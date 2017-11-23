// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package idna implements IDNA2008 using the compatibility processing
// defined by UTS (Unicode Technical Standard) #46, which defines a standard to
// deal with the transition from IDNA2003.
//
// IDNA2008 (Internationalized Domain Names for Applications), is defined in RFC
// 5890, RFC 5891, RFC 5892, RFC 5893 and RFC 5894.
// UTS #46 is defined in http://www.unicode.org/reports/tr46.
// See http://unicode.org/cldr/utility/idna.jsp for a visualization of the
// differences between these two standards.
package idna

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/secure/bidirule"
	"golang.org/x/text/unicode/bidi"
	"golang.org/x/text/unicode/norm"
)

// NOTE: Unlike common practice in Go APIs, the functions will return a
// sanitized domain name in case of errors. Browsers sometimes use a partially
// evaluated string as lookup.
// TODO: the current error handling is, in my opinion, the least opinionated.
// Other strategies are also viable, though:
// Option 1) Return an empty string in case of error, but allow the user to
//    specify explicitly which errors to ignore.
// Option 2) Return the partially evaluated string if it is itself a valid
//    string, otherwise return the empty string in case of error.
// Option 3) Option 1 and 2.
// Option 4) Always return an empty string for now and implement Option 1 as
//    needed, and document that the return string may not be empty in case of
//    error in the future.
// I think Option 1 is best, but it is quite opinionated.

// ToASCII is a wrapper for Punycode.ToASCII.
func ToASCII(s string) (string, error) ***REMOVED***
	return Punycode.process(s, true)
***REMOVED***

// ToUnicode is a wrapper for Punycode.ToUnicode.
func ToUnicode(s string) (string, error) ***REMOVED***
	return Punycode.process(s, false)
***REMOVED***

// An Option configures a Profile at creation time.
type Option func(*options)

// Transitional sets a Profile to use the Transitional mapping as defined in UTS
// #46. This will cause, for example, "ß" to be mapped to "ss". Using the
// transitional mapping provides a compromise between IDNA2003 and IDNA2008
// compatibility. It is used by most browsers when resolving domain names. This
// option is only meaningful if combined with MapForLookup.
func Transitional(transitional bool) Option ***REMOVED***
	return func(o *options) ***REMOVED*** o.transitional = true ***REMOVED***
***REMOVED***

// VerifyDNSLength sets whether a Profile should fail if any of the IDN parts
// are longer than allowed by the RFC.
func VerifyDNSLength(verify bool) Option ***REMOVED***
	return func(o *options) ***REMOVED*** o.verifyDNSLength = verify ***REMOVED***
***REMOVED***

// RemoveLeadingDots removes leading label separators. Leading runes that map to
// dots, such as U+3002 IDEOGRAPHIC FULL STOP, are removed as well.
//
// This is the behavior suggested by the UTS #46 and is adopted by some
// browsers.
func RemoveLeadingDots(remove bool) Option ***REMOVED***
	return func(o *options) ***REMOVED*** o.removeLeadingDots = remove ***REMOVED***
***REMOVED***

// ValidateLabels sets whether to check the mandatory label validation criteria
// as defined in Section 5.4 of RFC 5891. This includes testing for correct use
// of hyphens ('-'), normalization, validity of runes, and the context rules.
func ValidateLabels(enable bool) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		// Don't override existing mappings, but set one that at least checks
		// normalization if it is not set.
		if o.mapping == nil && enable ***REMOVED***
			o.mapping = normalize
		***REMOVED***
		o.trie = trie
		o.validateLabels = enable
		o.fromPuny = validateFromPunycode
	***REMOVED***
***REMOVED***

// StrictDomainName limits the set of permissible ASCII characters to those
// allowed in domain names as defined in RFC 1034 (A-Z, a-z, 0-9 and the
// hyphen). This is set by default for MapForLookup and ValidateForRegistration.
//
// This option is useful, for instance, for browsers that allow characters
// outside this range, for example a '_' (U+005F LOW LINE). See
// http://www.rfc-editor.org/std/std3.txt for more details This option
// corresponds to the UseSTD3ASCIIRules option in UTS #46.
func StrictDomainName(use bool) Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.trie = trie
		o.useSTD3Rules = use
		o.fromPuny = validateFromPunycode
	***REMOVED***
***REMOVED***

