package pktline

import (
	"fmt"
	"io"
)

// Encode returns payload encoded in pkt-line format.
func Encode(payload []byte) ([]byte, error) {
	if payload == nil {
		return []byte("0000"), nil
	}
	if len(payload)+headLen > maxLen {
		return nil, ErrTooLong
	}
	head := []byte(fmt.Sprintf("%04x", len(payload)+headLen))
	return append(head, payload...), nil
}

// Encoder encodes payloads in pkt-line format.
type Encoder struct {
	w io.Writer
}

// NewEncoder constructs a pkt-line encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode encodes payload and writes it to encoder output.
// If payload is nil, writes flush-pkt.
func (e *Encoder) Encode(payload []byte) error {
	line, err := Encode(payload)
	if err != nil {
		return err
	}
	_, err = e.w.Write(line)
	if err != nil {
		return err
	}
	return nil
}
