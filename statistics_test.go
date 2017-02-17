package epubgo

import (
	"testing"

	// . "github.com/smartystreets/goconvey/convey"
	"fmt"
)

type CharactorStatistic struct {
	File   string
	Length int
}

var (
	gcdxy_epub = "testdata/gcdxy.epub"

	statistics = []CharactorStatistic{
		{
			File:   "b_content.xhtml",
			Length: 220,
		},
		{
			File:   "chapter_00001.xhtml",
			Length: 3,
		},
		{
			File:   "chapter_00002.xhtml",
			Length: 685,
		},
		{
			File:   "chapter_00003.xhtml",
			Length: 1304,
		},
	}
)

func TestStatistics(t *testing.T) {
	f, err := Open(gcdxy_epub)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	total_count := 0
	for _, ele := range f.opf.Manifest {
		total_count += ele.CharactorCount
		if ele.CharactorCount > 0 {
			fmt.Println(ele.Href, " : ", ele.CharactorCount)
		}
	}
	fmt.Println("tatal: ", total_count)

	for _, statistic := range statistics {
		manifest := f.FileManifest(statistic.File)
		if manifest == nil {
			t.Fatal("file ", statistic.File, " expected")
		}
		if manifest.CharactorCount != statistic.Length {
			t.Fatalf("file %s expect %d and acturlly %d",
				statistic.File, statistic.Length, manifest.CharactorCount)
		}
	}
}

// func TestGetHtmlContent(t *testing.T) {

// 	Convey("get html content test", t, func() {
// 		for _, content_count := range statistics {
// 			file, err := os.Open(fmt.Sprintf("./testdata/gcdxy/ops/%s", content_count.File))
// 			if err != nil {
// 				t.Errorf("err: %s", err)
// 			}

// 			result, err := GetHtmlContent(file)
// 			// fmt.Println(len(result), " => ", string(result))
// 			So(len(result), ShouldEqual, content_count.Length)
// 		}
// 	})
// }
