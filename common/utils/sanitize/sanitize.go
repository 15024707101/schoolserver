package sanitize

import (
	"bytes"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"log"
	"reflect"
	"regexp"
	"strings"
)

const (
	// Can be more then NumCPU, because workers are blocked on net IO.
	numWorkers        = 10
	numWorkerChannels = numWorkers * 2
)

var (
	prefixRx = regexp.MustCompile("^.*<body>")
	suffixRx = regexp.MustCompile("</body>.*$")
)

var strict *bluemonday.Policy

func init() {
	strict = bluemonday.StrictPolicy()
}

func cleanupHTML(b []byte) []byte {
	node, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		log.Println(err)
		return b
	}
	var bb bytes.Buffer
	err = html.Render(&bb, node)
	if err != nil {
		log.Println(err)
		return b
	}
	return bb.Bytes()
}

// Sanitize 去掉html
func Sanitize(htmlInput string) string {
	str := strings.TrimSpace(strict.Sanitize(html.UnescapeString(htmlInput)))
	str = strings.Replace(str, "\n", "", -1)
	return str
}

// customBluemondayPolicy  自定义过滤策略
func customBluemondayPolicy() *bluemonday.Policy {
	p := bluemonday.StrictPolicy()
	p.AllowStandardAttributes()
	p.AllowStandardURLs()
	p.AllowElements("article", "aside")
	p.AllowElements("section")
	p.AllowElements("summary")
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("hgroup")
	p.AllowAttrs("cite").OnElements("blockquote")
	p.AllowElements("br", "div", "hr", "p", "span", "wbr")
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em",
		"figcaption", "mark", "s", "samp", "strong", "sub", "sup", "var")
	p.AllowAttrs("cite").OnElements("q")
	p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")
	p.AllowLists()
	p.AllowTables()
	p.SkipElementsContent("select")
	return p
}

//  去掉html属性
func HTMLPurifier(htmlInput string) string {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(htmlInput)
	return html
}

// 去掉所有html
func HTMLPurifierAll(htmlInput string) string {
	p := bluemonday.NewPolicy()
	html := p.Sanitize(htmlInput)
	return html
}

// SanitizeAllStringFields 净化结构体的字符串字段
func SanitizeAllStringFields(i interface{}) error {
	ifv := reflect.ValueOf(i)
	if ifv.Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}
	ift := reflect.Indirect(ifv).Type()
	if ift.Kind() != reflect.Struct {
		return errors.New("not point to a struct")
	}
	for i := 0; i < ift.NumField(); i++ {
		v := ift.Field(i)
		el := reflect.Indirect(ifv.Elem().FieldByName(v.Name))
		switch el.Kind() {
		case reflect.Slice:
			continue
		case reflect.Struct:
			continue
		case reflect.String:
			if el.CanSet() {
				el.SetString(Sanitize(el.String()))
			}
		}
	}
	return nil
}
