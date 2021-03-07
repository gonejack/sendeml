package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "embed"

	"github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	argConfDir  *string
	argFrom     *string
	argTo       *string
	argVerbose  = false
	argTemplate = false

	//go:embed smtp.json.example
	smtpTPL  string
	smtpConf = "smtp.json"
	sentDir  = "sent"

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

			logrus.Infof("config dir is %s", *argConfDir)
			{
				err := os.MkdirAll(*argConfDir, 0766)
				if err != nil {
					logrus.WithError(err).Fatalf("can not create config directory")
					return
				}
				smtpConf = filepath.Join(*argConfDir, smtpConf)
			}

			// create sent dir
			err := os.MkdirAll(sentDir, 0766)
			if err != nil {
				logrus.WithError(err).Fatalf("can not create sent directory")
				return
			}

			// parse smtp.json
			bytes, err := ioutil.ReadFile(smtpConf)
			if len(bytes) > 0 {
				if string(bytes) == smtpTPL {
					logrus.Infof("please edit %s", smtpConf)
					return
				}
				err = json.Unmarshal(bytes, &send)
			}
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					logrus.Errorf("smtp.json %s not found", smtpConf)
					_ = ioutil.WriteFile(smtpConf, []byte(smtpTPL), 0766)
					logrus.Infof("smtp.json %s created", smtpConf)
				} else {
					logrus.WithError(err).Errorf("parse smtp.json failed")
				}
				logrus.Infof("please edit %s", smtpConf)
				return
			}

			if len(args) == 0 {
				logrus.Fatalf("no .eml files given")
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
	argConfDir = cmd.PersistentFlags().StringP(
		"config-dir",
		"c",
		defaultConfigDir(),
		"config directory",
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
