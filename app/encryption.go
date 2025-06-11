package app

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"
)

// EncryptedWriter wraps w with an CTR cipher stream.
func EncryptedWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {

	// generate random initial value
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	// write clear IV to allow for decryption
	n, err := w.Write(iv)
	if err != nil || n != len(iv) {
		return nil, errors.New("could not write initial value")
	}

	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// EncryptedReader wraps r with an CTR cipher stream.
func EncryptedReader(key string, r io.Reader) (*cipher.StreamReader, error) {

	// read initial value
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if err != nil || n != len(iv) {
		return nil, errors.New("could not read initial value")
	}

	block, err := newBlock(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	return &cipher.StreamReader{S: stream, R: r}, nil
}

func newBlock(key string) (cipher.Block, error) {
	hash := md5.Sum([]byte(key))
	block, err := aes.NewCipher(hash[:])
	if err != nil {
		return nil, err
	}
	return block, nil
}
