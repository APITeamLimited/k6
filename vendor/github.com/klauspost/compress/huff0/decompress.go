package huff0

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/klauspost/compress/fse"
)

type dTable struct ***REMOVED***
	single []dEntrySingle
	double []dEntryDouble
***REMOVED***

// single-symbols decoding
type dEntrySingle struct ***REMOVED***
	entry uint16
***REMOVED***

// double-symbols decoding
type dEntryDouble struct ***REMOVED***
	seq   [4]byte
	nBits uint8
	len   uint8
***REMOVED***

// Uses special code for all tables that are < 8 bits.
const use8BitTables = true

// ReadTable will read a table from the input.
// The size of the input may be larger than the table definition.
// Any content remaining after the table definition will be returned.
// If no Scratch is provided a new one is allocated.
// The returned Scratch can be used for encoding or decoding input using this table.
func ReadTable(in []byte, s *Scratch) (s2 *Scratch, remain []byte, err error) ***REMOVED***
	s, err = s.prepare(in)
	if err != nil ***REMOVED***
		return s, nil, err
	***REMOVED***
	if len(in) <= 1 ***REMOVED***
		return s, nil, errors.New("input too small for table")
	***REMOVED***
	iSize := in[0]
	in = in[1:]
	if iSize >= 128 ***REMOVED***
		// Uncompressed
		oSize := iSize - 127
		iSize = (oSize + 1) / 2
		if int(iSize) > len(in) ***REMOVED***
			return s, nil, errors.New("input too small for table")
		***REMOVED***
		for n := uint8(0); n < oSize; n += 2 ***REMOVED***
			v := in[n/2]
			s.huffWeight[n] = v >> 4
			s.huffWeight[n+1] = v & 15
		***REMOVED***
		s.symbolLen = uint16(oSize)
		in = in[iSize:]
	***REMOVED*** else ***REMOVED***
		if len(in) < int(iSize) ***REMOVED***
			return s, nil, fmt.Errorf("input too small for table, want %d bytes, have %d", iSize, len(in))
		***REMOVED***
		// FSE compressed weights
		s.fse.DecompressLimit = 255
		hw := s.huffWeight[:]
		s.fse.Out = hw
		b, err := fse.Decompress(in[:iSize], s.fse)
		s.fse.Out = nil
		if err != nil ***REMOVED***
			return s, nil, err
		***REMOVED***
		if len(b) > 255 ***REMOVED***
			return s, nil, errors.New("corrupt input: output table too large")
		***REMOVED***
		s.symbolLen = uint16(len(b))
		in = in[iSize:]
	***REMOVED***

	// collect weight stats
	var rankStats [16]uint32
	weightTotal := uint32(0)
	for _, v := range s.huffWeight[:s.symbolLen] ***REMOVED***
		if v > tableLogMax ***REMOVED***
			return s, nil, errors.New("corrupt input: weight too large")
		***REMOVED***
		v2 := v & 15
		rankStats[v2]++
		// (1 << (v2-1)) is slower since the compiler cannot prove that v2 isn't 0.
		weightTotal += (1 << v2) >> 1
	***REMOVED***
	if weightTotal == 0 ***REMOVED***
		return s, nil, errors.New("corrupt input: weights zero")
	***REMOVED***

	// get last non-null symbol weight (implied, total must be 2^n)
	***REMOVED***
		tableLog := highBit32(weightTotal) + 1
		if tableLog > tableLogMax ***REMOVED***
			return s, nil, errors.New("corrupt input: tableLog too big")
		***REMOVED***
		s.actualTableLog = uint8(tableLog)
		// determine last weight
		***REMOVED***
			total := uint32(1) << tableLog
			rest := total - weightTotal
			verif := uint32(1) << highBit32(rest)
			lastWeight := highBit32(rest) + 1
			if verif != rest ***REMOVED***
				// last value must be a clean power of 2
				return s, nil, errors.New("corrupt input: last value not power of two")
			***REMOVED***
			s.huffWeight[s.symbolLen] = uint8(lastWeight)
			s.symbolLen++
			rankStats[lastWeight]++
		***REMOVED***
	***REMOVED***

	if (rankStats[1] < 2) || (rankStats[1]&1 != 0) ***REMOVED***
		// by construction : at least 2 elts of rank 1, must be even
		return s, nil, errors.New("corrupt input: min elt size, even check failed ")
	***REMOVED***

	// TODO: Choose between single/double symbol decoding

	// Calculate starting value for each rank
	***REMOVED***
		var nextRankStart uint32
		for n := uint8(1); n < s.actualTableLog+1; n++ ***REMOVED***
			current := nextRankStart
			nextRankStart += rankStats[n] << (n - 1)
			rankStats[n] = current
		***REMOVED***
	***REMOVED***

	// fill DTable (always full size)
	tSize := 1 << tableLogMax
	if len(s.dt.single) != tSize ***REMOVED***
		s.dt.single = make([]dEntrySingle, tSize)
	***REMOVED***
	cTable := s.prevTable
	if cap(cTable) < maxSymbolValue+1 ***REMOVED***
		cTable = make([]cTableEntry, 0, maxSymbolValue+1)
	***REMOVED***
	cTable = cTable[:maxSymbolValue+1]
	s.prevTable = cTable[:s.symbolLen]
	s.prevTableLog = s.actualTableLog

	for n, w := range s.huffWeight[:s.symbolLen] ***REMOVED***
		if w == 0 ***REMOVED***
			cTable[n] = cTableEntry***REMOVED***
				val:   0,
				nBits: 0,
			***REMOVED***
			continue
		***REMOVED***
		length := (uint32(1) << w) >> 1
		d := dEntrySingle***REMOVED***
			entry: uint16(s.actualTableLog+1-w) | (uint16(n) << 8),
		***REMOVED***

		rank := &rankStats[w]
		cTable[n] = cTableEntry***REMOVED***
			val:   uint16(*rank >> (w - 1)),
			nBits: uint8(d.entry),
		***REMOVED***

		single := s.dt.single[*rank : *rank+length]
		for i := range single ***REMOVED***
			single[i] = d
		***REMOVED***
		*rank += length
	***REMOVED***

	return s, in, nil
