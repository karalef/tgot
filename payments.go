package tgot

import "github.com/karalef/tgot/tg"

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	tg.InputInvoiceMessageContent
	StartParameter string
}

func (i Invoice) params(p params) {
	p.set("title", i.Title)
	p.set("description", i.Description)
	p.set("payload", i.Payload)
	p.set("provider_token", i.ProviderToken)
	p.set("currency", i.Currency)
	p.setJSON("prices", i.Prices)
	p.setInt("max_tip_amount", i.MaxTipAmount)
	p.setJSON("suggested_tip_amounts", i.SuggestedTipAmounts)
	p.set("start_parameter", i.StartParameter)
	p.set("provider_data", i.ProviderData)
	p.set("photo_url", i.PhotoURL)
	p.setInt("photo_size", i.PhotoSize)
	p.setInt("photo_width", i.PhotoWidth)
	p.setInt("photo_height", i.PhotoHeight)
	p.setBool("need_name", i.NeedName)
	p.setBool("need_phone_number", i.NeedPhoneNumber)
	p.setBool("need_email", i.NeedEmail)
	p.setBool("need_shipping_address", i.NeedShippingAddress)
	p.setBool("send_phone_number_to_provider", i.SendPhoneNumberToProvider)
	p.setBool("send_email_to_provider", i.SendEmailToProvider)
	p.setBool("is_flexible", i.IsFlexible)
}

// CreateInvoiceLink contains parameters for creating an invoice link.
type CreateInvoiceLink struct {
	tg.InputInvoiceMessageContent
}

func (c CreateInvoiceLink) params(p params) {
	// since the parameters are exactly the same, except for
	// StartParameter (since it is empty, it will not be included),
	// it is possible to use the already existing serialization function.
	i := Invoice{c.InputInvoiceMessageContent, ""}
	i.params(p)
}

// CreateInvoiceLink creates a link for an invoice.
func (c Context) CreateInvoiceLink(l CreateInvoiceLink) (string, error) {
	p := params{}
	l.params(p)
	return api[string](c, "createInvoiceLink", p)
}

var _ QueryContext[ShippingAnswer]

// ShippingAnswer represents an answer to shipping query.
type ShippingAnswer struct {
	OK              bool
	ShippingOptions []tg.ShippingOption
	ErrorMessage    string
}

func (a ShippingAnswer) answerData(p params, queryID string) string {
	p.set("shipping_query_id", queryID)
	p.setBool("ok", a.OK)
	p.setJSON("shipping_options", a.ShippingOptions)
	p.set("error_message", a.ErrorMessage)
	return "answerShippingQuery"
}

var _ QueryContext[PreCheckoutAnswer]

// PreCheckoutAnswer represents an answer to pre-checkout query.
type PreCheckoutAnswer struct {
	OK           bool
	ErrorMessage string
}

func (a PreCheckoutAnswer) answerData(p params, queryID string) string {
	p.set("pre_checkout_query_id", queryID)
	p.setBool("ok", a.OK)
	p.set("error_message", a.ErrorMessage)
	return "answerPreCheckoutQuery"
}
