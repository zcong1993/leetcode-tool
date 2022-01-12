package leetcode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/zcong1993/leetcode-tool/internal/config"

	"github.com/tidwall/gjson"
)

type Meta struct {
	Index        string
	Title        string
	Difficulty   string
	Tags         []string
	Link         string
	Content      string
	Code         string
	CodeSnippets string
}

type Tag struct {
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	TranslatedName string `json:"translatedName"`
}

var (
	difficultyMap = map[string]string{
		"Easy":   "简单",
		"Medium": "中等",
		"Hard":   "困难",
	}
)

type Leetcode struct {
	Config *config.Config
}

func NewLeetcode(config *config.Config) *Leetcode {
	return &Leetcode{Config: config}
}

func (l *Leetcode) getAllProblem() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, "https://leetcode-cn.com/api/problems/all/", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", l.Config.Cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (l *Leetcode) findProblemSlugByNumber(problems []byte, number string) string {
	return gjson.GetBytes(problems, fmt.Sprintf("stat_status_pairs.#(stat.frontend_question_id=\"%s\").stat.question__title_slug", number)).String()
}

func (l *Leetcode) getDetail(slug string) (*Meta, error) {
	if slug == "" {
		return nil, nil
	}
	req, err := http.NewRequest("POST", "https://leetcode-cn.com/graphql/", strings.NewReader(fmt.Sprintf(`{"operationName":"questionData","variables":{"titleSlug": "%s"},"query":"query questionData($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    questionId\n    questionFrontendId\n    boundTopicId\n    title\n    titleSlug\n    content\n    translatedTitle\n    translatedContent\n    isPaidOnly\n    difficulty\n    likes\n    dislikes\n    isLiked\n    similarQuestions\n    contributors {\n      username\n      profileUrl\n      avatarUrl\n      __typename\n    }\n    langToValidPlayground\n    topicTags {\n      name\n      slug\n      translatedName\n      __typename\n    }\n    companyTagStats\n    codeSnippets {\n      lang\n      langSlug\n      code\n      __typename\n    }\n    stats\n    hints\n    solution {\n      id\n      canSeeDetail\n      __typename\n    }\n    status\n    sampleTestCase\n    metaData\n    judgerAvailable\n    judgeType\n    mysqlSchemas\n    enableRunCode\n    envInfo\n    book {\n      id\n      bookName\n      pressName\n      source\n      shortDescription\n      fullDescription\n      bookImgUrl\n      pressImgUrl\n      productUrl\n      __typename\n    }\n    isSubscribed\n    isDailyQuestion\n    dailyRecordStatus\n    editorType\n    ugcQuestionId\n    style\n    __typename\n  }\n}\n"}`, slug)))
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("User-Agent", "axios/0.19.2")
	req.Header.Add("Host", "leetcode-cn.com")
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tagsResult := gjson.GetBytes(content, "data.question.topicTags.#.slug").Array()
	tags := make([]string, len(tagsResult))
	for i, t := range tagsResult {
		tags[i] = t.String()
	}

	codeSnippets := gjson.GetBytes(content, "data.question.codeSnippets").String()

	return &Meta{
		Index:        gjson.GetBytes(content, "data.question.questionId").String(),
		Title:        gjson.GetBytes(content, "data.question.translatedTitle").String(),
		Difficulty:   difficultyMap[gjson.GetBytes(content, "data.question.difficulty").String()],
		Tags:         tags,
		Link:         fmt.Sprintf("https://leetcode-cn.com/problems/%s/", gjson.GetBytes(content, "data.question.titleSlug").String()),
		Content:      gjson.GetBytes(content, "data.question.translatedContent").String(),
		Code:         gjson.GetBytes(content, "data.question.codeSnippets.#(lang=Go).code").String(),
		CodeSnippets: codeSnippets,
	}, nil
}

func (l *Leetcode) GetMetaByNumber(number string) (*Meta, error) {
	problems, err := l.getAllProblem()
	if err != nil {
		return nil, err
	}
	slug := l.findProblemSlugByNumber(problems, number)
	return l.getDetail(slug)
}

func (l *Leetcode) GetTags() ([]Tag, error) {
	resp, err := http.Get("https://leetcode-cn.com/problems/api/tags/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := make([]Tag, 0)
	err = json.Unmarshal([]byte(gjson.GetBytes(bt, "topics").Raw), &res)
	return res, err
}