***REMOVED***

// Decompress1X will decompress a 1X encoded stream.
// The length of the supplied input must match the end of a block exactly.
// Before this is called, the table must be initialized with ReadTable unless
// the encoder re-used the table.
// deprecated: Use the stateless Decoder() to get a concurrent version.
func (s *Scratch) Decompress1X(in []byte) (out []byte, err error) ***REMOVED***
	if cap(s.Out) < s.MaxDecodedSize ***REMOVED***
		s.Out = make([]byte, s.MaxDecodedSize)
	***REMOVED***
	s.Out = s.Out[:0:s.MaxDecodedSize]
	s.Out, err = s.Decoder().Decompress1X(s.Out, in)
	return s.Out, err
***REMOVED***

// Decompress4X will decompress a 4X encoded stream.
// Before this is called, the table must be initialized with ReadTable unless
// the encoder re-used the table.
// The length of the supplied input must match the end of a block exactly.
// The destination size of the uncompressed data must be known and provided.
// deprecated: Use the stateless Decoder() to get a concurrent version.
func (s *Scratch) Decompress4X(in []byte, dstSize int) (out []byte, err error) ***REMOVED***
	if dstSize > s.MaxDecodedSize ***REMOVED***
		return nil, ErrMaxDecodedSizeExceeded
	***REMOVED***
	if cap(s.Out) < dstSize ***REMOVED***
		s.Out = make([]byte, s.MaxDecodedSize)
	***REMOVED***
	s.Out = s.Out[:0:dstSize]
	s.Out, err = s.Decoder().Decompress4X(s.Out, in)
	return s.Out, err
***REMOVED***

// Decoder will return a stateless decoder that can be used by multiple
// decompressors concurrently.
// Before this is called, the table must be initialized with ReadTable.
// The Decoder is still linked to the scratch buffer so that cannot be reused.
// However, it is safe to discard the scratch.
func (s *Scratch) Decoder() *Decoder ***REMOVED***
	return &Decoder***REMOVED***
		dt:             s.dt,
		actualTableLog: s.actualTableLog,
		bufs:           &s.decPool,
	***REMOVED***
***REMOVED***

// Decoder provides stateless decoding.
type Decoder struct ***REMOVED***
	dt             dTable
	actualTableLog uint8
	bufs           *sync.Pool
