package reader

import "testing"
import "fmt"

const (
	bookPath       = "../testdata/a_dogs_tale.epub"
	noNCXPath      = "../testdata/noncx.epub"
	invalidNCXPath = "../testdata/invalidncx.epub"
	fileCapsPath   = "../testdata/fileCaps.epub"
)

func TestZipReader(t *testing.T) {
	zipReader, err := NewZipReader(bookPath)
	if err != nil {
		t.Fatalf("Open(%v) return an error: %v", bookPath, err)
	}
	// spew.Dump(zipReader)
	for index, zf := range zipReader.zip.File {
		fmt.Println(index, " : ", zf.Name)
	}
	zipReader.Close()
}
