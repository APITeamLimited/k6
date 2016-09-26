package js

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSleep(t *testing.T) ***REMOVED***
	start := time.Now()
	JSAPI***REMOVED******REMOVED***.Sleep(0.2)
	assert.True(t, time.Since(start) > 200*time.Millisecond)
	assert.True(t, time.Since(start) < 1*time.Second)
***REMOVED***

func TestDoGroup(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** group ***REMOVED*** from "speedboat";
	export default function() ***REMOVED***
		group("test", fn);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "test", vu.group.Name)
	***REMOVED***)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoGroupNested(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** group ***REMOVED*** from "speedboat";
	export default function() ***REMOVED***
		group("outer", function() ***REMOVED***
			group("inner", fn);
		***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	vu.vm.Set("fn", func() ***REMOVED***
		assert.Equal(t, "inner", vu.group.Name)
		assert.Equal(t, "outer", vu.group.Parent.Name)
	***REMOVED***)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)
***REMOVED***

func TestDoTest(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** test ***REMOVED*** from "speedboat";
	export default function() ***REMOVED***
		test(3, ***REMOVED*** "v === 3": (v) => v === 3 ***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)

	if !assert.Len(t, r.Tests, 1) ***REMOVED***
		return
	***REMOVED***
	ts := r.Tests[0]
	assert.Equal(t, "v === 3", ts.Name)
	assert.Equal(t, r.DefaultGroup, ts.Group)
	assert.Equal(t, int64(1), ts.Passes)
	assert.Equal(t, int64(0), ts.Fails)
***REMOVED***

func TestTestInGroup(t *testing.T) ***REMOVED***
	r, err := newSnippetRunner(`
	import ***REMOVED*** group, test ***REMOVED*** from "speedboat";
	export default function() ***REMOVED***
		group("group", function() ***REMOVED***
			test(3, ***REMOVED*** "v === 3": (v) => v === 3 ***REMOVED***);
		***REMOVED***);
	***REMOVED***`)
	assert.NoError(t, err)

	vu_, err := r.NewVU()
	assert.NoError(t, err)
	vu := vu_.(*VU)

	_, err = vu.RunOnce(context.Background())
	assert.NoError(t, err)

	assert.Len(t, r.Groups, 2)
	g := r.Groups[1]
	assert.Equal(t, "group", g.Name)

	assert.Len(t, r.Tests, 1)
	ts := r.Tests[0]
	assert.Equal(t, "v === 3", ts.Name)
	assert.Equal(t, g, ts.Group)
	assert.Equal(t, int64(1), ts.Passes)
	assert.Equal(t, int64(0), ts.Fails)
***REMOVED***
