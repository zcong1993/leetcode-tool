package update

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v2"
)

var (
	indexRegex      = regexp.MustCompile("@index (.+)")
	titleRegex      = regexp.MustCompile("@title (.+)")
	difficultyRegex = regexp.MustCompile("@difficulty (.+)")
	tagsRegex       = regexp.MustCompile("@tags (.+)")
	draftRegex      = regexp.MustCompile("@draft (.+)")
	linkRegex       = regexp.MustCompile("@link (.+)")
)

var (
	tableTpl = template.Must(template.New("table").Parse(tableStr))
)

type Meta struct {
	Index      string
	Title      string
	Difficulty string
	Tags       []string
	Draft      bool
	Fp         string
	Link       string
}

type Metas []*Meta

func (a Metas) Len() int      { return len(a) }
func (a Metas) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a Metas) Less(i, j int) bool {
	iIndex, _ := strconv.Atoi(a[i].Index)
	jIndex, _ := strconv.Atoi(a[j].Index)
	return iIndex < jIndex
}

type TableData struct {
	Metas Metas
	Total int
}

type TagMetas map[string](Metas)

func addMeta(tagMetas TagMetas, meta *Meta) {
	if meta == nil {
		return
	}
	for _, tag := range meta.Tags {
		if _, ok := tagMetas[tag]; !ok {
			tagMetas[tag] = make(Metas, 0)
		}
		tagMetas[tag] = append(tagMetas[tag], meta)
	}
	tagMetas["all"] = append(tagMetas["all"], meta)
}

func findTag(content []byte, r *regexp.Regexp) string {
	res := r.FindSubmatch(content)
	if len(res) < 2 {
		return ""
	}
	return string(res[1])
}

func findMeta(content []byte, fp string) *Meta {
	draft := findTag(content, draftRegex) == "" || findTag(content, draftRegex) == "true"
	if draft {
		return nil
	}
	tags := strings.Split(findTag(content, tagsRegex), ",")
	return &Meta{
		Index:      findTag(content, indexRegex),
		Title:      findTag(content, titleRegex),
		Difficulty: findTag(content, difficultyRegex),
		Tags:       tags,
		Draft:      draft,
		Fp:         filepath.Dir(fp),
		Link:       findTag(content, linkRegex),
	}
}

func genTable(data *TableData) string {
	var bf bytes.Buffer
	sort.Sort(data.Metas)
	tableTpl.Execute(&bf, data)
	return bf.String()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func Run() {
	files, err := doublestar.Glob("./solve/**/*")
	if err != nil {
		log.Fatal(err)
	}
	tagMetas := make(TagMetas, 0)
	tagMetas["all"] = make(Metas, 0)
	wg := sync.WaitGroup{}
	var lock sync.Mutex
	for _, fp := range files {
		if isFolder, _ := isDirectory(fp); isFolder {
			continue
		}

		if strings.HasSuffix(fp, ".md") {
			continue
		}
		wg.Add(1)
		fp := fp
		go func() {
			content, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			meta := findMeta(content, fp)
			if meta != nil {
				lock.Lock()
				addMeta(tagMetas, meta)
				lock.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for tag, metas := range tagMetas {
		fp := fmt.Sprintf("./toc/%s.md", tag)
		if !fileExists(fp) {
			continue
		}
		wg.Add(1)
		metas := metas
		go func() {
			content, err := ioutil.ReadFile(fp)
			if err != nil {
				log.Fatal(err)
			}
			table := genTable(&TableData{
				Metas: metas,
				Total: len(metas),
			})
			contents := strings.Split(string(content), "<!--- table -->")
			contents[1] = "<!--- table -->\n" + table
			newContent := strings.Join(contents, "")
			ioutil.WriteFile(fp, []byte(newContent), 0644)
			wg.Done()
		}()
	}
	wg.Wait()
}

var tableStr = `

总计: {{ .Total }}

| 序号 | 难度 | 题目                    | 解答                      |
| ---- | ---- | ------------------ | ---------------- |{{ range .Metas }}
| {{ .Index }} | {{ .Difficulty }} | [{{ .Title }}]({{ .Link }}) | [{{ .Fp }}](../{{ .Fp }})|{{ end }}
`
