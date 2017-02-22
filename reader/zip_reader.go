package reader

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"strings"
)

// 打开一个 epub 文件, 能够读取其中的文件

type ZipReader struct {
	file *os.File
	zip  *zip.Reader
}

// func (zr *ZipReader) GetFile(filePath string) ([]byte, error) {
// 	return nil, nil
// }

func (zr *ZipReader) contains(filePath string) bool {
	if zr.zip == nil {
		return false
	}
	for _, f := range zr.zip.File {
		if f.Name == filePath {
			return true
		}
	}
	return false
}

// Close closes the epub file
func (e *ZipReader) Close() {
	if e.file != nil {
		e.file.Close()
	}
}

// OpenFile opens a file inside the epub
func (e *ZipReader) OpenFile(name string) (io.ReadCloser, error) {
	return openFile(e.zip, name)
}

// Open opens an existing epub
func NewZipReader(path string) (e *ZipReader, err error) {
	e = new(ZipReader)
	e.file, err = os.Open(path)
	if err != nil {
		return
	}
	fileInfo, err := e.file.Stat()
	if err != nil {
		return
	}
	err = e.load(e.file, fileInfo.Size())
	if err != nil {
		return
	}

	if e.contains("mimetype") == false {
		err = errors.New("epub format error, no mimetype file")
		return
	}

	if e.contains("META-INF/container.xml") == false {
		err = errors.New("epub format error, no META-INF/container.xml file")
		return
	}
	return
}

// // Load loads an epub from an io.ReaderAt
// func Load(r io.ReaderAt, size int64) (e *Epub, err error) {
// 	e = new(Epub)
// 	e.file = nil
// 	// e.CharactorStatistics = CharactorStatisticArray{}
// 	err = e.load(r, size)
// 	return
// }

func (e *ZipReader) load(r io.ReaderAt, size int64) (err error) {
	e.zip, err = zip.NewReader(r, size)
	if err != nil {
		return
	}

	// e.rootPath, err = getRootPath(e.zip)
	// if err != nil {
	// 	return
	// }
	return
}

func openFile(file *zip.Reader, path string) (io.ReadCloser, error) {
	for _, f := range file.File {
		if f.Name == path {
			return f.Open()
		}
	}

	pathLower := strings.ToLower(path)
	for _, f := range file.File {
		if strings.ToLower(f.Name) == pathLower {
			return f.Open()
		}
	}

	return nil, errors.New("File " + path + " not found")
}

// // OpenFileId opens a file from it's id
// //
// // The id of the files often appears on metadata fields
// func (e Epub) OpenFileId(id string) (io.ReadCloser, error) {
// 	path := e.opf.filePath(id)
// 	return openFile(e.zip, e.rootPath+path)
// }

// type containerXML struct {
// 	// FIXME: only support for one rootfile, can it be more than one?
// 	Rootfile rootfile `xml:"rootfiles>rootfile"`
// }
// type rootfile struct {
// 	Path string `xml:"full-path,attr"`
// }

// func openOPF(file *zip.Reader) (io.ReadCloser, error) {

// 	path, err := getOpfPath(file)

// 	if err != nil {
// 		return nil, err
// 	}
// 	log.Println("opf path => ", path)
// 	return openFile(file, path)
// }

// func getRootPath(file *zip.Reader) (string, error) {
// 	opfPath, err := getOpfPath(file)
// 	if err != nil {
// 		return "", err
// 	}
// 	pathDir := path.Dir(opfPath)
// 	if pathDir == "." {
// 		return "", nil
// 	} else {
// 		return path.Dir(opfPath) + "/", nil
// 	}
// }

// func getOpfPath(file *zip.Reader) (string, error) {
// 	f, err := openFile(file, "META-INF/container.xml")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer f.Close()

// 	var c containerXML
// 	err = decodeXML(f, &c)
// 	return c.Rootfile.Path, err
// }

// func decodeXML(file io.Reader, v interface{}) error {
// 	decoder := xml.NewDecoder(file)
// 	decoder.Entity = xml.HTMLEntity
// 	decoder.CharsetReader = charset.NewReaderLabel
// 	return decoder.Decode(v)
// }
