// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"github.com/klauspost/compress/huff0"
)

// history contains the information transferred between blocks.
type history struct ***REMOVED***
	// Literal decompression
	huffTree *huff0.Scratch

	// Sequence decompression
	decoders      sequenceDecs
	recentOffsets [3]int

	// History buffer...
	b []byte

	// ignoreBuffer is meant to ignore a number of bytes
	// when checking for matches in history
	ignoreBuffer int

	windowSize       int
	allocFrameBuffer int // needed?
	error            bool
	dict             *dict
***REMOVED***

// reset will reset the history to initial state of a frame.
// The history must already have been initialized to the desired size.
func (h *history) reset() ***REMOVED***
	h.b = h.b[:0]
	h.ignoreBuffer = 0
	h.error = false
	h.recentOffsets = [3]int***REMOVED***1, 4, 8***REMOVED***
	if f := h.decoders.litLengths.fse; f != nil && !f.preDefined ***REMOVED***
		fseDecoderPool.Put(f)
	***REMOVED***
	if f := h.decoders.offsets.fse; f != nil && !f.preDefined ***REMOVED***
		fseDecoderPool.Put(f)
	***REMOVED***
	if f := h.decoders.matchLengths.fse; f != nil && !f.preDefined ***REMOVED***
		fseDecoderPool.Put(f)
	***REMOVED***
	h.decoders = sequenceDecs***REMOVED***br: h.decoders.br***REMOVED***
	if h.huffTree != nil ***REMOVED***
		if h.dict == nil || h.dict.litEnc != h.huffTree ***REMOVED***
			huffDecoderPool.Put(h.huffTree)
		***REMOVED***
	***REMOVED***
	h.huffTree = nil
	h.dict = nil
	//printf("history created: %+v (l: %d, c: %d)", *h, len(h.b), cap(h.b))
***REMOVED***

func (h *history) setDict(dict *dict) ***REMOVED***
	if dict == nil ***REMOVED***
		return
	***REMOVED***
	h.dict = dict
	h.decoders.litLengths = dict.llDec
	h.decoders.offsets = dict.ofDec
	h.decoders.matchLengths = dict.mlDec
	h.decoders.dict = dict.content
	h.recentOffsets = dict.offsets
	h.huffTree = dict.litEnc
***REMOVED***

// append bytes to history.
// This function will make sure there is space for it,
// if the buffer has been allocated with enough extra space.
func (h *history) append(b []byte) ***REMOVED***
	if len(b) >= h.windowSize ***REMOVED***
		// Discard all history by simply overwriting
		h.b = h.b[:h.windowSize]
		copy(h.b, b[len(b)-h.windowSize:])
		return
	***REMOVED***

	// If there is space, append it.
	if len(b) < cap(h.b)-len(h.b) ***REMOVED***
		h.b = append(h.b, b...)
		return
	***REMOVED***

	// Move data down so we only have window size left.
	// We know we have less than window size in b at this point.
	discard := len(b) + len(h.b) - h.windowSize
	copy(h.b, h.b[discard:])
	h.b = h.b[:h.windowSize]
	copy(h.b[h.windowSize-len(b):], b)
***REMOVED***

// ensureBlock will ensure there is space for at least one block...
func (h *history) ensureBlock() ***REMOVED***
	if cap(h.b) < h.allocFrameBuffer ***REMOVED***
		h.b = make([]byte, 0, h.allocFrameBuffer)
		return
	***REMOVED***

	avail := cap(h.b) - len(h.b)
	if avail >= h.windowSize || avail > maxCompressedBlockSize ***REMOVED***
		return
	***REMOVED***
	// Move data down so we only have window size left.
	// We know we have less than window size in b at this point.
	discard := len(h.b) - h.windowSize
	copy(h.b, h.b[discard:])
	h.b = h.b[:h.windowSize]
***REMOVED***

// append bytes to history without ever discarding anything.
func (h *history) appendKeep(b []byte) ***REMOVED***
	h.b = append(h.b, b...)
***REMOVED***
