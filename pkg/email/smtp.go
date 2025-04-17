package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
	UseTLS   bool
}

type SMTPClient struct {
	config SMTPConfig
}

type Message struct {
	FromEmail string
	FromName  string
	To        string
	Subject   string
	HTML      string
	Text      string
	Headers   map[string]string
}

func NewSMTPClient(config SMTPConfig) *SMTPClient {
	return &SMTPClient{
		config: config,
	}
}

func (c *SMTPClient) Send(message Message) (string, error) {
	if message.To == "" {
		return "", errors.New("recipient email is required")
	}
	if message.HTML == "" && message.Text == "" {
		return "", errors.New("either HTML or text content is required")
	}

	if message.FromEmail == "" {
		message.FromEmail = c.config.FromAddr
	}
	if message.FromName == "" {
		message.FromName = c.config.FromName
	}

	messageID := generateMessageID(message.To)
	emailContent, err := c.buildMIMEMessage(message, messageID)
	if err != nil {
		return "", err
	}

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	auth := smtp.PlainAuth("", c.config.Username, c.config.Password, c.config.Host)

	if c.config.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: c.config.Host,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return "", fmt.Errorf("TLS connection error: %w", err)
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, c.config.Host)
		if err != nil {
			return "", fmt.Errorf("SMTP client creation error: %w", err)
		}
		defer client.Close()

		if err = client.Auth(auth); err != nil {
			return "", fmt.Errorf("SMTP authentication error: %w", err)
		}

		if err = client.Mail(message.FromEmail); err != nil {
			return "", fmt.Errorf("SMTP FROM error: %w", err)
		}

		if err = client.Rcpt(message.To); err != nil {
			return "", fmt.Errorf("SMTP RCPT error: %w", err)
		}

		w, err := client.Data()
		if err != nil {
			return "", fmt.Errorf("SMTP DATA error: %w", err)
		}

		_, err = w.Write([]byte(emailContent))
		if err != nil {
			return "", fmt.Errorf("SMTP write error: %w", err)
		}

		err = w.Close()
		if err != nil {
			return "", fmt.Errorf("SMTP data close error: %w", err)
		}

		client.Quit()
	} else {
		err = smtp.SendMail(
			addr,
			auth,
			message.FromEmail,
			[]string{message.To},
			[]byte(emailContent),
		)
		if err != nil {
			return "", fmt.Errorf("SMTP send error: %w", err)
		}
	}

	return messageID, nil
}

func (c *SMTPClient) buildMIMEMessage(message Message, messageID string) (string, error) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\n", message.FromName, message.FromEmail))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", message.To))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", message.Subject))
	buf.WriteString(fmt.Sprintf("Message-ID: <%s>\r\n", messageID))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString("MIME-Version: 1.0\r\n")

	for name, value := range message.Headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", name, value))
	}

	boundary := "boundary_" + messageID[:8]
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", boundary))

	if message.Text != "" {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		buf.WriteString(message.Text)
		buf.WriteString("\r\n\r\n")
	}

	if message.HTML != "" {
		buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
		buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
		buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
		buf.WriteString(message.HTML)
		buf.WriteString("\r\n\r\n")
	}

	buf.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return buf.String(), nil
}

func generateMessageID(recipient string) string {
	domain := strings.Split(recipient, "@")[1]
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d.%d@%s", timestamp, time.Now().Unix(), domain)
}
