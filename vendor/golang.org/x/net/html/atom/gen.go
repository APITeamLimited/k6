// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

//go:generate go run gen.go
//go:generate go run gen.go -test

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// identifier converts s to a Go exported identifier.
// It converts "div" to "Div" and "accept-charset" to "AcceptCharset".
func identifier(s string) string ***REMOVED***
	b := make([]byte, 0, len(s))
	cap := true
	for _, c := range s ***REMOVED***
		if c == '-' ***REMOVED***
			cap = true
			continue
		***REMOVED***
		if cap && 'a' <= c && c <= 'z' ***REMOVED***
			c -= 'a' - 'A'
		***REMOVED***
		cap = false
		b = append(b, byte(c))
	***REMOVED***
	return string(b)
***REMOVED***

var test = flag.Bool("test", false, "generate table_test.go")

func genFile(name string, buf *bytes.Buffer) ***REMOVED***
	b, err := format.Source(buf.Bytes())
	if err != nil ***REMOVED***
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	***REMOVED***
	if err := ioutil.WriteFile(name, b, 0644); err != nil ***REMOVED***
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	***REMOVED***
***REMOVED***

func main() ***REMOVED***
	flag.Parse()

	var all []string
	all = append(all, elements...)
	all = append(all, attributes...)
	all = append(all, eventHandlers...)
	all = append(all, extra...)
	sort.Strings(all)

	// uniq - lists have dups
	w := 0
	for _, s := range all ***REMOVED***
		if w == 0 || all[w-1] != s ***REMOVED***
			all[w] = s
			w++
		***REMOVED***
	***REMOVED***
	all = all[:w]

	if *test ***REMOVED***
		var buf bytes.Buffer
		fmt.Fprintln(&buf, "// Code generated by go generate gen.go; DO NOT EDIT.\n")
		fmt.Fprintln(&buf, "//go:generate go run gen.go -test\n")
		fmt.Fprintln(&buf, "package atom\n")
		fmt.Fprintln(&buf, "var testAtomList = []string***REMOVED***")
		for _, s := range all ***REMOVED***
			fmt.Fprintf(&buf, "\t%q,\n", s)
		***REMOVED***
		fmt.Fprintln(&buf, "***REMOVED***")

		genFile("table_test.go", &buf)
		return
	***REMOVED***

	// Find hash that minimizes table size.
	var best *table
	for i := 0; i < 1000000; i++ ***REMOVED***
		if best != nil && 1<<(best.k-1) < len(all) ***REMOVED***
			break
		***REMOVED***
		h := rand.Uint32()
		for k := uint(0); k <= 16; k++ ***REMOVED***
			if best != nil && k >= best.k ***REMOVED***
				break
			***REMOVED***
			var t table
			if t.init(h, k, all) ***REMOVED***
				best = &t
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if best == nil ***REMOVED***
		fmt.Fprintf(os.Stderr, "failed to construct string table\n")
		os.Exit(1)
	***REMOVED***

	// Lay out strings, using overlaps when possible.
	layout := append([]string***REMOVED******REMOVED***, all...)

	// Remove strings that are substrings of other strings
	for changed := true; changed; ***REMOVED***
		changed = false
		for i, s := range layout ***REMOVED***
			if s == "" ***REMOVED***
				continue
			***REMOVED***
			for j, t := range layout ***REMOVED***
				if i != j && t != "" && strings.Contains(s, t) ***REMOVED***
					changed = true
					layout[j] = ""
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Join strings where one suffix matches another prefix.
	for ***REMOVED***
		// Find best i, j, k such that layout[i][len-k:] == layout[j][:k],
		// maximizing overlap length k.
		besti := -1
		bestj := -1
		bestk := 0
		for i, s := range layout ***REMOVED***
			if s == "" ***REMOVED***
				continue
			***REMOVED***
			for j, t := range layout ***REMOVED***
				if i == j ***REMOVED***
					continue
				***REMOVED***
				for k := bestk + 1; k <= len(s) && k <= len(t); k++ ***REMOVED***
					if s[len(s)-k:] == t[:k] ***REMOVED***
						besti = i
						bestj = j
						bestk = k
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if bestk > 0 ***REMOVED***
			layout[besti] += layout[bestj][bestk:]
			layout[bestj] = ""
			continue
		***REMOVED***
		break
	***REMOVED***

	text := strings.Join(layout, "")

	atom := map[string]uint32***REMOVED******REMOVED***
	for _, s := range all ***REMOVED***
		off := strings.Index(text, s)
		if off < 0 ***REMOVED***
			panic("lost string " + s)
		***REMOVED***
		atom[s] = uint32(off<<8 | len(s))
	***REMOVED***

	var buf bytes.Buffer
	// Generate the Go code.
	fmt.Fprintln(&buf, "// Code generated by go generate gen.go; DO NOT EDIT.\n")
	fmt.Fprintln(&buf, "//go:generate go run gen.go\n")
	fmt.Fprintln(&buf, "package atom\n\nconst (")

	// compute max len
	maxLen := 0
	for _, s := range all ***REMOVED***
		if maxLen < len(s) ***REMOVED***
			maxLen = len(s)
		***REMOVED***
		fmt.Fprintf(&buf, "\t%s Atom = %#x\n", identifier(s), atom[s])
	***REMOVED***
	fmt.Fprintln(&buf, ")\n")

	fmt.Fprintf(&buf, "const hash0 = %#x\n\n", best.h0)
	fmt.Fprintf(&buf, "const maxAtomLen = %d\n\n", maxLen)

	fmt.Fprintf(&buf, "var table = [1<<%d]Atom***REMOVED***\n", best.k)
	for i, s := range best.tab ***REMOVED***
		if s == "" ***REMOVED***
			continue
		***REMOVED***
		fmt.Fprintf(&buf, "\t%#x: %#x, // %s\n", i, atom[s], s)
	***REMOVED***
	fmt.Fprintf(&buf, "***REMOVED***\n")
	datasize := (1 << best.k) * 4

	fmt.Fprintln(&buf, "const atomText =")
	textsize := len(text)
	for len(text) > 60 ***REMOVED***
		fmt.Fprintf(&buf, "\t%q +\n", text[:60])
		text = text[60:]
	***REMOVED***
	fmt.Fprintf(&buf, "\t%q\n\n", text)

	genFile("table.go", &buf)

	fmt.Fprintf(os.Stdout, "%d atoms; %d string bytes + %d tables = %d total data\n", len(all), textsize, datasize, textsize+datasize)
