# sendeml
Command line tool to send eml files.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/sendeml)
![Build](https://github.com/gonejack/sendeml/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/sendeml.svg?color=blue)](LICENSE)

### Install
```shell
> go get github.com/gonejack/sendeml
```

### Config
Edit `~/.sendeml/smtp.json`
```shell
# new config
> sendeml -p > ~/.sendeml/smtp.json

# edit
> vi ~/.sendeml/smtp.json
```

### Usage

```shell
> sendeml [-c ~/.sendeml/smtp.json] [-f from] [-t address] *.eml
```

### requirement
- Go 1.16
