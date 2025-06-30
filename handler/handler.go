package handler

import (
	"errors"

	"github.com/karalef/tgot"
	"github.com/karalef/tgot/api/tg"
)

// ErrNotSpecified is returned by Handle if the update contains an object not specified in Handler.
var ErrNotSpecified = errors.New("handler is not specified")

var _ tgot.Handler = (*Handler)(nil)

// Handler contains all available updates handler functions.
type Handler struct {
	OnMessage                 func(*tgot.Message, *tg.Message)
	OnEditedMessage           func(*tgot.Message, *tg.Message)
	OnChannelPost             func(*tgot.Message, *tg.Message)
	OnEditedChannelPost       func(*tgot.Message, *tg.Message)
	OnBusinessConnnection     func(*tgot.User, *tg.BusinessConnection)
	OnBusinessMessage         func(*tgot.Message, *tg.Message)
	OnEditedBusinessMessage   func(*tgot.Message, *tg.Message)
	OnDeletedBusinessMessages func(*tgot.Chat, *tg.BusinessMessagesDeleted)
	OnCallbackQuery           func(tgot.Query[tgot.CallbackAnswer], *tg.CallbackQuery)
	OnMessageReaction         func(*tgot.Message, *tg.MessageReactionUpdated)
	OnMessageReactionCount    func(*tgot.Message, *tg.MessageReactionCountUpdated)
	OnInlineQuery             func(tgot.Query[tgot.InlineAnswer], *tg.InlineQuery)
	OnInlineChosen            func(*tgot.Message, *tg.InlineChosen)
	OnShippingQuery           func(tgot.Query[tgot.ShippingAnswer], *tg.ShippingQuery)
	OnPreCheckoutQuery        func(tgot.Query[tgot.PreCheckoutAnswer], *tg.PreCheckoutQuery)
	OnPurchasedPaidMedia      func(*tgot.User, *tg.PaidMediaPurchased)
	OnPoll                    func(tgot.BaseContext, *tg.Poll)
	OnPollAnswer              func(tgot.BaseContext, *tg.PollAnswer)
	OnMyChatMember            func(*tgot.Chat, *tg.ChatMemberUpdated)
	OnChatMember              func(*tgot.ChatMember, *tg.ChatMemberUpdated)
	OnChatJoinRequest         func(*tgot.ChatMember, *tg.ChatJoinRequest)
	OnChatBoost               func(*tgot.ChatMember, *tg.ChatBoostUpdated)
	OnChatBoostRemoved        func(*tgot.ChatMember, *tg.ChatBoostRemoved)
}

// Allowed returns list of allowed updates.
// If there is any change in the handler, this function must be called to get a new list,
// otherwise it may cause a panic.
func (h *Handler) Allowed() []string {
	list := make([]string, 0, 23)
	add := func(a bool, s string) {
		if a {
			list = append(list, s)
		}
	}
	add(h.OnMessage != nil, tg.UpdateTypeMessage)
	add(h.OnEditedMessage != nil, tg.UpdateTypeEditedMessage)
	add(h.OnChannelPost != nil, tg.UpdateTypeChannelPost)
	add(h.OnEditedChannelPost != nil, tg.UpdateTypeEditedChannelPost)
	add(h.OnBusinessConnnection != nil, tg.UpdateTypeBusinessConnection)
	add(h.OnBusinessMessage != nil, tg.UpdateTypeBusinessMessage)
	add(h.OnEditedBusinessMessage != nil, tg.UpdateTypeEditedBusinessMessage)
	add(h.OnDeletedBusinessMessages != nil, tg.UpdateTypeDeletedBusinessMessages)
	add(h.OnCallbackQuery != nil, tg.UpdateTypeCallbackQuery)
	add(h.OnMessageReaction != nil, tg.UpdateTypeMessageReaction)
	add(h.OnMessageReactionCount != nil, tg.UpdateTypeMessageReactionCount)
	add(h.OnInlineQuery != nil, tg.UpdateTypeInlineQuery)
	add(h.OnInlineChosen != nil, tg.UpdateTypeChosenInlineQuery)
	add(h.OnShippingQuery != nil, tg.UpdateTypeShippingQuery)
	add(h.OnPreCheckoutQuery != nil, tg.UpdateTypePreCheckoutQuery)
	add(h.OnPurchasedPaidMedia != nil, tg.UpdateTypePurchasedPaidMedia)
	add(h.OnPoll != nil, tg.UpdateTypePoll)
	add(h.OnPollAnswer != nil, tg.UpdateTypePollAnswer)
	add(h.OnMyChatMember != nil, tg.UpdateTypeMyChatMember)
	add(h.OnChatMember != nil, tg.UpdateTypeChatMember)
	add(h.OnChatJoinRequest != nil, tg.UpdateTypeChatJoinRequest)
	add(h.OnChatBoost != nil, tg.UpdateTypeChatBoost)
	add(h.OnChatBoostRemoved != nil, tg.UpdateTypeRemovedChatBoost)
	return list
}

