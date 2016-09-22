// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package epubgo

import (
	"io"
)

type XmlNCX struct {
	NavMap NavPointArray `xml:"navMap>navPoint"`
}
type NavPoint struct {
	Text      string        `xml:"navLabel>text"`
	Content   content       `xml:"content"`
	NavPoints NavPointArray `xml:"navPoint"`
}
type content struct {
	Src string `xml:"src,attr"`
}

type NavPointArray []*NavPoint

func parseNCX(ncx io.Reader) (*XmlNCX, error) {
	var n XmlNCX
	err := decodeXML(ncx, &n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (ncx XmlNCX) navMap() NavPointArray {
	return ncx.NavMap
}

func (point NavPoint) Title() string {
	return point.Text
}

func (point NavPoint) URL() string {
	return point.Content.Src
}

func (point NavPoint) Children() NavPointArray {
	return point.NavPoints
}
