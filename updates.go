package tgot

import (
	"context"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// GetUpdates contains parameters for getUpdates method.
type GetUpdates struct {
	Offset  int
	Limit   int
	Timeout int
	Allowed []string
}

// GetUpdates receives incoming updates using long polling.
func (b *Bot) GetUpdates(ctx context.Context, p GetUpdates) ([]tg.Update, error) {
	d := api.NewData()
	d.SetInt("limit", p.Limit)
	d.SetInt("timeout", p.Timeout)
	d.SetJSON("allowed", p.Allowed)
	d.SetInt("offset", p.Offset)
	return api.RequestContext[[]tg.Update](ctx, b.api, "getUpdates", d)
}

// WebhookData contains parameters for setWebhook method.
type WebhookData struct {
	URL            string
	Certificate    *tg.InputFile
	IPAddress      string
	MaxConnections int
	AllowedUpdates []string
	DropPending    bool
	SecretToken    string
}

// SetWebhook specifies a webhook URL.
// Use this method to specify a URL and receive incoming updates via an outgoing webhook.
func (b *Bot) SetWebhook(s WebhookData) (bool, error) {
	d := api.NewData().Set("url", s.URL)
	d.AddFile("certificate", s.Certificate)
	d.Set("ip_address", s.IPAddress)
	d.SetInt("max_connections", s.MaxConnections)
	d.SetJSON("allowed_updates", s.AllowedUpdates)
	d.SetBool("drop_pending_updates", s.DropPending)
	d.Set("secret_token", s.SecretToken)
	return api.Request[bool](b.api, "setWebhook", d)
}

// DeleteWebhook removes webhook integration if you decide to switch back to getUpdates.
func (b *Bot) DeleteWebhook(dropPending bool) (bool, error) {
	d := api.NewData().SetBool("drop_pending_updates", dropPending)
	return api.Request[bool](b.api, "deleteWebhook", d)
}

// GetWebhookInfo returns current webhook status.
func (b *Bot) GetWebhookInfo() (*tg.WebhookInfo, error) {
	return api.Request[*tg.WebhookInfo](b.api, "getWebhookInfo", nil)
}
