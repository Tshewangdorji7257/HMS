package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

// EmailConfig holds email configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// GetEmailConfig returns email configuration from environment variables
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getEnv("SMTP_PORT", "587"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@hostelmgmt.com"),
		FromName:     getEnv("FROM_NAME", "Hostel Management System"),
	}
}

// BookingConfirmationData holds data for booking confirmation email
type BookingConfirmationData struct {
	StudentName  string
	BuildingName string
	RoomNumber   string
	BedNumber    int
	BookingDate  string
	BookingID    string
}

// BookingCancellationData holds data for booking cancellation email
type BookingCancellationData struct {
	StudentName  string
	BuildingName string
	RoomNumber   string
	BedNumber    int
	CancelDate   string
	BookingID    string
}

// SendBookingConfirmationEmail sends a booking confirmation email
func SendBookingConfirmationEmail(toEmail string, data BookingConfirmationData) error {
	config := GetEmailConfig()

	// Skip if email credentials are not configured
	if config.SMTPUser == "" || config.SMTPPassword == "" {
		log.Println("⚠️  Email notifications disabled: SMTP credentials not configured")
		return nil
	}

	subject := "🎉 Booking Confirmed - Your Hostel Room is Reserved!"
	body := generateBookingConfirmationHTML(data)

	return sendEmail(config, toEmail, subject, body)
}

// SendBookingCancellationEmail sends a booking cancellation email
func SendBookingCancellationEmail(toEmail string, data BookingCancellationData) error {
	config := GetEmailConfig()

	// Skip if email credentials are not configured
	if config.SMTPUser == "" || config.SMTPPassword == "" {
		log.Println("⚠️  Email notifications disabled: SMTP credentials not configured")
		return nil
	}

	subject := "❌ Booking Cancelled - Your Reservation has been Cancelled"
	body := generateBookingCancellationHTML(data)

	return sendEmail(config, toEmail, subject, body)
}

// sendEmail sends an email using SMTP
func sendEmail(config *EmailConfig, to, subject, body string) error {
	// Email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build email message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// SMTP authentication
	auth := smtp.PlainAuth("", config.SMTPUser, config.SMTPPassword, config.SMTPHost)

	// Send email
	addr := fmt.Sprintf("%s:%s", config.SMTPHost, config.SMTPPort)
	err := smtp.SendMail(addr, auth, config.FromEmail, []string{to}, []byte(message))
	
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return err
	}

	log.Printf("✅ Email sent successfully to %s", to)
	return nil
}

// generateBookingConfirmationHTML generates HTML for booking confirmation email
func generateBookingConfirmationHTML(data BookingConfirmationData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .booking-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #eee; }
        .detail-label { font-weight: bold; color: #667eea; }
        .button { display: inline-block; background: #667eea; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Booking Confirmed!</h1>
            <p>Your hostel room has been successfully reserved</p>
        </div>
        <div class="content">
            <p>Dear {{.StudentName}},</p>
            <p>Congratulations! Your booking has been confirmed. Here are your reservation details:</p>
            
            <div class="booking-details">
                <h3 style="color: #667eea; margin-top: 0;">Booking Details</h3>
                <div class="detail-row">
                    <span class="detail-label">Booking ID:</span>
                    <span>{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Building:</span>
                    <span>{{.BuildingName}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Room Number:</span>
                    <span>{{.RoomNumber}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Bed Number:</span>
                    <span>{{.BedNumber}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Booking Date:</span>
                    <span>{{.BookingDate}}</span>
                </div>
            </div>

            <h3>📋 Next Steps:</h3>
            <ul>
                <li>Report to the hostel office within 7 days with your ID proof</li>
                <li>Complete the payment and documentation process</li>
                <li>Collect your room keys and access card</li>
                <li>Review the hostel rules and regulations</li>
            </ul>

            <p><strong>Important:</strong> Please keep this email for your records. You will need your Booking ID when visiting the hostel office.</p>
            
            <div style="text-align: center; margin-top: 30px;">
                <p style="color: #666;">If you have any questions, please contact the hostel administration.</p>
            </div>
        </div>
        <div class="footer">
            <p>This is an automated email from Hostel Management System</p>
            <p>Please do not reply to this email</p>
        </div>
    </div>
</body>
</html>
`

	t := template.Must(template.New("booking-confirmation").Parse(tmpl))
	var body bytes.Buffer
	t.Execute(&body, data)
	return body.String()
}

// generateBookingCancellationHTML generates HTML for booking cancellation email
func generateBookingCancellationHTML(data BookingCancellationData) string {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
        .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
        .booking-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #eee; }
        .detail-label { font-weight: bold; color: #f5576c; }
        .button { display: inline-block; background: #667eea; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Booking Cancelled</h1>
            <p>Your reservation has been cancelled</p>
        </div>
        <div class="content">
            <p>Dear {{.StudentName}},</p>
            <p>This email confirms that your hostel booking has been cancelled.</p>
            
            <div class="booking-details">
                <h3 style="color: #f5576c; margin-top: 0;">Cancelled Booking Details</h3>
                <div class="detail-row">
                    <span class="detail-label">Booking ID:</span>
                    <span>{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Building:</span>
                    <span>{{.BuildingName}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Room Number:</span>
                    <span>{{.RoomNumber}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Bed Number:</span>
                    <span>{{.BedNumber}}</span>
                </div>
                <div class="detail-row">
                    <span class="detail-label">Cancellation Date:</span>
                    <span>{{.CancelDate}}</span>
                </div>
            </div>

            <p>The bed is now available for other students to book.</p>

            <h3>📋 What's Next:</h3>
            <ul>
                <li>You can browse and book other available rooms</li>
                <li>Visit our hostel management system to explore options</li>
                <li>If you have any refund queries, contact the hostel office</li>
            </ul>
            
            <div style="text-align: center; margin-top: 30px;">
                <p style="color: #666;">If you have any questions or concerns, please contact the hostel administration.</p>
            </div>
        </div>
        <div class="footer">
            <p>This is an automated email from Hostel Management System</p>
            <p>Please do not reply to this email</p>
        </div>
    </div>
</body>
</html>
`

	t := template.Must(template.New("booking-cancellation").Parse(tmpl))
	var body bytes.Buffer
	t.Execute(&body, data)
	return body.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
