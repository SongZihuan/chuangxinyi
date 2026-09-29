package utils

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	errors "github.com/wuntsong-org/wterrors"
	"io"
)

func JsonMarshal(v interface{}) ([]byte, errors.WTError) {
	buf := bytes.NewBuffer(make([]byte, 0))
	e := jsoniter.NewEncoder(buf)
	err := e.Encode(v)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	res, err := io.ReadAll(buf)
	if err != nil {
		return nil, errors.WarpQuick(err)
	}

	return res, nil
}

func JsonUnmarshal(data []byte, v interface{}) errors.WTError {
	d := jsoniter.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	err := d.Decode(v)
	if err != nil {
		return errors.WarpQuick(err)
	}

	return nil
}
