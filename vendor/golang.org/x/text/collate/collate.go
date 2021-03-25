// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: remove hard-coded versions when we have implemented fractional weights.
// The current implementation is incompatible with later CLDR versions.
//go:generate go run maketables.go -cldr=23 -unicode=6.2.0

// Package collate contains types for comparing and sorting Unicode strings
// according to a given collation order.
package collate // import "golang.org/x/text/collate"

import (
	"bytes"
	"strings"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/language"
)

// Collator provides functionality for comparing strings for a given
// collation order.
type Collator struct ***REMOVED***
	options

	sorter sorter

	_iter [2]iter
***REMOVED***

func (c *Collator) iter(i int) *iter ***REMOVED***
	// TODO: evaluate performance for making the second iterator optional.
	return &c._iter[i]
***REMOVED***

// Supported returns the list of languages for which collating differs from its parent.
func Supported() []language.Tag ***REMOVED***
	// TODO: use language.Coverage instead.

	t := make([]language.Tag, len(tags))
	copy(t, tags)
	return t
***REMOVED***

func init() ***REMOVED***
	ids := strings.Split(availableLocales, ",")
	tags = make([]language.Tag, len(ids))
	for i, s := range ids ***REMOVED***
		tags[i] = language.Raw.MustParse(s)
	***REMOVED***
***REMOVED***

var tags []language.Tag

// New returns a new Collator initialized for the given locale.
func New(t language.Tag, o ...Option) *Collator ***REMOVED***
	index := colltab.MatchLang(t, tags)
	c := newCollator(getTable(locales[index]))

	// Set options from the user-supplied tag.
	c.setFromTag(t)

	// Set the user-supplied options.
	c.setOptions(o)

	c.init()
	return c
***REMOVED***

// NewFromTable returns a new Collator for the given Weighter.
func NewFromTable(w colltab.Weighter, o ...Option) *Collator ***REMOVED***
	c := newCollator(w)
	c.setOptions(o)
	c.init()
	return c
***REMOVED***

func (c *Collator) init() ***REMOVED***
	if c.numeric ***REMOVED***
		c.t = colltab.NewNumericWeighter(c.t)
	***REMOVED***
	c._iter[0].init(c)
	c._iter[1].init(c)
***REMOVED***

// Buffer holds keys generated by Key and KeyString.
type Buffer struct ***REMOVED***
	buf [4096]byte
	key []byte
***REMOVED***

func (b *Buffer) init() ***REMOVED***
	if b.key == nil ***REMOVED***
		b.key = b.buf[:0]
	***REMOVED***
***REMOVED***

// Reset clears the buffer from previous results generated by Key and KeyString.
func (b *Buffer) Reset() ***REMOVED***
	b.key = b.key[:0]
***REMOVED***

// Compare returns an integer comparing the two byte slices.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (c *Collator) Compare(a, b []byte) int ***REMOVED***
	// TODO: skip identical prefixes once we have a fast way to detect if a rune is
	// part of a contraction. This would lead to roughly a 10% speedup for the colcmp regtest.
	c.iter(0).SetInput(a)
	c.iter(1).SetInput(b)
	if res := c.compare(); res != 0 ***REMOVED***
		return res
	***REMOVED***
	if !c.ignore[colltab.Identity] ***REMOVED***
		return bytes.Compare(a, b)
	***REMOVED***
	return 0
***REMOVED***

