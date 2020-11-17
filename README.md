# leetcode-tool ![goreleaser](https://github.com/zcong1993/leetcode-tool/workflows/goreleaser/badge.svg)

<!--
[![Go Report Card](https://goreportcard.com/badge/github.com/zcong1993/leetcode-tool)](https://goreportcard.com/report/github.com/zcong1993/leetcode-tool)
-->

> 一个让你更方便刷题的工具

## Install

```bash
$ brew tap zcong1993/homebrew-tap
$ brew install zcong1993/homebrew-tap/leetcode-tool

# show help
$ leetcode-tool help
```

## 主要功能

### 新建题目代码

`leetcode-tool new [<flags>] <number>`

number 为网页上显示的题目序号, 例如: `leetcode-tool new 1` 创建第一题的代码.

此命令会创建一个文件夹和三个文件, 分别是: 代码文件, 测试文件, 题目描述文件.

题目信息从 leetcode 拉取, 并生成一些元信息供更新命令分类更新题目完成状态.

### 查看题目信息

`leetcode-tool meta <number>`

展示某道题目一些信息.

### 更新题目状态

`leetcode-tool update`

根据已完成题目类别更新 `toc` 文件夹下的算法分类状态.

### 更新 leetcode 分类

`leetcode-tool tags`

从 leetcode 网站拉取最新分类, 并创建 `toc` 对应文件, 一般不需要运行.

## Workflow

如何刷题?

```bash
# 1. 新建一个题目
$ leetcode-tool new 1
# 2. 写出题解, 和测试代码
...
# 3. 更新 toc 文件
$ leetcode-tool update
# 4. 提交代码
```

## 为什么需要这个工具

1. leetcode 网页 IDE 太弱了, 我需要本地解题
1. 网页题解难以沉淀, 例如一道题我希望整理多种解法
1. GitHub 易于分享题解
1. 根据自己需要, 组织题目做专题之类的总结

## 模板项目

- Typescript [zcong1993/leetcode-ts-template](https://github.com/zcong1993/leetcode-ts-template)

## 使用此工具的项目

- [zcong1993/algo-go](https://github.com/zcong1993/algo-go)

## License

MIT &copy; zcong1993
