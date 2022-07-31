package cmd

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.k6.io/k6/errext"
	"go.k6.io/k6/errext/exitcodes"
	"go.k6.io/k6/js"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/loader"
	"go.k6.io/k6/metrics"
)

const (
	testTypeJS      = "js"
	testTypeArchive = "archive"
)

// loadedTest contains all of data, details and dependencies of a fully-loaded
// and configured k6 test.
type loadedTest struct ***REMOVED***
	sourceRootPath string // contains the raw string the user supplied
	pwd            string
	source         *loader.SourceData
	fs             afero.Fs
	fileSystems    map[string]afero.Fs
	runtimeState   *lib.RuntimeState
	initRunner     lib.Runner // TODO: rename to something more appropriate

	// Only set if cliConfigGetter is supplied to loadAndConfigureTest() or if
	// consolidateDeriveAndValidateConfig() is manually called.
	consolidatedConfig Config
	derivedConfig      Config
***REMOVED***

func loadAndConfigureTest(
	gs *globalState, cmd *cobra.Command, args []string,
	// supply this if you want the test config consolidated and validated
	cliConfigGetter func(flags *pflag.FlagSet) (Config, error), // TODO: obviate
) (*loadedTest, error) ***REMOVED***
	if len(args) < 1 ***REMOVED***
		return nil, fmt.Errorf("k6 needs at least one argument to load the test")
	***REMOVED***

	sourceRootPath := args[0]
	gs.logger.Debugf("Resolving and reading test '%s'...", sourceRootPath)
	src, fileSystems, pwd, err := readSource(gs, sourceRootPath)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resolvedPath := src.URL.String()
	gs.logger.Debugf(
		"'%s' resolved to '%s' and successfully loaded %d bytes!",
		sourceRootPath, resolvedPath, len(src.Data),
	)

	gs.logger.Debugf("Gathering k6 runtime options...")
	runtimeOptions, err := getRuntimeOptions(cmd.Flags(), gs.envVars)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	registry := metrics.NewRegistry()
	state := &lib.RuntimeState***REMOVED***
		Logger:         gs.logger,
		RuntimeOptions: runtimeOptions,
		Registry:       registry,
		BuiltinMetrics: metrics.RegisterBuiltinMetrics(registry),
	***REMOVED***

	test := &loadedTest***REMOVED***
		pwd:            pwd,
		sourceRootPath: sourceRootPath,
		source:         src,
		fs:             gs.fs,
		fileSystems:    fileSystems,
		runtimeState:   state,
	***REMOVED***

	gs.logger.Debugf("Initializing k6 runner for '%s' (%s)...", sourceRootPath, resolvedPath)
	if err := test.initializeFirstRunner(gs); err != nil ***REMOVED***
		return nil, fmt.Errorf("could not initialize '%s': %w", sourceRootPath, err)
	***REMOVED***
	gs.logger.Debug("Runner successfully initialized!")

	if cliConfigGetter != nil ***REMOVED***
		if err := test.consolidateDeriveAndValidateConfig(gs, cmd, cliConfigGetter); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	return test, nil
***REMOVED***