***REMOVED***

func (d *Decoder) buffer() *[4][256]byte ***REMOVED***
	buf, ok := d.bufs.Get().(*[4][256]byte)
	if ok ***REMOVED***
		return buf
	***REMOVED***
	return &[4][256]byte***REMOVED******REMOVED***
***REMOVED***

// Decompress1X will decompress a 1X encoded stream.
// The cap of the output buffer will be the maximum decompressed size.
// The length of the supplied input must match the end of a block exactly.
func (d *Decoder) Decompress1X(dst, src []byte) ([]byte, error) ***REMOVED***
	if len(d.dt.single) == 0 ***REMOVED***
		return nil, errors.New("no table loaded")
	***REMOVED***
	if use8BitTables && d.actualTableLog <= 8 ***REMOVED***
		return d.decompress1X8Bit(dst, src)
	***REMOVED***
	var br bitReaderShifted
	err := br.init(src)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	maxDecodedSize := cap(dst)
	dst = dst[:0]

	// Avoid bounds check by always having full sized table.
	const tlSize = 1 << tableLogMax
	const tlMask = tlSize - 1
	dt := d.dt.single[:tlSize]

	// Use temp table to avoid bound checks/append penalty.
	bufs := d.buffer()
	buf := &bufs[0]
	var off uint8

	for br.off >= 8 ***REMOVED***
		br.fillFast()
		v := dt[br.peekBitsFast(d.actualTableLog)&tlMask]
		br.advance(uint8(v.entry))
		buf[off+0] = uint8(v.entry >> 8)

		v = dt[br.peekBitsFast(d.actualTableLog)&tlMask]
		br.advance(uint8(v.entry))
		buf[off+1] = uint8(v.entry >> 8)

		// Refill
		br.fillFast()

		v = dt[br.peekBitsFast(d.actualTableLog)&tlMask]
		br.advance(uint8(v.entry))
		buf[off+2] = uint8(v.entry >> 8)

		v = dt[br.peekBitsFast(d.actualTableLog)&tlMask]
		br.advance(uint8(v.entry))
		buf[off+3] = uint8(v.entry >> 8)

		off += 4
		if off == 0 ***REMOVED***
			if len(dst)+256 > maxDecodedSize ***REMOVED***
				br.close()
				d.bufs.Put(bufs)
				return nil, ErrMaxDecodedSizeExceeded
			***REMOVED***
			dst = append(dst, buf[:]...)
		***REMOVED***
	***REMOVED***

	if len(dst)+int(off) > maxDecodedSize ***REMOVED***
		d.bufs.Put(bufs)
		br.close()
		return nil, ErrMaxDecodedSizeExceeded
	***REMOVED***
	dst = append(dst, buf[:off]...)

	// br < 8, so uint8 is fine
	bitsLeft := uint8(br.off)*8 + 64 - br.bitsRead
	for bitsLeft > 0 ***REMOVED***
		br.fill()
		if false && br.bitsRead >= 32 ***REMOVED***
			if br.off >= 4 ***REMOVED***
				v := br.in[br.off-4:]
				v = v[:4]
				low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
				br.value = (br.value << 32) | uint64(low)
				br.bitsRead -= 32
				br.off -= 4
			***REMOVED*** else ***REMOVED***
				for br.off > 0 ***REMOVED***
					br.value = (br.value << 8) | uint64(br.in[br.off-1])
					br.bitsRead -= 8
					br.off--
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if len(dst) >= maxDecodedSize ***REMOVED***
			d.bufs.Put(bufs)
			br.close()
			return nil, ErrMaxDecodedSizeExceeded
		***REMOVED***
		v := d.dt.single[br.peekBitsFast(d.actualTableLog)&tlMask]
		nBits := uint8(v.entry)
		br.advance(nBits)
		bitsLeft -= nBits
		dst = append(dst, uint8(v.entry>>8))
	***REMOVED***
	d.bufs.Put(bufs)
	return dst, br.close()
***REMOVED***

