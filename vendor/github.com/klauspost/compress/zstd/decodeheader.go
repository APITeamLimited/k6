// Copyright 2020+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.

package zstd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// HeaderMaxSize is the maximum size of a Frame and Block Header.
// If less is sent to Header.Decode it *may* still contain enough information.
const HeaderMaxSize = 14 + 3

// Header contains information about the first frame and block within that.
type Header struct ***REMOVED***
	// SingleSegment specifies whether the data is to be decompressed into a
	// single contiguous memory segment.
	// It implies that WindowSize is invalid and that FrameContentSize is valid.
	SingleSegment bool

	// WindowSize is the window of data to keep while decoding.
	// Will only be set if SingleSegment is false.
	WindowSize uint64

	// Dictionary ID.
	// If 0, no dictionary.
	DictionaryID uint32

	// HasFCS specifies whether FrameContentSize has a valid value.
	HasFCS bool

	// FrameContentSize is the expected uncompressed size of the entire frame.
	FrameContentSize uint64

	// Skippable will be true if the frame is meant to be skipped.
	// This implies that FirstBlock.OK is false.
	Skippable bool

	// SkippableID is the user-specific ID for the skippable frame.
	// Valid values are between 0 to 15, inclusive.
	SkippableID int

	// SkippableSize is the length of the user data to skip following
	// the header.
	SkippableSize uint32

	// HeaderSize is the raw size of the frame header.
	//
	// For normal frames, it includes the size of the magic number and
	// the size of the header (per section 3.1.1.1).
	// It does not include the size for any data blocks (section 3.1.1.2) nor
	// the size for the trailing content checksum.
	//
	// For skippable frames, this counts the size of the magic number
	// along with the size of the size field of the payload.
	// It does not include the size of the skippable payload itself.
	// The total frame size is the HeaderSize plus the SkippableSize.
	HeaderSize int

	// First block information.
	FirstBlock struct ***REMOVED***
		// OK will be set if first block could be decoded.
		OK bool

		// Is this the last block of a frame?
		Last bool

		// Is the data compressed?
		// If true CompressedSize will be populated.
		// Unfortunately DecompressedSize cannot be determined
		// without decoding the blocks.
		Compressed bool

		// DecompressedSize is the expected decompressed size of the block.
		// Will be 0 if it cannot be determined.
		DecompressedSize int

		// CompressedSize of the data in the block.
		// Does not include the block header.
		// Will be equal to DecompressedSize if not Compressed.
		CompressedSize int
	***REMOVED***

	// If set there is a checksum present for the block content.
	// The checksum field at the end is always 4 bytes long.
	HasCheckSum bool
***REMOVED***

// Decode the header from the beginning of the stream.
// This will decode the frame header and the first block header if enough bytes are provided.
// It is recommended to provide at least HeaderMaxSize bytes.
// If the frame header cannot be read an error will be returned.
// If there isn't enough input, io.ErrUnexpectedEOF is returned.
// The FirstBlock.OK will indicate if enough information was available to decode the first block header.
func (h *Header) Decode(in []byte) error ***REMOVED***
	*h = Header***REMOVED******REMOVED***
	if len(in) < 4 ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	h.HeaderSize += 4
	b, in := in[:4], in[4:]
	if !bytes.Equal(b, frameMagic) ***REMOVED***
		if !bytes.Equal(b[1:4], skippableFrameMagic) || b[0]&0xf0 != 0x50 ***REMOVED***
			return ErrMagicMismatch
		***REMOVED***
		if len(in) < 4 ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		h.HeaderSize += 4
		h.Skippable = true
		h.SkippableID = int(b[0] & 0xf)
		h.SkippableSize = binary.LittleEndian.Uint32(in)
		return nil
	***REMOVED***

	// Read Window_Descriptor
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#window_descriptor
	if len(in) < 1 ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	fhd, in := in[0], in[1:]
	h.HeaderSize++
	h.SingleSegment = fhd&(1<<5) != 0
	h.HasCheckSum = fhd&(1<<2) != 0
	if fhd&(1<<3) != 0 ***REMOVED***
		return errors.New("reserved bit set on frame header")
	***REMOVED***

	if !h.SingleSegment ***REMOVED***
		if len(in) < 1 ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		var wd byte
		wd, in = in[0], in[1:]
		h.HeaderSize++
		windowLog := 10 + (wd >> 3)
		windowBase := uint64(1) << windowLog
		windowAdd := (windowBase / 8) * uint64(wd&0x7)
		h.WindowSize = windowBase + windowAdd
	***REMOVED***

	// Read Dictionary_ID
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#dictionary_id
	if size := fhd & 3; size != 0 ***REMOVED***
		if size == 3 ***REMOVED***
			size = 4
		***REMOVED***
		if len(in) < int(size) ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		b, in = in[:size], in[size:]
		h.HeaderSize += int(size)
		switch size ***REMOVED***
		case 1:
			h.DictionaryID = uint32(b[0])
		case 2:
			h.DictionaryID = uint32(b[0]) | (uint32(b[1]) << 8)
		case 4:
			h.DictionaryID = uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
		***REMOVED***
	***REMOVED***

	// Read Frame_Content_Size
	// https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md#frame_content_size
	var fcsSize int
	v := fhd >> 6
	switch v ***REMOVED***
	case 0:
		if h.SingleSegment ***REMOVED***
			fcsSize = 1
		***REMOVED***
	default:
		fcsSize = 1 << v
	***REMOVED***

	if fcsSize > 0 ***REMOVED***
		h.HasFCS = true
		if len(in) < fcsSize ***REMOVED***
			return io.ErrUnexpectedEOF
		***REMOVED***
		b, in = in[:fcsSize], in[fcsSize:]
		h.HeaderSize += int(fcsSize)
		switch fcsSize ***REMOVED***
		case 1:
			h.FrameContentSize = uint64(b[0])
		case 2:
			// When FCS_Field_Size is 2, the offset of 256 is added.
			h.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) + 256
		case 4:
			h.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24)
		case 8:
			d1 := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
			d2 := uint32(b[4]) | (uint32(b[5]) << 8) | (uint32(b[6]) << 16) | (uint32(b[7]) << 24)
			h.FrameContentSize = uint64(d1) | (uint64(d2) << 32)
		***REMOVED***
	***REMOVED***

	// Frame Header done, we will not fail from now on.
	if len(in) < 3 ***REMOVED***
		return nil
	***REMOVED***
	tmp := in[:3]
	bh := uint32(tmp[0]) | (uint32(tmp[1]) << 8) | (uint32(tmp[2]) << 16)
	h.FirstBlock.Last = bh&1 != 0
	blockType := blockType((bh >> 1) & 3)
	// find size.
	cSize := int(bh >> 3)
	switch blockType ***REMOVED***
	case blockTypeReserved:
		return nil
	case blockTypeRLE:
		h.FirstBlock.Compressed = true
		h.FirstBlock.DecompressedSize = cSize
		h.FirstBlock.CompressedSize = 1
	case blockTypeCompressed:
		h.FirstBlock.Compressed = true
		h.FirstBlock.CompressedSize = cSize
	case blockTypeRaw:
		h.FirstBlock.DecompressedSize = cSize
		h.FirstBlock.CompressedSize = cSize
	default:
		panic("Invalid block type")
	***REMOVED***

	h.FirstBlock.OK = true
	return nil
***REMOVED***