func (lt *loadedTest) initializeFirstRunner(gs *globalState) error ***REMOVED***
	testPath := lt.source.URL.String()
	logger := gs.logger.WithField("test_path", testPath)

	testType := lt.runtimeState.RuntimeOptions.TestType.String
	if testType == "" ***REMOVED***
		logger.Debug("Detecting test type for...")
		testType = detectTestType(lt.source.Data)
	***REMOVED***

	if lt.runtimeState.RuntimeOptions.KeyWriter.Valid ***REMOVED***
		logger.Warnf("SSLKEYLOGFILE was specified, logging TLS connection keys to '%s'...",
			lt.runtimeState.RuntimeOptions.KeyWriter.String)
		keylogFilename := lt.runtimeState.RuntimeOptions.KeyWriter.String
		// if path is absolute - no point doing anything
		if !filepath.IsAbs(keylogFilename) ***REMOVED***
			// filepath.Abs could be used but it will get the pwd from `os` package instead of what is in lt.pwd
			// this is against our general approach of not using `os` directly and makes testing harder
			keylogFilename = filepath.Join(lt.pwd, keylogFilename)
		***REMOVED***
		f, err := lt.fs.OpenFile(keylogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
		if err != nil ***REMOVED***
			return fmt.Errorf("couldn't get absolute path for keylog file: %w", err)
		***REMOVED***
		lt.runtimeState.KeyLogger = &syncWriteCloser***REMOVED***w: f***REMOVED***
	***REMOVED***
	switch testType ***REMOVED***
	case testTypeJS:
		logger.Debug("Trying to load as a JS test...")
		runner, err := js.New(lt.runtimeState, lt.source, lt.fileSystems)
		// TODO: should we use common.UnwrapGojaInterruptedError() here?
		if err != nil ***REMOVED***
			return fmt.Errorf("could not load JS test '%s': %w", testPath, err)
		***REMOVED***
		lt.initRunner = runner
		return nil

	case testTypeArchive:
		logger.Debug("Trying to load test as an archive bundle...")

		var arc *lib.Archive
		arc, err := lib.ReadArchive(bytes.NewReader(lt.source.Data))
		if err != nil ***REMOVED***
			return fmt.Errorf("could not load test archive bundle '%s': %w", testPath, err)
		***REMOVED***
		logger.Debugf("Loaded test as an archive bundle with type '%s'!", arc.Type)

		switch arc.Type ***REMOVED***
		case testTypeJS:
			logger.Debug("Evaluating JS from archive bundle...")
			lt.initRunner, err = js.NewFromArchive(lt.runtimeState, arc)
			if err != nil ***REMOVED***
				return fmt.Errorf("could not load JS from test archive bundle '%s': %w", testPath, err)
			***REMOVED***
			return nil
		default:
			return fmt.Errorf("archive '%s' has an unsupported test type '%s'", testPath, arc.Type)
		***REMOVED***
	default:
		return fmt.Errorf("unknown or unspecified test type '%s' for '%s'", testType, testPath)
	***REMOVED***
***REMOVED***

// readSource is a small wrapper around loader.ReadSource returning
// result of the load and filesystems map
func readSource(globalState *globalState, filename string) (*loader.SourceData, map[string]afero.Fs, string, error) ***REMOVED***
	pwd, err := globalState.getwd()
	if err != nil ***REMOVED***
		return nil, nil, "", err
	***REMOVED***

	filesystems := loader.CreateFilesystems(globalState.fs)
	src, err := loader.ReadSource(globalState.logger, filename, pwd, filesystems, globalState.stdIn)
	return src, filesystems, pwd, err
***REMOVED***

func detectTestType(data []byte) string ***REMOVED***
	if _, err := tar.NewReader(bytes.NewReader(data)).Next(); err == nil ***REMOVED***
		return testTypeArchive
	***REMOVED***
	return testTypeJS
***REMOVED***

func (lt *loadedTest) consolidateDeriveAndValidateConfig(
	gs *globalState, cmd *cobra.Command,
	cliConfGetter func(flags *pflag.FlagSet) (Config, error), // TODO: obviate
) error ***REMOVED***
	var cliConfig Config
	if cliConfGetter != nil ***REMOVED***
		gs.logger.Debug("Parsing CLI flags...")
		var err error
		cliConfig, err = cliConfGetter(cmd.Flags())
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	gs.logger.Debug("Consolidating config layers...")
	consolidatedConfig, err := getConsolidatedConfig(gs, cliConfig, lt.initRunner.GetOptions())
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	gs.logger.Debug("Parsing thresholds and validating config...")
	// Parse the thresholds, only if the --no-threshold flag is not set.
	// If parsing the threshold expressions failed, consider it as an
	// invalid configuration error.
	if !lt.runtimeState.RuntimeOptions.NoThresholds.Bool ***REMOVED***
		for metricName, thresholdsDefinition := range consolidatedConfig.Options.Thresholds ***REMOVED***
			err = thresholdsDefinition.Parse()
			if err != nil ***REMOVED***
				return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***

			err = thresholdsDefinition.Validate(metricName, lt.runtimeState.Registry)
			if err != nil ***REMOVED***
				return errext.WithExitCodeIfNone(err, exitcodes.InvalidConfig)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	derivedConfig, err := deriveAndValidateConfig(consolidatedConfig, lt.initRunner.IsExecutable, gs.logger)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	lt.consolidatedConfig = consolidatedConfig
	lt.derivedConfig = derivedConfig

	return nil
***REMOVED***

type syncWriteCloser struct ***REMOVED***
	w io.WriteCloser
	m sync.Mutex
***REMOVED***

func (cw *syncWriteCloser) Write(b []byte) (int, error) ***REMOVED***
	cw.m.Lock()
	defer cw.m.Unlock()
	return cw.w.Write(b)
***REMOVED***

func (cw *syncWriteCloser) Close() error ***REMOVED***
	cw.m.Lock()
	defer cw.m.Unlock()
	return cw.w.Close()
***REMOVED***
