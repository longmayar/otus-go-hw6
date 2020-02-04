package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestCopyWithLimitAndOffset(t *testing.T) {
	fileFrom, err := ioutil.TempFile(os.TempDir(), "fileFrom")
	if err != nil {
		log.Fatal("Cannot create temporary fileFrom", err)
	}

	text := []byte("0123456789")
	if _, err = fileFrom.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	fileTo, err := ioutil.TempFile(os.TempDir(), "fileTo")
	if err != nil {
		log.Fatal("Cannot create temporary fileTo", err)
	}

	resultFileSize, err := Copy(fileFrom.Name(), fileTo.Name(), 2, 5)
	fileToStat, _ := fileTo.Stat()
	fileToSize := fileToStat.Size()
	fileToContent, _ := ioutil.ReadAll(fileTo)

	assert.Equal(t, int64(2), resultFileSize)
	assert.Equal(t, int64(2), fileToSize)
	assert.Equal(t, "56", string(fileToContent))

	defer func() {
		_ = os.Remove(fileFrom.Name())
		_ = os.Remove(fileTo.Name())
	}()
}

func TestCopyWithoutLimitAndOffset(t *testing.T) {
	fileFrom, err := ioutil.TempFile(os.TempDir(), "fileFrom")
	if err != nil {
		log.Fatal("Cannot create temporary fileFrom", err)
	}

	text := []byte("0123456789")
	if _, err = fileFrom.Write(text); err != nil {
		log.Fatal("Failed to write to temporary file", err)
	}

	fileTo, err := ioutil.TempFile(os.TempDir(), "fileTo")
	if err != nil {
		log.Fatal("Cannot create temporary fileTo", err)
	}

	resultFileSize, err := Copy(fileFrom.Name(), fileTo.Name(), 0, 0)
	fileToStat, _ := fileTo.Stat()
	fileToSize := fileToStat.Size()
	fileToContent, _ := ioutil.ReadAll(fileTo)

	assert.Equal(t, int64(10), resultFileSize)
	assert.Equal(t, int64(10), fileToSize)
	assert.Equal(t, "0123456789", string(fileToContent))

	defer func() {
		_ = os.Remove(fileFrom.Name())
		_ = os.Remove(fileTo.Name())
	}()
}