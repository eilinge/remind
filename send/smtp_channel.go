package send

import (
	"bytes"

	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type Smtp struct {
	Host string
	Port int

	From     string
	Username string
	Password string
}

/*	Send get new SMTP Dialer, build message and send message(d.DialAndSend(m))
	m.SetHeader("From", s.Username)
	m.SetHeader("From_Name", s.From)
	m.SetHeader("To", rs...)
	m.SetHeader("Subject", n.Title)
	m.SetBody("text/plain", n.Content)
*/
func (s Smtp) Send(n *Notify, tos []*User) error {
	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	m := gomail.NewMessage()
	m.SetHeader("From", s.Username)
	m.SetHeader("From_Name", s.From)

	rs := make([]string, 0, len(tos))
	for _, u := range tos {
		if u.Email == "" {
			continue
		}
		rs = append(rs, u.Email)
	}

	l := logrus.WithField("notify", n).WithField("tos", tos)
	if len(rs) <= 0 {
		l.Info("no email address")
		return nil
	}
	m.SetHeader("To", rs...)
	m.SetHeader("Subject", n.Title)

	buffer := bytes.NewBuffer(nil)
	body := "Content: " + string([]rune(n.Content.Message.Cont)) + "\r\n"
	body += "Level: " + n.Content.Message.Level
	buffer.WriteString(body)

	m.SetBody("text/plain", buffer.String())

	return d.DialAndSend(m)
}