// NOTE: the following options pull in tables. The tables should not be linked
// in as long as the options are not used.

// BidiRule enables the Bidi rule as defined in RFC 5893. Any application
// that relies on proper validation of labels should include this rule.
func BidiRule() Option ***REMOVED***
	return func(o *options) ***REMOVED*** o.bidirule = bidirule.ValidString ***REMOVED***
***REMOVED***

// ValidateForRegistration sets validation options to verify that a given IDN is
// properly formatted for registration as defined by Section 4 of RFC 5891.
func ValidateForRegistration() Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.mapping = validateRegistration
		StrictDomainName(true)(o)
		ValidateLabels(true)(o)
		VerifyDNSLength(true)(o)
		BidiRule()(o)
	***REMOVED***
***REMOVED***

// MapForLookup sets validation and mapping options such that a given IDN is
// transformed for domain name lookup according to the requirements set out in
// Section 5 of RFC 5891. The mappings follow the recommendations of RFC 5894,
// RFC 5895 and UTS 46. It does not add the Bidi Rule. Use the BidiRule option
// to add this check.
//
// The mappings include normalization and mapping case, width and other
// compatibility mappings.
func MapForLookup() Option ***REMOVED***
	return func(o *options) ***REMOVED***
		o.mapping = validateAndMap
		StrictDomainName(true)(o)
		ValidateLabels(true)(o)
	***REMOVED***
***REMOVED***

type options struct ***REMOVED***
	transitional      bool
	useSTD3Rules      bool
	validateLabels    bool
	verifyDNSLength   bool
	removeLeadingDots bool

	trie *idnaTrie

	// fromPuny calls validation rules when converting A-labels to U-labels.
	fromPuny func(p *Profile, s string) error

	// mapping implements a validation and mapping step as defined in RFC 5895
	// or UTS 46, tailored to, for example, domain registration or lookup.
	mapping func(p *Profile, s string) (mapped string, isBidi bool, err error)

	// bidirule, if specified, checks whether s conforms to the Bidi Rule
	// defined in RFC 5893.
	bidirule func(s string) bool
***REMOVED***

// A Profile defines the configuration of an IDNA mapper.
type Profile struct ***REMOVED***
	options
***REMOVED***

func apply(o *options, opts []Option) ***REMOVED***
	for _, f := range opts ***REMOVED***
		f(o)
	***REMOVED***
***REMOVED***

// New creates a new Profile.
//
// With no options, the returned Profile is the most permissive and equals the
// Punycode Profile. Options can be passed to further restrict the Profile. The
// MapForLookup and ValidateForRegistration options set a collection of options,
// for lookup and registration purposes respectively, which can be tailored by
// adding more fine-grained options, where later options override earlier
// options.
func New(o ...Option) *Profile ***REMOVED***
	p := &Profile***REMOVED******REMOVED***
	apply(&p.options, o)
	return p
***REMOVED***

// ToASCII converts a domain or domain label to its ASCII form. For example,
// ToASCII("bücher.example.com") is "xn--bcher-kva.example.com", and
// ToASCII("golang") is "golang". If an error is encountered it will return
// an error and a (partially) processed result.
func (p *Profile) ToASCII(s string) (string, error) ***REMOVED***
	return p.process(s, true)
***REMOVED***

// ToUnicode converts a domain or domain label to its Unicode form. For example,
// ToUnicode("xn--bcher-kva.example.com") is "bücher.example.com", and
// ToUnicode("golang") is "golang". If an error is encountered it will return
// an error and a (partially) processed result.
func (p *Profile) ToUnicode(s string) (string, error) ***REMOVED***
	pp := *p
	pp.transitional = false
	return pp.process(s, false)
***REMOVED***

// String reports a string with a description of the profile for debugging
// purposes. The string format may change with different versions.
func (p *Profile) String() string ***REMOVED***
	s := ""
	if p.transitional ***REMOVED***
		s = "Transitional"
	***REMOVED*** else ***REMOVED***
		s = "NonTransitional"
	***REMOVED***
	if p.useSTD3Rules ***REMOVED***
		s += ":UseSTD3Rules"
	***REMOVED***
	if p.validateLabels ***REMOVED***
		s += ":ValidateLabels"
	***REMOVED***
	if p.verifyDNSLength ***REMOVED***
		s += ":VerifyDNSLength"
	***REMOVED***
	return s
