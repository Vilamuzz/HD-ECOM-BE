package services

import (
    "bytes"
    "context"
    "fmt"
    "html/template"
    "os"
    "time"

    "github.com/mailgun/mailgun-go/v4"
)

type TicketCommentEmailData struct {
    UserName     string
    TicketId     string
    TicketTitle  string
    Date         string
    Resolution   string
    CurrentYear  int
}

func (s *appService) SendTicketCommentEmail(toEmail, userName, ticketId, ticketTitle, resolution string) error {
    domain := os.Getenv("MAILGUN_DOMAIN")
    apiKey := os.Getenv("MAILGUN_API_KEY")
    fromEmail := os.Getenv("MAILGUN_FROM_EMAIL")

    if domain == "" || apiKey == "" || fromEmail == "" {
        return fmt.Errorf("mailgun configuration missing")
    }

    mg := mailgun.NewMailgun(domain, apiKey)

    // Prepare template data
    data := TicketCommentEmailData{
        UserName:    userName,
        TicketId:    ticketId,
        TicketTitle: ticketTitle,
        Date:        time.Now().Format("02 January 2006, 15:04"),
        Resolution:  resolution,
        CurrentYear: time.Now().Year(),
    }

    // Parse and execute template
    htmlBody, err := s.renderTicketCommentEmailTemplate(data)
    if err != nil {
        return fmt.Errorf("failed to render email template: %v", err)
    }

    // Create message
    message := mg.NewMessage(
        fromEmail,
        "Tiket Anda Diselesaikan ✔",
        "", // Plain text version (optional)
        toEmail,
    )
    message.SetHtml(htmlBody)

    // Send email
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
    defer cancel()

    _, _, err = mg.Send(ctx, message)
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }

    return nil
}

func (s *appService) renderTicketCommentEmailTemplate(data TicketCommentEmailData) (string, error) {
    tmpl := `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Tiket Anda Diselesaikan</title>
</head>

<body style="margin:0; padding:0; background-color:#f3f4f6; font-family:Arial,Helvetica,sans-serif;">

    <center style="width:100%; padding:20px 0; background-color:#f3f4f6;">

        <table width="600" style="width:600px; max-width:600px; background:#ffffff; border-radius:6px; border-collapse:collapse;">
            
            <!-- HEADER -->
            <tr>
                <td style="background:#f59e0b; padding:30px; text-align:center; color:#ffffff;">
                    <h1 style="margin:0; font-size:24px; font-weight:bold;">Tiket Diselesaikan ✔</h1>
                    <p style="margin:8px 0 0; font-size:14px;">Respons dari tim support kami</p>
                </td>
            </tr>

            <!-- BODY -->
            <tr>
                <td style="padding:30px; font-size:14px; color:#374151;">

                    <!-- GREETING -->
                    <p style="margin:0 0 15px;">
                        Halo <strong>{{.UserName}}</strong>,
                    </p>

                    <p style="margin:0 0 20px; line-height:1.6;">
                        Terima kasih telah menghubungi kami. Tim support telah menyelesaikan tiket Anda dan memberikan respons berikut:
                    </p>

                    <!-- TICKET INFO BOX -->
                    <table width="100%" style="width:100%; border-collapse:collapse; background:#f9fafb; border:1px solid #e5e7eb; border-radius:4px;">
                        <tr>
                            <td style="padding:12px; border-bottom:1px solid #e5e7eb;">
                                <strong style="color:#6b7280;">Nomor Tiket:</strong>
                                <span style="float:right; color:#111827;">#{{.TicketId}}</span>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding:12px; border-bottom:1px solid #e5e7eb;">
                                <strong style="color:#6b7280;">Judul:</strong>
                                <span style="float:right; color:#111827;">{{.TicketTitle}}</span>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding:12px; border-bottom:1px solid #e5e7eb;">
                                <strong style="color:#6b7280;">Status:</strong>
                                <span style="float:right; background:#10b981; color:#ffffff; padding:4px 10px; border-radius:12px; font-size:12px;">
                                    RESOLVED
                                </span>
                            </td>
                        </tr>

                        <tr>
                            <td style="padding:12px;">
                                <strong style="color:#6b7280;">Tanggal:</strong>
                                <span style="float:right; color:#111827;">{{.Date}}</span>
                            </td>
                        </tr>
                    </table>

                    <!-- RESOLUTION -->
                    <h3 style="margin:30px 0 12px; font-size:16px; color:#111827;">
                        Balasan dari Tim Support
                    </h3>

                    <div style="background:#ecfdf5; border:1px solid #d1fae5; padding:15px; border-radius:4px; color:#065f46; line-height:1.6; white-space:pre-wrap;">{{.Resolution}}</div>

                    <!-- MESSAGE BOX -->
                    <div style="margin-top:25px; background:#fef3c7; border:1px solid #fde68a; padding:15px; color:#92400e; border-radius:4px;">
                        <strong>Catatan:</strong> Jika Anda membutuhkan bantuan tambahan, jangan ragu untuk membuka tiket baru.
                    </div>

                    <!-- BUTTON -->
                    <div style="text-align:center; margin:30px 0;">
                        <a href="https://helpdesk.magangslab.store/contact-support"
                            style="background:#f59e0b; padding:12px 24px; display:inline-block; text-decoration:none; color:#ffffff; font-weight:bold; border-radius:4px;">
                            Lihat Tiket Lengkap
                        </a>
                    </div>

                </td>
            </tr>

            <!-- FOOTER -->
            <tr>
                <td style="padding:25px 30px; text-align:center; font-size:12px; color:#6b7280;">
                    <strong>SecondCycle Help Center</strong>
                    <br><br>
                    Kami siap membantu Anda kapan saja.<br>
                    Kunjungi:
                    <a href="https://helpdesk.magangslab.store" style="color:#f59e0b; text-decoration:none;">Pusat Bantuan</a>
                    <br><br>
                    <a href="https://magangslab.store" style="color:#f59e0b; margin:0 6px;">Website</a> |
                    <a href="https://helpdesk.magangslab.store/help" style="color:#f59e0b; margin:0 6px;">Bantuan</a> |
                    <a href="https://helpdesk.magangslab.store/privacy" style="color:#f59e0b; margin:0 6px;">Privasi</a>

                    <br><br>
                    <span style="color:#9ca3af;">
                        © {{.CurrentYear}} SecondCycle. Email ini dikirim karena Anda membuat tiket support.
                    </span>
                </td>
            </tr>

        </table>

    </center>

</body>
</html>`

    t, err := template.New("email").Parse(tmpl)
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := t.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}