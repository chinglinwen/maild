package p

import (
	"crypto/tls"
	"log"
	"net/smtp"
)

func sendtls(subject, body, to string) (err error) {
	if to == "" {
		to = defaultreceiver
	}
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	servername := smtpaddr + ":" + smtpport
	log.Printf("send to %v\n", servername)
	auth := smtp.PlainAuth("", from, pass, smtpaddr)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpaddr,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return
	}

	c, err := smtp.NewClient(conn, smtpaddr)
	if err != nil {
		return
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return
	}

	// To && From
	if err = c.Mail(from); err != nil {
		return
	}

	if err = c.Rcpt(to); err != nil {
		return
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return
	}
	log.Printf("log writed msg\n")

	err = w.Close()
	if err != nil {
		return
	}

	return c.Quit()
}
