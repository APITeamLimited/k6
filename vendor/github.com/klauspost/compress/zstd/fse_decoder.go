// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"fmt"
)

const (
	tablelogAbsoluteMax = 9
)

const (
	/*!MEMORY_USAGE :
	 *  Memory usage formula : N->2^N Bytes (examples : 10 -> 1KB; 12 -> 4KB ; 16 -> 64KB; 20 -> 1MB; etc.)
	 *  Increasing memory usage improves compression ratio
	 *  Reduced memory usage can improve speed, due to cache effect
	 *  Recommended max value is 14, for 16KB, which nicely fits into Intel x86 L1 cache */
	maxMemoryUsage = tablelogAbsoluteMax + 2

	maxTableLog    = maxMemoryUsage - 2
	maxTablesize   = 1 << maxTableLog
	maxTableMask   = (1 << maxTableLog) - 1
	minTablelog    = 5
	maxSymbolValue = 255
)

// fseDecoder provides temporary storage for compression and decompression.
type fseDecoder struct ***REMOVED***
	dt             [maxTablesize]decSymbol // Decompression table.
	symbolLen      uint16                  // Length of active part of the symbol table.
	actualTableLog uint8                   // Selected tablelog.
	maxBits        uint8                   // Maximum number of additional bits

	// used for table creation to avoid allocations.
	stateTable [256]uint16
	norm       [maxSymbolValue + 1]int16
	preDefined bool
***REMOVED***

// tableStep returns the next table index.
func tableStep(tableSize uint32) uint32 ***REMOVED***
	return (tableSize >> 1) + (tableSize >> 3) + 3
***REMOVED***

// readNCount will read the symbol distribution so decoding tables can be constructed.
func (s *fseDecoder) readNCount(b *byteReader, maxSymbol uint16) error ***REMOVED***
	var (
		charnum   uint16
		previous0 bool
	)
	if b.remain() < 4 ***REMOVED***
		return errors.New("input too small")
	***REMOVED***
	bitStream := b.Uint32NC()
	nbBits := uint((bitStream & 0xF) + minTablelog) // extract tableLog
	if nbBits > tablelogAbsoluteMax ***REMOVED***
		println("Invalid tablelog:", nbBits)
		return errors.New("tableLog too large")
	***REMOVED***
	bitStream >>= 4
	bitCount := uint(4)

	s.actualTableLog = uint8(nbBits)
	remaining := int32((1 << nbBits) + 1)
	threshold := int32(1 << nbBits)
	gotTotal := int32(0)
	nbBits++

	for remaining > 1 && charnum <= maxSymbol ***REMOVED***
		if previous0 ***REMOVED***
			//println("prev0")
			n0 := charnum
			for (bitStream & 0xFFFF) == 0xFFFF ***REMOVED***
				//println("24 x 0")
				n0 += 24
				if r := b.remain(); r > 5 ***REMOVED***
					b.advance(2)
					// The check above should make sure we can read 32 bits
					bitStream = b.Uint32NC() >> bitCount
				***REMOVED*** else ***REMOVED***
					// end of bit stream
					bitStream >>= 16
					bitCount += 16
				***REMOVED***
			***REMOVED***
			//printf("bitstream: %d, 0b%b", bitStream&3, bitStream)
			for (bitStream & 3) == 3 ***REMOVED***
				n0 += 3
				bitStream >>= 2
				bitCount += 2
			***REMOVED***
			n0 += uint16(bitStream & 3)
			bitCount += 2

			if n0 > maxSymbolValue ***REMOVED***
				return errors.New("maxSymbolValue too small")
			***REMOVED***
			//println("inserting ", n0-charnum, "zeroes from idx", charnum, "ending before", n0)
			for charnum < n0 ***REMOVED***
				s.norm[uint8(charnum)] = 0
				charnum++
			***REMOVED***

			if r := b.remain(); r >= 7 || r-int(bitCount>>3) >= 4 ***REMOVED***
				b.advance(bitCount >> 3)
				bitCount &= 7
				// The check above should make sure we can read 32 bits
				bitStream = b.Uint32NC() >> bitCount
			***REMOVED*** else ***REMOVED***
				bitStream >>= 2
			***REMOVED***
		***REMOVED***

		max := (2*threshold - 1) - remaining
		var count int32

		if int32(bitStream)&(threshold-1) < max ***REMOVED***
			count = int32(bitStream) & (threshold - 1)
			if debugAsserts && nbBits < 1 ***REMOVED***
				panic("nbBits underflow")
			***REMOVED***
			bitCount += nbBits - 1
		***REMOVED*** else ***REMOVED***
			count = int32(bitStream) & (2*threshold - 1)
			if count >= threshold ***REMOVED***
				count -= max
			***REMOVED***
			bitCount += nbBits
		***REMOVED***

		// extra accuracy
		count--
		if count < 0 ***REMOVED***
			// -1 means +1
			remaining += count
			gotTotal -= count
		***REMOVED*** else ***REMOVED***
			remaining -= count
			gotTotal += count
		***REMOVED***
		s.norm[charnum&0xff] = int16(count)
		charnum++
		previous0 = count == 0
		for remaining < threshold ***REMOVED***
			nbBits--
			threshold >>= 1
		***REMOVED***

		if r := b.remain(); r >= 7 || r-int(bitCount>>3) >= 4 ***REMOVED***
			b.advance(bitCount >> 3)
			bitCount &= 7
			// The check above should make sure we can read 32 bits
			bitStream = b.Uint32NC() >> (bitCount & 31)
		***REMOVED*** else ***REMOVED***
			bitCount -= (uint)(8 * (len(b.b) - 4 - b.off))
			b.off = len(b.b) - 4
			bitStream = b.Uint32() >> (bitCount & 31)
		***REMOVED***
	***REMOVED***
	s.symbolLen = charnum
	if s.symbolLen <= 1 ***REMOVED***
		return fmt.Errorf("symbolLen (%d) too small", s.symbolLen)
	***REMOVED***
	if s.symbolLen > maxSymbolValue+1 ***REMOVED***
		return fmt.Errorf("symbolLen (%d) too big", s.symbolLen)
	***REMOVED***
	if remaining != 1 ***REMOVED***
		return fmt.Errorf("corruption detected (remaining %d != 1)", remaining)
	***REMOVED***
	if bitCount > 32 ***REMOVED***
		return fmt.Errorf("corruption detected (bitCount %d > 32)", bitCount)
	***REMOVED***
	if gotTotal != 1<<s.actualTableLog ***REMOVED***
		return fmt.Errorf("corruption detected (total %d != %d)", gotTotal, 1<<s.actualTableLog)
	***REMOVED***
	b.advance((bitCount + 7) >> 3)
	// println(s.norm[:s.symbolLen], s.symbolLen)
	return s.buildDtable()
