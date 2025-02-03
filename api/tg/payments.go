package tg

import (
	"github.com/karalef/tgot/api/internal/oneof"
)

// LabeledPrice represents a portion of the price for goods or services.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

// Invoice contains basic information about an invoice..
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"` // Three-letter ISO 4217 currency code
	TotalAmount    int    `json:"total_amount"`
}

// ShippingAddress represents a shipping address.
type ShippingAddress struct {
	CountryCode string `json:"country_code"` // Two-letter ISO 3166-1 alpha-2 country code
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo represents information about an order.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// ShippingOption represents one shipping option.
type ShippingOption struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

// SuccessfulPayment contains basic information about a successful payment.
type SuccessfulPayment struct {
	Curency                 string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id,omitempty"`
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// StarTransactions contains a list of Telegram Star transactions.
type StarTransactions struct {
	Transactions []StarTransaction `json:"transactions"`
}

// StarTransaction describes a Telegram Star transaction.
type StarTransaction struct {
	ID       string              `json:"id"`
	Amount   uint                `json:"amount"`
	Date     int64               `json:"date"`
	Source   *TransactionPartner `json:"source"`
	Receiver *TransactionPartner `json:"receiver"`
}

// TransactionPartnerType represents the type of a transaction partner.
type TransactionPartnerType string

// all available transaction partner types.
const (
	TransactionPartnerTypeUser        TransactionPartnerType = "user"
	TransactionPartnerTypeFragment    TransactionPartnerType = "fragment"
	TransactionPartnerTypeTelegramAds TransactionPartnerType = "telegram_ads"
	TransactionPartnerTypeTelegramApi TransactionPartnerType = "telegram_api"
	TransactionPartnerTypeOther       TransactionPartnerType = "other"
)

var transactionPartnerTypes = oneof.NewMap[TransactionPartnerType](
	TransactionPartnerUser{},
	TransactionPartnerFragment{},
	TransactionPartnerTelegramAds{},
	TransactionPartnerTelegramApi{},
	TransactionPartnerOther{},
)

func (TransactionPartnerType) TypeFor(t TransactionPartnerType) oneof.Type {
	return transactionPartnerTypes.TypeFor(t)
}

// TransactionPartner describes the source of a transaction, or its recipient for outgoing transactions.
type TransactionPartner = oneof.Object[TransactionPartnerType, oneof.IDTypeType]

// TransactionPartnerUser describes a transaction with a user.
type TransactionPartnerUser struct {
	User             *User       `json:"user"`
	InvoicePayload   string      `json:"invoice_payload,omitempty"`
	PaidMedia        []PaidMedia `json:"paid_media,omitempty"`
	PaidMediaPayload string      `json:"paid_media_payload"`
}

func (TransactionPartnerUser) Type() TransactionPartnerType { return TransactionPartnerTypeUser }

// TransactionPartnerFragment describes a transaction with Fragment.
type TransactionPartnerFragment struct {
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
}

func (TransactionPartnerFragment) Type() TransactionPartnerType {
	return TransactionPartnerTypeFragment
}

// TransactionPartnerTelegramAds describes a withdrawal transaction to the Telegram Ads platform.
type TransactionPartnerTelegramAds struct{}

func (TransactionPartnerTelegramAds) Type() TransactionPartnerType {
	return TransactionPartnerTypeTelegramAds
}

// TransactionPartnerTelegramApi describes a transaction with payment for paid broadcasting.
type TransactionPartnerTelegramApi struct{}

func (TransactionPartnerTelegramApi) Type() TransactionPartnerType {
	return TransactionPartnerTypeTelegramApi
}

// TransactionPartnerOther describes a transaction with an unknown source or recipient.
type TransactionPartnerOther struct{}

func (TransactionPartnerOther) Type() TransactionPartnerType { return TransactionPartnerTypeOther }

// RevenueWithdrawalStateType represents the type of a revenue withdrawal state.
type RevenueWithdrawalStateType string

// all available revenue withdrawal state types.
const (
	RevenueWithdrawalStateTypePending   RevenueWithdrawalStateType = "pending"
	RevenueWithdrawalStateTypeSucceeded RevenueWithdrawalStateType = "succeeded"
	RevenueWithdrawalStateTypeFailed    RevenueWithdrawalStateType = "failed"
)

var revenueWithdrawalStateTypes = oneof.NewMap[RevenueWithdrawalStateType](
	RevenueWithdrawalStatePending{},
	RevenueWithdrawalStateSucceeded{},
	RevenueWithdrawalStateFailed{},
)

func (RevenueWithdrawalStateType) TypeFor(t RevenueWithdrawalStateType) oneof.Type {
	return revenueWithdrawalStateTypes.TypeFor(t)
}

// RevenueWithdrawalState describes the state of a revenue withdrawal operation.
type RevenueWithdrawalState = oneof.Object[RevenueWithdrawalStateType, oneof.IDTypeType]

// RevenueWithdrawalStatePending means the withdrawal is in progress.
type RevenueWithdrawalStatePending struct{}

func (RevenueWithdrawalStatePending) Type() RevenueWithdrawalStateType {
	return RevenueWithdrawalStateTypePending
}

// RevenueWithdrawalStateSucceeded means the withdrawal was successed.
type RevenueWithdrawalStateSucceeded struct {
	Date int64  `json:"date"`
	URL  string `json:"url"`
}

func (RevenueWithdrawalStateSucceeded) Type() RevenueWithdrawalStateType {
	return RevenueWithdrawalStateTypeSucceeded
}

// RevenueWithdrawalStateFailed means the withdrawal failed and the transaction was refunded.
type RevenueWithdrawalStateFailed struct{}

func (RevenueWithdrawalStateFailed) Type() RevenueWithdrawalStateType {
	return RevenueWithdrawalStateTypeFailed
}

// RefundedPayment contains basic information about a refunded payment.
type RefundedPayment struct {
	Currency                string `json:"currency"`
	TotalAmount             uint   `json:"total_amount"`
	InvoicePayload          string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
}
