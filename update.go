package telebot

import (
	"net/http"
	"regexp"
	"strings"
)

// Update object represents an incoming update.
type Update struct {
	ID int `json:"update_id"`

	Message                 *Message                 `json:"message,omitempty"`
	EditedMessage           *Message                 `json:"edited_message,omitempty"`
	ChannelPost             *Message                 `json:"channel_post,omitempty"`
	EditedChannelPost       *Message                 `json:"edited_channel_post,omitempty"`
	MessageReaction         *MessageReaction         `json:"message_reaction"`
	MessageReactionCount    *MessageReactionCount    `json:"message_reaction_count"`
	Callback                *Callback                `json:"callback_query,omitempty"`
	Query                   *Query                   `json:"inline_query,omitempty"`
	InlineResult            *InlineResult            `json:"chosen_inline_result,omitempty"`
	ShippingQuery           *ShippingQuery           `json:"shipping_query,omitempty"`
	PreCheckoutQuery        *PreCheckoutQuery        `json:"pre_checkout_query,omitempty"`
	Poll                    *Poll                    `json:"poll,omitempty"`
	PollAnswer              *PollAnswer              `json:"poll_answer,omitempty"`
	MyChatMember            *ChatMemberUpdate        `json:"my_chat_member,omitempty"`
	ChatMember              *ChatMemberUpdate        `json:"chat_member,omitempty"`
	ChatJoinRequest         *ChatJoinRequest         `json:"chat_join_request,omitempty"`
	Boost                   *BoostUpdated            `json:"chat_boost"`
	BoostRemoved            *BoostRemoved            `json:"removed_chat_boost"`
	BusinessConnection      *BusinessConnection      `json:"business_connection"`
	BusinessMessage         *Message                 `json:"business_message"`
	EditedBusinessMessage   *Message                 `json:"edited_business_message"`
	DeletedBusinessMessages *BusinessMessagesDeleted `json:"deleted_business_messages"`

	// Secret is a secret header passed in http request along with update.
	Secret string
	// Args is a custom map of arguments that can be passed to the update.
	Args    map[string]string
	Request *http.Request

	// Error field is used to store an error that occurred while processing the update.
	Error error
}

// ProcessUpdate processes an update with registered handlers
func (hm *Handler) ProcessUpdate(b *Bot, u Update) {
	hm.ProcessContext(NewContext(b, u))
}

var (
	cmdRx   = regexp.MustCompile(`^(/\w+)(@(\w+))?(\s|$)(.+)?`)
	cbackRx = regexp.MustCompile(`^\f([-\w]+)(\|(.+))?$`)
)

