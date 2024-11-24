package telebot

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// required to test send and edit methods
	token = os.Getenv("TELEBOT_SECRET")
	b, _  = newTestBot() // cached bot instance to avoid getMe method flooding

	chatID, _    = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	userID, _    = strconv.ParseInt(os.Getenv("USER_ID"), 10, 64)
	channelID, _ = strconv.ParseInt(os.Getenv("CHANNEL_ID"), 10, 64)

	to      = &Chat{ID: chatID}    // to chat recipient for send and edit methods
	user    = &User{ID: userID}    // to user recipient for some special cases
	channel = &Chat{ID: channelID} // to channel recipient for some special cases

	logo  = FromURL("https://telegra.ph/file/c95b8fe46dd3df15d12e5.png")
	thumb = FromURL("https://telegra.ph/file/fe28e378784b3a4e367fb.png")
)

func defaultHandler() *Handler {
	return NewHandler(HandlerSettings{})
}

func defaultSettings() Settings {
	return Settings{Token: token, Handler: defaultHandler()}
}

func newTestBot() (*Bot, error) {
	return NewBot(defaultSettings())
}

func TestNewBot(t *testing.T) {
	var pref Settings
	_, err := NewBot(pref)
	assert.Error(t, err)

	pref.Token = "BAD TOKEN"
	_, err = NewBot(pref)
	assert.Error(t, err)

	pref.URL = "BAD URL"
	_, err = NewBot(pref)
	assert.Error(t, err)

	b, err := NewBot(Settings{Handler: NewHandler(HandlerSettings{Offline: true})})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, b.Me)
	assert.Equal(t, 100, cap(b.Updates))

	pref = defaultSettings()
	pref.URL = "http://api.telegram.org" // not https
	pref.Poller = &LongPoller{Timeout: time.Second}
	pref.Updates = 50

	b, err = NewBot(pref)
	require.NoError(t, err)
	assert.Equal(t, pref.Poller, b.Poller)
	assert.Equal(t, 50, cap(b.Updates))
}

