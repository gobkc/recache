package recache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func GzipEncode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, err := gzipWriter.Write(input)
	if err != nil {
		_ = gzipWriter.Close()
		return nil, err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GzipDecode(input []byte) ([]byte, error) {
	bytesReader := bytes.NewReader(input)
	gzipReader, err := gzip.NewReader(bytesReader)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = gzipReader.Close()
	}()
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(gzipReader); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Marshal(data interface{}) ([]byte, error) {
	marshalData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	gzipData, err := GzipEncode(marshalData)
	if err != nil {
		return nil, err
	}
	return gzipData, err
}

func Unmarshal(input []byte, output interface{}) error {
	decodeData, err := GzipDecode(input)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodeData, output)
	if err != nil {
		return err
	}
	return nil
}
