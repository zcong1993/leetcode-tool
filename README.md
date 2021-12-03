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

## 使用说明

[https://blog.cong.moe/post/2020-11-30-leetcode_tool](https://blog.cong.moe/post/2020-11-30-leetcode_tool)

## 支持语言

- Golang go
- Typescript ts
- Javascript js
- Python3 py3

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

### 重置项目

假如想要重新从零开始, 或者你的项目是基于别人项目 fork 的, 可以使用如下方式清理已有题解:

```shell
# 1. 删除所有题解文件
rm -rf solve/*
# 2. 重新构建 toc 文件, -f 参数会强制覆盖
leetcode-tool tags -f
# 2.1 假如你还保留了部分题解, 还需要更新下题目状态
leetcode-tool update
```

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

## 常见问题

### 1. 报错 panic: runtime error: invalid memory address or nil pointer dereference

因为 LeetCode 网站现在某些请求会校验 cookie, 采取的修复方法是请求增加了 cookie, 但是内置 cookie 没法确保长期有效.

所以 `.leetcode.json` 配置文件中支持 cookie 配置, 后续请访问此链接 https://leetcode-cn.com/api/problems/all 拿到 cookie 自行更新配置文件.

## License

MIT &copy; zcong1993
