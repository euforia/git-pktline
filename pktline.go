// Package pktline implements pkt-line format encoding used by Git's transfer protocol.
// https://github.com/git/git/blob/master/Documentation/technical/protocol-common.txt
package pktline

import (
	"errors"
	"io"
)

// Errors returned by methods in package pktline.
var (
	ErrShortRead   = errors.New("input is too short")
	ErrInputExcess = errors.New("input is too long")
	ErrTooLong     = errors.New("too long payload")
	ErrInvalidLen  = errors.New("invalid length")
)

const (
	headLen = 4
	maxLen  = 65524 // 65520 bytes of data
)

// EncoderDecoder serves as both Encoder and Decoder.
type EncoderDecoder struct {
	Encoder
	Decoder
}

// NewEncoderDecoder constructs pkt-line encoder/decoder.
func NewEncoderDecoder(rw io.ReadWriter) *EncoderDecoder {
	return &EncoderDecoder{*NewEncoder(rw), *NewDecoder(rw)}
}
