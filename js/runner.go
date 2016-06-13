package js

import (
	"github.com/GeertJohan/go.rice"
	log "github.com/Sirupsen/logrus"
	"github.com/loadimpact/speedboat"
	"golang.org/x/net/context"
	"gopkg.in/olebedev/go-duktape.v2"
	"os"
)

type Runner struct ***REMOVED***
	Filename string
	Source   string
***REMOVED***

func New(filename, src string) *Runner ***REMOVED***
	return &Runner***REMOVED***
		Filename: filename,
		Source:   src,
	***REMOVED***
***REMOVED***

func (r *Runner) RunVU(ctx context.Context, t speedboat.Test, id int) ***REMOVED***
	js := duktape.New()
	setupGlobalObject(js, t, id)
	bridgeAPI(js, contextForAPI(ctx))

	if err := putScript(js, r.Filename, r.Source); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't compile script")
		return
	***REMOVED***

	vendor, err := rice.FindBox("vendor")
	if err != nil ***REMOVED***
		log.WithError(err).Error("Script vendor files missing; try `git submodule update --init`")
		return
	***REMOVED***
	vendorFiles := []string***REMOVED***"lodash/dist/lodash.min.js"***REMOVED***
	for _, filename := range vendorFiles ***REMOVED***
		src, err := vendor.String(filename)
		if err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't read dependent script")
			return
		***REMOVED***
		if err := loadScript(js, filename, src); err != nil ***REMOVED***
			log.WithError(err).Error("Couldn't load dependency")
			return
		***REMOVED***
	***REMOVED***

	lib, err := rice.FindBox("lib")
	if err != nil ***REMOVED***
		log.WithError(err).Error("Script support files absent")
		return
	***REMOVED***
	if err := lib.Walk("/", func(path string, info os.FileInfo, err error) error ***REMOVED***
		if info.IsDir() ***REMOVED***
			return nil
		***REMOVED***
		src, err := lib.String(path)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		return loadScript(js, path, src)
	***REMOVED***); err != nil ***REMOVED***
		log.WithError(err).Error("Couldn't load support file")
		return
	***REMOVED***

	js.PushGlobalObject()
	js.PushString(scriptProp)
	for ***REMOVED***
		js.DupTop()
		if js.PcallProp(-3, 0) != duktape.ErrNone ***REMOVED***
			err := getJSError(js)
			log.WithFields(log.Fields***REMOVED***
				"file":  err.Filename,
				"line":  err.Line,
				"error": err.Message,
			***REMOVED***).Error("Script Error")
		***REMOVED***
		js.Pop()

		select ***REMOVED***
		case <-ctx.Done():
			return
		default:
		***REMOVED***
	***REMOVED***
***REMOVED***

func bridgeAPI(js *duktape.Context, ctx context.Context) ***REMOVED***
	api := map[string]map[string]APIFunc***REMOVED***
		"http": map[string]APIFunc***REMOVED***
			"request":            apiHTTPRequest,
			"setMaxConnsPerHost": apiHTTPSetMaxConnsPerHost,
		***REMOVED***,
		"log": map[string]APIFunc***REMOVED***
			"log": apiLogLog,
		***REMOVED***,
		"test": map[string]APIFunc***REMOVED***
			"abort": apiTestAbort,
		***REMOVED***,
	***REMOVED***
	global := map[string]APIFunc***REMOVED***
		"sleep": apiSleep,
	***REMOVED***

	js.PushGlobalObject()
	defer js.Pop()

	for fname, fn := range global ***REMOVED***
		fn := fn
		js.PushGoFunction(func(js *duktape.Context) int ***REMOVED***
			return fn(js, ctx)
		***REMOVED***)
		js.PutPropString(-2, fname)
	***REMOVED***

	js.GetPropString(-1, "__modules__")
	defer js.Pop()

	for modname, mod := range api ***REMOVED***
		js.PushObject()
		for fname, fn := range mod ***REMOVED***
			fn := fn
			js.PushGoFunction(func(js *duktape.Context) int ***REMOVED***
				return fn(js, ctx)
			***REMOVED***)
			js.PutPropString(-2, fname)
		***REMOVED***
		js.PutPropString(-2, modname)
	***REMOVED***
***REMOVED***