***REMOVED***

var (
	// Punycode is a Profile that does raw punycode processing with a minimum
	// of validation.
	Punycode *Profile = punycode

	// Lookup is the recommended profile for looking up domain names, according
	// to Section 5 of RFC 5891. The exact configuration of this profile may
	// change over time.
	Lookup *Profile = lookup

	// Display is the recommended profile for displaying domain names.
	// The configuration of this profile may change over time.
	Display *Profile = display

	// Registration is the recommended profile for checking whether a given
	// IDN is valid for registration, according to Section 4 of RFC 5891.
	Registration *Profile = registration

	punycode = &Profile***REMOVED******REMOVED***
	lookup   = &Profile***REMOVED***options***REMOVED***
		transitional:   true,
		useSTD3Rules:   true,
		validateLabels: true,
		trie:           trie,
		fromPuny:       validateFromPunycode,
		mapping:        validateAndMap,
		bidirule:       bidirule.ValidString,
	***REMOVED******REMOVED***
	display = &Profile***REMOVED***options***REMOVED***
		useSTD3Rules:   true,
		validateLabels: true,
		trie:           trie,
		fromPuny:       validateFromPunycode,
		mapping:        validateAndMap,
		bidirule:       bidirule.ValidString,
	***REMOVED******REMOVED***
	registration = &Profile***REMOVED***options***REMOVED***
		useSTD3Rules:    true,
		validateLabels:  true,
		verifyDNSLength: true,
		trie:            trie,
		fromPuny:        validateFromPunycode,
		mapping:         validateRegistration,
		bidirule:        bidirule.ValidString,
	***REMOVED******REMOVED***

	// TODO: profiles
	// Register: recommended for approving domain names: don't do any mappings
	// but rather reject on invalid input. Bundle or block deviation characters.
)

type labelError struct***REMOVED*** label, code_ string ***REMOVED***

func (e labelError) code() string ***REMOVED*** return e.code_ ***REMOVED***
func (e labelError) Error() string ***REMOVED***
	return fmt.Sprintf("idna: invalid label %q", e.label)
***REMOVED***

type runeError rune

func (e runeError) code() string ***REMOVED*** return "P1" ***REMOVED***
func (e runeError) Error() string ***REMOVED***
	return fmt.Sprintf("idna: disallowed rune %U", e)
***REMOVED***

// process implements the algorithm described in section 4 of UTS #46,
// see http://www.unicode.org/reports/tr46.
func (p *Profile) process(s string, toASCII bool) (string, error) ***REMOVED***
	var err error
	var isBidi bool
	if p.mapping != nil ***REMOVED***
		s, isBidi, err = p.mapping(p, s)
	***REMOVED***
	// Remove leading empty labels.
	if p.removeLeadingDots ***REMOVED***
		for ; len(s) > 0 && s[0] == '.'; s = s[1:] ***REMOVED***
		***REMOVED***
	***REMOVED***
	// TODO: allow for a quick check of the tables data.
	// It seems like we should only create this error on ToASCII, but the
	// UTS 46 conformance tests suggests we should always check this.
	if err == nil && p.verifyDNSLength && s == "" ***REMOVED***
		err = &labelError***REMOVED***s, "A4"***REMOVED***
	***REMOVED***
	labels := labelIter***REMOVED***orig: s***REMOVED***
	for ; !labels.done(); labels.next() ***REMOVED***
		label := labels.label()
		if label == "" ***REMOVED***
			// Empty labels are not okay. The label iterator skips the last
			// label if it is empty.
			if err == nil && p.verifyDNSLength ***REMOVED***
				err = &labelError***REMOVED***s, "A4"***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		if strings.HasPrefix(label, acePrefix) ***REMOVED***
			u, err2 := decode(label[len(acePrefix):])
			if err2 != nil ***REMOVED***
				if err == nil ***REMOVED***
					err = err2
				***REMOVED***
				// Spec says keep the old label.
				continue
			***REMOVED***
			isBidi = isBidi || bidirule.DirectionString(u) != bidi.LeftToRight
			labels.set(u)
			if err == nil && p.validateLabels ***REMOVED***
				err = p.fromPuny(p, u)
			***REMOVED***
			if err == nil ***REMOVED***
				// This should be called on NonTransitional, according to the
				// spec, but that currently does not have any effect. Use the
				// original profile to preserve options.
				err = p.validateLabel(u)
			***REMOVED***
		***REMOVED*** else if err == nil ***REMOVED***
			err = p.validateLabel(label)
		***REMOVED***
	***REMOVED***
	if isBidi && p.bidirule != nil && err == nil ***REMOVED***
		for labels.reset(); !labels.done(); labels.next() ***REMOVED***
			if !p.bidirule(labels.label()) ***REMOVED***
				err = &labelError***REMOVED***s, "B"***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if toASCII ***REMOVED***
		for labels.reset(); !labels.done(); labels.next() ***REMOVED***
			label := labels.label()
			if !ascii(label) ***REMOVED***
				a, err2 := encode(acePrefix, label)
				if err == nil ***REMOVED***
					err = err2
				***REMOVED***
				label = a
				labels.set(a)
			***REMOVED***
			n := len(label)
			if p.verifyDNSLength && err == nil && (n == 0 || n > 63) ***REMOVED***
				err = &labelError***REMOVED***label, "A4"***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	s = labels.result()
	if toASCII && p.verifyDNSLength && err == nil ***REMOVED***
		// Compute the length of the domain name minus the root label and its dot.
		n := len(s)
		if n > 0 && s[n-1] == '.' ***REMOVED***
			n--
		***REMOVED***
		if len(s) < 1 || n > 253 ***REMOVED***
			err = &labelError***REMOVED***s, "A4"***REMOVED***
		***REMOVED***
	***REMOVED***
	return s, err