// decompress1X8Bit will decompress a 1X encoded stream with tablelog <= 8.
// The cap of the output buffer will be the maximum decompressed size.
// The length of the supplied input must match the end of a block exactly.
func (d *Decoder) decompress1X8Bit(dst, src []byte) ([]byte, error) ***REMOVED***
	if d.actualTableLog == 8 ***REMOVED***
		return d.decompress1X8BitExactly(dst, src)
	***REMOVED***
	var br bitReaderBytes
	err := br.init(src)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	maxDecodedSize := cap(dst)
	dst = dst[:0]

	// Avoid bounds check by always having full sized table.
	dt := d.dt.single[:256]

	// Use temp table to avoid bound checks/append penalty.
	bufs := d.buffer()
	buf := &bufs[0]
	var off uint8

	switch d.actualTableLog ***REMOVED***
	case 8:
		const shift = 8 - 8
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					br.close()
					d.bufs.Put(bufs)
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 7:
		const shift = 8 - 7
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					br.close()
					d.bufs.Put(bufs)
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 6:
		const shift = 8 - 6
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 5:
		const shift = 8 - 5
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 4:
		const shift = 8 - 4
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 3:
		const shift = 8 - 3
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 2:
		const shift = 8 - 2
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	case 1:
		const shift = 8 - 1
		for br.off >= 4 ***REMOVED***
			br.fillFast()
			v := dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+0] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+1] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+2] = uint8(v.entry >> 8)

			v = dt[uint8(br.value>>(56+shift))]
			br.advance(uint8(v.entry))
			buf[off+3] = uint8(v.entry >> 8)

			off += 4
			if off == 0 ***REMOVED***
				if len(dst)+256 > maxDecodedSize ***REMOVED***
					d.bufs.Put(bufs)
					br.close()
					return nil, ErrMaxDecodedSizeExceeded
				***REMOVED***
				dst = append(dst, buf[:]...)
			***REMOVED***
		***REMOVED***
	default:
		d.bufs.Put(bufs)
		return nil, fmt.Errorf("invalid tablelog: %d", d.actualTableLog)
	***REMOVED***

	if len(dst)+int(off) > maxDecodedSize ***REMOVED***
		d.bufs.Put(bufs)
		br.close()
		return nil, ErrMaxDecodedSizeExceeded
	***REMOVED***
	dst = append(dst, buf[:off]...)

	// br < 4, so uint8 is fine
	bitsLeft := int8(uint8(br.off)*8 + (64 - br.bitsRead))
	shift := (8 - d.actualTableLog) & 7

	for bitsLeft > 0 ***REMOVED***
		if br.bitsRead >= 64-8 ***REMOVED***
			for br.off > 0 ***REMOVED***
				br.value |= uint64(br.in[br.off-1]) << (br.bitsRead - 8)
				br.bitsRead -= 8
				br.off--
			***REMOVED***
		***REMOVED***
		if len(dst) >= maxDecodedSize ***REMOVED***
			br.close()
			d.bufs.Put(bufs)
			return nil, ErrMaxDecodedSizeExceeded
		***REMOVED***
		v := dt[br.peekByteFast()>>shift]
		nBits := uint8(v.entry)
		br.advance(nBits)
		bitsLeft -= int8(nBits)
		dst = append(dst, uint8(v.entry>>8))
	***REMOVED***
	d.bufs.Put(bufs)
	return dst, br.close()
***REMOVED***