// ProcessContext processes the given context.
// A started bot calls this function automatically.
func (hm *Handler) ProcessContext(c Context) {
	u := c.Update()

	if u.Message != nil {
		m := u.Message

		if m.PinnedMessage != nil {
			hm.handle(OnPinned, c)
			return
		}

		if m.Origin != nil {
			hm.handle(OnForward, c)
		}

		// Commands
		if m.Text != "" {
			// Filtering malicious messages
			if m.Text[0] == '\a' {
				return
			}

			match := cmdRx.FindAllStringSubmatch(m.Text, -1)
			if match != nil {
				// Syntax: "</command>@<bot> <payload>"
				command, botName := match[0][1], match[0][3]

				if botName != "" && !strings.EqualFold(c.Bot().Me().Username, botName) {
					return
				}

				m.Payload = match[0][5]
				if hm.handle(command, c) {
					return
				}
			}

			// 1:1 satisfaction
			if hm.handle(m.Text, c) {
				return
			}

			if m.ReplyTo != nil {
				hm.handle(OnReply, c)
			}

			hm.handle(OnText, c)
			return
		}

		if hm.handleMedia(c) {
			return
		}

		if m.Contact != nil {
			hm.handle(OnContact, c)
			return
		}
		if m.Location != nil {
			hm.handle(OnLocation, c)
			return
		}
		if m.Venue != nil {
			hm.handle(OnVenue, c)
			return
		}
		if m.Game != nil {
			hm.handle(OnGame, c)
			return
		}
		if m.Dice != nil {
			hm.handle(OnDice, c)
			return
		}
		if m.Invoice != nil {
			hm.handle(OnInvoice, c)
			return
		}
		if m.Payment != nil {
			hm.handle(OnPayment, c)
			return
		}
		if m.RefundedPayment != nil {
			hm.handle(OnRefund, c)
			return
		}
		if m.TopicCreated != nil {
			hm.handle(OnTopicCreated, c)
			return
		}
		if m.TopicReopened != nil {
			hm.handle(OnTopicReopened, c)
			return
		}
		if m.TopicClosed != nil {
			hm.handle(OnTopicClosed, c)
			return
		}
		if m.TopicEdited != nil {
			hm.handle(OnTopicEdited, c)
			return
		}
		if m.GeneralTopicHidden != nil {
			hm.handle(OnGeneralTopicHidden, c)
			return
		}
		if m.GeneralTopicUnhidden != nil {
			hm.handle(OnGeneralTopicUnhidden, c)
			return
		}
		if m.WriteAccessAllowed != nil {
			hm.handle(OnWriteAccessAllowed, c)
			return
		}

		wasAdded := (m.UserJoined != nil && m.UserJoined.ID == c.Bot().Me().ID) ||
			(m.UsersJoined != nil && isUserInList(c.Bot().Me(), m.UsersJoined))
		if m.GroupCreated || m.SuperGroupCreated || wasAdded {
			hm.handle(OnAddedToGroup, c)
			return
		}

		if m.UserJoined != nil {
			hm.handle(OnUserJoined, c)
			return
		}
		if m.UsersJoined != nil {
			for _, user := range m.UsersJoined {
				m.UserJoined = &user
				hm.handle(OnUserJoined, c)
			}
			return
		}
		if m.UserLeft != nil {
			hm.handle(OnUserLeft, c)
			return
		}

		if m.UserShared != nil {
			hm.handle(OnUserShared, c)
			return
		}
		if m.ChatShared != nil {
			hm.handle(OnChatShared, c)
			return
		}

		if m.NewGroupTitle != "" {
			hm.handle(OnNewGroupTitle, c)
			return
		}
		if m.NewGroupPhoto != nil {
			hm.handle(OnNewGroupPhoto, c)
			return
		}
		if m.GroupPhotoDeleted {
			hm.handle(OnGroupPhotoDeleted, c)
			return
		}

		if m.GroupCreated {
			hm.handle(OnGroupCreated, c)
			return
		}
		if m.SuperGroupCreated {
			hm.handle(OnSuperGroupCreated, c)
			return
		}
		if m.ChannelCreated {
			hm.handle(OnChannelCreated, c)
			return
		}

		if m.MigrateTo != 0 {
			m.MigrateFrom = m.Chat.ID
			hm.handle(OnMigration, c)
			return
		}

		if m.VideoChatStarted != nil {
			hm.handle(OnVideoChatStarted, c)
			return
		}
		if m.VideoChatEnded != nil {
			hm.handle(OnVideoChatEnded, c)
			return
		}
		if m.VideoChatParticipants != nil {
			hm.handle(OnVideoChatParticipants, c)
			return
		}
		if m.VideoChatScheduled != nil {
			hm.handle(OnVideoChatScheduled, c)
			return
		}

		if m.WebAppData != nil {
			hm.handle(OnWebApp, c)
			return
		}

		if m.ProximityAlert != nil {
			hm.handle(OnProximityAlert, c)
			return
		}
		if m.AutoDeleteTimer != nil {
			hm.handle(OnAutoDeleteTimer, c)
			return
		}
	}

	if u.EditedMessage != nil {
		hm.handle(OnEdited, c)
		return
	}

	if u.ChannelPost != nil {
		m := u.ChannelPost

		if m.PinnedMessage != nil {
			hm.handle(OnPinned, c)
			return
		}

		hm.handle(OnChannelPost, c)
		return
	}

	if u.EditedChannelPost != nil {
		hm.handle(OnEditedChannelPost, c)
		return
	}

	if u.Callback != nil {
		if data := u.Callback.Data; data != "" && data[0] == '\f' {
			match := cbackRx.FindAllStringSubmatch(data, -1)
			if match != nil {
				unique, payload := match[0][1], match[0][3]
				if handler, ok := hm.handlers["\f"+unique]; ok {
					u.Callback.Unique = unique
					u.Callback.Data = payload
					hm.runHandler(handler, c)
					return
				}
			}
		}

		hm.handle(OnCallback, c)
		return
	}

	if u.Query != nil {
		hm.handle(OnQuery, c)
		return
	}

	if u.InlineResult != nil {
		hm.handle(OnInlineResult, c)
		return
	}

	if u.ShippingQuery != nil {
		hm.handle(OnShipping, c)
		return
	}

	if u.PreCheckoutQuery != nil {
		hm.handle(OnCheckout, c)
		return
	}

	if u.Poll != nil {
		hm.handle(OnPoll, c)
		return
	}
	if u.PollAnswer != nil {
		hm.handle(OnPollAnswer, c)
		return
	}

	if u.MyChatMember != nil {
		hm.handle(OnMyChatMember, c)
		return
	}
	if u.ChatMember != nil {
		hm.handle(OnChatMember, c)
		return
	}
	if u.ChatJoinRequest != nil {
		hm.handle(OnChatJoinRequest, c)
		return
	}

	if u.Boost != nil {
		hm.handle(OnBoost, c)
		return
	}
	if u.BoostRemoved != nil {
		hm.handle(OnBoostRemoved, c)
		return
	}

	if u.BusinessConnection != nil {
		hm.handle(OnBusinessConnection, c)
		return
	}
	if u.BusinessMessage != nil {
		hm.handle(OnBusinessMessage, c)
		return
	}
	if u.EditedBusinessMessage != nil {
		hm.handle(OnEditedBusinessMessage, c)
		return
	}
	if u.DeletedBusinessMessages != nil {
		hm.handle(OnDeletedBusinessMessages, c)
		return
	}
}

