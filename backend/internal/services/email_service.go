package services

import (
	"fmt"
	"net/smtp"

	"backend/internal/models"
)

type EmailJob struct {
	User         *models.User
	SmtpHost     string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
	FrontendURL  string
}

func SendVerificationEmail(job *EmailJob) error {
	verificationURL := fmt.Sprintf("%s/verify_user?username=%s", job.FrontendURL, job.User.Username)

	to := job.User.Email
	subject := "Email Verification Code on Relaxio"

	body := fmt.Sprintf(`
    <html>
    <body>
        <p>Your verification code is: <strong>%s</strong></p>
        <p>Click the following link to verify your account:</p>
        <p><a href="%s">Verify Account</a></p>
    </body>
    </html>
    `, job.User.VerificationCode, verificationURL)

	msg := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", job.SmtpUser, job.SmtpPassword, job.SmtpHost)
	err := smtp.SendMail(job.SmtpHost+":"+job.SmtpPort, auth, job.SmtpUser, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}