// CompareString returns an integer comparing the two strings.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (c *Collator) CompareString(a, b string) int ***REMOVED***
	// TODO: skip identical prefixes once we have a fast way to detect if a rune is
	// part of a contraction. This would lead to roughly a 10% speedup for the colcmp regtest.
	c.iter(0).SetInputString(a)
	c.iter(1).SetInputString(b)
	if res := c.compare(); res != 0 ***REMOVED***
		return res
	***REMOVED***
	if !c.ignore[colltab.Identity] ***REMOVED***
		if a < b ***REMOVED***
			return -1
		***REMOVED*** else if a > b ***REMOVED***
			return 1
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func compareLevel(f func(i *iter) int, a, b *iter) int ***REMOVED***
	a.pce = 0
	b.pce = 0
	for ***REMOVED***
		va := f(a)
		vb := f(b)
		if va != vb ***REMOVED***
			if va < vb ***REMOVED***
				return -1
			***REMOVED***
			return 1
		***REMOVED*** else if va == 0 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (c *Collator) compare() int ***REMOVED***
	ia, ib := c.iter(0), c.iter(1)
	// Process primary level
	if c.alternate != altShifted ***REMOVED***
		// TODO: implement script reordering
		if res := compareLevel((*iter).nextPrimary, ia, ib); res != 0 ***REMOVED***
			return res
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// TODO: handle shifted
	***REMOVED***
	if !c.ignore[colltab.Secondary] ***REMOVED***
		f := (*iter).nextSecondary
		if c.backwards ***REMOVED***
			f = (*iter).prevSecondary
		***REMOVED***
		if res := compareLevel(f, ia, ib); res != 0 ***REMOVED***
			return res
		***REMOVED***
	***REMOVED***
	// TODO: special case handling (Danish?)
	if !c.ignore[colltab.Tertiary] || c.caseLevel ***REMOVED***
		if res := compareLevel((*iter).nextTertiary, ia, ib); res != 0 ***REMOVED***
			return res
		***REMOVED***
		if !c.ignore[colltab.Quaternary] ***REMOVED***
			if res := compareLevel((*iter).nextQuaternary, ia, ib); res != 0 ***REMOVED***
				return res
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

// Key returns the collation key for str.
// Passing the buffer buf may avoid memory allocations.
// The returned slice will point to an allocation in Buffer and will remain
// valid until the next call to buf.Reset().
func (c *Collator) Key(buf *Buffer, str []byte) []byte ***REMOVED***
	// See https://www.unicode.org/reports/tr10/#Main_Algorithm for more details.
	buf.init()
	return c.key(buf, c.getColElems(str))
***REMOVED***

// KeyFromString returns the collation key for str.
// Passing the buffer buf may avoid memory allocations.
// The returned slice will point to an allocation in Buffer and will retain
// valid until the next call to buf.ResetKeys().
func (c *Collator) KeyFromString(buf *Buffer, str string) []byte ***REMOVED***
	// See https://www.unicode.org/reports/tr10/#Main_Algorithm for more details.
	buf.init()
	return c.key(buf, c.getColElemsString(str))
***REMOVED***

func (c *Collator) key(buf *Buffer, w []colltab.Elem) []byte ***REMOVED***
	processWeights(c.alternate, c.t.Top(), w)
	kn := len(buf.key)
	c.keyFromElems(buf, w)
	return buf.key[kn:]
***REMOVED***

func (c *Collator) getColElems(str []byte) []colltab.Elem ***REMOVED***
	i := c.iter(0)
	i.SetInput(str)
	for i.Next() ***REMOVED***
	***REMOVED***
	return i.Elems
***REMOVED***

func (c *Collator) getColElemsString(str string) []colltab.Elem ***REMOVED***
	i := c.iter(0)
	i.SetInputString(str)
	for i.Next() ***REMOVED***
	***REMOVED***
	return i.Elems
***REMOVED***

type iter struct ***REMOVED***
	wa [512]colltab.Elem

	colltab.Iter
	pce int
***REMOVED***

func (i *iter) init(c *Collator) ***REMOVED***
	i.Weighter = c.t
	i.Elems = i.wa[:0]
***REMOVED***

func (i *iter) nextPrimary() int ***REMOVED***
	for ***REMOVED***
		for ; i.pce < i.N; i.pce++ ***REMOVED***
			if v := i.Elems[i.pce].Primary(); v != 0 ***REMOVED***
				i.pce++
				return v
			***REMOVED***
		***REMOVED***
		if !i.Next() ***REMOVED***
			return 0
		***REMOVED***
	***REMOVED***
	panic("should not reach here")
***REMOVED***

func (i *iter) nextSecondary() int ***REMOVED***
	for ; i.pce < len(i.Elems); i.pce++ ***REMOVED***
		if v := i.Elems[i.pce].Secondary(); v != 0 ***REMOVED***
			i.pce++
			return v
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (i *iter) prevSecondary() int ***REMOVED***
	for ; i.pce < len(i.Elems); i.pce++ ***REMOVED***
		if v := i.Elems[len(i.Elems)-i.pce-1].Secondary(); v != 0 ***REMOVED***
			i.pce++
			return v
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (i *iter) nextTertiary() int ***REMOVED***
	for ; i.pce < len(i.Elems); i.pce++ ***REMOVED***
		if v := i.Elems[i.pce].Tertiary(); v != 0 ***REMOVED***
			i.pce++
			return int(v)
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func (i *iter) nextQuaternary() int ***REMOVED***
	for ; i.pce < len(i.Elems); i.pce++ ***REMOVED***
		if v := i.Elems[i.pce].Quaternary(); v != 0 ***REMOVED***
			i.pce++
			return v
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func appendPrimary(key []byte, p int) []byte ***REMOVED***
	// Convert to variable length encoding; supports up to 23 bits.
	if p <= 0x7FFF ***REMOVED***
		key = append(key, uint8(p>>8), uint8(p))
	***REMOVED*** else ***REMOVED***
		key = append(key, uint8(p>>16)|0x80, uint8(p>>8), uint8(p))
	***REMOVED***
	return key