***REMOVED***

type byLen []string

func (x byLen) Less(i, j int) bool ***REMOVED*** return len(x[i]) > len(x[j]) ***REMOVED***
func (x byLen) Swap(i, j int)      ***REMOVED*** x[i], x[j] = x[j], x[i] ***REMOVED***
func (x byLen) Len() int           ***REMOVED*** return len(x) ***REMOVED***

// fnv computes the FNV hash with an arbitrary starting value h.
func fnv(h uint32, s string) uint32 ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		h ^= uint32(s[i])
		h *= 16777619
	***REMOVED***
	return h
***REMOVED***

// A table represents an attempt at constructing the lookup table.
// The lookup table uses cuckoo hashing, meaning that each string
// can be found in one of two positions.
type table struct ***REMOVED***
	h0   uint32
	k    uint
	mask uint32
	tab  []string
***REMOVED***

// hash returns the two hashes for s.
func (t *table) hash(s string) (h1, h2 uint32) ***REMOVED***
	h := fnv(t.h0, s)
	h1 = h & t.mask
	h2 = (h >> 16) & t.mask
	return
***REMOVED***

// init initializes the table with the given parameters.
// h0 is the initial hash value,
// k is the number of bits of hash value to use, and
// x is the list of strings to store in the table.
// init returns false if the table cannot be constructed.
func (t *table) init(h0 uint32, k uint, x []string) bool ***REMOVED***
	t.h0 = h0
	t.k = k
	t.tab = make([]string, 1<<k)
	t.mask = 1<<k - 1
	for _, s := range x ***REMOVED***
		if !t.insert(s) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// insert inserts s in the table.
func (t *table) insert(s string) bool ***REMOVED***
	h1, h2 := t.hash(s)
	if t.tab[h1] == "" ***REMOVED***
		t.tab[h1] = s
		return true
	***REMOVED***
	if t.tab[h2] == "" ***REMOVED***
		t.tab[h2] = s
		return true
	***REMOVED***
	if t.push(h1, 0) ***REMOVED***
		t.tab[h1] = s
		return true
	***REMOVED***
	if t.push(h2, 0) ***REMOVED***
		t.tab[h2] = s
		return true
	***REMOVED***
	return false