func (hm *Handler) handle(end string, c Context) bool {
	if handler, ok := hm.handlers[end]; ok {
		hm.runHandler(handler, c)
		return true
	}
	return false
}

func (hm *Handler) handleMedia(c Context) bool {
	var (
		m     = c.Message()
		fired = true
	)

	switch {
	case m.Photo != nil:
		fired = hm.handle(OnPhoto, c)
	case m.Voice != nil:
		fired = hm.handle(OnVoice, c)
	case m.Audio != nil:
		fired = hm.handle(OnAudio, c)
	case m.Animation != nil:
		fired = hm.handle(OnAnimation, c)
	case m.Document != nil:
		fired = hm.handle(OnDocument, c)
	case m.Sticker != nil:
		fired = hm.handle(OnSticker, c)
	case m.Video != nil:
		fired = hm.handle(OnVideo, c)
	case m.VideoNote != nil:
		fired = hm.handle(OnVideoNote, c)
	default:
		return false
	}

	if !fired {
		return hm.handle(OnMedia, c)
	}

	return true
}

func (hm *Handler) runHandler(h HandlerFunc, c Context) {
	f := func() {
		if err := h(c); err != nil {
			hm.onError(err, c)
		}
	}
	if hm.synchronous {
		f()
	} else {
		go f()
	}
}

func isUserInList(user *User, list []User) bool {
	for _, user2 := range list {
		if user.ID == user2.ID {
			return true
		}
	}
	return false
}
