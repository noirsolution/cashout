package cashout

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Cashout struct {
	Webhook string
	Proxies []string
}

// NewCashout create a new cashout object, with the Discord webhook link for now
func NewCashout(webhook string, proxies []string) *Cashout {
	m := &Cashout{
		Webhook: webhook,
		Proxies: proxies,
	}

	return m
}

// SendDiscordWebhook is used to send a new embed to the Monitor's webhook
func (c *Cashout) SendDiscordWebhook(webhook Cashout) error {
	webhookJSON, err := json.Marshal(webhook)

	if err != nil {
		log.Fatal(err)
	}

	_, err = http.Post(c.Webhook, "application/json", bytes.NewBuffer(webhookJSON))

	if err != nil {
		return err
	}

	return nil
}
