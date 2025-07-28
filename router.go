package tgot

import "github.com/karalef/tgot/api/tg"

var _ Handler = (*Router)(nil)

// Router contains all available updates handler functions.
type Router struct {
	OnMessage                 func(*Message, *tg.Message)
	OnEditedMessage           func(*Message, *tg.Message)
	OnChannelPost             func(*Message, *tg.Message)
	OnEditedChannelPost       func(*Message, *tg.Message)
	OnBusinessConnnection     func(*User, *tg.BusinessConnection)
	OnBusinessMessage         func(*Message, *tg.Message)
	OnEditedBusinessMessage   func(*Message, *tg.Message)
	OnDeletedBusinessMessages func(*Chat, *tg.BusinessMessagesDeleted)
	OnCallbackQuery           func(Query[CallbackAnswer], *tg.CallbackQuery)
	OnMessageReaction         func(*Message, *tg.MessageReactionUpdated)
	OnMessageReactionCount    func(*Message, *tg.MessageReactionCountUpdated)
	OnInlineQuery             func(Query[InlineAnswer], *tg.InlineQuery)
	OnInlineChosen            func(*Message, *tg.InlineChosen)
	OnShippingQuery           func(Query[ShippingAnswer], *tg.ShippingQuery)
	OnPreCheckoutQuery        func(Query[PreCheckoutAnswer], *tg.PreCheckoutQuery)
	OnPurchasedPaidMedia      func(*User, *tg.PaidMediaPurchased)
	OnPoll                    func(BaseContext, *tg.Poll)
	OnPollAnswer              func(BaseContext, *tg.PollAnswer)
	OnMyChatMember            func(*Chat, *tg.ChatMemberUpdated)
	OnChatMember              func(*ChatMember, *tg.ChatMemberUpdated)
	OnChatJoinRequest         func(*ChatMember, *tg.ChatJoinRequest)
	OnChatBoost               func(*ChatMember, *tg.ChatBoostUpdated)
	OnChatBoostRemoved        func(*ChatMember, *tg.ChatBoostRemoved)
}

