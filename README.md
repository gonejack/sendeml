# sendeml

Command line tool to send eml files

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
