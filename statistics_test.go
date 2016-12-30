package epubgo

import (
	"testing"
	// . "github.com/smartystreets/goconvey/convey"
)

var (
	gcdxy_epub = "testdata/gcdxy.epub"

	statistics = []CharactorStatistic{
		{
			File:   "b_content_xhtml",
			Length: 220,
		},
		{
			File:   "chapter_00001_xhtml",
			Length: 3,
		},
		{
			File:   "chapter_00002_xhtml",
			Length: 685,
		},
		{
			File:   "chapter_00003_xhtml",
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

	if len(f.CharactorStatistics) < 2 {
		t.Errorf("count err")
	}

	test_statistic := func(sta CharactorStatistic) (bool, int) {
		for _, statistic := range f.CharactorStatistics {
			if statistic.File == sta.File {
				if statistic.Length == sta.Length {
					return true, sta.Length
				} else {
					return false, statistic.Length
				}
			}
		}
		return false, -1
	}

	for _, statistic := range statistics {
		b, n := test_statistic(statistic)
		if b == false {
			t.Errorf("file %s expect %d and acturlly %d",
				statistic.File, statistic.Length, n)
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
