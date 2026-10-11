package main

import (
	"context"
	"strconv"
)

type Payment struct {
	ID          string
	AmountCents int
	Recipient   string
}

func NeedsConfirmation(p Payment) bool { return p.AmountCents >= 100000 }

func NotifyPayment(ctx context.Context, c *Client, p Payment, templateID string) (string, error) {
	if !NeedsConfirmation(p) {
		return "", nil
	}
	return c.SendSMS(ctx, p.Recipient, templateID, map[string]string{"payment_id": p.ID, "amount_cents": strconv.Itoa(p.AmountCents)})
}