***REMOVED***

// push attempts to push aside the entry in slot i.
func (t *table) push(i uint32, depth int) bool ***REMOVED***
	if depth > len(t.tab) ***REMOVED***
		return false
	***REMOVED***
	s := t.tab[i]
	h1, h2 := t.hash(s)
	j := h1 + h2 - i
	if t.tab[j] != "" && !t.push(j, depth+1) ***REMOVED***
		return false
	***REMOVED***
	t.tab[j] = s
	return true
***REMOVED***

// The lists of element names and attribute keys were taken from
// https://html.spec.whatwg.org/multipage/indices.html#index
// as of the "HTML Living Standard - Last Updated 18 September 2017" version.

// "command", "keygen" and "menuitem" have been removed from the spec,
// but are kept here for backwards compatibility.
var elements = []string***REMOVED***
	"a",
	"abbr",
	"address",
	"area",
	"article",
	"aside",
	"audio",
	"b",
	"base",
	"bdi",
	"bdo",
	"blockquote",
	"body",
	"br",
	"button",
	"canvas",
	"caption",
	"cite",
	"code",
	"col",
	"colgroup",
	"command",
	"data",
	"datalist",
	"dd",
	"del",
	"details",
	"dfn",
	"dialog",
	"div",
	"dl",
	"dt",
	"em",
	"embed",
	"fieldset",
	"figcaption",
	"figure",
	"footer",
	"form",
	"h1",
	"h2",
	"h3",
	"h4",
	"h5",
	"h6",
	"head",
	"header",
	"hgroup",
	"hr",
	"html",
	"i",
	"iframe",
	"img",
	"input",
	"ins",
	"kbd",
	"keygen",
	"label",
	"legend",
	"li",
	"link",
	"main",
	"map",
	"mark",
	"menu",
	"menuitem",
	"meta",
	"meter",
	"nav",
	"noscript",
	"object",
	"ol",
	"optgroup",
	"option",
	"output",
	"p",
	"param",
	"picture",
	"pre",
	"progress",
	"q",
	"rp",
	"rt",
	"ruby",
	"s",
	"samp",
	"script",
	"section",
	"select",
	"slot",
	"small",
	"source",
	"span",
	"strong",
	"style",
	"sub",
	"summary",
	"sup",
	"table",
	"tbody",
	"td",
	"template",
	"textarea",
	"tfoot",
	"th",
	"thead",
	"time",
	"title",
	"tr",
	"track",
	"u",
	"ul",
	"var",
	"video",
	"wbr",
***REMOVED***

// https://html.spec.whatwg.org/multipage/indices.html#attributes-3
//
// "challenge", "command", "contextmenu", "dropzone", "icon", "keytype", "mediagroup",
// "radiogroup", "spellcheck", "scoped", "seamless", "sortable" and "sorted" have been removed from the spec,
// but are kept here for backwards compatibility.
var attributes = []string***REMOVED***
	"abbr",
	"accept",
	"accept-charset",
	"accesskey",
	"action",
	"allowfullscreen",
	"allowpaymentrequest",
	"allowusermedia",
	"alt",
	"as",
	"async",
	"autocomplete",
	"autofocus",
	"autoplay",
	"challenge",
	"charset",
	"checked",
	"cite",
	"class",
	"color",
	"cols",
	"colspan",
	"command",
	"content",
	"contenteditable",
	"contextmenu",
	"controls",
	"coords",
	"crossorigin",
	"data",
	"datetime",
	"default",
	"defer",
	"dir",
	"dirname",
	"disabled",
	"download",
	"draggable",
	"dropzone",
	"enctype",
	"for",
	"form",
	"formaction",
	"formenctype",
	"formmethod",
	"formnovalidate",
	"formtarget",
	"headers",
	"height",
	"hidden",
	"high",
	"href",
	"hreflang",
	"http-equiv",
	"icon",
	"id",
	"inputmode",
	"integrity",
	"is",
	"ismap",
	"itemid",
	"itemprop",
	"itemref",
	"itemscope",
	"itemtype",
	"keytype",
	"kind",
	"label",
	"lang",
	"list",
	"loop",
	"low",
	"manifest",
	"max",
	"maxlength",
	"media",
	"mediagroup",
	"method",
	"min",
	"minlength",
	"multiple",
	"muted",
	"name",
	"nomodule",
	"nonce",
	"novalidate",
	"open",
	"optimum",
	"pattern",
	"ping",
	"placeholder",
	"playsinline",
	"poster",
	"preload",
	"radiogroup",
	"readonly",
	"referrerpolicy",
	"rel",
	"required",
	"reversed",
	"rows",
	"rowspan",
	"sandbox",
	"spellcheck",
	"scope",
	"scoped",
	"seamless",
	"selected",
	"shape",
	"size",
	"sizes",
	"sortable",
	"sorted",
	"slot",
	"span",
	"spellcheck",
	"src",
	"srcdoc",
	"srclang",
	"srcset",
	"start",
	"step",
	"style",
	"tabindex",
	"target",
	"title",
	"translate",
	"type",
	"typemustmatch",
	"updateviacache",
	"usemap",
	"value",
	"width",
	"workertype",
	"wrap",
