package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	//go:embed smtp.json.example
	smtpTPL string
	sentDir = "sent"

	argFrom     *string
	argTo       *string
	argSMTP     *string
	argVerbose  = false
	argTemplate = false

	send sender
	cmd  = &cobra.Command{
		Short: "Command line tool to send eml files",
		Use:   "sendeml [-c smtp.json] [-f from] [-t address] *.eml",
		Run: func(cmd *cobra.Command, args []string) {
			if argTemplate {
				fmt.Print(smtpTPL)
				return
			}

			if argVerbose {
				logrus.SetLevel(logrus.DebugLevel)
			}

			// parse smtp.json
			cpath := *argSMTP
			fd, err := os.Open(cpath)
			if err == nil {
				err = json.NewDecoder(fd).Decode(&send.config)
			}
			if errors.Is(err, fs.ErrNotExist) {
				logrus.Errorf("%s not found", cpath)
				_ = ioutil.WriteFile(cpath, []byte(smtpTPL), 0766)
				logrus.Infof("%s created", cpath)
				return
			}
			if err != nil {
				logrus.WithError(err).Errorf("parse %s failed", cpath)
				return
			}

			if len(args) == 0 {
				args, _ = filepath.Glob("*.eml")
			}
			if len(args) == 0 {
				logrus.Fatalf("no .eml files given")
				return
			}

			// create sent dir
			err = os.MkdirAll(sentDir, 0766)
			if err != nil {
				logrus.WithError(err).Fatalf("can not create sent directory")
				return
			}

			send.sendAndMove(args)
		},
	}
)

func defaultConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".sendeml")
}
func init() {
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false
	argFrom = cmd.PersistentFlags().StringP(
		"from",
		"",
		"",
		"email address from",
	)
	argTo = cmd.PersistentFlags().StringP(
		"to",
		"",
		"",
		"email address to",
	)
	argSMTP = cmd.PersistentFlags().StringP(
		"smtp",
		"c",
		filepath.Join(defaultConfigDir(), "smtp.json"),
		"smtp config",
	)
	cmd.PersistentFlags().BoolVarP(
		&argTemplate,
		"print-template",
		"p",
		false,
		"print smtp.json template",
	)
	cmd.PersistentFlags().BoolVarP(
		&argVerbose,
		"verbose",
		"v",
		false,
		"verbose",
	)
	logrus.SetFormatter(&formatter.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		//NoColors:        true,
		HideKeys:    true,
		CallerFirst: true,
		FieldsOrder: []string{"feed", "article", "source"},
	})
}

func main() {
	_ = cmd.Execute()
}
