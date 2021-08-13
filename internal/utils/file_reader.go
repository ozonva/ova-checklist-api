package utils

import (
	"io"
	"os"
)

type FileReadingVisitor interface {
	OnOpenFail(filename string, err error)
	OnOpenSuccess(filename string, reader io.Reader)
	OnCloseFail(filename string, err error)
	OnCloseSuccess(filename string)
}

type FileOpeningStrategy interface {
	Open(filename string) (io.ReadCloser, error)
}

// FileSystemOpeningStrategy implements FileOpeningStrategy
type FileSystemOpeningStrategy struct {}

func (fs *FileSystemOpeningStrategy) Open(filename string) (io.ReadCloser, error) {
	return os.Open(filename)
}

func closeFile(visitor FileReadingVisitor, file io.ReadCloser, filename string) {
	if err := file.Close(); err != nil {
		visitor.OnCloseFail(filename, err)
		return
	}
	visitor.OnCloseSuccess(filename)
}

func ReadFiles(fs FileOpeningStrategy, visitor FileReadingVisitor, filenames ...string) {
	for _, filename := range filenames {
		func() {
			file, err := fs.Open(filename)
			if err != nil {
				visitor.OnOpenFail(filename, err)
				return
			}
			defer closeFile(visitor, file, filename)
			visitor.OnOpenSuccess(filename, file)
		}()
	}
}
