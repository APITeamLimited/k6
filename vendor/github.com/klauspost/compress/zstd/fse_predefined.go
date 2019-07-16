// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"fmt"
	"math"
)

var (
	// fsePredef are the predefined fse tables as defined here:
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#default-distributions
	// These values are already transformed.
	fsePredef [3]fseDecoder

	// fsePredefEnc are the predefined encoder based on fse tables as defined here:
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#default-distributions
	// These values are already transformed.
	fsePredefEnc [3]fseEncoder

	// symbolTableX contain the transformations needed for each type as defined in
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#the-codes-for-literals-lengths-match-lengths-and-offsets
	symbolTableX [3][]baseOffset

	// maxTableSymbol is the biggest supported symbol for each table type
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#the-codes-for-literals-lengths-match-lengths-and-offsets
	maxTableSymbol = [3]uint8***REMOVED***tableLiteralLengths: maxLiteralLengthSymbol, tableOffsets: maxOffsetLengthSymbol, tableMatchLengths: maxMatchLengthSymbol***REMOVED***

	// bitTables is the bits table for each table.
	bitTables = [3][]byte***REMOVED***tableLiteralLengths: llBitsTable[:], tableOffsets: nil, tableMatchLengths: mlBitsTable[:]***REMOVED***
)

type tableIndex uint8

const (
	// indexes for fsePredef and symbolTableX
	tableLiteralLengths tableIndex = 0
	tableOffsets        tableIndex = 1
	tableMatchLengths   tableIndex = 2

	maxLiteralLengthSymbol = 35
	maxOffsetLengthSymbol  = 30
	maxMatchLengthSymbol   = 52
)

// baseOffset is used for calculating transformations.
type baseOffset struct ***REMOVED***
	baseLine uint32
	addBits  uint8
***REMOVED***

// fillBase will precalculate base offsets with the given bit distributions.
func fillBase(dst []baseOffset, base uint32, bits ...uint8) ***REMOVED***
	if len(bits) != len(dst) ***REMOVED***
		panic(fmt.Sprintf("len(dst) (%d) != len(bits) (%d)", len(dst), len(bits)))
	***REMOVED***
	for i, bit := range bits ***REMOVED***
		if base > math.MaxInt32 ***REMOVED***
			panic(fmt.Sprintf("invalid decoding table, base overflows int32"))
		***REMOVED***

		dst[i] = baseOffset***REMOVED***
			baseLine: base,
			addBits:  bit,
		***REMOVED***
		base += 1 << bit
	***REMOVED***
***REMOVED***

func init() ***REMOVED***
	// Literals length codes
	tmp := make([]baseOffset, 36)
	for i := range tmp[:16] ***REMOVED***
		tmp[i] = baseOffset***REMOVED***
			baseLine: uint32(i),
			addBits:  0,
		***REMOVED***
	***REMOVED***
	fillBase(tmp[16:], 16, 1, 1, 1, 1, 2, 2, 3, 3, 4, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
	symbolTableX[tableLiteralLengths] = tmp

	// Match length codes
	tmp = make([]baseOffset, 53)
	for i := range tmp[:32] ***REMOVED***
		tmp[i] = baseOffset***REMOVED***
			// The transformation adds the 3 length.
			baseLine: uint32(i) + 3,
			addBits:  0,
		***REMOVED***
	***REMOVED***
	fillBase(tmp[32:], 35, 1, 1, 1, 1, 2, 2, 3, 3, 4, 4, 5, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
	symbolTableX[tableMatchLengths] = tmp

	// Offset codes
	tmp = make([]baseOffset, maxOffsetBits+1)
	tmp[1] = baseOffset***REMOVED***
		baseLine: 1,
		addBits:  1,
	***REMOVED***
	fillBase(tmp[2:], 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30)
	symbolTableX[tableOffsets] = tmp

	// Fill predefined tables and transform them.
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#default-distributions
	for i := range fsePredef[:] ***REMOVED***
		f := &fsePredef[i]
		switch tableIndex(i) ***REMOVED***
		case tableLiteralLengths:
			// https://github.com/facebook/zstd/blob/ededcfca57366461021c922720878c81a5854a0a/lib/decompress/zstd_decompress_block.c#L243
			f.actualTableLog = 6
			copy(f.norm[:], []int16***REMOVED***4, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1,
				2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 2, 1, 1, 1, 1, 1,
				-1, -1, -1, -1***REMOVED***)
			f.symbolLen = 36
		case tableOffsets:
			// https://github.com/facebook/zstd/blob/ededcfca57366461021c922720878c81a5854a0a/lib/decompress/zstd_decompress_block.c#L281
			f.actualTableLog = 5
			copy(f.norm[:], []int16***REMOVED***
				1, 1, 1, 1, 1, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, -1, -1***REMOVED***)
			f.symbolLen = 29
		case tableMatchLengths:
			//https://github.com/facebook/zstd/blob/ededcfca57366461021c922720878c81a5854a0a/lib/decompress/zstd_decompress_block.c#L304
			f.actualTableLog = 6
			copy(f.norm[:], []int16***REMOVED***
				1, 4, 3, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
				1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1,
				-1, -1, -1, -1, -1***REMOVED***)
			f.symbolLen = 53
		***REMOVED***
		if err := f.buildDtable(); err != nil ***REMOVED***
			panic(fmt.Errorf("building table %v: %v", tableIndex(i), err))
		***REMOVED***
		if err := f.transform(symbolTableX[i]); err != nil ***REMOVED***
			panic(fmt.Errorf("building table %v: %v", tableIndex(i), err))
		***REMOVED***
		f.preDefined = true

		// Create encoder as well
		enc := &fsePredefEnc[i]
		copy(enc.norm[:], f.norm[:])
		enc.symbolLen = f.symbolLen
		enc.actualTableLog = f.actualTableLog
		if err := enc.buildCTable(); err != nil ***REMOVED***
			panic(fmt.Errorf("building encoding table %v: %v", tableIndex(i), err))
		***REMOVED***
		enc.setBits(bitTables[i])
		enc.preDefined = true
	***REMOVED***
***REMOVED***
