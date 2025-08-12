package service

import (
	"fmt"

	"github.com/mucunga90/ecommerce/config"
	"github.com/mucunga90/ecommerce/internal"

	"gopkg.in/gomail.v2"
)

type notifier struct {
	cfg *config.Config
	sms service
}

func NewNotifier(cfg *config.Config) *notifier {
	sms := NewSMS(cfg.SMS.AfricaSTUser, cfg.SMS.AfricaSTKey, cfg.SMS.SMTPHost)
	return &notifier{cfg: cfg, sms: sms}
}

func (n *notifier) SendOrderSMS(evt internal.OrderCreatedEvent) error {
	message := fmt.Sprintf("Hello %s, your order #%d has been placed. Total: %.2f",
		evt.CustomerName, evt.OrderID, evt.Total)
	_, err := n.sms.Send(n.cfg.SMS.AfricaSTUser, evt.CustomerPhone, message)
	return err
}

func (n *notifier) SendOrderEmail(evt internal.OrderCreatedEvent) error {
	m := gomail.NewMessage()
	m.SetHeader("From", n.cfg.Email.EmailFrom)
	m.SetHeader("To", n.cfg.Email.AdminEmail)
	m.SetHeader("Subject", fmt.Sprintf("New Order #%d", evt.OrderID))

	body := fmt.Sprintf("Customer: %s\nOrder ID: %d\nTotal: %.2f\n\nItems:\n",
		evt.CustomerName, evt.OrderID, evt.Total)
	for _, item := range evt.Items {
		body += fmt.Sprintf("Product ID: %d | Qty: %d | Price: %.2f\n", item.ProductID, item.Quantity, item.UnitPrice)
	}
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(n.cfg.Email.SMTPHost, n.cfg.Email.SMTPPort, n.cfg.Email.EmailFrom, n.cfg.Email.EmailPassword)
	return d.DialAndSend(m)
}