// Allowed returns list of allowed updates.
// If there is any change in the handler, this function must be called to get a new list,
// otherwise it may cause a panic.
func (h *Router) Allowed() []string {
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

func (h *Router) Handle(ctx Empty, upd *tg.Update) {
	switch {

	case upd.Message != nil && h.OnMessage != nil:
		h.OnMessage(WithMessage(ctx, ChatMsgID(upd.Message)), upd.Message)

	case upd.EditedMessage != nil && h.OnEditedMessage != nil:
		h.OnEditedMessage(WithMessage(ctx, ChatMsgID(upd.EditedMessage)), upd.EditedMessage)

	case upd.ChannelPost != nil && h.OnChannelPost != nil:
		h.OnChannelPost(WithMessage(ctx, ChatMsgID(upd.ChannelPost)), upd.ChannelPost)

	case upd.EditedChannelPost != nil && h.OnEditedChannelPost != nil:
		h.OnEditedChannelPost(WithMessage(ctx, ChatMsgID(upd.EditedChannelPost)), upd.EditedChannelPost)

	case upd.BusinessConnection != nil && h.OnBusinessConnnection != nil:
		h.OnBusinessConnnection(WithUser(ctx, upd.BusinessConnection.User.ID), upd.BusinessConnection)

	case upd.BusinessMessage != nil && h.OnBusinessMessage != nil:
		h.OnBusinessMessage(WithMessage(ctx, ChatMsgID(upd.BusinessMessage)), upd.BusinessMessage)

	case upd.EditedBusinessMessage != nil && h.OnEditedBusinessMessage != nil:
		h.OnEditedBusinessMessage(WithMessage(ctx, ChatMsgID(upd.EditedBusinessMessage)), upd.EditedBusinessMessage)

	case upd.DeletedBusinessMessages != nil && h.OnDeletedBusinessMessages != nil:
		chatid := NewChatID(upd.DeletedBusinessMessages.Chat.ID, upd.DeletedBusinessMessages.ConnID)
		h.OnDeletedBusinessMessages(WithChatID(ctx, chatid), upd.DeletedBusinessMessages)

	case upd.CallbackQuery != nil && h.OnCallbackQuery != nil:
		qctx := WithQuery[CallbackAnswer](ctx, upd.CallbackQuery.ID, upd.CallbackQuery.From)
		h.OnCallbackQuery(qctx, upd.CallbackQuery)

	case upd.MessageReaction != nil && h.OnMessageReaction != nil:
		msgid := MsgID(upd.MessageReaction.Chat.ID, upd.MessageReaction.MessageID)
		h.OnMessageReaction(WithMessage(ctx, msgid), upd.MessageReaction)

	case upd.MessageReactionCount != nil && h.OnMessageReactionCount != nil:
		msgid := MsgID(upd.MessageReactionCount.Chat.ID, upd.MessageReactionCount.MessageID)
		h.OnMessageReactionCount(WithMessage(ctx, msgid), upd.MessageReactionCount)

	case upd.InlineQuery != nil && h.OnInlineQuery != nil:
		qctx := WithQuery[InlineAnswer](ctx, upd.InlineQuery.ID, upd.InlineQuery.From)
		h.OnInlineQuery(qctx, upd.InlineQuery)

	case upd.InlineChosen != nil && h.OnInlineChosen != nil:
		h.OnInlineChosen(WithMessage(ctx, InlineMsgID(upd.InlineChosen)), upd.InlineChosen)

	case upd.ShippingQuery != nil && h.OnShippingQuery != nil:
		qctx := WithQuery[ShippingAnswer](ctx, upd.ShippingQuery.ID, upd.ShippingQuery.From)
		h.OnShippingQuery(qctx, upd.ShippingQuery)

	case upd.PreCheckoutQuery != nil && h.OnPreCheckoutQuery != nil:
		qctx := WithQuery[PreCheckoutAnswer](ctx, upd.PreCheckoutQuery.ID, upd.PreCheckoutQuery.From)
		h.OnPreCheckoutQuery(qctx, upd.PreCheckoutQuery)

	case upd.PurchasedPaidMedia != nil && h.OnPurchasedPaidMedia != nil:
		h.OnPurchasedPaidMedia(WithUser(ctx, upd.PurchasedPaidMedia.From.ID), upd.PurchasedPaidMedia)

	case upd.Poll != nil && h.OnPoll != nil:
		h.OnPoll(ctx, upd.Poll)

	case upd.PollAnswer != nil && h.OnPollAnswer != nil:
		h.OnPollAnswer(ctx, upd.PollAnswer)

	case upd.MyChatMember != nil && h.OnMyChatMember != nil:
		h.OnMyChatMember(WithChatID(ctx, NewChatID(upd.MyChatMember.Chat.ID)), upd.MyChatMember)

	case upd.ChatMember != nil && h.OnChatMember != nil:
		chatid := NewChatID(upd.ChatMember.Chat.ID)
		h.OnChatMember(WithChatMember(ctx, chatid, upd.ChatMember.New.User.ID), upd.ChatMember)

	case upd.ChatJoinRequest != nil && h.OnChatJoinRequest != nil:
		chatid := NewChatID(upd.ChatMember.Chat.ID)
		h.OnChatJoinRequest(WithChatMember(ctx, chatid, upd.ChatJoinRequest.From.ID), upd.ChatJoinRequest)

	case upd.ChatBoost != nil && h.OnChatBoost != nil:
		chatid := NewChatID(upd.ChatBoost.Chat.ID)
		userid := upd.ChatBoost.Boost.Source.User.ID
		h.OnChatBoost(WithChatMember(ctx, chatid, userid), upd.ChatBoost)

	case upd.RemovedChatBoost != nil && h.OnChatBoostRemoved != nil:
		chatid := NewChatID(upd.RemovedChatBoost.Chat.ID)
		userid := upd.RemovedChatBoost.Source.User.ID
		h.OnChatBoostRemoved(WithChatMember(ctx, chatid, userid), upd.RemovedChatBoost)
	}
}
