// Package p contains an HTTP Cloud Function.
package p

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"
)

var (
	from     string // from env
	pass     string // from env
	smtpaddr = "smtp.163.com"
	smtpport = "25"
	// smtpport        = "465"
	defaultreceiver = "w@clwen.com"
)

func init() {
	log.Println("starting...")
	Init()
}
func Init() {
	from = os.Getenv("MAIL_USER")
	if from == "" {
		log.Fatal("email addrs not set from env")
	}
	pass = os.Getenv("MAIL_PASS")
	if from == "" {
		log.Fatal("email pass not set from env")
	}

	if x := os.Getenv("MAIL_SMTP_ADDR"); x != "" {
		smtpaddr = x
	}
	if x := os.Getenv("MAIL_SMTP_PORT"); x != "" {
		smtpport = x
	}
	if x := os.Getenv("MAIL_DEFAULT_RECEIVER"); x != "" {
		defaultreceiver = x
	}
	log.Printf("set default receiver to %v\n", defaultreceiver)
	testsend()
}

var once sync.Once

func F(w http.ResponseWriter, r *http.Request) {
	once.Do(Init)

	d := Mail{
		Subject: r.FormValue("subject"),
		Body:    r.FormValue("body"),
		To:      r.FormValue("to"),
	}

	if d.Subject == "" && d.Body == "" {
		log.Printf("no form provided, try decode body")
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
			fmt.Fprintf(w, "decode body err: %v\n", err)
			return
		}
	}

	if d.Subject == "" && d.Body == "" {
		fmt.Fprintf(w, "subject and body can't both be empty\n")
		return
	}

	n := len(d.Body)
	if n > 100 {
		n = 100
	}
	digest := d.Body[:n] + " ..."
	if d.Subject == "" {
		d.Subject = digest
	}
	err := send(d.Subject, d.Body, d.To)
	if err != nil {
		fmt.Fprintf(w, "send err: %v\n", err)
		log.Printf("send err: %v\n", err)
		return
	}
	fmt.Fprintf(w, "send okay, to: %v, subject: %v, body: %v\n", d.To, d.Subject, digest)
	log.Printf("send okay, to: %v, subject: %v, body: %v\n", d.To, d.Subject, digest)
}

type Mail struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

func send(subject, body, to string) error {
	if to == "" {
		to = defaultreceiver
	}
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	return smtp.SendMail(smtpaddr+":"+smtpport,
		smtp.PlainAuth("", from, pass, smtpaddr),
		from, []string{to}, []byte(msg))
}

func testsend() {
	s := fmt.Sprintf("mail service started at %v\n", time.Now())
	b := fmt.Sprintf("mail service envs:\n %v\n", os.Environ())
	err := send(s, b, "")
	if err != nil {
		log.Fatalf("send err: %v\n", err)
	}
	log.Printf("sent start msg ok\n")
}
