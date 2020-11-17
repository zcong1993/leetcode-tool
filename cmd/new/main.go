package new

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	config2 "github.com/zcong1993/leetcode-tool/internal/config"

	"github.com/tidwall/gjson"

	"github.com/zcong1993/leetcode-tool/pkg/leetcode"
)

type LanguageConfig struct {
	LeetcodeLang   string
	CodeTplStr     string
	TestCodeTplStr string
	CodeFileName   string
	TestFileName   string
}

const (
	folder = "solve"
	prefix = "solve"
)

var (
	languageConfigs = map[string]LanguageConfig{
		"go": {
			CodeTplStr:     codeStrGo,
			TestCodeTplStr: testCodeStrGo,
			LeetcodeLang:   "Go",
			CodeFileName:   "solve_%s.go",
			TestFileName:   "solve_%s_test.go",
		},
		"ts": {
			CodeTplStr:     codeStrTs,
			TestCodeTplStr: testCodeStrTs,
			LeetcodeLang:   "TypeScript",
			CodeFileName:   "solve_%s.ts",
			TestFileName:   "solve_%s.test.ts",
		},
		"js": {
			CodeTplStr:     codeStrJs,
			TestCodeTplStr: testCodeStrJs,
			LeetcodeLang:   "JavaScript",
			CodeFileName:   "solve_%s.js",
			TestFileName:   "solve_%s.test.js",
		},
	}
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func normalizeNumber(number string) string {
	if len(number) >= 4 {
		return number
	}
	return strings.Repeat("0", 4-len(number)) + number
}

func mustExecuteTemplate(name string, str string, data interface{}) []byte {
	tpl := template.Must(template.New(name).Parse(str))
	var bf bytes.Buffer
	err := tpl.Execute(&bf, data)
	if err != nil {
		log.Fatalf("mustExecuteTemplate %s error: %s\n", name, err.Error())
	}
	return bf.Bytes()
}

type MetaWithFolder struct {
	leetcode.Meta
	Folder     string
	TagStr     string
	FrontendId string
}

func Run(n string, lang string) {
	if lang == "" {
		lang = config2.GetLang()
	}
	config, ok := languageConfigs[lang]
	if !ok {
		supportLangs := make([]string, len(languageConfigs))
		i := 0
		for l := range languageConfigs {
			supportLangs[i] = l
			i++
		}
		log.Fatalf("invalid lang, now support %s\n", strings.Join(supportLangs, ","))
	}
	meta, err := leetcode.GetMetaByNumber(n)
	if err != nil || meta == nil {
		log.Fatal(err, meta)
	}
	number := normalizeNumber(meta.Index)
	folderName := prefix + number
	fp := filepath.Join(folder, folderName)
	os.MkdirAll(fp, 0755)
	codeFp := filepath.Join(fp, fmt.Sprintf(config.CodeFileName, number))
	codeTestFp := filepath.Join(fp, fmt.Sprintf(config.TestFileName, number))
	problemFp := filepath.Join(fp, "problem.md")
	metaf := &MetaWithFolder{
		*meta,
		folderName,
		strings.Join(meta.Tags, ","),
		n,
	}
	metaf.Meta.Content = strings.ReplaceAll(metaf.Meta.Content, "↵", "")
	metaf.Meta.Code = gjson.Get(metaf.CodeSnippets, fmt.Sprintf("#(lang=%s).code", config.LeetcodeLang)).String()

	if !fileExists(codeFp) {
		bf := mustExecuteTemplate("code", config.CodeTplStr, metaf)
		ioutil.WriteFile(codeFp, bf, 0644)
	}

	if !fileExists(codeTestFp) {
		bf := mustExecuteTemplate("test", config.TestCodeTplStr, metaf)
		ioutil.WriteFile(codeTestFp, bf, 0644)
	}

	if !fileExists(problemFp) {
		bf := mustExecuteTemplate("problem", problemStr, metaf)
		ioutil.WriteFile(problemFp, bf, 0644)
	}
	fmt.Printf("Done: %s\n", fp)
}

var (
	codeStrGo = `package {{ .Folder }}

/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
*/

{{ .Code }}
`

	testCodeStrGo = `package {{ .Folder }}_test

`

	problemStr = `# [{{ .Index }}. {{ .Title }}]({{ .Link }})

{{ .Content }}
`
)

var (
	codeStrTs = `/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
*/

export {{ .Code }}
`
	testCodeStrTs = `
it('solve_{{ .Index }} should pass', () => {})
`
)

var (
	codeStrJs = `/**
 * @index {{ .Index }}
 * @title {{ .Title }}
 * @difficulty {{ .Difficulty }}
 * @tags {{ .TagStr }}
 * @draft false
 * @link {{ .Link }}
 * @frontendId {{ .FrontendId }}
*/

{{ .Code }}
`
	testCodeStrJs = `
it('solve_{{ .Index }} should pass', () => {})
`
)
