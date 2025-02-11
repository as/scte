// This bitreader was copied from github.com/as/bit
// its not worth being a dependency
package scte

import (
	"encoding/binary"
	"fmt"
	"io"
)

// NewReader returns a new bitstream Reader with b as the
// bytestream
func NewReader(b []byte) *Reader {
	return &Reader{
		b:   b,
		at:  0,
		len: len(b)*8 + 7, // fixed point bit offset
	}
}

// Reader reads bits from a byte stream
type Reader struct {
	b   []byte
	at  int
	len int
	err error
}

// ReadPrint reads a named symbol of n bytes from the underlying
// reader, returning it as a uint64
func (r *Reader) ReadPrint(name string, n int) (val uint64) {
	val = r.Read(n)
	fmt.Printf("Read %d (%q) = %x\n", n, name, val)
	return val
}

// Decode decodes an aribtrary n-bit big-endian number into dst
// dst should be a pointer to any integer or boolean value up
// to 64 bits wide
func (r *Reader) Decode(dst any, n int) (val uint64) {
	val = r.Read(n)
	switch p := dst.(type) {
	case *uint64:
		*p = val
	case *uint32:
		*p = uint32(val)
	case *int:
		*p = int(val)
	case *uint:
		*p = uint(val)
	case *bool:
		*p = val != 0
	case *uint16:
		*p = uint16(val)
	case *uint8:
		*p = uint8(val)
	case *int64:
		*p = int64(val)
	case *int32:
		*p = int32(val)
	case *int16:
		*p = int16(val)
	case *int8:
		*p = int8(val)
	}
	return
}

func (r *Reader) Ignore(n int) (val uint64) {
	return r.ReadPrint("Ignored", n)
}

func (r *Reader) ok() bool {
	return r.err == nil
}

func (r *Reader) Err() error {
	return r.err
}

// Read reads n bits and returns it as a uint64, if the
// read is not byte-aligned, up to calls to read may be issued
// recursively
func (r *Reader) Read(n int) (val uint64) {
	if n < 0 || r.at+n > r.len {
		r.err = io.EOF
		return 0
	}
	i := r.at / 8 // byte offset
	m := r.at % 8 // bit offset

	val = readBE(r.b[i:]) << m
	r.at += n
	if extra := 64 - int(n+m); extra < 0 {
		// read over 64+7 bits, so issue another read call
		// to get the rest of the data if there's room
		r.at += extra
		val |= r.Read(-extra)
	} else {
		// read a value less than 64 bits, shift it into
		// its intended representation
		val >>= 64 - n
	}

	return val
}

// Peek looks ahead up to 64 bits in the reader without advancing it
func (r *Reader) Peek(n int) (val uint64) {
	at, err := r.at, r.err
	val = r.Read(n)
	r.at, r.err = at, err
	return
}

// Offset returns the current bit offset of the reader, or the
// number of bits read. Divide by 8 for the byte offset.
func (r *Reader) Offset() int {
	return r.at
}

// readBE reads bytes in the buffer into a uint64
// the buffer p can be less than 64-bits
func readBE(p []byte) (n uint64) {
	if len(p) >= 8 {
		return binary.BigEndian.Uint64(p)
	}
	for i := 0; i < len(p); i++ {
		n |= uint64(p[i]) << (8 * (7 - i))
	}
	return n
}
