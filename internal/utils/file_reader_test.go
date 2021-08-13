package utils

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// implements io.ReadCloser
type inMemoryFile struct {
	content string
	readOffset int
	isCloseable bool
}

// implements FileOpeningStrategy
type inMemoryFileSystem struct {
	files map[string]inMemoryFile
}

// implements FileReadingVisitor
type inMemoryFileVisitor struct {
	openedFiles []string
	readContent []string
	closedFiles []string

	notOpenedFiles []string
	notClosedFiles []string
}

// methods for inMemoryFile
func newInMemoryFile(content string) inMemoryFile {
	return inMemoryFile{content, 0, true}
}

func newNotCloseableInMemoryFile(content string) inMemoryFile {
	return inMemoryFile{content, 0, false}
}

func (f *inMemoryFile) Read(dst []byte) (int, error) {
	if f.readOffset >= len(f.content) {
		return 0, nil
	}
	lenToRead := min(len(dst), len(f.content) - f.readOffset)
	bytes := []byte(f.content)
	end := f.readOffset + lenToRead
	src := bytes[f.readOffset:end]
	f.readOffset = end
	copy(dst, src)
	return lenToRead, nil
}

func (f *inMemoryFile) Close() error {
	if f.isCloseable {
		return nil
	}
	return os.ErrDeadlineExceeded
}

// methods for inMemoryFileSystem
func (fs *inMemoryFileSystem) Open(filename string) (io.ReadCloser, error) {
	if file, exists := fs.files[filename]; exists {
		return &file, nil
	}
	return nil, os.ErrNotExist
}

// methods for inMemoryFileVisitor
func (f *inMemoryFileVisitor) OnOpenFail(filename string, _ error) {
	f.notOpenedFiles = append(f.notOpenedFiles, filename)
}

func (f *inMemoryFileVisitor) OnOpenSuccess(filename string, file io.Reader) {
	f.openedFiles = append(f.openedFiles, filename)
	content := make([]byte, 0)
	buf := make([]byte, 100)
	for {
		bytesRead, err := file.Read(buf)
		if bytesRead == 0 || err != nil {
			break
		}
		content = append(content, buf[:bytesRead]...)
	}
	f.readContent = append(f.readContent, string(content))
}

func (f *inMemoryFileVisitor) OnCloseFail(filename string, _ error) {
	f.notClosedFiles = append(f.notClosedFiles, filename)
}

func (f *inMemoryFileVisitor) OnCloseSuccess(filename string) {
	f.closedFiles = append(f.closedFiles, filename)
}

func TestReadConfigFiles(t *testing.T) {
	var visitor inMemoryFileVisitor
	fs := inMemoryFileSystem{
		map[string]inMemoryFile{
			"hello.txt":              newInMemoryFile("Hello, world!"),
			"world.jpg":              newInMemoryFile("Never gonna give you up"),
			"i_am_not_closeable.exe": newNotCloseableInMemoryFile("it's gonna blow!"),
		},
	}

	ReadFiles(&fs, &visitor, "hello.txt", "world.jpg", "not_exists.py", "i_am_not_closeable.exe")
	assert.Equal(t, []string{"hello.txt", "world.jpg", "i_am_not_closeable.exe"}, visitor.openedFiles)
	assert.Equal(t, []string{"Hello, world!", "Never gonna give you up", "it's gonna blow!"}, visitor.readContent)
	assert.Equal(t, []string{"hello.txt", "world.jpg"}, visitor.closedFiles)
	assert.Equal(t, []string{"not_exists.py"}, visitor.notOpenedFiles)
	assert.Equal(t, []string{"i_am_not_closeable.exe"}, visitor.notClosedFiles)
}