// decompress1X8Bit will decompress a 1X encoded stream with tablelog <= 8.
// The cap of the output buffer will be the maximum decompressed size.
// The length of the supplied input must match the end of a block exactly.
func (d *Decoder) decompress1X8BitExactly(dst, src []byte) ([]byte, error) ***REMOVED***
	var br bitReaderBytes
	err := br.init(src)
	if err != nil ***REMOVED***
		return dst, err
	***REMOVED***
	maxDecodedSize := cap(dst)
	dst = dst[:0]

	// Avoid bounds check by always having full sized table.
	dt := d.dt.single[:256]

	// Use temp table to avoid bound checks/append penalty.
	bufs := d.buffer()
	buf := &bufs[0]
	var off uint8

	const shift = 56

	//fmt.Printf("mask: %b, tl:%d\n", mask, d.actualTableLog)
	for br.off >= 4 ***REMOVED***
		br.fillFast()
		v := dt[uint8(br.value>>shift)]
		br.advance(uint8(v.entry))
		buf[off+0] = uint8(v.entry >> 8)

		v = dt[uint8(br.value>>shift)]
		br.advance(uint8(v.entry))
		buf[off+1] = uint8(v.entry >> 8)

		v = dt[uint8(br.value>>shift)]
		br.advance(uint8(v.entry))
		buf[off+2] = uint8(v.entry >> 8)

		v = dt[uint8(br.value>>shift)]
		br.advance(uint8(v.entry))
		buf[off+3] = uint8(v.entry >> 8)

		off += 4
		if off == 0 ***REMOVED***
			if len(dst)+256 > maxDecodedSize ***REMOVED***
				d.bufs.Put(bufs)
				br.close()
				return nil, ErrMaxDecodedSizeExceeded
			***REMOVED***
			dst = append(dst, buf[:]...)
		***REMOVED***
	***REMOVED***

	if len(dst)+int(off) > maxDecodedSize ***REMOVED***
		d.bufs.Put(bufs)
		br.close()
		return nil, ErrMaxDecodedSizeExceeded
	***REMOVED***
	dst = append(dst, buf[:off]...)

	// br < 4, so uint8 is fine
	bitsLeft := int8(uint8(br.off)*8 + (64 - br.bitsRead))
	for bitsLeft > 0 ***REMOVED***
		if br.bitsRead >= 64-8 ***REMOVED***
			for br.off > 0 ***REMOVED***
				br.value |= uint64(br.in[br.off-1]) << (br.bitsRead - 8)
				br.bitsRead -= 8
				br.off--
			***REMOVED***
		***REMOVED***
		if len(dst) >= maxDecodedSize ***REMOVED***
			d.bufs.Put(bufs)
			br.close()
			return nil, ErrMaxDecodedSizeExceeded
		***REMOVED***
		v := dt[br.peekByteFast()]
		nBits := uint8(v.entry)
		br.advance(nBits)
		bitsLeft -= int8(nBits)
		dst = append(dst, uint8(v.entry>>8))
	***REMOVED***
	d.bufs.Put(bufs)
	return dst, br.close()
***REMOVED***

