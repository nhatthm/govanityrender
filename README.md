# Go Vanity HTML Renderer

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/govanityrender)](https://github.com/nhatthm/govanityrender/releases/latest)
[![Build Status](https://github.com/nhatthm/govanityrender/actions/workflows/docker-dev.yaml/badge.svg)](https://github.com/nhatthm/govanityrender/actions/workflows/docker-dev.yaml)
[![codecov](https://codecov.io/gh/nhatthm/govanityrender/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/govanityrender)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/vanityrender)
[![Donate](https://img.shields.io/badge/%20-Donate-%20?style=flat&logo=githubsponsors&color=E5E4E2)](http://donate.nhat.me)

Go Vanity HTML Renderer generates the html pages to use with GitHub Pages that allows you to set custom import paths for your Go packages.

## Prerequisites

- `Go >= 1.27`

## Install

```bash
go get go.nhat.io/vanityrender
```

## Usage

```shell
$ vanityrender --help
  -config string
    	config file (default "config.json")
  -homepage-tpl string
    	template file
  -modules string
    	rebuild only the listed modules, comma separated
  -out string
    	output path (default "build")
```

**Examples**

```text
$ vanityrender -config config.json -out build
```

## Donation

If this project saved you some development time, buy me a cup of coffee :)

[![donate](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](http://donate.nhat.me)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://github.com/nhatthm/donate.nhat.me/blob/master/images/qr_sponsor.png" width="147px" />

[<sub><sup>[table of contents]</sup></sub>](#table-of-contents)
