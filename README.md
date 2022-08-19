# Go Vanity HTML Renderer

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/govanityrender)](https://github.com/nhatthm/govanityrender/releases/latest)
[![Build Status](https://github.com/nhatthm/govanityrender/actions/workflows/docker-dev.yaml/badge.svg)](https://github.com/nhatthm/govanityrender/actions/workflows/docker-dev.yaml)
[![codecov](https://codecov.io/gh/nhatthm/govanityrender/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/govanityrender)
[![Go Report Card](https://goreportcard.com/badge/go.nhat.io/vanityrender)](https://goreportcard.com/report/go.nhat.io/vanityrender)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/go.nhat.io/vanityrender)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

Go Vanity HTML Renderer generates the html pages to use with GitHub Pages that allows you to set custom import paths for your Go packages.

## Prerequisites

- `Go >= 1.18`

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
  -out string
    	output path (default "build")
```

**Examples**

```text
$ vanityrender -config config.json -out build
```

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