// Decompress4X will decompress a 4X encoded stream.
// The length of the supplied input must match the end of a block exactly.
// The *capacity* of the dst slice must match the destination size of
// the uncompressed data exactly.
func (d *Decoder) decompress4X8bit(dst, src []byte) ([]byte, error) ***REMOVED***
	if d.actualTableLog == 8 ***REMOVED***
		return d.decompress4X8bitExactly(dst, src)
	***REMOVED***

	var br [4]bitReaderBytes
	start := 6
	for i := 0; i < 3; i++ ***REMOVED***
		length := int(src[i*2]) | (int(src[i*2+1]) << 8)
		if start+length >= len(src) ***REMOVED***
			return nil, errors.New("truncated input (or invalid offset)")
		***REMOVED***
		err := br[i].init(src[start : start+length])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		start += length
	***REMOVED***
	err := br[3].init(src[start:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// destination, offset to match first output
	dstSize := cap(dst)
	dst = dst[:dstSize]
	out := dst
	dstEvery := (dstSize + 3) / 4

	shift := (56 + (8 - d.actualTableLog)) & 63

	const tlSize = 1 << 8
	single := d.dt.single[:tlSize]

	// Use temp table to avoid bound checks/append penalty.
	buf := d.buffer()
	var off uint8
	var decoded int

	// Decode 4 values from each decoder/loop.
	const bufoff = 256
	for ***REMOVED***
		if br[0].off < 4 || br[1].off < 4 || br[2].off < 4 || br[3].off < 4 ***REMOVED***
			break
		***REMOVED***

		***REMOVED***
			// Interleave 2 decodes.
			const stream = 0
			const stream2 = 1
			br1 := &br[stream]
			br2 := &br[stream2]
			br1.fillFast()
			br2.fillFast()

			v := single[uint8(br1.value>>shift)].entry
			v2 := single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off] = uint8(v >> 8)
			buf[stream2][off] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+1] = uint8(v >> 8)
			buf[stream2][off+1] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+2] = uint8(v >> 8)
			buf[stream2][off+2] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+3] = uint8(v >> 8)
			buf[stream2][off+3] = uint8(v2 >> 8)
		***REMOVED***

		***REMOVED***
			const stream = 2
			const stream2 = 3
			br1 := &br[stream]
			br2 := &br[stream2]
			br1.fillFast()
			br2.fillFast()

			v := single[uint8(br1.value>>shift)].entry
			v2 := single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off] = uint8(v >> 8)
			buf[stream2][off] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+1] = uint8(v >> 8)
			buf[stream2][off+1] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+2] = uint8(v >> 8)
			buf[stream2][off+2] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+3] = uint8(v >> 8)
			buf[stream2][off+3] = uint8(v2 >> 8)
		***REMOVED***

		off += 4

		if off == 0 ***REMOVED***
			if bufoff > dstEvery ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 1")
			***REMOVED***
			copy(out, buf[0][:])
			copy(out[dstEvery:], buf[1][:])
			copy(out[dstEvery*2:], buf[2][:])
			copy(out[dstEvery*3:], buf[3][:])
			out = out[bufoff:]
			decoded += bufoff * 4
			// There must at least be 3 buffers left.
			if len(out) < dstEvery*3 ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 2")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if off > 0 ***REMOVED***
		ioff := int(off)
		if len(out) < dstEvery*3+ioff ***REMOVED***
			d.bufs.Put(buf)
			return nil, errors.New("corruption detected: stream overrun 3")
		***REMOVED***
		copy(out, buf[0][:off])
		copy(out[dstEvery:], buf[1][:off])
		copy(out[dstEvery*2:], buf[2][:off])
		copy(out[dstEvery*3:], buf[3][:off])
		decoded += int(off) * 4
		out = out[off:]
	***REMOVED***

	// Decode remaining.
	// Decode remaining.
	remainBytes := dstEvery - (decoded / 4)
	for i := range br ***REMOVED***
		offset := dstEvery * i
		endsAt := offset + remainBytes
		if endsAt > len(out) ***REMOVED***
			endsAt = len(out)
		***REMOVED***
		br := &br[i]
		bitsLeft := br.remaining()
		for bitsLeft > 0 ***REMOVED***
			if br.finished() ***REMOVED***
				d.bufs.Put(buf)
				return nil, io.ErrUnexpectedEOF
			***REMOVED***
			if br.bitsRead >= 56 ***REMOVED***
				if br.off >= 4 ***REMOVED***
					v := br.in[br.off-4:]
					v = v[:4]
					low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
					br.value |= uint64(low) << (br.bitsRead - 32)
					br.bitsRead -= 32
					br.off -= 4
				***REMOVED*** else ***REMOVED***
					for br.off > 0 ***REMOVED***
						br.value |= uint64(br.in[br.off-1]) << (br.bitsRead - 8)
						br.bitsRead -= 8
						br.off--
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// end inline...
			if offset >= endsAt ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 4")
			***REMOVED***

			// Read value and increment offset.
			v := single[uint8(br.value>>shift)].entry
			nBits := uint8(v)
			br.advance(nBits)
			bitsLeft -= uint(nBits)
			out[offset] = uint8(v >> 8)
			offset++
		***REMOVED***
		if offset != endsAt ***REMOVED***
			d.bufs.Put(buf)
			return nil, fmt.Errorf("corruption detected: short output block %d, end %d != %d", i, offset, endsAt)
		***REMOVED***
		decoded += offset - dstEvery*i
		err = br.close()
		if err != nil ***REMOVED***
			d.bufs.Put(buf)
			return nil, err
		***REMOVED***
	***REMOVED***
	d.bufs.Put(buf)
	if dstSize != decoded ***REMOVED***
		return nil, errors.New("corruption detected: short output block")
	***REMOVED***
	return dst, nil
***REMOVED***

