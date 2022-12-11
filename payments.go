package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// CreateInvoiceLink contains parameters for creating an invoice link.
type CreateInvoiceLink = tg.InputInvoiceMessageContent

// CreateInvoiceLink creates a link for an invoice.
func (c Context) CreateInvoiceLink(l CreateInvoiceLink) (string, error) {
	d := api.NewData()
	Invoice{
		Title:                     l.Title,
		Description:               l.Description,
		Payload:                   l.Payload,
		ProviderToken:             l.ProviderToken,
		Currency:                  l.Currency,
		Prices:                    l.Prices,
		MaxTipAmount:              l.MaxTipAmount,
		SuggestedTipAmounts:       l.SuggestedTipAmounts,
		ProviderData:              l.ProviderData,
		PhotoURL:                  l.PhotoURL,
		PhotoSize:                 l.PhotoSize,
		PhotoWidth:                l.PhotoWidth,
		PhotoHeight:               l.PhotoHeight,
		NeedName:                  l.NeedName,
		NeedPhoneNumber:           l.NeedPhoneNumber,
		NeedEmail:                 l.NeedEmail,
		NeedShippingAddress:       l.NeedShippingAddress,
		SendPhoneNumberToProvider: l.SendPhoneNumberToProvider,
		SendEmailToProvider:       l.SendEmailToProvider,
		IsFlexible:                l.IsFlexible,
	}.sendData(d)
	return method[string](c, "createInvoiceLink", d)
}

// ShippingContext type.
type ShippingContext = QueryContext[ShippingAnswer]

// ShippingAnswer represents an answer to shipping query.
type ShippingAnswer struct {
	OK              bool
	ShippingOptions []tg.ShippingOption
	ErrorMessage    string
}

func (a ShippingAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("shipping_query_id", queryID)
	d.SetBool("ok", a.OK)
	d.SetJSON("shipping_options", a.ShippingOptions)
	d.Set("error_message", a.ErrorMessage)
	return "answerShippingQuery"
}

// PreCheckoutContext type.
type PreCheckoutContext = QueryContext[PreCheckoutAnswer]

// PreCheckoutAnswer represents an answer to pre-checkout query.
type PreCheckoutAnswer struct {
	OK           bool
	ErrorMessage string
}

func (a PreCheckoutAnswer) answerData(d *api.Data, queryID string) string {
	d.Set("pre_checkout_query_id", queryID)
	d.SetBool("ok", a.OK)
	d.Set("error_message", a.ErrorMessage)
	return "answerPreCheckoutQuery"
}
