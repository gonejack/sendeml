package main

import (
	"encoding/json"
	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var (
	flagVerbose = false

	argFrom *string
	argTo   *string

	sendFile = "send.json"
	sentDir  = "sent"
	send     Send

	cmd = &cobra.Command{
		Short: "Send eml files",
		Use:   "sendeml [-f from] [-t address] *.eml",
		Run: func(cmd *cobra.Command, args []string) {
			if flagVerbose {
				logrus.SetLevel(logrus.DebugLevel)
			}

			// create sent dir
			err := os.MkdirAll(sentDir, 0777)
			if err != nil {
				logrus.WithError(err).Fatalf("can not create sent directory")
				return
			}

			// parse send
			bytes, err := ioutil.ReadFile(sendFile)
			if err == nil && len(bytes) > 0 {
				err = json.Unmarshal(bytes, &send)
				if err != nil {
					logrus.WithError(err).Fatalf("parse %s failed", sendFile)
				}
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
	argFrom = cmd.PersistentFlags().StringP(
		"from",
		"f",
		"",
		"email from",
	)
	argTo = cmd.PersistentFlags().StringP(
		"to",
		"t",
		"",
		"email to",
	)
	cmd.PersistentFlags().BoolVarP(
		&flagVerbose,
		"verbose",
		"v",
		false,
		"Verbose",
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