// Decompress4X will decompress a 4X encoded stream.
// The length of the supplied input must match the end of a block exactly.
// The *capacity* of the dst slice must match the destination size of
// the uncompressed data exactly.
func (d *Decoder) decompress4X8bitExactly(dst, src []byte) ([]byte, error) ***REMOVED***
	var br [4]bitReaderBytes
	start := 6
	for i := 0; i < 3; i++ ***REMOVED***
		length := int(src[i*2]) | (int(src[i*2+1]) << 8)
		if start+length >= len(src) ***REMOVED***
			return nil, errors.New("truncated input (or invalid offset)")
		***REMOVED***
		err := br[i].init(src[start : start+length])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		start += length
	***REMOVED***
	err := br[3].init(src[start:])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// destination, offset to match first output
	dstSize := cap(dst)
	dst = dst[:dstSize]
	out := dst
	dstEvery := (dstSize + 3) / 4

	const shift = 56
	const tlSize = 1 << 8
	const tlMask = tlSize - 1
	single := d.dt.single[:tlSize]

	// Use temp table to avoid bound checks/append penalty.
	buf := d.buffer()
	var off uint8
	var decoded int

	// Decode 4 values from each decoder/loop.
	const bufoff = 256
	for ***REMOVED***
		if br[0].off < 4 || br[1].off < 4 || br[2].off < 4 || br[3].off < 4 ***REMOVED***
			break
		***REMOVED***

		***REMOVED***
			// Interleave 2 decodes.
			const stream = 0
			const stream2 = 1
			br1 := &br[stream]
			br2 := &br[stream2]
			br1.fillFast()
			br2.fillFast()

			v := single[uint8(br1.value>>shift)].entry
			v2 := single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off] = uint8(v >> 8)
			buf[stream2][off] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+1] = uint8(v >> 8)
			buf[stream2][off+1] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+2] = uint8(v >> 8)
			buf[stream2][off+2] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+3] = uint8(v >> 8)
			buf[stream2][off+3] = uint8(v2 >> 8)
		***REMOVED***

		***REMOVED***
			const stream = 2
			const stream2 = 3
			br1 := &br[stream]
			br2 := &br[stream2]
			br1.fillFast()
			br2.fillFast()

			v := single[uint8(br1.value>>shift)].entry
			v2 := single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off] = uint8(v >> 8)
			buf[stream2][off] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+1] = uint8(v >> 8)
			buf[stream2][off+1] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+2] = uint8(v >> 8)
			buf[stream2][off+2] = uint8(v2 >> 8)

			v = single[uint8(br1.value>>shift)].entry
			v2 = single[uint8(br2.value>>shift)].entry
			br1.bitsRead += uint8(v)
			br1.value <<= v & 63
			br2.bitsRead += uint8(v2)
			br2.value <<= v2 & 63
			buf[stream][off+3] = uint8(v >> 8)
			buf[stream2][off+3] = uint8(v2 >> 8)
		***REMOVED***

		off += 4

		if off == 0 ***REMOVED***
			if bufoff > dstEvery ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 1")
			***REMOVED***
			copy(out, buf[0][:])
			copy(out[dstEvery:], buf[1][:])
			copy(out[dstEvery*2:], buf[2][:])
			copy(out[dstEvery*3:], buf[3][:])
			out = out[bufoff:]
			decoded += bufoff * 4
			// There must at least be 3 buffers left.
			if len(out) < dstEvery*3 ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 2")
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if off > 0 ***REMOVED***
		ioff := int(off)
		if len(out) < dstEvery*3+ioff ***REMOVED***
			return nil, errors.New("corruption detected: stream overrun 3")
		***REMOVED***
		copy(out, buf[0][:off])
		copy(out[dstEvery:], buf[1][:off])
		copy(out[dstEvery*2:], buf[2][:off])
		copy(out[dstEvery*3:], buf[3][:off])
		decoded += int(off) * 4
		out = out[off:]
	***REMOVED***

	// Decode remaining.
	remainBytes := dstEvery - (decoded / 4)
	for i := range br ***REMOVED***
		offset := dstEvery * i
		endsAt := offset + remainBytes
		if endsAt > len(out) ***REMOVED***
			endsAt = len(out)
		***REMOVED***
		br := &br[i]
		bitsLeft := br.remaining()
		for bitsLeft > 0 ***REMOVED***
			if br.finished() ***REMOVED***
				d.bufs.Put(buf)
				return nil, io.ErrUnexpectedEOF
			***REMOVED***
			if br.bitsRead >= 56 ***REMOVED***
				if br.off >= 4 ***REMOVED***
					v := br.in[br.off-4:]
					v = v[:4]
					low := (uint32(v[0])) | (uint32(v[1]) << 8) | (uint32(v[2]) << 16) | (uint32(v[3]) << 24)
					br.value |= uint64(low) << (br.bitsRead - 32)
					br.bitsRead -= 32
					br.off -= 4
				***REMOVED*** else ***REMOVED***
					for br.off > 0 ***REMOVED***
						br.value |= uint64(br.in[br.off-1]) << (br.bitsRead - 8)
						br.bitsRead -= 8
						br.off--
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// end inline...
			if offset >= endsAt ***REMOVED***
				d.bufs.Put(buf)
				return nil, errors.New("corruption detected: stream overrun 4")
			***REMOVED***

			// Read value and increment offset.
			v := single[br.peekByteFast()].entry
			nBits := uint8(v)
			br.advance(nBits)
			bitsLeft -= uint(nBits)
			out[offset] = uint8(v >> 8)
			offset++
		***REMOVED***
		if offset != endsAt ***REMOVED***
			d.bufs.Put(buf)
			return nil, fmt.Errorf("corruption detected: short output block %d, end %d != %d", i, offset, endsAt)
		***REMOVED***

		decoded += offset - dstEvery*i
		err = br.close()
		if err != nil ***REMOVED***
			d.bufs.Put(buf)
			return nil, err
		***REMOVED***
	***REMOVED***
	d.bufs.Put(buf)
	if dstSize != decoded ***REMOVED***
		return nil, errors.New("corruption detected: short output block")
	***REMOVED***
	return dst, nil
