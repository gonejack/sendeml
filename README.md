# sendeml

Command line tool to send eml files

### Install
```shell
go get github.com/gonejack/sendeml
```

### Edit mail server config
Edit `smtp.json.example` into `smtp.json`

### Usage

```shell
sendeml [-c smtp.json] [-f from] [-t address] *.eml
```