***REMOVED***

// decSymbol contains information about a state entry,
// Including the state offset base, the output symbol and
// the number of bits to read for the low part of the destination state.
// Using a composite uint64 is faster than a struct with separate members.
type decSymbol uint64

func newDecSymbol(nbits, addBits uint8, newState uint16, baseline uint32) decSymbol ***REMOVED***
	return decSymbol(nbits) | (decSymbol(addBits) << 8) | (decSymbol(newState) << 16) | (decSymbol(baseline) << 32)
***REMOVED***

func (d decSymbol) nbBits() uint8 ***REMOVED***
	return uint8(d)
***REMOVED***

func (d decSymbol) addBits() uint8 ***REMOVED***
	return uint8(d >> 8)
***REMOVED***

func (d decSymbol) newState() uint16 ***REMOVED***
	return uint16(d >> 16)
***REMOVED***

func (d decSymbol) baseline() uint32 ***REMOVED***
	return uint32(d >> 32)
***REMOVED***

func (d decSymbol) baselineInt() int ***REMOVED***
	return int(d >> 32)
***REMOVED***

func (d *decSymbol) set(nbits, addBits uint8, newState uint16, baseline uint32) ***REMOVED***
	*d = decSymbol(nbits) | (decSymbol(addBits) << 8) | (decSymbol(newState) << 16) | (decSymbol(baseline) << 32)
***REMOVED***

func (d *decSymbol) setNBits(nBits uint8) ***REMOVED***
	const mask = 0xffffffffffffff00
	*d = (*d & mask) | decSymbol(nBits)
***REMOVED***

func (d *decSymbol) setAddBits(addBits uint8) ***REMOVED***
	const mask = 0xffffffffffff00ff
	*d = (*d & mask) | (decSymbol(addBits) << 8)
***REMOVED***

func (d *decSymbol) setNewState(state uint16) ***REMOVED***
	const mask = 0xffffffff0000ffff
	*d = (*d & mask) | decSymbol(state)<<16
***REMOVED***

func (d *decSymbol) setBaseline(baseline uint32) ***REMOVED***
	const mask = 0xffffffff
	*d = (*d & mask) | decSymbol(baseline)<<32
***REMOVED***

func (d *decSymbol) setExt(addBits uint8, baseline uint32) ***REMOVED***
	const mask = 0xffff00ff
	*d = (*d & mask) | (decSymbol(addBits) << 8) | (decSymbol(baseline) << 32)
***REMOVED***

// decSymbolValue returns the transformed decSymbol for the given symbol.
func decSymbolValue(symb uint8, t []baseOffset) (decSymbol, error) ***REMOVED***
	if int(symb) >= len(t) ***REMOVED***
		return 0, fmt.Errorf("rle symbol %d >= max %d", symb, len(t))
	***REMOVED***
	lu := t[symb]
	return newDecSymbol(0, lu.addBits, 0, lu.baseLine), nil
***REMOVED***

// setRLE will set the decoder til RLE mode.
func (s *fseDecoder) setRLE(symbol decSymbol) ***REMOVED***
	s.actualTableLog = 0
	s.maxBits = symbol.addBits()
	s.dt[0] = symbol