***REMOVED***

// matches will compare a decoding table to a coding table.
// Errors are written to the writer.
// Nothing will be written if table is ok.
func (s *Scratch) matches(ct cTable, w io.Writer) ***REMOVED***
	if s == nil || len(s.dt.single) == 0 ***REMOVED***
		return
	***REMOVED***
	dt := s.dt.single[:1<<s.actualTableLog]
	tablelog := s.actualTableLog
	ok := 0
	broken := 0
	for sym, enc := range ct ***REMOVED***
		errs := 0
		broken++
		if enc.nBits == 0 ***REMOVED***
			for _, dec := range dt ***REMOVED***
				if uint8(dec.entry>>8) == byte(sym) ***REMOVED***
					fmt.Fprintf(w, "symbol %x has decoder, but no encoder\n", sym)
					errs++
					break
				***REMOVED***
			***REMOVED***
			if errs == 0 ***REMOVED***
				broken--
			***REMOVED***
			continue
		***REMOVED***
		// Unused bits in input
		ub := tablelog - enc.nBits
		top := enc.val << ub
		// decoder looks at top bits.
		dec := dt[top]
		if uint8(dec.entry) != enc.nBits ***REMOVED***
			fmt.Fprintf(w, "symbol 0x%x bit size mismatch (enc: %d, dec:%d).\n", sym, enc.nBits, uint8(dec.entry))
			errs++
		***REMOVED***
		if uint8(dec.entry>>8) != uint8(sym) ***REMOVED***
			fmt.Fprintf(w, "symbol 0x%x decoder output mismatch (enc: %d, dec:%d).\n", sym, sym, uint8(dec.entry>>8))
			errs++
		***REMOVED***
		if errs > 0 ***REMOVED***
			fmt.Fprintf(w, "%d errros in base, stopping\n", errs)
			continue
		***REMOVED***
		// Ensure that all combinations are covered.
		for i := uint16(0); i < (1 << ub); i++ ***REMOVED***
			vval := top | i
			dec := dt[vval]
			if uint8(dec.entry) != enc.nBits ***REMOVED***
				fmt.Fprintf(w, "symbol 0x%x bit size mismatch (enc: %d, dec:%d).\n", vval, enc.nBits, uint8(dec.entry))
				errs++
			***REMOVED***
			if uint8(dec.entry>>8) != uint8(sym) ***REMOVED***
				fmt.Fprintf(w, "symbol 0x%x decoder output mismatch (enc: %d, dec:%d).\n", vval, sym, uint8(dec.entry>>8))
				errs++
			***REMOVED***
			if errs > 20 ***REMOVED***
				fmt.Fprintf(w, "%d errros, stopping\n", errs)
				break
			***REMOVED***
		***REMOVED***
		if errs == 0 ***REMOVED***
			ok++
			broken--
		***REMOVED***
	***REMOVED***
	if broken > 0 ***REMOVED***
		fmt.Fprintf(w, "%d broken, %d ok\n", broken, ok)
	***REMOVED***
***REMOVED***
