package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/tg"
)

// Invoice contains information about the invoice to be sent.
type Invoice struct {
	tg.InputInvoiceMessageContent
	StartParameter string
}

func (i Invoice) data() api.Data {
	d := api.NewData()
	d.Set("title", i.Title)
	d.Set("description", i.Description)
	d.Set("payload", i.Payload)
	d.Set("provider_token", i.ProviderToken)
	d.Set("currency", i.Currency)
	d.SetJSON("prices", i.Prices)
	d.SetInt("max_tip_amount", i.MaxTipAmount)
	d.SetJSON("suggested_tip_amounts", i.SuggestedTipAmounts)
	d.Set("start_parameter", i.StartParameter)
	d.Set("provider_data", i.ProviderData)
	d.Set("photo_url", i.PhotoURL)
	d.SetInt("photo_size", i.PhotoSize)
	d.SetInt("photo_width", i.PhotoWidth)
	d.SetInt("photo_height", i.PhotoHeight)
	d.SetBool("need_name", i.NeedName)
	d.SetBool("need_phone_number", i.NeedPhoneNumber)
	d.SetBool("need_email", i.NeedEmail)
	d.SetBool("need_shipping_address", i.NeedShippingAddress)
	d.SetBool("send_phone_number_to_provider", i.SendPhoneNumberToProvider)
	d.SetBool("send_email_to_provider", i.SendEmailToProvider)
	d.SetBool("is_flexible", i.IsFlexible)
	return d
}

// CreateInvoiceLink contains parameters for creating an invoice link.
type CreateInvoiceLink struct {
	tg.InputInvoiceMessageContent
}

func (c CreateInvoiceLink) data() api.Data {
	// since the parameters are exactly the same, except for
	// StartParameter (since it is empty, it will not be included),
	// it is possible to use the already existing serialization function.
	return Invoice{c.InputInvoiceMessageContent, ""}.data()
}

// CreateInvoiceLink creates a link for an invoice.
func (c Context) CreateInvoiceLink(l CreateInvoiceLink) (string, error) {
	return method[string](c, "createInvoiceLink", l.data())
}

// ShippingContext type.
type ShippingContext = QueryContext[ShippingAnswer]

// ShippingAnswer represents an answer to shipping query.
type ShippingAnswer struct {
	OK              bool
	ShippingOptions []tg.ShippingOption
	ErrorMessage    string
}

func (a ShippingAnswer) answerData(queryID string) (string, api.Data) {
	d := api.NewData()
	d.Set("shipping_query_id", queryID)
	d.SetBool("ok", a.OK)
	d.SetJSON("shipping_options", a.ShippingOptions)
	d.Set("error_message", a.ErrorMessage)
	return "answerShippingQuery", d
}

// PreCheckoutContext type.
type PreCheckoutContext = QueryContext[PreCheckoutAnswer]

// PreCheckoutAnswer represents an answer to pre-checkout query.
type PreCheckoutAnswer struct {
	OK           bool
	ErrorMessage string
}

func (a PreCheckoutAnswer) answerData(queryID string) (string, api.Data) {
	d := api.NewData()
	d.Set("pre_checkout_query_id", queryID)
	d.SetBool("ok", a.OK)
	d.Set("error_message", a.ErrorMessage)
	return "answerPreCheckoutQuery", d
}