***REMOVED***

// buildDtable will build the decoding table.
func (s *fseDecoder) buildDtable() error ***REMOVED***
	tableSize := uint32(1 << s.actualTableLog)
	highThreshold := tableSize - 1
	symbolNext := s.stateTable[:256]

	// Init, lay down lowprob symbols
	***REMOVED***
		for i, v := range s.norm[:s.symbolLen] ***REMOVED***
			if v == -1 ***REMOVED***
				s.dt[highThreshold].setAddBits(uint8(i))
				highThreshold--
				symbolNext[i] = 1
			***REMOVED*** else ***REMOVED***
				symbolNext[i] = uint16(v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	// Spread symbols
	***REMOVED***
		tableMask := tableSize - 1
		step := tableStep(tableSize)
		position := uint32(0)
		for ss, v := range s.norm[:s.symbolLen] ***REMOVED***
			for i := 0; i < int(v); i++ ***REMOVED***
				s.dt[position].setAddBits(uint8(ss))
				position = (position + step) & tableMask
				for position > highThreshold ***REMOVED***
					// lowprob area
					position = (position + step) & tableMask
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if position != 0 ***REMOVED***
			// position must reach all cells once, otherwise normalizedCounter is incorrect
			return errors.New("corrupted input (position != 0)")
		***REMOVED***
	***REMOVED***

	// Build Decoding table
	***REMOVED***
		tableSize := uint16(1 << s.actualTableLog)
		for u, v := range s.dt[:tableSize] ***REMOVED***
			symbol := v.addBits()
			nextState := symbolNext[symbol]
			symbolNext[symbol] = nextState + 1
			nBits := s.actualTableLog - byte(highBits(uint32(nextState)))
			s.dt[u&maxTableMask].setNBits(nBits)
			newState := (nextState << nBits) - tableSize
			if newState > tableSize ***REMOVED***
				return fmt.Errorf("newState (%d) outside table size (%d)", newState, tableSize)
			***REMOVED***
			if newState == uint16(u) && nBits == 0 ***REMOVED***
				// Seems weird that this is possible with nbits > 0.
				return fmt.Errorf("newState (%d) == oldState (%d) and no bits", newState, u)
			***REMOVED***
			s.dt[u&maxTableMask].setNewState(newState)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// transform will transform the decoder table into a table usable for
// decoding without having to apply the transformation while decoding.
// The state will contain the base value and the number of bits to read.
func (s *fseDecoder) transform(t []baseOffset) error ***REMOVED***
	tableSize := uint16(1 << s.actualTableLog)
	s.maxBits = 0
	for i, v := range s.dt[:tableSize] ***REMOVED***
		add := v.addBits()
		if int(add) >= len(t) ***REMOVED***
			return fmt.Errorf("invalid decoding table entry %d, symbol %d >= max (%d)", i, v.addBits(), len(t))
		***REMOVED***
		lu := t[add]
		if lu.addBits > s.maxBits ***REMOVED***
			s.maxBits = lu.addBits
		***REMOVED***
		v.setExt(lu.addBits, lu.baseLine)
		s.dt[i] = v
	***REMOVED***
	return nil
***REMOVED***

type fseState struct ***REMOVED***
	dt    []decSymbol
	state decSymbol
***REMOVED***

// Initialize and decodeAsync first state and symbol.
func (s *fseState) init(br *bitReader, tableLog uint8, dt []decSymbol) ***REMOVED***
	s.dt = dt
	br.fill()
	s.state = dt[br.getBits(tableLog)]
***REMOVED***

// next returns the current symbol and sets the next state.
// At least tablelog bits must be available in the bit reader.
func (s *fseState) next(br *bitReader) ***REMOVED***
	lowBits := uint16(br.getBits(s.state.nbBits()))
	s.state = s.dt[s.state.newState()+lowBits]
***REMOVED***

// finished returns true if all bits have been read from the bitstream
// and the next state would require reading bits from the input.
func (s *fseState) finished(br *bitReader) bool ***REMOVED***
	return br.finished() && s.state.nbBits() > 0
***REMOVED***

// final returns the current state symbol without decoding the next.
func (s *fseState) final() (int, uint8) ***REMOVED***
	return s.state.baselineInt(), s.state.addBits()
***REMOVED***

// final returns the current state symbol without decoding the next.
func (s decSymbol) final() (int, uint8) ***REMOVED***
	return s.baselineInt(), s.addBits()
***REMOVED***

// nextFast returns the next symbol and sets the next state.
// This can only be used if no symbols are 0 bits.
// At least tablelog bits must be available in the bit reader.
func (s *fseState) nextFast(br *bitReader) (uint32, uint8) ***REMOVED***
	lowBits := br.get16BitsFast(s.state.nbBits())
	s.state = s.dt[s.state.newState()+lowBits]
	return s.state.baseline(), s.state.addBits()
***REMOVED***
