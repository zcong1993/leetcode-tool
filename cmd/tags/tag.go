package tags

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/zcong1993/leetcode-tool/pkg/leetcode"
)

var tagTpl = template.Must(template.New("tag").Parse(tagStr))

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func Run(force bool) {
	tags, err := leetcode.GetTags()
	if err != nil {
		log.Fatal(err)
	}

	tags = append([]leetcode.Tag{{
		Name:           "all",
		Slug:           "all",
		TranslatedName: "汇总",
	}}, tags...)

	wg := sync.WaitGroup{}
	sb := strings.Builder{}
	for _, tag := range tags {
		if tag.TranslatedName == "" {
			tag.TranslatedName = tag.Name
		}
		fp := filepath.Join("./toc", tag.Slug+".md")

		sb.WriteString(fmt.Sprintf("- [%s](%s)\n", tag.TranslatedName, fp))

		tag := tag
		wg.Add(1)
		go func() {
			if !force && fileExists(fp) {
				wg.Done()
				return
			}
			var content bytes.Buffer
			err := tagTpl.Execute(&content, tag)
			if err != nil {
				log.Fatal(err)
			}
			err = ioutil.WriteFile(fp, content.Bytes(), 0644)
			if err != nil {
				log.Printf("write file %s error, %s\n", fp, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println(sb.String())
}

var tagStr = `# {{ .TranslatedName }}

<!--- table -->

`
