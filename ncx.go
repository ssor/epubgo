// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package epubgo

import (
	"io"
	"strings"
)

type XmlNCX struct {
	NavMap NavPointArray `xml:"navMap>navPoint"`
}
type NavPoint struct {
	Text      string        `xml:"navLabel>text"`
	Content   content       `xml:"content"`
	NavPoints NavPointArray `xml:"navPoint"`
	Level     int
}

// func (np *NavPoint) SetContentCount(file string, count int) bool {
func (np *NavPoint) SetContentCount(getCount func(string) int) {
	if len(np.Content.Src) > 0 {
		np.Content.Count = getCount(np.Content.Src)
	}
	if np.NavPoints != nil && len(np.NavPoints) > 0 {
		np.NavPoints.SetContentCount(getCount)
	}
}

func (np *NavPoint) ResetLevel(level int) {
	np.Level = level
	if np.NavPoints != nil && len(np.NavPoints) > 0 {
		np.NavPoints.ResetLevel(level + 1)
	}
}

func (np *NavPoint) TotalContentLength() int {
	total := np.Content.Count

	if np.NavPoints != nil {
		total += np.NavPoints.TotalContentLength()
	}
	return total
}

type content struct {
	Src   string `xml:"src,attr"`
	Count int
}

type NavPointArray []*NavPoint

func (nps NavPointArray) ResetLevel(level int) {
	for _, np := range nps {
		np.ResetLevel(level + 1)
	}
}

func (nps NavPointArray) SetContentCount(getCount func(string) int) {
	if nps != nil && len(nps) > 0 {
		for _, np := range nps {
			np.SetContentCount(getCount)
		}
	}
}

func (nps NavPointArray) TotalContentLength() int {
	total := 0

	if nps != nil && len(nps) > 0 {
		for _, np := range nps {
			total += np.TotalContentLength()
		}
	}
	return total
}

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

func (point *NavPoint) Title() string {
	return point.Text
}

func (np *NavPoint) LevelTitle() string {
	return strings.Repeat(" ", np.Level*4) + np.Text
}

func (point NavPoint) URL() string {
	return point.Content.Src
}

func (point NavPoint) Children() NavPointArray {
	return point.NavPoints
}