***REMOVED***

func normalize(p *Profile, s string) (mapped string, isBidi bool, err error) ***REMOVED***
	// TODO: consider first doing a quick check to see if any of these checks
	// need to be done. This will make it slower in the general case, but
	// faster in the common case.
	mapped = norm.NFC.String(s)
	isBidi = bidirule.DirectionString(mapped) == bidi.RightToLeft
	return mapped, isBidi, nil
***REMOVED***

func validateRegistration(p *Profile, s string) (idem string, bidi bool, err error) ***REMOVED***
	// TODO: filter need for normalization in loop below.
	if !norm.NFC.IsNormalString(s) ***REMOVED***
		return s, false, &labelError***REMOVED***s, "V1"***REMOVED***
	***REMOVED***
	for i := 0; i < len(s); ***REMOVED***
		v, sz := trie.lookupString(s[i:])
		if sz == 0 ***REMOVED***
			return s, bidi, runeError(utf8.RuneError)
		***REMOVED***
		bidi = bidi || info(v).isBidi(s[i:])
		// Copy bytes not copied so far.
		switch p.simplify(info(v).category()) ***REMOVED***
		// TODO: handle the NV8 defined in the Unicode idna data set to allow
		// for strict conformance to IDNA2008.
		case valid, deviation:
		case disallowed, mapped, unknown, ignored:
			r, _ := utf8.DecodeRuneInString(s[i:])
			return s, bidi, runeError(r)
		***REMOVED***
		i += sz
	***REMOVED***
	return s, bidi, nil
***REMOVED***

func (c info) isBidi(s string) bool ***REMOVED***
	if !c.isMapped() ***REMOVED***
		return c&attributesMask == rtl
	***REMOVED***
	// TODO: also store bidi info for mapped data. This is possible, but a bit
	// cumbersome and not for the common case.
	p, _ := bidi.LookupString(s)
	switch p.Class() ***REMOVED***
	case bidi.R, bidi.AL, bidi.AN:
		return true
	***REMOVED***
	return false
***REMOVED***

