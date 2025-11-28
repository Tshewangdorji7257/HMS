# Email Notifications Setup Guide

The booking service now supports automatic email notifications for booking confirmations and cancellations. Students will receive professional HTML emails when they book or cancel their hostel rooms.

## Features

✨ **Booking Confirmation Emails**
- Sent automatically when a student successfully books a room
- Contains booking details (Building, Room, Bed, Booking ID)
- Professional HTML template with gradient header
- Includes next steps and instructions

✨ **Cancellation Confirmation Emails**
- Sent automatically when a booking is cancelled
- Confirms cancellation with booking details
- Professional HTML template
- Includes information about rebooking

## Email Templates

Both email templates include:
- 🎨 Professional HTML design with gradients and styling
- 📋 Complete booking details (Building, Room, Bed, Booking ID)
- 📅 Timestamp of booking/cancellation
- 📝 Helpful instructions and next steps
- 🔒 Automated system notification footer

## Configuration

### Option 1: Gmail SMTP (Recommended for Testing)

1. **Create a Gmail App Password:**
   - Go to your Google Account: https://myaccount.google.com/
   - Click "Security" → "2-Step Verification" (enable if not already)
   - Scroll to "App passwords"
   - Create a new app password for "Mail"
   - Copy the 16-character password

2. **Update docker-compose.yml:**
   ```yaml
   booking-service:
     environment:
       # ... other environment variables ...
       SMTP_HOST: smtp.gmail.com
       SMTP_PORT: 587
       SMTP_USER: your-gmail@gmail.com
       SMTP_PASSWORD: your-16-char-app-password
       FROM_EMAIL: your-gmail@gmail.com
       FROM_NAME: Hostel Management System
   ```

### Option 2: Other Email Providers

**SendGrid:**
```yaml
SMTP_HOST: smtp.sendgrid.net
SMTP_PORT: 587
SMTP_USER: apikey
SMTP_PASSWORD: your-sendgrid-api-key
FROM_EMAIL: noreply@yourdomain.com
FROM_NAME: Hostel Management System
```

**Outlook/Hotmail:**
```yaml
SMTP_HOST: smtp-mail.outlook.com
SMTP_PORT: 587
SMTP_USER: your-email@outlook.com
SMTP_PASSWORD: your-password
FROM_EMAIL: your-email@outlook.com
FROM_NAME: Hostel Management System
```

**AWS SES:**
```yaml
SMTP_HOST: email-smtp.us-east-1.amazonaws.com
SMTP_PORT: 587
SMTP_USER: your-aws-ses-username
SMTP_PASSWORD: your-aws-ses-password
FROM_EMAIL: verified-email@yourdomain.com
FROM_NAME: Hostel Management System
```

## Default Behavior (No Configuration)

If SMTP credentials are not configured:
- The booking service will work normally
- Email notifications will be **disabled automatically**
- A warning will be logged: `⚠️ Email notifications disabled: SMTP credentials not configured`
- Bookings and cancellations will still work perfectly
- No errors will be thrown

## How It Works

### Booking Flow:
1. Student books a room through the frontend
2. Frontend sends `user_email` in the booking request
3. Booking is created in the database
4. Bed occupancy is updated
5. **Email is sent asynchronously** (non-blocking)
6. Success response is returned immediately
7. Email delivery happens in the background

### Cancellation Flow:
1. Student cancels a booking
2. Frontend sends `user_email` as query parameter
3. Booking status is updated to 'cancelled'
4. Bed is freed up
5. **Email is sent asynchronously** (non-blocking)
6. Success response is returned immediately
7. Email delivery happens in the background

## Testing Email Notifications

### 1. Configure SMTP (using Gmail example):

Edit `Backend/docker-compose.yml`:
```yaml
booking-service:
  environment:
    # Uncomment and configure these lines:
    SMTP_HOST: smtp.gmail.com
    SMTP_PORT: 587
    SMTP_USER: your-test-email@gmail.com
    SMTP_PASSWORD: your-app-password-here
    FROM_EMAIL: your-test-email@gmail.com
    FROM_NAME: HMS Test System
```

### 2. Restart the booking service:

```powershell
cd Backend
docker-compose down booking-service
docker-compose up -d booking-service
```

### 3. Test booking:

```powershell
# Make sure user has email in their profile
# Book a room through the frontend
# Check the recipient's email inbox for confirmation
```

### 4. Test cancellation:

```powershell
# Cancel the booking through the frontend
# Check the recipient's email inbox for cancellation confirmation
```

### 5. Check logs:

```powershell
docker-compose logs -f booking-service
```

