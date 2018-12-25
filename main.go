package main

import (
	"github.com/namsral/flag"
	"fmt"
	"net/http"
	"strings"
	"crypto/tls"

	//"github.com/natefinch/lumberjack"
	gomail "gopkg.in/gomail.v2"
	"github.com/chinglinwen/log"
)

var (
	port        = flag.String("p", "3001", "listening port")
	smtpAddr    = flag.String("smtpaddr", "", "smtp address")
	smtpPort    = flag.Int("port", 25, "smtp port")
	smtpUser    = flag.String("user", "", "smtp user")
	smtpPass    = flag.String("pass", "", "smtp pass")
	defaultfrom = flag.String("from", "", "from (default=user)")
)

func main() {
	http.HandleFunc("/", mailHandler)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}

func mailHandler(w http.ResponseWriter, r *http.Request) {
	receivers := r.FormValue("receiver")
	if receivers == "" {
		fmt.Fprintf(w, "receiver is empty\n")
		return
	}
	subject := r.FormValue("subject")
	if subject == "" {
		subject = "empty"
	}
	body := r.FormValue("body")
	if body == "" {
		body = "empty"
	}
	from := r.FormValue("from")
	if from == "" {
		from = *defaultfrom
	} else {
		from += "<" + *defaultfrom + ">"
	}

	rs := strings.Split(receivers, ",")
	err := mail(from, subject, body, rs)
	if err != nil {
		fmt.Fprintf(w, "send error: %v\n", err)
		log.Printf("send error: %v\n", err)
		return
	}
	fmt.Fprintf(w, "sent okay")
	log.Printf("sent to %v, %v, %v", rs, subject, body)
}

func mail(from, subject, body string, receivers []string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", receivers...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	d := gomail.NewDialer(*smtpAddr, *smtpPort, *smtpUser, *smtpPass)
        d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return d.DialAndSend(m)
}

func init() {
	flag.Parse()
	if *defaultfrom == "" {
		*defaultfrom = *smtpUser
	}
	/*
	log.SetOutput(&lumberjack.Logger{
		Filename:   "maild.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})
	*/
	log.Println("starting...")
}
