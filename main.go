package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var from, to string
var limit, offset int

func init() {
	flag.StringVar(&from, "from", "", "copy from")
	flag.StringVar(&to, "to", "", "copy to")
	flag.IntVar(&limit, "limit", 0, "limit in input file")
	flag.IntVar(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	validationErr := validateInputParams()
	if validationErr != nil {
		fmt.Println(validationErr)
		return
	}

	written, err := Copy(from, to, limit, offset)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Done! %d bytes successfully copied\n", written)
}

func Copy(from string, to string, limit int, offset int) (int64, error) {
	var written int64
	var err error

	fileFrom, err := os.Open(from)
	if err != nil {
		return 0, err
	}

	fileTo, err := os.Create(to)
	if err != nil {
		return 0, err
	}

	defer func() {
		_ = fileFrom.Close()
		_ = fileTo.Close()
	}()

	fileToSize, err := getResultFileSize(fileFrom, offset, limit)
	if err != nil {
		return 0, err
	}

	if offset > 0 {
		_, err := fileFrom.Seek(int64(offset), 0)
		if err != nil {
			return 0, err
		}
	}

	bar := pb.Full.Start64(fileToSize)
	barReader := bar.NewProxyReader(fileFrom)

	if limit > 0 {
		written, err = io.CopyN(fileTo, barReader, int64(limit))
	} else {
		written, err = io.Copy(fileTo, barReader)
	}
	
	if err != nil {
		return 0, err
	}

	bar.Finish()

	return written, nil
}

func validateInputParams() error {
	if from == "" {
		return errors.New("-from param required")
	}

	if to == "" {
		return errors.New("-to param required")
	}

	return nil
}

func getResultFileSize(file *os.File, offset int, limit int) (int64, error){
	if limit > 0 {
		return int64(limit), nil
	}
	
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	fileSize := fileStat.Size()
	
	if int64(offset) > fileSize {
		return 0, errors.New(fmt.Sprintf("Offset %d is larger than fileFrom size %d", offset, fileSize))
	}
	
	return fileStat.Size() - int64(offset), nil
}
