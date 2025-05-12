package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	config EmailConfig
}

func NewEmailService(config EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// SendPasswordResetEmail отправляет электронное письмо с ссылкой для сброса пароля
func (s *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	from := s.config.Username
	password := s.config.Password

	// Заголовки письма
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = "Восстановление пароля в ZLib"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Тело письма
	body := fmt.Sprintf(
		`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Восстановление пароля</title>
</head>
<body>
    <h2>Восстановление пароля в ZLib</h2>
    <p>Для восстановления пароля перейдите по ссылке ниже:</p>
    <p><a href="http://localhost:3000/auth?token=%s">Восстановить пароль</a></p>
    <p>Если вы не запрашивали восстановление пароля, проигнорируйте это письмо.</p>
    <p>С уважением, команда ZLib</p>
</body>
</html>
`, resetToken,
	)

	// Формирование полного сообщения
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Настройка TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         s.config.Host,
	}

	// Подключение к серверу
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Аутентификация
	auth := smtp.PlainAuth("", from, password, s.config.Host)
	if err = client.Auth(auth); err != nil {
		return err
	}

	// Отправитель и получатель
	if err = client.Mail(from); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}

	// Отправка сообщения
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