***REMOVED***

// keyFromElems converts the weights ws to a compact sequence of bytes.
// The result will be appended to the byte buffer in buf.
func (c *Collator) keyFromElems(buf *Buffer, ws []colltab.Elem) ***REMOVED***
	for _, v := range ws ***REMOVED***
		if w := v.Primary(); w > 0 ***REMOVED***
			buf.key = appendPrimary(buf.key, w)
		***REMOVED***
	***REMOVED***
	if !c.ignore[colltab.Secondary] ***REMOVED***
		buf.key = append(buf.key, 0, 0)
		// TODO: we can use one 0 if we can guarantee that all non-zero weights are > 0xFF.
		if !c.backwards ***REMOVED***
			for _, v := range ws ***REMOVED***
				if w := v.Secondary(); w > 0 ***REMOVED***
					buf.key = append(buf.key, uint8(w>>8), uint8(w))
				***REMOVED***
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for i := len(ws) - 1; i >= 0; i-- ***REMOVED***
				if w := ws[i].Secondary(); w > 0 ***REMOVED***
					buf.key = append(buf.key, uint8(w>>8), uint8(w))
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if c.caseLevel ***REMOVED***
		buf.key = append(buf.key, 0, 0)
	***REMOVED***
	if !c.ignore[colltab.Tertiary] || c.caseLevel ***REMOVED***
		buf.key = append(buf.key, 0, 0)
		for _, v := range ws ***REMOVED***
			if w := v.Tertiary(); w > 0 ***REMOVED***
				buf.key = append(buf.key, uint8(w))
			***REMOVED***
		***REMOVED***
		// Derive the quaternary weights from the options and other levels.
		// Note that we represent MaxQuaternary as 0xFF. The first byte of the
		// representation of a primary weight is always smaller than 0xFF,
		// so using this single byte value will compare correctly.
		if !c.ignore[colltab.Quaternary] && c.alternate >= altShifted ***REMOVED***
			if c.alternate == altShiftTrimmed ***REMOVED***
				lastNonFFFF := len(buf.key)
				buf.key = append(buf.key, 0)
				for _, v := range ws ***REMOVED***
					if w := v.Quaternary(); w == colltab.MaxQuaternary ***REMOVED***
						buf.key = append(buf.key, 0xFF)
					***REMOVED*** else if w > 0 ***REMOVED***
						buf.key = appendPrimary(buf.key, w)
						lastNonFFFF = len(buf.key)
					***REMOVED***
				***REMOVED***
				buf.key = buf.key[:lastNonFFFF]
			***REMOVED*** else ***REMOVED***
				buf.key = append(buf.key, 0)
				for _, v := range ws ***REMOVED***
					if w := v.Quaternary(); w == colltab.MaxQuaternary ***REMOVED***
						buf.key = append(buf.key, 0xFF)
					***REMOVED*** else if w > 0 ***REMOVED***
						buf.key = appendPrimary(buf.key, w)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func processWeights(vw alternateHandling, top uint32, wa []colltab.Elem) ***REMOVED***
	ignore := false
	vtop := int(top)
	switch vw ***REMOVED***
	case altShifted, altShiftTrimmed:
		for i := range wa ***REMOVED***
			if p := wa[i].Primary(); p <= vtop && p != 0 ***REMOVED***
				wa[i] = colltab.MakeQuaternary(p)
				ignore = true
			***REMOVED*** else if p == 0 ***REMOVED***
				if ignore ***REMOVED***
					wa[i] = colltab.Ignore
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				ignore = false
			***REMOVED***
		***REMOVED***
	case altBlanked:
		for i := range wa ***REMOVED***
			if p := wa[i].Primary(); p <= vtop && (ignore || p != 0) ***REMOVED***
				wa[i] = colltab.Ignore
				ignore = true
			***REMOVED*** else ***REMOVED***
				ignore = false
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
