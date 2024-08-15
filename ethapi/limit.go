package ethapi

import (
	"bytes"
	"encoding/json"
	"errors"
)

var (
	// responseSizeLimit limits maximum size of RPC response
	responseSizeLimit int
	// ErrMaxResponseSize is returned if size of the RPC call
	// response is over defined limit
	ErrMaxResponseSize = errors.New("response size exceeded")
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
func NewJsonResultBuffer() *JsonResultBuffer {
	b := JsonResultBuffer{}
	b.WriteString("[")
	return &b
}

// GetResult returns finalized byte array
func (b *JsonResultBuffer) GetResult() []byte {
	b.WriteString("]")
	return b.Bytes()
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
		b.WriteString(",")
	}
	b.Write(res)
	return nil
}
