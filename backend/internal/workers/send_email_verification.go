package workers

import (
	"backend/internal/models"
	"fmt"
	"log"
	"net/smtp"
	"time"
)

func SendVerificationEmail(user *models.User, smtpHost, smtpPort, smtpUser, smtpPassword, frontendURL string) error {
	verificationURL := fmt.Sprintf("%s/verify_user?username=%s", frontendURL, user.Username)

	to := user.Email
	subject := "Email Verification Code on Relaxio"
	body := fmt.Sprintf(`
    <html>
    <body>
        <p>Your verification code is: <strong>%s</strong></p>
        <p>Click the following link to verify your account:</p>
        <p><a href="%s">Verify Account</a></p>
    </body>
    </html>
    `, user.VerificationCode, verificationURL)

	msg := []byte("MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	for attempts := 0; attempts < 3; attempts++ {
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
		if err != nil {
			log.Printf("Attempt %d - Error sending email: %s\n", attempts+1, err)
			time.Sleep(2 * time.Second)
		} else {
			log.Printf("Verification email sent to %s\n", to)
			return nil
		}
	}
	return fmt.Errorf("failed to send verification email after multiple attempts")
}
