/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package js

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/netext"
	"github.com/loadimpact/k6/stats"
	"github.com/oxtoacart/bpool"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestInitContextRequire(t *testing.T) ***REMOVED***
	t.Run("Modules", func(t *testing.T) ***REMOVED***
		t.Run("Nonexistent", func(t *testing.T) ***REMOVED***
			_, err := getSimpleBundle("/script.js", `import "k6/NONEXISTENT";`)
			assert.EqualError(t, err, "GoError: unknown builtin module: k6/NONEXISTENT")
		***REMOVED***)

		t.Run("k6", func(t *testing.T) ***REMOVED***
			b, err := getSimpleBundle("/script.js", `
					import k6 from "k6";
					export let _k6 = k6;
					export let dummy = "abc123";
					export default function() ***REMOVED******REMOVED***
			`)
			if !assert.NoError(t, err, "bundle error") ***REMOVED***
				return
			***REMOVED***

			bi, err := b.Instantiate()
			if !assert.NoError(t, err, "instance error") ***REMOVED***
				return
			***REMOVED***

			exports := bi.Runtime.Get("exports").ToObject(bi.Runtime)
			if assert.NotNil(t, exports) ***REMOVED***
				_, defaultOk := goja.AssertFunction(exports.Get("default"))
				assert.True(t, defaultOk, "default export is not a function")
				assert.Equal(t, "abc123", exports.Get("dummy").String())
			***REMOVED***

			k6 := bi.Runtime.Get("_k6").ToObject(bi.Runtime)
			if assert.NotNil(t, k6) ***REMOVED***
				_, groupOk := goja.AssertFunction(k6.Get("group"))
				assert.True(t, groupOk, "k6.group is not a function")
			***REMOVED***

			t.Run("group", func(t *testing.T) ***REMOVED***
				b, err := getSimpleBundle("/script.js", `
						import ***REMOVED*** group ***REMOVED*** from "k6";
						export let _group = group;
						export let dummy = "abc123";
						export default function() ***REMOVED******REMOVED***
				`)
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				bi, err := b.Instantiate()
				if !assert.NoError(t, err) ***REMOVED***
					return
				***REMOVED***

				exports := bi.Runtime.Get("exports").ToObject(bi.Runtime)
				if assert.NotNil(t, exports) ***REMOVED***
					_, defaultOk := goja.AssertFunction(exports.Get("default"))
					assert.True(t, defaultOk, "default export is not a function")
					assert.Equal(t, "abc123", exports.Get("dummy").String())
				***REMOVED***

				_, groupOk := goja.AssertFunction(exports.Get("_group"))
				assert.True(t, groupOk, "***REMOVED*** group ***REMOVED*** is not a function")
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)

	t.Run("Files", func(t *testing.T) ***REMOVED***
		t.Run("Nonexistent", func(t *testing.T) ***REMOVED***
			path := filepath.FromSlash("/nonexistent.js")
			_, err := getSimpleBundle("/script.js", `import "/nonexistent.js"; export default function() ***REMOVED******REMOVED***`)
			assert.EqualError(t, err, fmt.Sprintf("GoError: open %s: file does not exist", path))
		***REMOVED***)
		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			fs := afero.NewMemMapFs()
			assert.NoError(t, afero.WriteFile(fs, "/file.js", []byte***REMOVED***0x00***REMOVED***, 0755))
			_, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data:     []byte(`import "/file.js"; export default function() ***REMOVED******REMOVED***`),
			***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
			assert.Contains(t, err.Error(), "SyntaxError: /file.js: Unexpected character '\x00' (1:0)\n> 1 | \x00\n")
		***REMOVED***)
		t.Run("Error", func(t *testing.T) ***REMOVED***
			fs := afero.NewMemMapFs()
			assert.NoError(t, afero.WriteFile(fs, "/file.js", []byte(`throw new Error("aaaa")`), 0755))
			_, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data:     []byte(`import "/file.js"; export default function() ***REMOVED******REMOVED***`),
			***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
			assert.EqualError(t, err, "Error: aaaa at /file.js:2:7(4)")
		***REMOVED***)

		imports := map[string]struct ***REMOVED***
			LibPath    string
			ConstPaths map[string]string
		***REMOVED******REMOVED***
			"./lib.js": ***REMOVED***"/path/to/lib.js", map[string]string***REMOVED***
				"":               "",
				"./const.js":     "/path/to/const.js",
				"../const.js":    "/path/const.js",
				"./sub/const.js": "/path/to/sub/const.js",
			***REMOVED******REMOVED***,
			"../lib.js": ***REMOVED***"/path/lib.js", map[string]string***REMOVED***
				"":               "",
				"./const.js":     "/path/const.js",
				"../const.js":    "/const.js",
				"./sub/const.js": "/path/sub/const.js",
			***REMOVED******REMOVED***,
			"./dir/lib.js": ***REMOVED***"/path/to/dir/lib.js", map[string]string***REMOVED***
				"":               "",
				"./const.js":     "/path/to/dir/const.js",
				"../const.js":    "/path/to/const.js",
				"./sub/const.js": "/path/to/dir/sub/const.js",
			***REMOVED******REMOVED***,
			"/path/to/lib.js": ***REMOVED***"/path/to/lib.js", map[string]string***REMOVED***
				"":               "",
				"./const.js":     "/path/to/const.js",
				"../const.js":    "/path/const.js",
				"./sub/const.js": "/path/to/sub/const.js",
			***REMOVED******REMOVED***,
		***REMOVED***
		for libName, data := range imports ***REMOVED***
			t.Run("lib=\""+libName+"\"", func(t *testing.T) ***REMOVED***
				for constName, constPath := range data.ConstPaths ***REMOVED***
					name := "inline"
					if constName != "" ***REMOVED***
						name = "const=\"" + constName + "\""
					***REMOVED***
					t.Run(name, func(t *testing.T) ***REMOVED***
						fs := afero.NewMemMapFs()
						src := &lib.SourceData***REMOVED***
							Filename: `/path/to/script.js`,
							Data: []byte(fmt.Sprintf(`
								import fn from "%s";
								let v = fn();
								export default function() ***REMOVED***
								***REMOVED***;
							`, libName)),
						***REMOVED***

						jsLib := `export default function() ***REMOVED*** return 12345; ***REMOVED***`
						if constName != "" ***REMOVED***
							jsLib = fmt.Sprintf(
								`import ***REMOVED*** c ***REMOVED*** from "%s"; export default function() ***REMOVED*** return c; ***REMOVED***`,
								constName,
							)

							constsrc := `export let c = 12345;`
							assert.NoError(t, fs.MkdirAll(filepath.Dir(constPath), 0755))
							assert.NoError(t, afero.WriteFile(fs, constPath, []byte(constsrc), 0644))
						***REMOVED***

						assert.NoError(t, fs.MkdirAll(filepath.Dir(data.LibPath), 0755))
						assert.NoError(t, afero.WriteFile(fs, data.LibPath, []byte(jsLib), 0644))

						b, err := NewBundle(src, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
						if !assert.NoError(t, err) ***REMOVED***
							return
						***REMOVED***
						if constPath != "" ***REMOVED***
							assert.Contains(t, b.BaseInitContext.programs, constPath)
						***REMOVED***

						_, err = b.Instantiate()
						if !assert.NoError(t, err) ***REMOVED***
							return
						***REMOVED***
					***REMOVED***)
				***REMOVED***
			***REMOVED***)
		***REMOVED***

		t.Run("Isolation", func(t *testing.T) ***REMOVED***
			fs := afero.NewMemMapFs()
			assert.NoError(t, afero.WriteFile(fs, "/a.js", []byte(`const myvar = "a";`), 0644))
			assert.NoError(t, afero.WriteFile(fs, "/b.js", []byte(`const myvar = "b";`), 0644))
			b, err := NewBundle(&lib.SourceData***REMOVED***
				Filename: "/script.js",
				Data: []byte(`
				import "./a.js";
				import "./b.js";
				export default function() ***REMOVED***
					if (typeof myvar != "undefined") ***REMOVED***
						throw new Error("myvar is set in global scope");
					***REMOVED***
				***REMOVED***;
				`),
			***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***

			bi, err := b.Instantiate()
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			_, err = bi.Default(goja.Undefined())
			assert.NoError(t, err)
		***REMOVED***)
	***REMOVED***)
***REMOVED***

func createAndReadFile(t *testing.T, file string, content []byte, expectedLength int, binary bool) (*BundleInstance, error) ***REMOVED***
	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/path/to", 0755))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/"+file, []byte(content), 0644))

	afero.WriteFile(fs, "/path/to/"+file, []byte(content), 0644)
	binaryArg := ""
	if binary ***REMOVED***
		binaryArg = ",\"b\""
	***REMOVED***

	b, err := NewBundle(&lib.SourceData***REMOVED***
		Filename: "/path/to/script.js",
		Data: []byte(fmt.Sprintf(`
				export let data = open("/path/to/%s"%s);
				var expectedLength = %d;
				if (data.length != expectedLength) ***REMOVED***
					throw new Error("Length not equal, expected: " + expectedLength + ", actual: " + data.length);
				***REMOVED***
				export default function() ***REMOVED******REMOVED***
				`, file, binaryArg, expectedLength)),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)

	if !assert.NoError(t, err) ***REMOVED***
		return nil, err
	***REMOVED***

	bi, err := b.Instantiate()
	if !assert.NoError(t, err) ***REMOVED***
		return nil, err
	***REMOVED***
	return bi, nil
***REMOVED***

func TestInitContextOpen(t *testing.T) ***REMOVED***

	testCases := []struct ***REMOVED***
		content []byte
		file    string
		length  int
	***REMOVED******REMOVED***
		***REMOVED***[]byte("hello world!"), "ascii", 12***REMOVED***,
		***REMOVED***[]byte("?((¯°·._.• ţ€$ţɨɲǥ µɲɨȼ๏ď€ΣSЫ ɨɲ Ќ6 •._.·°¯))؟•"), "utf", 47***REMOVED***,
		***REMOVED***[]byte***REMOVED***044, 226, 130, 172***REMOVED***, "utf-8", 2***REMOVED***, // $€
		//***REMOVED***[]byte***REMOVED***00, 36, 32, 127***REMOVED***, "utf-16", 2***REMOVED***,   // $€
	***REMOVED***
	for _, tc := range testCases ***REMOVED***
		t.Run(tc.file, func(t *testing.T) ***REMOVED***
			bi, err := createAndReadFile(t, tc.file, tc.content, tc.length, false)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
			assert.Equal(t, string(tc.content), bi.Runtime.Get("data").Export())
		***REMOVED***)
	***REMOVED***

	t.Run("Binary", func(t *testing.T) ***REMOVED***
		bi, err := createAndReadFile(t, "/path/to/file.bin", []byte("hi!\x0f\xff\x01"), 6, true)
		if !assert.NoError(t, err) ***REMOVED***
			return
		***REMOVED***
		bytes := []byte***REMOVED***104, 105, 33, 15, 255, 1***REMOVED***
		assert.Equal(t, bytes, bi.Runtime.Get("data").Export())
	***REMOVED***)

	testdata := map[string]string***REMOVED***
		"Absolute": "/path/to/file",
		"Relative": "./file",
	***REMOVED***

	for name, loadPath := range testdata ***REMOVED***
		t.Run(name, func(t *testing.T) ***REMOVED***
			_, err := createAndReadFile(t, loadPath, []byte("content"), 7, false)
			if !assert.NoError(t, err) ***REMOVED***
				return
			***REMOVED***
		***REMOVED***)
	***REMOVED***

	t.Run("Nonexistent", func(t *testing.T) ***REMOVED***
		fs := afero.NewMemMapFs()
		path := filepath.FromSlash("/nonexistent.txt")
		_, err := NewBundle(&lib.SourceData***REMOVED***
			Filename: "/script.js",
			Data:     []byte(`open("/nonexistent.txt"); export default function() ***REMOVED******REMOVED***`),
		***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
		assert.EqualError(t, err, fmt.Sprintf("GoError: open %s: file does not exist", path))
	***REMOVED***)

***REMOVED***

func TestRequestWithBinaryFile(t *testing.T) ***REMOVED***
	t.Parallel()

	ch := make(chan bool, 1)

	h := func(w http.ResponseWriter, r *http.Request) ***REMOVED***
		defer func() ***REMOVED***
			ch <- true
		***REMOVED***()

		assert.NoError(t, r.ParseMultipartForm(32<<20))
		file, _, err := r.FormFile("file")
		assert.NoError(t, err)
		defer func() ***REMOVED***
			assert.NoError(t, file.Close())
		***REMOVED***()
		bytes := make([]byte, 3)
		_, err = file.Read(bytes)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hi!"), bytes)
		assert.Equal(t, "this is a standard form field", r.FormValue("field"))
	***REMOVED***

	srv := httptest.NewServer(http.HandlerFunc(h))
	defer srv.Close()

	fs := afero.NewMemMapFs()
	assert.NoError(t, fs.MkdirAll("/path/to", 0755))
	assert.NoError(t, afero.WriteFile(fs, "/path/to/file.bin", []byte("hi!"), 0644))

	b, err := NewBundle(&lib.SourceData***REMOVED***
		Filename: "/path/to/script.js",
		Data: []byte(fmt.Sprintf(`
			import http from "k6/http";
			let binFile = open("/path/to/file.bin", "b");
			export default function() ***REMOVED***
				var data = ***REMOVED***
					field: "this is a standard form field",
					file: http.file(binFile, "test.bin")
				***REMOVED***;
				var res = http.post("%s", data);
				return true;
			***REMOVED***
			`, srv.URL)),
	***REMOVED***, fs, lib.RuntimeOptions***REMOVED******REMOVED***)
	assert.NoError(t, err)

	bi, err := b.Instantiate()
	assert.NoError(t, err)

	root, err := lib.NewGroup("", nil)
	assert.NoError(t, err)

	logger := log.New()
	logger.Level = log.DebugLevel
	logger.Out = ioutil.Discard

	state := &lib.State***REMOVED***
		Options: lib.Options***REMOVED******REMOVED***,
		Logger:  logger,
		Group:   root,
		Transport: &http.Transport***REMOVED***
			DialContext: (netext.NewDialer(net.Dialer***REMOVED***
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
				DualStack: true,
			***REMOVED***)).DialContext,
		***REMOVED***,
		BPool:   bpool.NewBufferPool(1),
		Samples: make(chan stats.SampleContainer, 500),
	***REMOVED***

	ctx := context.Background()
	ctx = lib.WithState(ctx, state)
	ctx = common.WithRuntime(ctx, bi.Runtime)
	*bi.Context = ctx

	v, err := bi.Default(goja.Undefined())
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v.Export())

	<-ch
***REMOVED***
