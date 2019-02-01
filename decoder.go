package pktline

import (
	"bytes"
	"io"
	"strconv"
)

// Decode parses a single pkt-line and returns it's payload.
// Input longer than a single pkt-line is considered an error.
func Decode(line []byte) (payload []byte, err error) {
	buffer := bytes.NewBuffer(line)
	decoder := NewDecoder(buffer)
	err = decoder.Decode(&payload)
	if err != nil {
		return
	}
	if buffer.Len() != 0 {
		err = ErrInputExcess
		return
	}
	return
}

// Decoder decodes input in pkt-line format.
type Decoder struct {
	r io.Reader
}

// NewDecoder constructs a new pkt-line decoder.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads a single pkt-line and stores its payload.
// Flush-pkt is causes *payload to be nil.
func (d *Decoder) Decode(payload *[]byte) error {
	head := make([]byte, headLen)
	_, err := io.ReadFull(d.r, head)
	if err == io.ErrUnexpectedEOF {
		return ErrShortRead
	}
	if err != nil {
		return err
	}
	lineLen, err := strconv.ParseInt(string(head), 16, 16)
	if err != nil {
		return err
	}
	if lineLen == 0 { // flush-pkt
		*payload = nil
		return nil
	}
	if lineLen < headLen {
		return ErrInvalidLen
	}
	*payload = make([]byte, lineLen-headLen)
	if lineLen == headLen { // empty line
		return nil
	}
	_, err = io.ReadFull(d.r, *payload)
	if err == io.ErrUnexpectedEOF {
		return ErrShortRead
	}
	if err != nil {
		return err
	}
	return nil
}

// DecodeUntilFlush decodes pkt-line messages until it encounters flush-pkt.
// The flush-pkt is not included in output.
// If error is not nil, output contains data that was read before the error occured.
func (d *Decoder) DecodeUntilFlush(lines *[][]byte) (err error) {
	*lines = make([][]byte, 0)
	for {
		var l []byte
		err = d.Decode(&l)
		if err != nil {
			return
		}
		if l == nil { // flush-pkt
			return
		}
		*lines = append(*lines, l)
	}
}
