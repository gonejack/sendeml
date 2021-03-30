# sendeml

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gonejack/sendeml)
![Build](https://github.com/gonejack/sendeml/actions/workflows/go.yml/badge.svg)
[![GitHub license](https://img.shields.io/github/license/gonejack/sendeml.svg?color=red)](LICENSE)

Command line tool to send eml files.

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
> sendeml *.eml
```
```
Usage:
  sendeml [-c smtp.json] [-f from] [-t address] *.eml [flags]

Flags:
      --from string      email address from
      --to string        email address to
  -c, --smtp string      smtp config (default "/Users/youi/.sendeml/smtp.json")
  -p, --print-template   print smtp.json template
  -v, --verbose          verbose
  -h, --help             help for sendeml
```