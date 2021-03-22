package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

type sender struct {
	From     string
	To       string
	Host     string
	Port     int
	Username string
	Password string

	auth smtp.Auth
}

func (s *sender) sendAndMove(emails []string) {
	for _, eml := range emails {
		log := logrus.WithField("email", eml)

		log.Debugf("sending")
		err := s.sendEmail(eml)
		if err != nil {
			log.WithError(err).Errorf("send failed")
			continue
		}
		log.Info("sent")

		rename := filepath.Join(sentDir, filepath.Base(eml))
		index := 1
		for {
			_, err := os.Stat(rename)
			if errors.Is(err, os.ErrNotExist) {
				break
			} else {
				rename = fmt.Sprintf("%s#%d", rename, index)
				index += 1
			}
		}
		err = os.Rename(eml, rename)
		if err != nil {
			log.WithError(err).Errorf("move failed")
		}
	}
}

func (s *sender) sendEmail(eml string) (err error) {
	file, err := os.Open(eml)
	if err != nil {
		return
	}
	defer file.Close()

	e, err := email.NewEmailFromReader(file)
	if err != nil {
		return
	}

	e.From = s.getFrom()
	e.To = []string{s.getTo()}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := s.getAuth()

	if s.Port == 465 {
		return e.SendWithTLS(addr, auth, &tls.Config{ServerName: s.Host})
	} else {
		return e.Send(addr, auth)
	}
}

func (s *sender) getFrom() string {
	if *argFrom != "" {
		return *argFrom
	}
	return s.From
}
func (s *sender) getTo() string {
	if *argTo != "" {
		return *argTo
	}
	return s.To
}
func (s *sender) getAuth() smtp.Auth {
	if s.auth == nil {
		s.auth = smtp.PlainAuth("", send.Username, send.Password, send.Host)
	}
	return s.auth
}
