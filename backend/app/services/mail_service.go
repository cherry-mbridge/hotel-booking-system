package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"
)

type MailService struct {
	enabled  bool
	host     string
	port     string
	username string
	password string
	from     string
}

func NewMailService() *MailService {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "myothandatoe@gmail.com"
	}

	return &MailService{
		enabled:  host != "" && username != "" && password != "",
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *MailService) sendMail(to, subject, body string) error {
	if !s.enabled {
		log.Printf("[EMAIL MOCK] To: %s | Subject: %s\n%s\n", to, subject, body)
		return nil
	}

	msg := []byte("To: " + to + "\r\n" +
		"From: " + s.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	addr := s.host + ":" + s.port
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(addr, auth, s.from, []string{to}, msg)
}

type bookingStatusData struct {
	CustomerName string
	RoomName     string
	Status       string
	StatusColor  string
	CheckIn      string
	CheckOut     string
	TotalPrice   string
	Message      string
	Year         int
}

func (s *MailService) SendBookingStatusNotification(toEmail, customerName, roomName, status, checkIn, checkOut string, totalPrice float64) error {
	data := bookingStatusData{
		CustomerName: customerName,
		RoomName:     roomName,
		Status:       status,
		CheckIn:      checkIn,
		CheckOut:     checkOut,
		TotalPrice:   fmt.Sprintf("%.2f", totalPrice),
		Year:         time.Now().Year(),
	}

	switch status {
	case "confirmed":
		data.StatusColor = "#22c55e"
		data.Message = "Your booking has been confirmed! We look forward to welcoming you. Please arrive at the check-in desk on your scheduled date."
	case "rejected":
		data.StatusColor = "#ef4444"
		data.Message = "We regret to inform you that your booking request could not be accommodated at this time. Please contact us or try another date/room."
	case "cancelled":
		data.StatusColor = "#ef4444"
		data.Message = "Your booking has been cancelled as requested. If you have any questions, please reach out to our support team."
	default:
		data.StatusColor = "#3b82f6"
		data.Message = "Your booking status has been updated. Please review the details below."
	}

	tmplStr := `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <style>
    body { font-family: Arial, sans-serif; background: #f8fafc; margin: 0; padding: 0; }
    .container { max-width: 600px; margin: 40px auto; background: #ffffff; border-radius: 24px; overflow: hidden; box-shadow: 0 10px 40px rgba(0,0,0,0.08); }
    .header { background: #0f172a; padding: 40px 30px; text-align: center; }
    .header h1 { color: #ffffff; margin: 0; font-size: 24px; }
    .header p { color: #94a3b8; margin: 8px 0 0; font-size: 14px; }
    .content { padding: 30px; }
    .status-badge { display: inline-block; padding: 10px 24px; border-radius: 999px; color: #ffffff; font-weight: bold; font-size: 14px; text-transform: uppercase; letter-spacing: 1px; }
    .detail-row { display: flex; justify-content: space-between; padding: 14px 0; border-bottom: 1px solid #f1f5f9; }
    .detail-row:last-child { border-bottom: none; }
    .detail-label { color: #64748b; font-size: 14px; }
    .detail-value { color: #0f172a; font-weight: bold; font-size: 14px; }
    .message { background: #f8fafc; border-radius: 16px; padding: 20px; margin-top: 24px; color: #334155; font-size: 14px; line-height: 1.6; }
    .footer { text-align: center; padding: 20px 30px; color: #94a3b8; font-size: 12px; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>Lumina Hotel & Resorts</h1>
      <p>Booking Status Update</p>
    </div>
    <div class="content">
      <p style="color:#334155; font-size:16px;">Hi {{.CustomerName}},</p>
      <p style="color:#64748b; font-size:14px; margin-top:8px;">Your reservation has been updated. Here are the latest details:</p>

      <div style="text-align:center; margin: 30px 0;">
        <span class="status-badge" style="background: {{.StatusColor}};">{{.Status}}</span>
      </div>

      <div class="detail-row">
        <span class="detail-label">Room : </span>
        <span class="detail-value">{{.RoomName}}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Check In : </span>
        <span class="detail-value">{{.CheckIn}}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Check Out : </span>
        <span class="detail-value">{{.CheckOut}}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Total Price : </span>
        <span class="detail-value">${{.TotalPrice}}</span>
      </div>

      <div class="message">
        {{.Message}}
      </div>
    </div>
    <div class="footer">
      &copy; {{.Year}} Lumina Hotel & Resorts. All rights reserved.<br>
      If you have any questions, please contact our support team.
    </div>
  </div>
</body>
</html>`

	tmpl, err := template.New("booking-status").Parse(tmplStr)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	subject := fmt.Sprintf("Lumina Hotel - Booking %s", capitalize(status))
	return s.sendMail(toEmail, subject, buf.String())
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
