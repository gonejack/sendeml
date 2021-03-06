# sendeml

Command line tool to send eml files

### Install
```shell
go get github.com/gonejack/sendeml
```

### Config
Create `smtp.json` by
```shell
sendeml -p
```

### Usage

```shell
sendeml [-c smtp.json] [-f from] [-t address] *.eml
```

### requirement
- Go 1.16
