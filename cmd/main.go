package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zcong1993/leetcode-tool/cmd/new"
	"github.com/zcong1993/leetcode-tool/cmd/tags"
	"github.com/zcong1993/leetcode-tool/cmd/update"
	"github.com/zcong1993/leetcode-tool/pkg/leetcode"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "master"
	commit  = ""
	date    = ""
	builtBy = ""
)

var (
	app = kingpin.New("algo", "A command-line tool for algo-go repo.")

	updateCmd = app.Command("update", "Update readme.")

	newCmd = app.Command("new", "Init a new leetcode problem.")
	number = newCmd.Arg("number", "problem number").Required().String()
	lang   = newCmd.Flag("lang", "language").String()

	metaCmd    = app.Command("meta", "Show problem meta by number.")
	metaNumber = metaCmd.Arg("number", "problem number").Required().String()

	tagsCmd   = app.Command("tags", "Update tag toc files.")
	tagsForce = tagsCmd.Flag("force", "force update file").Short('f').Bool()
)

func showMeta(number string) {
	meta, err := leetcode.GetMetaByNumber(number)
	if err != nil {
		log.Fatal(err)
	}
	if meta == nil {
		log.Fatal("mate not found")
	}
	meta.Content = ""
	meta.Code = ""
	meta.CodeSnippets = ""
	fmt.Printf("%+v\n", meta)
}

func main() {
	app.Version(buildVersion(version, commit, date, builtBy))
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case updateCmd.FullCommand():
		update.Run()
	case newCmd.FullCommand():
		new.Run(*number, *lang)
	case metaCmd.FullCommand():
		showMeta(*metaNumber)
	case tagsCmd.FullCommand():
		tags.Run(*tagsForce)
	}
}

func buildVersion(version, commit, date, builtBy string) string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
