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

	"github.com/tidwall/gjson"

	"github.com/zcong1993/leetcode-tool/pkg/leetcode"
)

type TplFile struct {
	Name     string
	FileName string
	TplStr   string
}

type LanguageConfig struct {
	LeetcodeLang string
	TplFiles     []TplFile
}

const (
	folder = "solve"
	prefix = "solve"
)

var (
	languageConfigs = map[string]LanguageConfig{
		"go": {
			LeetcodeLang: "Go",
			TplFiles:     []TplFile{{"code", "solve_%s.go", codeStrGo}, {"test", "solve_%s_test.go", testCodeStrGo}},
		},
		"ts": {
			LeetcodeLang: "TypeScript",
			TplFiles:     []TplFile{{"code", "solve_%s.ts", codeStrTs}, {"test", "solve_%s.test.ts", testCodeStrTs}},
		},
		"js": {
			LeetcodeLang: "JavaScript",
			TplFiles:     []TplFile{{"code", "solve_%s.js", codeStrJs}, {"test", "solve_%s.test.js", testCodeStrJs}},
		},
		"py3": {
			LeetcodeLang: "Python3",
			TplFiles:     []TplFile{{"code", "solve_%s.py", codeStrPy3}, {"test", "test_%s.py", testCodeStrPy3}, {"__init__", "__init__.py", ""}},
		},
		"java": {
			LeetcodeLang: "Java",
			TplFiles:     []TplFile{{"code", "solve_%s.java", codeStrJava}, {"test", "test_%s.java", testCodeStrJava}},
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

func Run(lc *leetcode.Leetcode, n string, lang string) {
	if lang == "" {
		lang = lc.Config.Lang
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
	meta, err := lc.GetMetaByNumber(n)
	if err != nil || meta == nil {
		log.Fatal(err, meta)
	}
	number := normalizeNumber(meta.Index)
	folderName := prefix + number
	fp := filepath.Join(folder, folderName)
	os.MkdirAll(fp, 0755)
	metaf := &MetaWithFolder{
		*meta,
		folderName,
		strings.Join(meta.Tags, ","),
		n,
	}
	metaf.Meta.Content = strings.ReplaceAll(metaf.Meta.Content, "↵", "")
	metaf.Meta.Code = gjson.Get(metaf.CodeSnippets, fmt.Sprintf("#(lang=%s).code", config.LeetcodeLang)).String()

	problemFp := filepath.Join(fp, "problem.md")
	if !fileExists(problemFp) {
		bf := mustExecuteTemplate("problem", problemStr, metaf)
		ioutil.WriteFile(problemFp, bf, 0644)
	}

	for _, tpl := range config.TplFiles {
		fileName := tpl.FileName
		if strings.Count(tpl.FileName, "%s") > 0 {
			fileName = fmt.Sprintf(tpl.FileName, number)
		}
		fp := filepath.Join(fp, fileName)
		if !fileExists(fp) {
			bf := mustExecuteTemplate(tpl.Name, tpl.TplStr, metaf)
			ioutil.WriteFile(fp, bf, 0644)
		}
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

var (
	codeStrPy3 = `'''
@index {{ .Index }}
@title {{ .Title }}
@difficulty {{ .Difficulty }}
@tags {{ .TagStr }}
@draft false
@link {{ .Link }}
@frontendId {{ .FrontendId }}
'''

{{ .Code }}
`
	testCodeStrPy3 = `def test_solve():
	pass
`
)

var (
	codeStrJava = `package {{ .Folder }};

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

	testCodeStrJava = `package {{ .Folder }};

public class test_{{ printf "%04s" .Index }} {
	public static void main(String[] args) {
		Solution solution = new Solution();
		// do some test
	}
}
`
)