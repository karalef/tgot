package api

import (
	"context"

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
func (a *API) GetUpdates(ctx context.Context, p GetUpdates) ([]tg.Update, error) {
	d := NewData()
	d.SetInt("limit", p.Limit)
	d.SetInt("timeout", p.Timeout)
	d.SetJSON("allowed", p.Allowed)
	d.SetInt("offset", p.Offset)
	return RequestContext[[]tg.Update](ctx, a, "getUpdates", d)
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
func (a *API) SetWebhook(s WebhookData) (bool, error) {
	d := NewData().Set("url", s.URL)
	d.AddFile("certificate", s.Certificate)
	d.Set("ip_address", s.IPAddress)
	d.SetInt("max_connections", s.MaxConnections)
	d.SetJSON("allowed_updates", s.AllowedUpdates)
	d.SetBool("drop_pending_updates", s.DropPending)
	d.Set("secret_token", s.SecretToken)
	return Request[bool](a, "setWebhook", d)
}

// DeleteWebhook removes webhook integration if you decide to switch back to getUpdates.
func (a *API) DeleteWebhook(dropPending bool) (bool, error) {
	d := NewData().SetBool("drop_pending_updates", dropPending)
	return Request[bool](a, "deleteWebhook", d)
}

// GetWebhookInfo returns current webhook status.
func (a *API) GetWebhookInfo() (*tg.WebhookInfo, error) {
	return Request[*tg.WebhookInfo](a, "getWebhookInfo", nil)
}
