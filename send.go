package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/fs"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/gonejack/email"
	"github.com/sirupsen/logrus"
)

type config struct {
	From     string
	To       string
	Host     string
	Port     int
	Username string
	Password string
}

type sender struct {
	config
	auth   smtp.Auth
	client *smtp.Client
}

func (s *sender) sendAndMove(emails []string) {
	client, err := s.getClient()
	if err != nil {
		logrus.WithError(err).Fatalf("get client failed")
		return
	}
	s.client = client
	defer func() {
		_ = s.client.Quit()
		_ = s.client.Close()
	}()

	for _, eml := range emails {
		log := logrus.WithField("email", eml)

		log.Infof("sending")
		err := s.sendEmail(eml)
		if err != nil {
			log.WithError(err).Errorf("send failed")
			continue
		}
		log.Debugf("sent")

		rename := filepath.Join(sentDir, filepath.Base(eml))
		index := 1
		for {
			if _, err := os.Open(rename); errors.Is(err, fs.ErrNotExist) {
				break
			}

			rename = fmt.Sprintf("%s#%d", rename, index)
			index += 1
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

	return e.SendWithClient(s.client)
}
func (s *sender) getClient() (c *smtp.Client, err error) {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	auth := s.getAuth()

	c, err = smtp.Dial(addr)
	if err != nil {
		return
	}

	err = c.Hello("localhost")
	if err != nil {
		return
	}

	ok, _ := c.Extension("STARTTLS")
	if ok {
		err = c.StartTLS(&tls.Config{ServerName: s.Host})
		if err != nil {
			return
		}
	}

	if auth != nil {
		ok, _ := c.Extension("AUTH")
		if ok {
			err = c.Auth(auth)
			if err != nil {
				return
			}
		}
	}

	return
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
