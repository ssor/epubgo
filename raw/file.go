// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package raw

import (
	"encoding/xml"

	"golang.org/x/net/html/charset"
	// "github.com/golang/net/html/charset"
	"io"
)

type containerXML struct {
	// FIXME: only support for one rootfile, can it be more than one?
	Rootfile rootfile `xml:"rootfiles>rootfile"`
}
type rootfile struct {
	Path string `xml:"full-path,attr"`
}

func decodeXML(file io.Reader, v interface{}) error {
	decoder := xml.NewDecoder(file)
	decoder.Entity = xml.HTMLEntity
	decoder.CharsetReader = charset.NewReaderLabel
	return decoder.Decode(v)
}

func isTextContent(mediaType string) bool {
	return mediaType == "application/xhtml+xml"
}

// func openFile(file *zip.Reader, path string) (io.ReadCloser, error) {
// 	for _, f := range file.File {
// 		if f.Name == path {
// 			return f.Open()
// 		}
// 	}

// 	pathLower := strings.ToLower(path)
// 	for _, f := range file.File {
// 		if strings.ToLower(f.Name) == pathLower {
// 			return f.Open()
// 		}
// 	}

// 	return nil, errors.New("File " + path + " not found")
// }
