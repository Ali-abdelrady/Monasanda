package osioutil

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"

	"github.com/klauspost/compress/zstd"
)

func GzipCompressData(parData []byte) ([]byte, error) {
	var pvBuffer bytes.Buffer
	pvGzipWriter := gzip.NewWriter(&pvBuffer)

	_, pvError := pvGzipWriter.Write(parData)

	if pvError == nil {
		pvError = pvGzipWriter.Close()
	}

	return pvBuffer.Next(pvBuffer.Len()), pvError
}

func GzipDecompressData(parData []byte) ([]byte, error) {
	var pvBuffer bytes.Buffer
	_, pvError := pvBuffer.Write(parData)
	if pvError != nil {
		return nil, pvError
	}

	pvGzipReader, pvError := gzip.NewReader(&pvBuffer)
	if pvError != nil {
		return nil, pvError
	}

	paContent, pvError := io.ReadAll(pvGzipReader)
	if pvError != nil {
		return nil, pvError
	}

	pvGzipReader.Close()

	return paContent, pvError
}

func ZLibCompressData(parData []byte) []byte {
	var b bytes.Buffer
	pvWriter := zlib.NewWriter(&b)
	pvWriter.Write(parData)
	pvWriter.Close()

	return b.Bytes()
}

func ZLibDecompressData(parData []byte) ([]byte, error) {
	b := bytes.NewReader(parData)

	pvReader, pvError := zlib.NewReader(b)
	if pvError != nil {
		return []byte{}, pvError
	}

	var b1 bytes.Buffer

	_, pvError = io.Copy(&b1, pvReader)
	if pvError != nil {
		return []byte{}, pvError
	}

	pvReader.Close()

	return b1.Bytes(), nil
}

func ZSTDCompressData(parData []byte) ([]byte, error) {
	// https://pkg.go.dev/github.com/klauspost/compress/zstd#example-ZipCompressor

	pvBuf := bytes.Buffer{}

	pvEncoder, pvError := zstd.NewWriter(&pvBuf)
	if pvError != nil {
		return []byte{}, fmt.Errorf("[CMPRS001] %v", pvError)
	}

	_, pvError = io.Copy(pvEncoder, bytes.NewBuffer(parData))
	if pvError != nil {
		pvEncoder.Close()
		return []byte{}, fmt.Errorf("[CMPRS002] %v", pvError)
	}

	pvError = pvEncoder.Close()
	if pvError != nil {
		return []byte{}, fmt.Errorf("[CMPRS003] %v", pvError)
	}

	return pvBuf.Bytes(), nil
}

func ZSTDDecompressData(parData []byte) ([]byte, error) {
	pvDecoder, pvError := zstd.NewReader(bytes.NewBuffer(parData))
	if pvError != nil {
		return []byte{}, fmt.Errorf("[CMPRS004] %v", pvError)
	}
	defer pvDecoder.Close()

	var pvBuf bytes.Buffer
	_, pvError = io.Copy(&pvBuf, pvDecoder)
	if pvError != nil {
		return []byte{}, fmt.Errorf("[CMPRS005] %v", pvError)
	}

	return pvBuf.Bytes(), nil
}