Look for:
- ✅ `Email sent successfully to user@example.com`
- ⚠️ `Email notifications disabled: SMTP credentials not configured`
- ❌ `Failed to send email to user@example.com: <error>`

## Email Templates Preview

### Booking Confirmation Email:
```
Subject: 🎉 Booking Confirmed - Your Hostel Room is Reserved!

Dear John Doe,

Congratulations! Your booking has been confirmed.

┌─────────────────────────────────┐
│ Booking Details                 │
├─────────────────────────────────┤
│ Booking ID: abc-123-xyz         │
│ Building: North Wing            │
│ Room Number: 101                │
│ Bed Number: 2                   │
│ Booking Date: January 15, 2024  │
└─────────────────────────────────┘

Next Steps:
• Report to hostel office within 7 days
• Complete payment and documentation
• Collect your room keys
```

### Cancellation Email:
```
Subject: ❌ Booking Cancelled - Your Reservation has been Cancelled

Dear John Doe,

This email confirms that your hostel booking has been cancelled.

┌─────────────────────────────────┐
│ Cancelled Booking Details       │
├─────────────────────────────────┤
│ Booking ID: abc-123-xyz         │
│ Building: North Wing            │
│ Room Number: 101                │
│ Bed Number: 2                   │
│ Cancellation Date: Jan 20, 2024 │
└─────────────────────────────────┘
```

## Important Notes

⚠️ **Security:**
- Never commit SMTP credentials to Git
- Use environment variables or secrets management
- Use app passwords instead of main passwords
- Enable 2FA on email accounts

⚠️ **Rate Limits:**
- Gmail: 500 emails/day for free accounts
- SendGrid: 100 emails/day on free tier
- Consider using a dedicated email service for production

⚠️ **Email Delivery:**
- Emails are sent asynchronously (non-blocking)
- Failed email delivery does NOT fail the booking
- Check spam/junk folders if emails don't arrive
- Monitor logs for delivery issues

⚠️ **Testing:**
- Test with real email addresses
- Check spam folders
- Verify HTML rendering in different email clients
- Test both booking and cancellation flows

## Troubleshooting

### Emails Not Being Sent

1. **Check SMTP Configuration:**
   ```powershell
   docker-compose logs booking-service | Select-String "SMTP"
   ```

2. **Verify Environment Variables:**
   ```powershell
   docker inspect hostel_booking_service | ConvertFrom-Json | Select-Object -ExpandProperty Config | Select-Object -ExpandProperty Env
   ```

3. **Check for Errors:**
   ```powershell
   docker-compose logs booking-service | Select-String "email"
   ```

### Gmail Specific Issues

**Error: "Username and Password not accepted"**
- Make sure you're using an App Password, not your regular password
- Enable 2-Step Verification first
- Generate a new App Password

**Error: "Less secure app access"**
- This is deprecated - use App Passwords instead
- Go to Google Account → Security → App Passwords

### Email Goes to Spam

- Add SPF/DKIM records if using custom domain
- Use a reputable email service provider
- Avoid spam trigger words in subject/content
- Ensure FROM_EMAIL matches SMTP_USER for Gmail

## Production Recommendations

For production deployments:

1. **Use a Dedicated Email Service:**
   - SendGrid (99.9% delivery rate)
   - AWS SES (scalable, cheap)
   - Mailgun (developer-friendly)
   - Postmark (transactional emails)

2. **Configure Proper DNS Records:**
   - SPF record for sender verification
   - DKIM for email authentication
   - DMARC for email policy

3. **Monitor Email Delivery:**
   - Track sent/delivered/bounced emails
   - Set up webhooks for delivery status
   - Log all email attempts

4. **Handle Failures Gracefully:**
   - Implement retry logic
   - Queue failed emails
   - Alert admins on repeated failures

5. **Use Email Templates:**
   - Store templates in database
   - Support multiple languages
   - A/B test email content

## Code Architecture

```
booking-service/
├── utils/
│   └── email.go           # Email service implementation
├── handlers/
│   └── booking.go         # Booking handlers with email calls
└── models/
    └── booking.go         # Updated with user_email field
```

The email service is:
- ✅ Non-blocking (uses goroutines)
- ✅ Fail-safe (errors logged, not thrown)
- ✅ Configurable (environment variables)
- ✅ Professional (HTML templates with styling)
- ✅ Informative (booking details included)

## Support

If you encounter issues:
1. Check the logs: `docker-compose logs -f booking-service`
2. Verify SMTP settings in docker-compose.yml
3. Test with a simple SMTP client first
4. Check firewall/network restrictions
5. Review the backend handlers code

---

**Ready to send emails?** Configure your SMTP settings in `docker-compose.yml` and restart the booking service!
