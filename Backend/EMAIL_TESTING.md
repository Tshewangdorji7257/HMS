# Email Notification Testing Guide

## Quick Start - Testing Email Notifications

### Step 1: Configure Gmail SMTP (5 minutes)

1. **Enable 2-Step Verification on Gmail:**
   - Go to: https://myaccount.google.com/security
   - Enable "2-Step Verification"

2. **Create App Password:**
   - Go to: https://myaccount.google.com/apppasswords
   - Select "Mail" and "Windows Computer"
   - Copy the 16-character password (e.g., `abcd efgh ijkl mnop`)

3. **Update docker-compose.yml:**
   ```yaml
   booking-service:
     environment:
       # ... keep existing variables ...
       # ADD these lines (remove the # at the start):
       SMTP_HOST: smtp.gmail.com
       SMTP_PORT: 587
       SMTP_USER: your-email@gmail.com
       SMTP_PASSWORD: abcdefghijklmnop  # paste your app password (remove spaces)
       FROM_EMAIL: your-email@gmail.com
       FROM_NAME: Hostel Management System
   ```

4. **Restart booking service:**
   ```powershell
   cd Backend
   docker-compose restart booking-service
   ```

### Step 2: Test Booking Email (2 minutes)

1. **Open your frontend application:**
   - Make sure you're logged in with an account that has an email address

2. **Book a room:**
   - Browse buildings
   - Select a room with available beds
   - Click "Book Now" on an available bed
   - Confirm the booking

3. **Check your email:**
   - Look for: "🎉 Booking Confirmed - Your Hostel Room is Reserved!"
   - Check spam folder if not in inbox
   - Email should arrive within 5-30 seconds

### Step 3: Test Cancellation Email (1 minute)

1. **Cancel the booking:**
   - Go to "My Bookings" page
   - Click "Cancel Booking" on your active booking
   - Confirm cancellation

2. **Check your email:**
   - Look for: "❌ Booking Cancelled - Your Reservation has been Cancelled"
   - Should arrive within 5-30 seconds

### Step 4: Verify in Logs

```powershell
cd Backend
docker-compose logs -f booking-service
```

Look for these messages:
- ✅ `Email sent successfully to your-email@gmail.com`
- Or: ⚠️ `Email notifications disabled: SMTP credentials not configured` (if not configured)
- Or: ❌ `Failed to send email: <error>` (if configuration is wrong)

## Testing Without Configuration

If you don't configure SMTP:
- ✅ Bookings work normally
- ✅ Cancellations work normally
- ⚠️ No emails are sent
- ⚠️ Warning logged: "Email notifications disabled"
- ✅ **No errors or failures**

## Troubleshooting

### "Username and Password not accepted"
**Solution:** You're using your regular Gmail password. You MUST use an App Password.
1. Enable 2-Step Verification first
2. Create an App Password at https://myaccount.google.com/apppasswords
3. Use that 16-character password (remove spaces)

### Emails go to spam
**Solution:** Mark as "Not Spam" in Gmail. Or use a dedicated email service like SendGrid.

### No email received
**Checklist:**
- [ ] SMTP credentials configured in docker-compose.yml?
- [ ] Booking service restarted after config change?
- [ ] Check spam/junk folder?
- [ ] User account has email address?
- [ ] Check logs: `docker-compose logs booking-service | Select-String "email"`

### Email takes long to arrive
**Normal:** 5-30 seconds delay is normal
**If longer:** Check your email provider's rate limits

## Sample Email Content

### Booking Confirmation:
```
From: Hostel Management System <your-email@gmail.com>
To: student@example.com
Subject: 🎉 Booking Confirmed - Your Hostel Room is Reserved!

[Beautiful HTML email with:]
- Booking ID
- Building name
- Room number
- Bed number  
- Booking date
- Next steps
```

### Cancellation Confirmation:
```
From: Hostel Management System <your-email@gmail.com>
To: student@example.com
Subject: ❌ Booking Cancelled - Your Reservation has been Cancelled

[Beautiful HTML email with:]
- Booking ID
- Building name
- Room number
- Bed number
- Cancellation date
- Rebooking information
```

## Production Setup

For production, use a professional email service:

### SendGrid (Recommended):
```yaml
SMTP_HOST: smtp.sendgrid.net
SMTP_PORT: 587
SMTP_USER: apikey
SMTP_PASSWORD: SG.xxxxx-your-api-key-xxxxx
FROM_EMAIL: noreply@yourdomain.com
FROM_NAME: Your Company Name
```

Benefits:
- 100 emails/day free
- 99.9% delivery rate
- Email analytics
- Better deliverability

### AWS SES:
```yaml
SMTP_HOST: email-smtp.us-east-1.amazonaws.com
SMTP_PORT: 587
SMTP_USER: your-aws-username
SMTP_PASSWORD: your-aws-password
FROM_EMAIL: verified@yourdomain.com
FROM_NAME: Your Company Name
```

Benefits:
- $0.10 per 1,000 emails
- Highly scalable
- AWS integration
- Reliable delivery

## Quick Commands

**View email-related logs:**
```powershell
docker-compose logs booking-service | Select-String "email"
```

**Check SMTP configuration:**
```powershell
docker inspect hostel_booking_service | ConvertFrom-Json | Select-Object -ExpandProperty Config | Select-Object -ExpandProperty Env | Select-String "SMTP"
```

**Restart after config change:**
```powershell
docker-compose restart booking-service
```

**View live logs:**
```powershell
docker-compose logs -f booking-service
```

## Need Help?

1. **Check the full documentation:** `EMAIL_NOTIFICATIONS.md`
2. **View logs:** `docker-compose logs booking-service`
3. **Verify configuration:** Check `docker-compose.yml`
4. **Test SMTP separately:** Use a tool like `telnet` or online SMTP tester

---

**Ready?** Configure your SMTP settings and test your first booking! 📧
