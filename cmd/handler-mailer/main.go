// Take well-formed json from a sensu check result a context rich html document to be mail to
// one or more addresses.
//
// LICENSE:
//   Copyright 2016 Yieldbot. <devops@yieldbot.com>
//   Released under the MIT License; see LICENSE
//   for details.

package main

import (
	"bytes"
	"fmt"
	"github.com/codegangsta/cli"
	// "github.com/yieldbot/sensumailer/lib"
	"github.com/yieldbot/sensuplugin/sensuhandler"
	"github.com/yieldbot/sensuplugin/sensuutil"
	// "log"
	"net/smtp"
	"os"
	// "time"
)

func main() {

	var emailAddress string
	var smtpHost string
	var smtpPort string
	var emailSender string
	var debug bool

	app := cli.NewApp()
	app.Name = "handler-mailer"
	app.Usage = "Send context rich html alert notifications via email"
	app.Action = func(c *cli.Context) {

		if debug {
			fmt.Printf("This is the sending address: %v \n", emailSender)
			fmt.Printf("This is the recieving address: %v\n", emailAddress)
			fmt.Printf("This is the smtp address: %v:%v\n", smtpHost, smtpPort)
			sensuutil.Exit("debug")
		}

		// Get the sensu event data
		sensuEvent := new(sensuhandler.SensuEvent)
		sensuEvent = sensuEvent.AcquireSensuEvent()

		// Connect to the remote SMTP server.
		s, err := smtp.Dial(smtpHost + ":" + smtpPort)
		if err != nil {
			sensuutil.EHndlr(err)
		}
		defer s.Close()

		// Set the sender and recipient.
		s.Mail(emailSender)
		s.Rcpt(emailAddress)

		// Send the email body.
		ws, err := s.Data()
		if err != nil {
			sensuutil.EHndlr(err)
		}
		defer ws.Close()
		buf := bytes.NewBufferString("This is the email body.")
		if _, err = buf.WriteTo(ws); err != nil {
			sensuutil.EHndlr(err)
		}

		fmt.Printf("Email sent to %s\n", emailAddress)
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "address",
			Value:       "mattjones@yieldbot.com",
			Usage:       "email address to send to",
			EnvVar:      "SENSU_HANDLER_EMAIL_ADDRESS",
			Destination: &emailAddress,
		},
		cli.StringFlag{
			Name:        "host",
			Value:       "localhost",
			Usage:       "smtp server",
			EnvVar:      "SENSU_HANDLER_EMAIL_HOST",
			Destination: &smtpHost,
		},
		cli.StringFlag{
			Name:        "port",
			Value:       "25",
			Usage:       "smtp port",
			EnvVar:      "SENSU_HANDLER_EMAIL_PORT",
			Destination: &smtpPort,
		},
		cli.StringFlag{
			Name:        "sender",
			Value:       "sensu@yieldbot.com",
			Usage:       "email sender",
			EnvVar:      "SENSU_HANDLER_EMAIL_SENDER",
			Destination: &emailSender,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Print debugging info, no alerts will be sent",
			Destination: &debug,
		},
	}
	app.Run(os.Args)
}