func validateAndMap(p *Profile, s string) (vm string, bidi bool, err error) ***REMOVED***
	var (
		b []byte
		k int
	)
	// combinedInfoBits contains the or-ed bits of all runes. We use this
	// to derive the mayNeedNorm bit later. This may trigger normalization
	// overeagerly, but it will not do so in the common case. The end result
	// is another 10% saving on BenchmarkProfile for the common case.
	var combinedInfoBits info
	for i := 0; i < len(s); ***REMOVED***
		v, sz := trie.lookupString(s[i:])
		if sz == 0 ***REMOVED***
			b = append(b, s[k:i]...)
			b = append(b, "\ufffd"...)
			k = len(s)
			if err == nil ***REMOVED***
				err = runeError(utf8.RuneError)
			***REMOVED***
			break
		***REMOVED***
		combinedInfoBits |= info(v)
		bidi = bidi || info(v).isBidi(s[i:])
		start := i
		i += sz
		// Copy bytes not copied so far.
		switch p.simplify(info(v).category()) ***REMOVED***
		case valid:
			continue
		case disallowed:
			if err == nil ***REMOVED***
				r, _ := utf8.DecodeRuneInString(s[start:])
				err = runeError(r)
			***REMOVED***
			continue
		case mapped, deviation:
			b = append(b, s[k:start]...)
			b = info(v).appendMapping(b, s[start:i])
		case ignored:
			b = append(b, s[k:start]...)
			// drop the rune
		case unknown:
			b = append(b, s[k:start]...)
			b = append(b, "\ufffd"...)
		***REMOVED***
		k = i
	***REMOVED***
	if k == 0 ***REMOVED***
		// No changes so far.
		if combinedInfoBits&mayNeedNorm != 0 ***REMOVED***
			s = norm.NFC.String(s)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		b = append(b, s[k:]...)
		if norm.NFC.QuickSpan(b) != len(b) ***REMOVED***
			b = norm.NFC.Bytes(b)
		***REMOVED***
		// TODO: the punycode converters require strings as input.
		s = string(b)
	***REMOVED***
	return s, bidi, err
***REMOVED***

// A labelIter allows iterating over domain name labels.
type labelIter struct ***REMOVED***
	orig     string
	slice    []string
	curStart int
	curEnd   int
	i        int
***REMOVED***

func (l *labelIter) reset() ***REMOVED***
	l.curStart = 0
	l.curEnd = 0
	l.i = 0
***REMOVED***

func (l *labelIter) done() bool ***REMOVED***
	return l.curStart >= len(l.orig)
***REMOVED***

func (l *labelIter) result() string ***REMOVED***
	if l.slice != nil ***REMOVED***
		return strings.Join(l.slice, ".")
	***REMOVED***
	return l.orig
***REMOVED***

func (l *labelIter) label() string ***REMOVED***
	if l.slice != nil ***REMOVED***
		return l.slice[l.i]
	***REMOVED***
	p := strings.IndexByte(l.orig[l.curStart:], '.')
	l.curEnd = l.curStart + p
	if p == -1 ***REMOVED***
		l.curEnd = len(l.orig)
	***REMOVED***
	return l.orig[l.curStart:l.curEnd]
***REMOVED***

// next sets the value to the next label. It skips the last label if it is empty.
func (l *labelIter) next() ***REMOVED***
	l.i++
	if l.slice != nil ***REMOVED***
		if l.i >= len(l.slice) || l.i == len(l.slice)-1 && l.slice[l.i] == "" ***REMOVED***
			l.curStart = len(l.orig)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		l.curStart = l.curEnd + 1
		if l.curStart == len(l.orig)-1 && l.orig[l.curStart] == '.' ***REMOVED***
			l.curStart = len(l.orig)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (l *labelIter) set(s string) ***REMOVED***
	if l.slice == nil ***REMOVED***
		l.slice = strings.Split(l.orig, ".")
	***REMOVED***
	l.slice[l.i] = s
***REMOVED***

// acePrefix is the ASCII Compatible Encoding prefix.
const acePrefix = "xn--"

func (p *Profile) simplify(cat category) category ***REMOVED***
	switch cat ***REMOVED***
	case disallowedSTD3Mapped:
		if p.useSTD3Rules ***REMOVED***
			cat = disallowed
		***REMOVED*** else ***REMOVED***
			cat = mapped
		***REMOVED***
	case disallowedSTD3Valid:
		if p.useSTD3Rules ***REMOVED***
			cat = disallowed
		***REMOVED*** else ***REMOVED***
			cat = valid
		***REMOVED***
	case deviation:
		if !p.transitional ***REMOVED***
			cat = valid
		***REMOVED***
	case validNV8, validXV8:
		// TODO: handle V2008
		cat = valid
	***REMOVED***
	return cat
***REMOVED***

