// Custom loader for execution from redis

package node

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/lib/executor"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
	"gopkg.in/guregu/null.v3"
)

func loadAndConfigureTest(
	gs *globalState,
	job map[string]string,
) (*nodeLoadedAndConfiguredTest, error) {
	test, err := loadTest(gs, job)
	if err != nil {
		return nil, err
	}

	return test.consolidateDeriveAndValidateConfig(gs, job)
}

func loadTest(gs *globalState, job map[string]string) (*nodeLoadedTest, error) {
	sourceName := job["sourceName"]

	if sourceName == "" {
		return nil, fmt.Errorf("sourceName not found on job, this is probably a bug")
	}

	stringSource := job["source"]

	if stringSource == "" {
		return nil, fmt.Errorf("source not found on job, this is probably a bug")
	}

	source := &loader.SourceData{
		URL:  &url.URL{Path: sourceName},
		Data: []byte(stringSource),
	}

	filesystems := map[string]afero.Fs{
		"file": afero.NewMemMapFs(),
	}

	f, err := afero.TempFile(filesystems["file"], "", sourceName)
	if err != nil {
		return nil, err
	}

	_, err = f.Write([]byte(stringSource))
	if err != nil {
		return nil, err
	}

	// Store the source in the filesystem
	sourceRootPath := sourceName

	// For now runtime options are constant for all tests
	// TODO: make this configurable
	runtimeOptions := lib.RuntimeOptions{
		TestType:             null.StringFrom(testTypeJS),
		IncludeSystemEnvVars: null.BoolFrom(false),
		CompatibilityMode:    null.StringFrom("extended"),
		NoThresholds:         null.BoolFrom(false),
		NoSummary:            null.BoolFrom(false),
		SummaryExport:        null.StringFrom(""),
		Env:                  make(map[string]string),
	}

	registry := metrics.NewRegistry()

	preInitState := &lib.TestPreInitState{
		// These gs will need to be changed as on the cloud
		Logger:         gs.logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
	}

	test := &nodeLoadedTest{
		sourceRootPath: sourceRootPath,
		source:         source,
		fs:             gs.fs,
		pwd:            "",
		fileSystems:    filesystems,
		preInitState:   preInitState,
	}

	gs.logger.Debugf("Initializing k6 runner for '%s' (%s)...", sourceRootPath)
	if err := test.initializeFirstRunner(gs); err != nil {
		return nil, fmt.Errorf("could not initialize '%s': %w", sourceRootPath, err)
	}
	gs.logger.Debug("Runner successfully initialized!")
	return test, nil
}

func detectTestType(data []byte) string {
	if _, err := tar.NewReader(bytes.NewReader(data)).Next(); err == nil {
		return testTypeArchive
	}
	return testTypeJS
}

func (lt *nodeLoadedTest) initializeFirstRunner(gs *globalState) error {
	testPath := lt.source.URL.String()
	logger := gs.logger.WithField("test_path", testPath)

	testType := lt.preInitState.RuntimeOptions.TestType.String

	if testType == "" {
		logger.Debug("Detecting test type for...")
		testType = detectTestType(lt.source.Data)
	}

	if lt.preInitState.RuntimeOptions.KeyWriter.Valid {
		logger.Warnf("SSLKEYLOGFILE was specified, logging TLS connection keys to '%s'...",
			lt.preInitState.RuntimeOptions.KeyWriter.String)
		keylogFilename := lt.preInitState.RuntimeOptions.KeyWriter.String
		// if path is absolute - no point doing anything
		if !filepath.IsAbs(keylogFilename) {
			// filepath.Abs could be used but it will get the pwd from `os` package instead of what is in lt.pwd
			// this is against our general approach of not using `os` directly and makes testing harder
			keylogFilename = filepath.Join(lt.pwd, keylogFilename)
		}
		f, err := lt.fs.OpenFile(keylogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil {
			return fmt.Errorf("couldn't get absolute path for keylog file: %w", err)
		}
		lt.keyLogger = f
		lt.preInitState.KeyLogger = &syncWriter{w: f}
	}

	switch testType {
	case testTypeJS:
		logger.Debug("Trying to load as a JS test...")
		runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems)
		// TODO: should we use common.UnwrapGojaInterruptedError() here?
		if err != nil {
			return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
		}
		lt.initRunner = runner
		return nil

	case testTypeArchive:
		logger.Debug("Trying to load test as an archive bundle...")

		var arc *lib.Archive
		arc, err := lib.ReadArchive(bytes.NewReader(lt.source.Data))
		if err != nil {
			return fmt.Errorf("could not load test archive bundle '%s': %w", testPath, err)
		}
		logger.Debugf("Loaded test as an archive bundle with type '%s'!", arc.Type)

		switch arc.Type {
		case testTypeJS:
			logger.Debug("Evaluating JS from archive bundle...")
			lt.initRunner, err = js.NewFromArchive(lt.preInitState, arc)
			if err != nil {
				return fmt.Errorf("could not load JS from test archive bundle '%s': %w", testPath, err)
			}
			return nil
		default:
			return fmt.Errorf("archive '%s' has an unsupported test type '%s'", testPath, arc.Type)
		}
	default:
		return fmt.Errorf("unknown or unspecified test type '%s' for '%s'", testType, testPath)
	}
}

