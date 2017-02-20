// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package epubgo

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// Epub holds all the data of the ebook
type Epub struct {
	file     *os.File
	zip      *zip.Reader
	rootPath string
	metadata MetaDataList
	opf      *xmlOPF
	NCX      *XmlNCX
}

type MetaDataList map[string][]MdataElement
type MdataElement struct {
	content string
	attr    map[string]string
}

// Open opens an existing epub
func Open(path string) (e *Epub, err error) {
	e = new(Epub)
	e.file, err = os.Open(path)
	if err != nil {
		return
	}
	fileInfo, err := e.file.Stat()
	if err != nil {
		return
	}
	err = e.load(e.file, fileInfo.Size())
	return
}

// Load loads an epub from an io.ReaderAt
func Load(r io.ReaderAt, size int64) (e *Epub, err error) {
	e = new(Epub)
	e.file = nil
	// e.CharactorStatistics = CharactorStatisticArray{}
	err = e.load(r, size)
	return
}

func (e *Epub) load(r io.ReaderAt, size int64) (err error) {
	e.zip, err = zip.NewReader(r, size)
	if err != nil {
		return
	}

	e.rootPath, err = getRootPath(e.zip)
	if err != nil {
		return
	}

	return e.parseFiles()
}

func isTextContent(mediaType string) bool {
	return mediaType == "application/xhtml+xml"
}

func (e *Epub) CountFileCharactor(ele *manifest) error {

	if isTextContent(ele.MediaType) {
		if len(ele.Href) > 0 {
			reader_closer, err := e.OpenFile(ele.Href)
			if err != nil {
				return err
			}
			defer reader_closer.Close()
			content, err := getHtmlContent(reader_closer)
			if err != nil {
				return err
			}
			ele.CharactorCount = len(content)
			return nil
		}
	}
	return nil
}

func (e *Epub) parseFiles() error {
	opfFile, err := openOPF(e.zip)
	if err != nil {
		return err
	}
	defer opfFile.Close()
	e.opf, err = parseOPF(opfFile)
	if err != nil {
		return err
	}

	if e.opf.Manifest != nil && len(e.opf.Manifest) > 0 {
		for _, ele := range e.opf.Manifest {
			err = e.CountFileCharactor(ele)
			if err != nil {
				return err
			}
		}
	}

	e.metadata = e.opf.toMData()
	ncxPath := e.opf.ncxPath()
	if ncxPath != "" {
		ncx, err := e.OpenFile(ncxPath)
		if err != nil {
			return errors.New("Can't open the NCX file")
		}
		defer ncx.Close()
		e.NCX, err = parseNCX(ncx)
		if err != nil {
			return err
		}
		e.NCX.NavMap.SetContentCount(func(file string) int {
			for _, ele := range e.opf.Manifest {
				if ele.Href == file {
					return ele.CharactorCount
				}
			}
			return 0
		})
	}
	return nil
}

// Close closes the epub file
func (e Epub) Close() {
	if e.file != nil {
		e.file.Close()
	}
}

// OpenFile opens a file inside the epub
func (e Epub) OpenFile(name string) (io.ReadCloser, error) {
	return openFile(e.zip, e.rootPath+name)
}

// OpenFileId opens a file from it's id
//
// The id of the files often appears on metadata fields
func (e Epub) OpenFileId(id string) (io.ReadCloser, error) {
	path := e.opf.filePath(id)
	return openFile(e.zip, e.rootPath+path)
}

// Navigation returns a navigation iterator
func (e Epub) Navigation() (*NavigationIterator, error) {
	if e.NCX == nil {
		return nil, errors.New("There is no NCX file on the epub")
	}
	return newNavigationIterator(e.NCX.navMap())
}

// Spine returns a spine iterator
func (e Epub) Spine() (*SpineIterator, error) {
	return newSpineIterator(&e)
}

// Metadata returns the values of a metadata field
//
// The valid field names are:
//    title, language, identifier, creator, subject, description, publisher,
//    contributor, date, type, format, source, relation, coverage, rights, meta
func (e Epub) Metadata(field string) ([]string, error) {
	elem, ok := e.metadata[field]
	if ok {
		cont := make([]string, len(elem))
		for i, e := range elem {
			cont[i] = e.content
		}
		return cont, nil
	}

	return nil, errors.New("Field " + field + " don't exists")
}

// MetadataFields retunrs the list of metadata fields pressent on the current epub
func (e Epub) MetadataFields() []string {
	fields := make([]string, len(e.metadata))
	i := 0
	for k, _ := range e.metadata {
		fields[i] = k
		i++
	}
	return fields
}

// MetadataAttr returns the metadata attributes
//
// Returns: an array of maps of each attribute and it's value.
// The array has the fields on the same order than the Metadata method.
func (e Epub) MetadataAttr(field string) ([]map[string]string, error) {
	elem, ok := e.metadata[field]
	if ok {
		attr := make([]map[string]string, len(elem))
		for i, e := range elem {
			attr[i] = e.attr
		}
		return attr, nil
	}

	return nil, errors.New("Field " + field + " don't exists")
}

func (e *Epub) FileManifest(file string) *manifest {
	for _, ele := range e.opf.Manifest {
		if ele.Href == file {
			return ele
		}
	}
	return nil
}

func getHtmlContent(reader io.Reader) ([]rune, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	src := string(b)
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	body_head := strings.Index(src, "<body>")
	body_tail := strings.Index(src, "</body>")
	if body_head >= 0 && body_tail >= 0 && body_tail > body_head {
		src = src[body_head+6 : body_tail]
	}
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	src = strings.Replace(src, "\n", "", -1)
	src = strings.TrimSpace(src)
	src_rune := []rune(src)

	return src_rune, nil
}
