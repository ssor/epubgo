package epubgo

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

// The charactor count of each chapter
type CharactorStatistic struct {
	File   string
	Length int
}

func GetHtmlContent(reader io.Reader) ([]rune, error) {
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
