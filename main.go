package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "embed"

	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	//go:embed smtp.json.example
	smtpTPL    string
	argFrom    *string
	argTo      *string
	argSMTP    *string
	argVerbose = false

	sentDir = "sent"
	send    Send

	cmd = &cobra.Command{
		Short: "Send eml files",
		Use:   "sendeml [-c smtp.json] [-f from] [-t address] *.eml",
		Run: func(cmd *cobra.Command, args []string) {
			if argVerbose {
				logrus.SetLevel(logrus.DebugLevel)
			}

			// create sent dir
			err := os.MkdirAll(sentDir, 0777)
			if err != nil {
				logrus.WithError(err).Fatalf("can not create sent directory")
				return
			}

			// parse send
			bytes, err := ioutil.ReadFile(*argSMTP)
			if len(bytes) > 0 {
				err = json.Unmarshal(bytes, &send)
			}
			if err != nil {
				if os.IsNotExist(err) {
					logrus.WithError(err).Errorf("smtp config not found")
				} else {
					logrus.WithError(err).Errorf("parse smtp config failed")
				}
				logrus.Infof("please create smtp.json from this template:")
				fmt.Println(smtpTPL)
				return
			}

			if len(args) > 0 {
				send.sendAndMove(args)
			} else {
				_ = cmd.Help()
			}
		},
	}
)

func init() {
	cmd.Flags().SortFlags = false
	cmd.PersistentFlags().SortFlags = false
	argSMTP = cmd.PersistentFlags().StringP(
		"smtp-config",
		"c",
		"smtp.json",
		"smtp config file",
	)
	argFrom = cmd.PersistentFlags().StringP(
		"from",
		"f",
		"",
		"email address from",
	)
	argTo = cmd.PersistentFlags().StringP(
		"to",
		"t",
		"",
		"email address to",
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
