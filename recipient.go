package tgot

import (
	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
)

// AsRecipientID returns RecipientID that contains chat ID.
func (c ChatID) AsRecipientID() RecipientID {
	return RecipientID{
		chatid:   c.id,
		username: c.username,
	}
}

// // AsRecipientID returns RecipientID that contains user ID.
func (u *User) AsRecipientID() RecipientID {
	return RecipientID{userid: u.id}
}

// RecipientID represents user id, chat id or channel username.
type RecipientID struct {
	chatid   int64
	userid   int64
	username string
}

func (r RecipientID) ChatID() int64    { return r.chatid }
func (r RecipientID) UserID() int64    { return r.userid }
func (r RecipientID) Username() string { return r.username }

func (r RecipientID) set(d *api.Data) {
	d.Set("chat_id", r.username)
	d.SetInt64("chat_id", r.chatid)
	d.SetInt64("user_id", r.userid)
}

// WithRecipient creates a new Recipient from context and recipient id.
// It copies the context but resets the context parameters.
func WithRecipient(ctx BaseContext, r RecipientID) *Recipient {
	d := api.NewData()
	r.set(d)
	return &Recipient{
		context: ctx.ctx().with(d),
		id:      r,
	}
}

// Recipient provides method to work with recipients.
type Recipient struct {
	*context
	id RecipientID
}

// Gift contains sendGift method parameters.
type Gift struct {
	ID            string             `tg:"gift_id"`
	PayForUpgrade bool               `tg:"pay_for_upgrade"`
	Text          string             `tg:"text"`
	ParseMode     tg.ParseMode       `tg:"text_parse_mode"`
	Entities      []tg.MessageEntity `tg:"text_entities"`
}

// SendGift sends a gift to the given user or channel chat.
func (r *Recipient) SendGift(g Gift) error {
	return r.method("sendGift", api.NewDataFrom(g))
}
