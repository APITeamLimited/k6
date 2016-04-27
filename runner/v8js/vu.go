package v8js

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/ry/v8worker"
	"reflect"
)

func (vu *VUContext) RegisterModules(w *v8worker.Worker) error ***REMOVED***
	vu.mods = map[string]Module***REMOVED***
		"global": Module***REMOVED***
			"sleep": vu.Sleep,
		***REMOVED***,
		"console": Module***REMOVED***
			"log": vu.ConsoleLog,
		***REMOVED***,
		"http": Module***REMOVED***
			"get": vu.HTTPGet,
		***REMOVED***,
	***REMOVED***

	for modname, mod := range vu.mods ***REMOVED***
		jsMod := fmt.Sprintf(`
		speedboat._modules["%s"] = ***REMOVED******REMOVED***;
		`, modname)
		for name, mem := range mod ***REMOVED***
			t := reflect.TypeOf(mem)

			if t.Kind() != reflect.Func ***REMOVED***
				return errors.New("Not a function: " + modname + "." + name)
			***REMOVED***

			jsFn := fmt.Sprintf(`speedboat._modules["%s"]["%s"] = function() ***REMOVED***
				var args = [];
			`, modname, name)

			numArgs := t.NumIn()
			if !t.IsVariadic() ***REMOVED***
				jsFn += fmt.Sprintf(`
					if (arguments.length != %d) ***REMOVED***
						throw new Error("wrong number of arguments");
					***REMOVED***
				`, t.NumIn())
			***REMOVED*** else ***REMOVED***
				numArgs--
			***REMOVED***

			for i := 0; i < numArgs; i++ ***REMOVED***
				aT := t.In(i)
				jsFn += fmt.Sprintf("args.push(speedboat._require.%s(arguments[%d]));", aT.Kind().String(), i)
			***REMOVED***
			if t.IsVariadic() ***REMOVED***
				varArg := t.In(numArgs)
				eT := varArg.Elem()
				jsFn += fmt.Sprintf(`
					for (var i = %d; i < arguments.length; i++) ***REMOVED***
						args.push(speedboat._require.%s(arguments[i]));
					***REMOVED***
				`, numArgs, eT.Kind().String())
			***REMOVED***

			jsFn += `
				$sendSync(JSON.stringify(***REMOVED***
					m: '` + modname + `',
					f: '` + name + `',
					a: args,
				***REMOVED***));
			***REMOVED***`
			jsMod += "\n\n" + jsFn
		***REMOVED***

		if err := w.Load("module:"+modname, jsMod); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Make functions in the "global" module global, preimport console
	makeGlobals := `
	for (key in speedboat._modules['global']) ***REMOVED***
		eval(key + " = speedboat._modules['global']['" + key + "'];");
	***REMOVED***
	var console = speedboat._modules['console'];
	`
	if err := w.Load("internal:preload", makeGlobals); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (vu *VUContext) Recv(raw string) ***REMOVED***
***REMOVED***

func (vu *VUContext) RecvSync(raw string) string ***REMOVED***
	call := struct ***REMOVED***
		Mod  string        `json:"m"`
		Fn   string        `json:"f"`
		Args []interface***REMOVED******REMOVED*** `json:"a"`
	***REMOVED******REMOVED******REMOVED***
	if err := json.Unmarshal([]byte(raw), &call); err != nil ***REMOVED***
		log.WithError(err).Error("Malformed call")
		return ""
	***REMOVED***
	log.WithFields(log.Fields***REMOVED***
		"mod":  call.Mod,
		"fn":   call.Fn,
		"args": call.Args,
	***REMOVED***).Debug("Sync call")

	mod, ok := vu.mods[call.Mod]
	if !ok ***REMOVED***
		log.WithField("mod", call.Mod).Error("Unknown module")
		return ""
	***REMOVED***

	fn, ok := mod[call.Fn]
	if !ok ***REMOVED***
		log.WithFields(log.Fields***REMOVED***
			"mod": call.Mod,
			"fn":  call.Fn,
		***REMOVED***).Error("Unknown function")
		return ""
	***REMOVED***

	args := make([]reflect.Value, len(call.Args))
	for i, arg := range call.Args ***REMOVED***
		args[i] = reflect.ValueOf(arg)
	***REMOVED***

	defer func() ***REMOVED***
		if err := recover(); err != nil ***REMOVED***
			log.WithField("error", err).Error("Go call panicked")
		***REMOVED***
	***REMOVED***()
	fnV := reflect.ValueOf(fn)
	log.WithField("T", fnV.Type().String()).Debug("Function")
	reflect.ValueOf(fn).Call(args)

	return ""
***REMOVED***
