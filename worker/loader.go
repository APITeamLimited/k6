// Custom loader for execution from redis

package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/APITeamLimited/k6-worker/errext"
	"github.com/APITeamLimited/k6-worker/errext/exitcodes"
	"github.com/APITeamLimited/k6-worker/js"
	"github.com/APITeamLimited/k6-worker/lib"
	"github.com/APITeamLimited/k6-worker/lib/executor"
	"github.com/APITeamLimited/k6-worker/lib/types"
	"github.com/APITeamLimited/k6-worker/loader"
	"github.com/APITeamLimited/k6-worker/metrics"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"gopkg.in/guregu/null.v3"
)

func loadAndConfigureTest(
	gs *globalState,
	job map[string]string,
	workerInfo *lib.WorkerInfo,
) (*workerLoadedAndConfiguredTest, error) {
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

	test := &workerLoadedTest{
		sourceRootPath: sourceRootPath,
		source:         source,
		fs:             gs.fs,
		pwd:            "",
		fileSystems:    filesystems,
		preInitState:   preInitState,
	}

	gs.logger.Debugf("Initializing k6 runner for '%s' (%s)...", sourceRootPath)
	if err := test.initializeFirstRunner(gs, workerInfo); err != nil {
		return nil, fmt.Errorf("could not initialize '%s': %w", sourceRootPath, err)
	}
	gs.logger.Debug("Runner successfully initialized!")

	return test.consolidateDeriveAndValidateConfig(gs, job)
}

func (lt *workerLoadedTest) initializeFirstRunner(gs *globalState, workerInfo *lib.WorkerInfo) error {
	testPath := lt.source.URL.String()
	logger := gs.logger.WithField("test_path", testPath)

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
		lt.preInitState.KeyLogger = &consoleWriter{
			ctx:      gs.ctx,
			client:   workerInfo.Client,
			jobId:    workerInfo.JobId,
			workerId: workerInfo.WorkerId,
		}
	}

	runner, err := js.New(lt.preInitState, lt.source, lt.fileSystems, workerInfo)
	// TODO: should we use common.UnwrapGojaInterruptedError() here?
	if err != nil {
		return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
	}
	lt.initRunner = runner
	return nil

}

func (lt *workerLoadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, job map[string]string,
) (*workerLoadedAndConfiguredTest, error) {

	// Options have already been determined by the orchestrator

	var redisOptions = lib.Options{}

	if job["options"] == "" {
		return nil, fmt.Errorf("options not found on job, this is probably a bug")
	}

	err := json.Unmarshal([]byte(job["options"]), &redisOptions)
	if err != nil {
		return nil, fmt.Errorf("could not parse options: %w", err)
	}

	consolidatedConfig := getConsolidatedConfig(redisOptions)

	// TODO: get other config sources eg

	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.preInitState.RuntimeOptions.NoThresholds.Bool {
		for metricName, thresholdsDefinition := range consolidatedConfig.Options.Thresholds {
			err := thresholdsDefinition.Parse()
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

	return &workerLoadedAndConfiguredTest{
		workerLoadedTest:   lt,
		consolidatedConfig: consolidatedConfig,
		derivedConfig:      derivedConfig,
	}, nil
}

func getConsolidatedConfig(parsedOptions lib.Options) Config {
	consolidatedConfig := Config{
		Options: parsedOptions,
	}

	consolidatedConfig = applyDefault(consolidatedConfig)

	return consolidatedConfig
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

func (lct *workerLoadedAndConfiguredTest) buildTestRunState(
	configToReinject lib.Options,
) (*lib.TestRunState, error) {
	// This might be the full derived or just the consolidated options
	if err := lct.initRunner.SetOptions(configToReinject); err != nil {
		return nil, err
	}

	// TODO: init atlas root worker, etc.

	return &lib.TestRunState{
		TestPreInitState: lct.preInitState,
		Runner:           lct.initRunner,
		Options:          lct.derivedConfig.Options, // we will always run with the derived options
	}, nil
}

func applyDefault(conf Config) Config {
	if conf.SystemTags == nil {
		conf.SystemTags = &metrics.DefaultSystemTagSet
	}
	if conf.SummaryTrendStats == nil {
		conf.SummaryTrendStats = lib.DefaultSummaryTrendStats
	}
	defDNS := types.DefaultDNSConfig()
	if !conf.DNS.TTL.Valid {
		conf.DNS.TTL = defDNS.TTL
	}
	if !conf.DNS.Select.Valid {
		conf.DNS.Select = defDNS.Select
	}
	if !conf.DNS.Policy.Valid {
		conf.DNS.Policy = defDNS.Policy
	}
	if !conf.SetupTimeout.Valid {
		conf.SetupTimeout.Duration = types.Duration(60 * time.Second)
	}
	if !conf.TeardownTimeout.Valid {
		conf.TeardownTimeout.Duration = types.Duration(60 * time.Second)
	}
	return conf
}