func (h *Handler) Handle(ctx tgot.Empty, upd *tg.Update) error {
	switch {
	default:
		return ErrNotSpecified

	case upd.Message != nil && h.OnMessage != nil:
		h.OnMessage(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.Message)), upd.Message)

	case upd.EditedMessage != nil && h.OnEditedMessage != nil:
		h.OnEditedMessage(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.EditedMessage)), upd.EditedMessage)

	case upd.ChannelPost != nil && h.OnChannelPost != nil:
		h.OnChannelPost(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.ChannelPost)), upd.ChannelPost)

	case upd.EditedChannelPost != nil && h.OnEditedChannelPost != nil:
		h.OnEditedChannelPost(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.EditedChannelPost)), upd.EditedChannelPost)

	case upd.BusinessConnection != nil && h.OnBusinessConnnection != nil:
		h.OnBusinessConnnection(tgot.WithUser(ctx, upd.BusinessConnection.User.ID), upd.BusinessConnection)

	case upd.BusinessMessage != nil && h.OnBusinessMessage != nil:
		h.OnBusinessMessage(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.BusinessMessage)), upd.BusinessMessage)

	case upd.EditedBusinessMessage != nil && h.OnEditedBusinessMessage != nil:
		h.OnEditedBusinessMessage(tgot.WithMessage(ctx, tgot.ChatMsgID(upd.EditedBusinessMessage)), upd.EditedBusinessMessage)

	case upd.DeletedBusinessMessages != nil && h.OnDeletedBusinessMessages != nil:
		chatid := tgot.NewChatID(upd.DeletedBusinessMessages.Chat.ID, upd.DeletedBusinessMessages.BusinessConnectionID)
		h.OnDeletedBusinessMessages(tgot.WithChatID(ctx, chatid), upd.DeletedBusinessMessages)

	case upd.CallbackQuery != nil && h.OnCallbackQuery != nil:
		qctx := tgot.WithQuery[tgot.CallbackAnswer](ctx, upd.CallbackQuery.ID, upd.CallbackQuery.From)
		h.OnCallbackQuery(qctx, upd.CallbackQuery)

	case upd.MessageReaction != nil && h.OnMessageReaction != nil:
		msgid := tgot.MsgID(upd.MessageReaction.Chat.ID, upd.MessageReaction.MessageID)
		h.OnMessageReaction(tgot.WithMessage(ctx, msgid), upd.MessageReaction)

	case upd.MessageReactionCount != nil && h.OnMessageReactionCount != nil:
		msgid := tgot.MsgID(upd.MessageReactionCount.Chat.ID, upd.MessageReactionCount.MessageID)
		h.OnMessageReactionCount(tgot.WithMessage(ctx, msgid), upd.MessageReactionCount)

	case upd.InlineQuery != nil && h.OnInlineQuery != nil:
		qctx := tgot.WithQuery[tgot.InlineAnswer](ctx, upd.InlineQuery.ID, upd.InlineQuery.From)
		h.OnInlineQuery(qctx, upd.InlineQuery)

	case upd.InlineChosen != nil && h.OnInlineChosen != nil:
		h.OnInlineChosen(tgot.WithMessage(ctx, tgot.InlineMsgID(upd.InlineChosen)), upd.InlineChosen)

	case upd.ShippingQuery != nil && h.OnShippingQuery != nil:
		qctx := tgot.WithQuery[tgot.ShippingAnswer](ctx, upd.ShippingQuery.ID, upd.ShippingQuery.From)
		h.OnShippingQuery(qctx, upd.ShippingQuery)

	case upd.PreCheckoutQuery != nil && h.OnPreCheckoutQuery != nil:
		qctx := tgot.WithQuery[tgot.PreCheckoutAnswer](ctx, upd.PreCheckoutQuery.ID, upd.PreCheckoutQuery.From)
		h.OnPreCheckoutQuery(qctx, upd.PreCheckoutQuery)

	case upd.PurchasedPaidMedia != nil && h.OnPurchasedPaidMedia != nil:
		h.OnPurchasedPaidMedia(tgot.WithUser(ctx, upd.PurchasedPaidMedia.From.ID), upd.PurchasedPaidMedia)

	case upd.Poll != nil && h.OnPoll != nil:
		h.OnPoll(ctx, upd.Poll)

	case upd.PollAnswer != nil && h.OnPollAnswer != nil:
		h.OnPollAnswer(ctx, upd.PollAnswer)

	case upd.MyChatMember != nil && h.OnMyChatMember != nil:
		h.OnMyChatMember(tgot.WithChatID(ctx, tgot.NewChatID(upd.MyChatMember.Chat.ID)), upd.MyChatMember)

	case upd.ChatMember != nil && h.OnChatMember != nil:
		chatid := tgot.NewChatID(upd.ChatMember.Chat.ID)
		h.OnChatMember(tgot.WithChatMember(ctx, chatid, upd.ChatMember.New.User.ID), upd.ChatMember)

	case upd.ChatJoinRequest != nil && h.OnChatJoinRequest != nil:
		chatid := tgot.NewChatID(upd.ChatMember.Chat.ID)
		h.OnChatJoinRequest(tgot.WithChatMember(ctx, chatid, upd.ChatJoinRequest.From.ID), upd.ChatJoinRequest)

	case upd.ChatBoost != nil && h.OnChatBoost != nil:
		chatid := tgot.NewChatID(upd.ChatBoost.Chat.ID)
		userid := upd.ChatBoost.Boost.Source.User.ID
		h.OnChatBoost(tgot.WithChatMember(ctx, chatid, userid), upd.ChatBoost)

	case upd.RemovedChatBoost != nil && h.OnChatBoostRemoved != nil:
		chatid := tgot.NewChatID(upd.RemovedChatBoost.Chat.ID)
		userid := upd.RemovedChatBoost.Source.User.ID
		h.OnChatBoostRemoved(tgot.WithChatMember(ctx, chatid, userid), upd.RemovedChatBoost)
	}

	return nil
}