***REMOVED***

// "onautocomplete", "onautocompleteerror", "onmousewheel",
// "onshow" and "onsort" have been removed from the spec,
// but are kept here for backwards compatibility.
var eventHandlers = []string***REMOVED***
	"onabort",
	"onautocomplete",
	"onautocompleteerror",
	"onauxclick",
	"onafterprint",
	"onbeforeprint",
	"onbeforeunload",
	"onblur",
	"oncancel",
	"oncanplay",
	"oncanplaythrough",
	"onchange",
	"onclick",
	"onclose",
	"oncontextmenu",
	"oncopy",
	"oncuechange",
	"oncut",
	"ondblclick",
	"ondrag",
	"ondragend",
	"ondragenter",
	"ondragexit",
	"ondragleave",
	"ondragover",
	"ondragstart",
	"ondrop",
	"ondurationchange",
	"onemptied",
	"onended",
	"onerror",
	"onfocus",
	"onhashchange",
	"oninput",
	"oninvalid",
	"onkeydown",
	"onkeypress",
	"onkeyup",
	"onlanguagechange",
	"onload",
	"onloadeddata",
	"onloadedmetadata",
	"onloadend",
	"onloadstart",
	"onmessage",
	"onmessageerror",
	"onmousedown",
	"onmouseenter",
	"onmouseleave",
	"onmousemove",
	"onmouseout",
	"onmouseover",
	"onmouseup",
	"onmousewheel",
	"onwheel",
	"onoffline",
	"ononline",
	"onpagehide",
	"onpageshow",
	"onpaste",
	"onpause",
	"onplay",
	"onplaying",
	"onpopstate",
	"onprogress",
	"onratechange",
	"onreset",
	"onresize",
	"onrejectionhandled",
	"onscroll",
	"onsecuritypolicyviolation",
	"onseeked",
	"onseeking",
	"onselect",
	"onshow",
	"onsort",
	"onstalled",
	"onstorage",
	"onsubmit",
	"onsuspend",
	"ontimeupdate",
	"ontoggle",
	"onunhandledrejection",
	"onunload",
	"onvolumechange",
	"onwaiting",
***REMOVED***

// extra are ad-hoc values not covered by any of the lists above.
var extra = []string***REMOVED***
	"align",
	"annotation",
	"annotation-xml",
	"applet",
	"basefont",
	"bgsound",
	"big",
	"blink",
	"center",
	"color",
	"desc",
	"face",
	"font",
	"foreignObject", // HTML is case-insensitive, but SVG-embedded-in-HTML is case-sensitive.
	"foreignobject",
	"frame",
	"frameset",
	"image",
	"isindex",
	"listing",
	"malignmark",
	"marquee",
	"math",
	"mglyph",
	"mi",
	"mn",
	"mo",
	"ms",
	"mtext",
	"nobr",
	"noembed",
	"noframes",
	"plaintext",
	"prompt",
	"public",
	"spacer",
	"strike",
	"svg",
	"system",
	"tt",
	"xmp",
***REMOVED***
