package raw

import (
	"errors"
	"io"
	"path"

	"github.com/ssor/epubgo/reader"
)

type Reader interface {
	OpenFile(name string) (io.ReadCloser, error)
	Close()
}

func NewEpub(path string) (*Epub, error) {
	epub_reader, err := reader.NewZipReader(path)
	if err != nil {
		return nil, err
	}
	e := &Epub{
		reader: epub_reader,
	}
	e.rootPath, err = e.getRootPath()
	if err != nil {
		return nil, err
	}

	err = e.parseFiles()
	if err != nil {
		return nil, err
	}
	return e, nil
}

// Epub holds all the data of the ebook
type Epub struct {
	rootPath string
	metadata MetaDataList
	opf      *xmlOPF
	NCX      *XmlNCX
	reader   Reader
}

type MetaDataList map[string][]MdataElement
type MdataElement struct {
	content string
	attr    map[string]string
}

// func (e *Epub) CountFileCharactor(ele *manifest) error {

// 	if isTextContent(ele.MediaType) {
// 		if len(ele.Href) > 0 {
// 			reader_closer, err := e.OpenFile(ele.Href)
// 			if err != nil {
// 				return err
// 			}
// 			defer reader_closer.Close()
// 			content, err := getHtmlContent(reader_closer)
// 			if err != nil {
// 				return err
// 			}
// 			ele.CharactorCount = len(content)
// 			return nil
// 		}
// 	}
// 	return nil
// }

func (e *Epub) parseFiles() error {
	opfFile, err := e.openOPF()
	if err != nil {
		return err
	}
	defer opfFile.Close()
	e.opf, err = parseOPF(opfFile)
	if err != nil {
		return err
	}

	// if e.opf.Manifest != nil && len(e.opf.Manifest) > 0 {
	// 	for _, ele := range e.opf.Manifest {
	// 		err = e.CountFileCharactor(ele)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// }

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
		// e.NCX.NavMap.SetContentCount(func(file string) int {
		// 	for _, ele := range e.opf.Manifest {
		// 		if ele.Href == file {
		// 			return ele.CharactorCount
		// 		}
		// 	}
		// 	return 0
		// })
	}
	return nil
}

func (e *Epub) openOPF() (io.ReadCloser, error) {

	path, err := e.getOpfPath()

	if err != nil {
		return nil, err
	}
	// log.Println("opf path => ", path)
	return e.reader.OpenFile(path)
}

func (e *Epub) getRootPath() (string, error) {
	opfPath, err := e.getOpfPath()
	if err != nil {
		return "", err
	}
	pathDir := path.Dir(opfPath)
	if pathDir == "." {
		return "", nil
	} else {
		return path.Dir(opfPath) + "/", nil
	}
}

func (e *Epub) getOpfPath() (string, error) {
	f, err := e.reader.OpenFile("META-INF/container.xml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	var c containerXML
	err = decodeXML(f, &c)
	return c.Rootfile.Path, err
}

// Close closes the epub file
func (e Epub) Close() {
	e.reader.Close()
}

// OpenFile opens a file inside the epub
func (e Epub) OpenFile(name string) (io.ReadCloser, error) {
	return e.reader.OpenFile(path.Join(e.rootPath, name))
}

func (e *Epub) Files() []string {
	files := []string{}

	for _, item := range e.opf.Manifest {
		files = append(files, item.Href)
	}
	return files
}

// OpenFileId opens a file from it's id
//
// The id of the files often appears on metadata fields
func (e Epub) OpenFileId(id string) (io.ReadCloser, error) {
	path := e.opf.filePath(id)
	return e.OpenFile(path)
}

func (e *Epub) GetFileHrefByID(id string) string {
	return e.opf.filePath(id)
}

func (e *Epub) NavPoints() NavPointArray {
	if e.NCX == nil {
		return NavPointArray{}
	}

	return e.NCX.navMap()
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
func (e *Epub) MetadataAttr(field string) ([]map[string]string, error) {
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

// func getHtmlContent(reader io.Reader) ([]rune, error) {
// 	b, err := ioutil.ReadAll(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	src := string(b)
// 	//将HTML标签全转换成小写
// 	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
// 	src = re.ReplaceAllStringFunc(src, strings.ToLower)

// 	body_head := strings.Index(src, "<body>")
// 	body_tail := strings.Index(src, "</body>")
// 	if body_head >= 0 && body_tail >= 0 && body_tail > body_head {
// 		src = src[body_head+6 : body_tail]
// 	}
// 	//去除STYLE
// 	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
// 	src = re.ReplaceAllString(src, "")

// 	//去除SCRIPT
// 	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
// 	src = re.ReplaceAllString(src, "")

// 	//去除所有尖括号内的HTML代码，并换成换行符
// 	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
// 	src = re.ReplaceAllString(src, "\n")

// 	//去除连续的换行符
// 	re, _ = regexp.Compile("\\s{2,}")
// 	src = re.ReplaceAllString(src, "\n")

// 	src = strings.Replace(src, "\n", "", -1)
// 	src = strings.TrimSpace(src)
// 	src_rune := []rune(src)

// 	return src_rune, nil
// }

// // Open opens an existing epub
// func Open(path string) (e *Epub, err error) {
// 	e = new(Epub)
// 	e.file, err = os.Open(path)
// 	if err != nil {
// 		return
// 	}
// 	fileInfo, err := e.file.Stat()
// 	if err != nil {
// 		return
// 	}
// 	err = e.load(e.file, fileInfo.Size())
// 	return
// }

// // Load loads an epub from an io.ReaderAt
// func Load(r io.ReaderAt, size int64) (e *Epub, err error) {
// 	e = new(Epub)
// 	e.file = nil
// 	// e.CharactorStatistics = CharactorStatisticArray{}
// 	err = e.load(r, size)
// 	return
// }

// func (e *Epub) load(r io.ReaderAt, size int64) (err error) {
// 	e.zip, err = zip.NewReader(r, size)
// 	if err != nil {
// 		return
// 	}

// 	e.rootPath, err = getRootPath(e.zip)
// 	if err != nil {
// 		return
// 	}

// 	return e.parseFiles()
// }