func validateFromPunycode(p *Profile, s string) error ***REMOVED***
	if !norm.NFC.IsNormalString(s) ***REMOVED***
		return &labelError***REMOVED***s, "V1"***REMOVED***
	***REMOVED***
	// TODO: detect whether string may have to be normalized in the following
	// loop.
	for i := 0; i < len(s); ***REMOVED***
		v, sz := trie.lookupString(s[i:])
		if sz == 0 ***REMOVED***
			return runeError(utf8.RuneError)
		***REMOVED***
		if c := p.simplify(info(v).category()); c != valid && c != deviation ***REMOVED***
			return &labelError***REMOVED***s, "V6"***REMOVED***
		***REMOVED***
		i += sz
	***REMOVED***
	return nil
***REMOVED***

const (
	zwnj = "\u200c"
	zwj  = "\u200d"
)

type joinState int8

const (
	stateStart joinState = iota
	stateVirama
	stateBefore
	stateBeforeVirama
	stateAfter
	stateFAIL
)

var joinStates = [][numJoinTypes]joinState***REMOVED***
	stateStart: ***REMOVED***
		joiningL:   stateBefore,
		joiningD:   stateBefore,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateVirama,
	***REMOVED***,
	stateVirama: ***REMOVED***
		joiningL: stateBefore,
		joiningD: stateBefore,
	***REMOVED***,
	stateBefore: ***REMOVED***
		joiningL:   stateBefore,
		joiningD:   stateBefore,
		joiningT:   stateBefore,
		joinZWNJ:   stateAfter,
		joinZWJ:    stateFAIL,
		joinVirama: stateBeforeVirama,
	***REMOVED***,
	stateBeforeVirama: ***REMOVED***
		joiningL: stateBefore,
		joiningD: stateBefore,
		joiningT: stateBefore,
	***REMOVED***,
	stateAfter: ***REMOVED***
		joiningL:   stateFAIL,
		joiningD:   stateBefore,
		joiningT:   stateAfter,
		joiningR:   stateStart,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateAfter, // no-op as we can't accept joiners here
	***REMOVED***,
	stateFAIL: ***REMOVED***
		0:          stateFAIL,
		joiningL:   stateFAIL,
		joiningD:   stateFAIL,
		joiningT:   stateFAIL,
		joiningR:   stateFAIL,
		joinZWNJ:   stateFAIL,
		joinZWJ:    stateFAIL,
		joinVirama: stateFAIL,
	***REMOVED***,
***REMOVED***

// validateLabel validates the criteria from Section 4.1. Item 1, 4, and 6 are
// already implicitly satisfied by the overall implementation.
func (p *Profile) validateLabel(s string) (err error) ***REMOVED***
	if s == "" ***REMOVED***
		if p.verifyDNSLength ***REMOVED***
			return &labelError***REMOVED***s, "A4"***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	if !p.validateLabels ***REMOVED***
		return nil
	***REMOVED***
	trie := p.trie // p.validateLabels is only set if trie is set.
	if len(s) > 4 && s[2] == '-' && s[3] == '-' ***REMOVED***
		return &labelError***REMOVED***s, "V2"***REMOVED***
	***REMOVED***
	if s[0] == '-' || s[len(s)-1] == '-' ***REMOVED***
		return &labelError***REMOVED***s, "V3"***REMOVED***
	***REMOVED***
	// TODO: merge the use of this in the trie.
	v, sz := trie.lookupString(s)
	x := info(v)
	if x.isModifier() ***REMOVED***
		return &labelError***REMOVED***s, "V5"***REMOVED***
	***REMOVED***
	// Quickly return in the absence of zero-width (non) joiners.
	if strings.Index(s, zwj) == -1 && strings.Index(s, zwnj) == -1 ***REMOVED***
		return nil
	***REMOVED***
	st := stateStart
	for i := 0; ; ***REMOVED***
		jt := x.joinType()
		if s[i:i+sz] == zwj ***REMOVED***
			jt = joinZWJ
		***REMOVED*** else if s[i:i+sz] == zwnj ***REMOVED***
			jt = joinZWNJ
		***REMOVED***
		st = joinStates[st][jt]
		if x.isViramaModifier() ***REMOVED***
			st = joinStates[st][joinVirama]
		***REMOVED***
		if i += sz; i == len(s) ***REMOVED***
			break
		***REMOVED***
		v, sz = trie.lookupString(s[i:])
		x = info(v)
	***REMOVED***
	if st == stateFAIL || st == stateAfter ***REMOVED***
		return &labelError***REMOVED***s, "C"***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func ascii(s string) bool ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] >= utf8.RuneSelf ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***