func TestBot(t *testing.T) {
	if b == nil {
		t.Skip("Cached bot instance is bad (probably wrong or empty TELEBOT_SECRET)")
	}
	if chatID == 0 {
		t.Skip("CHAT_ID is required for Bot methods test")
	}

	_, err := b.Send(to, nil)
	assert.Equal(t, ErrUnsupportedWhat, err)
	_, err = b.Edit(&Message{Chat: &Chat{}}, nil)
	assert.Equal(t, ErrUnsupportedWhat, err)

	_, err = b.Send(nil, "")
	assert.Equal(t, ErrBadRecipient, err)
	_, err = b.Forward(nil, nil)
	assert.Equal(t, ErrBadRecipient, err)

	photo := &Photo{
		File:    logo,
		Caption: t.Name(),
	}
	var msg *Message

	t.Run("Send(what=Sendable)", func(t *testing.T) {
		msg, err = b.Send(to, photo)
		require.NoError(t, err)
		assert.NotNil(t, msg.Photo)
		assert.Equal(t, photo.Caption, msg.Caption)
	})

	t.Run("SendAlbum()", func(t *testing.T) {
		_, err = b.SendAlbum(nil, nil)
		assert.Equal(t, ErrBadRecipient, err)

		_, err = b.SendAlbum(to, nil)
		assert.Error(t, err)

		photo2 := *photo
		photo2.Caption = ""

		msgs, err := b.SendAlbum(to, Album{photo, &photo2}, ModeHTML)
		require.NoError(t, err)
		assert.Len(t, msgs, 2)
		assert.NotEmpty(t, msgs[0].AlbumID)
	})

	t.Run("SendPaid()", func(t *testing.T) {
		_, err = b.SendPaid(nil, 0, nil)
		assert.Equal(t, ErrBadRecipient, err)

		_, err = b.SendPaid(channel, 0, nil)
		assert.Error(t, err)

		photo2 := *photo
		photo2.Caption = ""

		msg, err := b.SendPaid(channel, 1, PaidAlbum{photo, &photo2}, ModeHTML)
		require.NoError(t, err)
		require.NotNil(t, msg)
		assert.Equal(t, 1, msg.PaidMedia.Stars)
		assert.Equal(t, 2, len(msg.PaidMedia.PaidMedia))
	})

	t.Run("EditCaption()+ParseMode", func(t *testing.T) {
		b.handler.parseMode = "html"

		edited, err := b.EditCaption(msg, "<b>new caption with html</b>")
		require.NoError(t, err)
		assert.Equal(t, "new caption with html", edited.Caption)
		assert.Equal(t, EntityBold, edited.CaptionEntities[0].Type)

		sleep()

		edited, err = b.EditCaption(msg, "*new caption with markdown*", ModeMarkdown)
		require.NoError(t, err)
		assert.Equal(t, "new caption with markdown", edited.Caption)
		assert.Equal(t, EntityBold, edited.CaptionEntities[0].Type)

		sleep()

		edited, err = b.EditCaption(msg, "_new caption with markdown \\(V2\\)_", ModeMarkdownV2)
		require.NoError(t, err)
		assert.Equal(t, "new caption with markdown (V2)", edited.Caption)
		assert.Equal(t, EntityItalic, edited.CaptionEntities[0].Type)
	})

	t.Run("Edit(what=Media)", func(t *testing.T) {
		photo.Caption = "<code>new caption with html</code>"

		edited, err := b.Edit(msg, photo)
		require.NoError(t, err)
		assert.Equal(t, edited.Photo.UniqueID, photo.UniqueID)
		assert.Equal(t, EntityCode, edited.CaptionEntities[0].Type)

		resp, err := http.Get("https://telegra.ph/file/274e5eb26f348b10bd8ee.mp4")
		require.NoError(t, err)
		defer resp.Body.Close()

		file, err := ioutil.TempFile("", "")
		require.NoError(t, err)

		_, err = io.Copy(file, resp.Body)
		require.NoError(t, err)

		animation := &Animation{
			File:     FromDisk(file.Name()),
			Caption:  t.Name(),
			FileName: "animation.gif",
		}

		msg, err := b.Send(msg.Chat, animation)
		require.NoError(t, err)

		if msg.Animation != nil {
			assert.Equal(t, msg.Animation.FileID, animation.FileID)
		} else {
			assert.Equal(t, msg.Document.FileID, animation.FileID)
		}

		_, err = b.Edit(edited, animation)
		require.NoError(t, err)
	})

	t.Run("Edit(what=Animation)", func(t *testing.T) {})

	t.Run("Send(what=string)", func(t *testing.T) {
		msg, err = b.Send(to, t.Name())
		require.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		rpl, err := b.Reply(msg, t.Name())
		require.NoError(t, err)
		assert.Equal(t, rpl.Text, msg.Text)
		assert.NotNil(t, rpl.ReplyTo)
		assert.Equal(t, rpl.ReplyTo, msg)
		assert.True(t, rpl.IsReply())

		fwd, err := b.Forward(to, msg)
		require.NoError(t, err)
		assert.NotNil(t, msg, fwd)
		assert.True(t, fwd.IsForwarded())

		fwd.ID += 1 // nonexistent message
		_, err = b.Forward(to, fwd)
		assert.Equal(t, ErrNotFoundToForward, err)
	})

	t.Run("Edit(what=string)", func(t *testing.T) {
		msg, err = b.Edit(msg, t.Name())
		require.NoError(t, err)
		assert.Equal(t, t.Name(), msg.Text)

		_, err = b.Edit(msg, msg.Text)
		assert.Error(t, err) // message is not modified
	})

	t.Run("Edit(what=ReplyMarkup)", func(t *testing.T) {
		good := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{
				{{
					Data: "btn",
					Text: "Hi Telebot!",
				}},
			},
		}
		bad := &ReplyMarkup{
			InlineKeyboard: [][]InlineButton{
				{{
					Data: strings.Repeat("*", 65),
					Text: "Bad Button",
				}},
			},
		}

		edited, err := b.Edit(msg, good)
		require.NoError(t, err)
		assert.Equal(t, edited.ReplyMarkup.InlineKeyboard, good.InlineKeyboard)

		edited, err = b.EditReplyMarkup(edited, nil)
		require.NoError(t, err)
		assert.Nil(t, edited.ReplyMarkup)

		_, err = b.Edit(edited, bad)
		assert.Equal(t, ErrBadButtonData, err)
	})

	t.Run("Edit(what=Location)", func(t *testing.T) {
		loc := &Location{Lat: 42, Lng: 69, LivePeriod: 60}
		edited, err := b.Send(to, loc)
		require.NoError(t, err)
		assert.NotNil(t, edited.Location)

		loc = &Location{Lat: loc.Lng, Lng: loc.Lat}
		edited, err = b.Edit(edited, *loc)
		require.NoError(t, err)
		assert.NotNil(t, edited.Location)
	})

	// Should be after the Edit tests.
	t.Run("Delete()", func(t *testing.T) {
		require.NoError(t, b.Delete(msg))
	})

	t.Run("Notify()", func(t *testing.T) {
		assert.Equal(t, ErrBadRecipient, b.Notify(nil, Typing))
		require.NoError(t, b.Notify(to, Typing))
	})

	t.Run("Answer()", func(t *testing.T) {
		assert.Error(t, b.Answer(&Query{}, &QueryResponse{
			Results: Results{&ArticleResult{}},
		}))
	})

	t.Run("Respond()", func(t *testing.T) {
		assert.Error(t, b.Respond(&Callback{}, &CallbackResponse{}))
	})

	t.Run("Payments", func(t *testing.T) {
		assert.NotPanics(t, func() {
			b.Accept(&PreCheckoutQuery{})
			b.Accept(&PreCheckoutQuery{}, "error")
		})
		assert.NotPanics(t, func() {
			b.Ship(&ShippingQuery{})
			b.Ship(&ShippingQuery{}, "error")
			b.Ship(&ShippingQuery{}, ShippingOption{}, ShippingOption{})
			assert.Equal(t, ErrUnsupportedWhat, b.Ship(&ShippingQuery{}, 0))
		})
	})

	t.Run("Commands", func(t *testing.T) {
		var (
			set1 = []Command{{
				Text:        "test1",
				Description: "test command 1",
			}}
			set2 = []Command{{
				Text:        "test2",
				Description: "test command 2",
			}}
			scope = CommandScope{
				Type:   CommandScopeChat,
				ChatID: chatID,
			}
		)

		err := b.SetCommands(set1)
		require.NoError(t, err)

		cmds, err := b.Commands()
		require.NoError(t, err)
		assert.Equal(t, set1, cmds)

		err = b.SetCommands(set2, "en", scope)
		require.NoError(t, err)

		cmds, err = b.Commands()
		require.NoError(t, err)
		assert.Equal(t, set1, cmds)

		cmds, err = b.Commands("en", scope)
		require.NoError(t, err)
		assert.Equal(t, set2, cmds)

		require.NoError(t, b.DeleteCommands("en", scope))
		require.NoError(t, b.DeleteCommands())
	})

	t.Run("InviteLink", func(t *testing.T) {
		inviteLink, err := b.CreateInviteLink(&Chat{ID: chatID}, nil)
		require.NoError(t, err)
		assert.True(t, len(inviteLink.InviteLink) > 0)

		sleep()

		response, err := b.EditInviteLink(&Chat{ID: chatID}, &ChatInviteLink{InviteLink: inviteLink.InviteLink})
		require.NoError(t, err)
		assert.True(t, len(response.InviteLink) > 0)

		sleep()

		response, err = b.RevokeInviteLink(&Chat{ID: chatID}, inviteLink.InviteLink)
		require.Nil(t, err)
		assert.True(t, len(response.InviteLink) > 0)
	})
}

func sleep() {
	time.Sleep(time.Second)
}
