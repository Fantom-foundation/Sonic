package ethapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

const (
	// bufferStartSize default buffer start size in bytes
	bufferStartSize = 100000
)

var (
	// responseSizeLimit limits maximum size of RPC response
	responseSizeLimit int
	// ErrMaxResponseSize is returned if size of the RPC call
	// response is over defined limit
	ErrMaxResponseSize = errors.New("response size exceeded")
	// ErrResponseTooLarge is returned when the result buffer cannot be expanded
	ErrResponseTooLarge = errors.New("response too large")
	// ErrCannotInitBuffer is returned when initialization of the buffer failed
	ErrCannotInitBuffer = errors.New("cannot initialize write buffer")
)

// SetResponseSizeLimit sets maximum size of RPC response in bytes
func SetResponseSizeLimit(limit int) {
	responseSizeLimit = limit
}

// JsonResultBuffer is a bytes buffer for jsonRawMessage result
type JsonResultBuffer struct {
	bytes.Buffer
}

// NewJsonResultBuffer creates new bytes buffer
func NewJsonResultBuffer() (b *JsonResultBuffer, err error) {
	defer func() {
		if recover() != nil {
			err = ErrCannotInitBuffer
		}
	}()
	b = &JsonResultBuffer{}
	// grow buffer to default start size
	b.Grow(bufferStartSize)
	if err := b.writeString("["); err != nil {
		return nil, err
	}
	return b, nil
}

// GetResult returns finalized byte array
func (b *JsonResultBuffer) GetResult() ([]byte, error) {
	if err := b.writeString("]"); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// AddObject marshals and adds object into buffer.
// In case that buffer size is over defined limit, error is returned
func (b *JsonResultBuffer) AddObject(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if responseSizeLimit > 0 && b.Len()+len(res) > responseSizeLimit {
		return ErrMaxResponseSize
	}
	if b.Len() > 1 {
		if err := b.writeString(","); err != nil {
			return err
		}
	}
	if err = b.writeString(string(res)); err != nil {
		return err
	}

	return nil
}

// writeString appends the contents of s to the buffer and handle possible panic
func (b *JsonResultBuffer) writeString(s string) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrResponseTooLarge
		}
	}()
	b.WriteString(s)
	return nil
}
