package tg

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

// ShippingQuery contains information about an incoming shipping query.
type ShippingQuery struct {
	ID              string           `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery contains information about an incoming pre-checkout query.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             *User      `json:"from"`
	Curency          string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
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
	TransactionPartnerTypeOther       TransactionPartnerType = "other"
)

// TransactionPartner describes the source of a transaction, or its recipient for outgoing transactions.
type TransactionPartner struct {
	Type TransactionPartnerType `json:"type"`

	// user
	User           *User  `json:"user,omitempty"`
	InvoicePayload string `json:"invoice_payload,omitempty"`

	// fragment
	WithdrawalState *RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
}

// RevenueWithdrawalStateType represents the type of a revenue withdrawal state.
type RevenueWithdrawalStateType string

// all available revenue withdrawal state types.
const (
	RevenueWithdrawalStateTypePending   RevenueWithdrawalStateType = "pending"
	RevenueWithdrawalStateTypeSucceeded RevenueWithdrawalStateType = "succeeded"
	RevenueWithdrawalStateTypeFailed    RevenueWithdrawalStateType = "failed"
)

// RevenueWithdrawalState describes the state of a revenue withdrawal operation.
type RevenueWithdrawalState struct {
	Type RevenueWithdrawalStateType `json:"type"`

	// succeeded
	Date int64  `json:"date"`
	URL  string `json:"url"`
}