func (lt *nodeLoadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, job map[string]string,
) (*nodeLoadedAndConfiguredTest, error) {

	// TODO: implement consolidateDeriveAndValidateConfig behavior#

	var parsedOptions lib.Options
	err := json.Unmarshal([]byte(job["options"]), &parsedOptions)

	if err != nil {
		return nil, fmt.Errorf("could not parse options: %w", err)
	}

	consolidatedConfig := Config{
		Options: parsedOptions,
	}

	gs.logger.Debug("Parsing thresholds and validating config...")
	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.preInitState.RuntimeOptions.NoThresholds.Bool {
		for metricName, thresholdsDefinition := range consolidatedConfig.Options.Thresholds {
			err = thresholdsDefinition.Parse()
			if err != nil {
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			}

			err = thresholdsDefinition.Validate(metricName, lt.preInitState.Registry)
			if err != nil {
				return nil, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			}
		}
	}

	derivedConfig, err := deriveAndValidateConfig(consolidatedConfig, lt.initRunner.IsExecutable, gs.logger)
	if err != nil {
		return nil, err
	}

	return &nodeLoadedAndConfiguredTest{
		nodeLoadedTest:     lt,
		consolidatedConfig: consolidatedConfig,
		derivedConfig:      derivedConfig,
	}, nil
}

func deriveAndValidateConfig(
	conf Config, isExecutable func(string) bool, logger logrus.FieldLogger,
) (result Config, err error) {
	result = conf
	result.Options, err = executor.DeriveScenariosFromShortcuts(conf.Options, logger)
	if err == nil {
		err = validateConfig(result, isExecutable)
	}
	return result, errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
}

func validateConfig(conf Config, isExecutable func(string) bool) error {
	errList := conf.Validate()

	for _, ec := range conf.Scenarios {
		if err := validateScenarioConfig(ec, isExecutable); err != nil {
			errList = append(errList, err)
		}
	}

	return consolidateErrorMessage(errList, "There were problems with the specified script configuration:")
}

func consolidateErrorMessage(errList []error, title string) error {
	if len(errList) == 0 {
		return nil
	}

	errMsgParts := []string{title}
	for _, err := range errList {
		errMsgParts = append(errMsgParts, fmt.Sprintf("\t- %s", err.Error()))
	}

	return errors.New(strings.Join(errMsgParts, "\n"))
}

func validateScenarioConfig(conf lib.ExecutorConfig, isExecutable func(string) bool) error {
	execFn := conf.GetExec()
	if !isExecutable(execFn) {
		return fmt.Errorf("executor %s: function '%s' not found in exports", conf.GetName(), execFn)
	}
	return nil
}

func (lct *nodeLoadedAndConfiguredTest) buildTestRunState(
	configToReinject lib.Options,
) (*lib.TestRunState, error) {
	// This might be the full derived or just the consodlidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil {
		return nil, err
	}

	// TODO: init atlas root node, etc.

	return &lib.TestRunState{
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	}, nil
}
