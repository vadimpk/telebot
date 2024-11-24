package telebot

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Handle(t *testing.T) {
	h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})

	h.Handle("/start", func(c Context) error { return nil })

	reply := ReplyButton{Text: "reply"}
	h.Handle(&reply, func(c Context) error { return nil })

	inline := InlineButton{Unique: "inline"}
	h.Handle(&inline, func(c Context) error { return nil })

	btnReply := (&ReplyMarkup{}).Text("btnReply")
	h.Handle(&btnReply, func(c Context) error { return nil })

	btnInline := (&ReplyMarkup{}).Data("", "btnInline")
	h.Handle(&btnInline, func(c Context) error { return nil })

	assert.Contains(t, h.handlers, "/start")
	assert.Contains(t, h.handlers, btnReply.CallbackUnique())
	assert.Contains(t, h.handlers, btnInline.CallbackUnique())
	assert.Contains(t, h.handlers, reply.CallbackUnique())
	assert.Contains(t, h.handlers, inline.CallbackUnique())
}

func TestHandler_ProcessUpdate(t *testing.T) {
	h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})

	b, err := NewBot(Settings{Handler: h})
	if err != nil {
		t.Fatal(err)
	}

	h.Handle(OnMedia, func(c Context) error {
		assert.NotNil(t, c.Message().Photo)
		return nil
	})

	h.Handle("/start", func(c Context) error {
		assert.Equal(t, "/start", c.Text())
		return nil
	})
	h.Handle("hello", func(c Context) error {
		assert.Equal(t, "hello", c.Text())
		return nil
	})
	h.Handle(OnText, func(c Context) error {
		assert.Equal(t, "text", c.Text())
		return nil
	})
	h.Handle(OnPinned, func(c Context) error {
		assert.NotNil(t, c.Message())
		return nil
	})
	h.Handle(OnPhoto, func(c Context) error {
		assert.NotNil(t, c.Message().Photo)
		return nil
	})
	h.Handle(OnVoice, func(c Context) error {
		assert.NotNil(t, c.Message().Voice)
		return nil
	})
	h.Handle(OnAudio, func(c Context) error {
		assert.NotNil(t, c.Message().Audio)
		return nil
	})
	h.Handle(OnAnimation, func(c Context) error {
		assert.NotNil(t, c.Message().Animation)
		return nil
	})
	h.Handle(OnDocument, func(c Context) error {
		assert.NotNil(t, c.Message().Document)
		return nil
	})
	h.Handle(OnSticker, func(c Context) error {
		assert.NotNil(t, c.Message().Sticker)
		return nil
	})
	h.Handle(OnVideo, func(c Context) error {
		assert.NotNil(t, c.Message().Video)
		return nil
	})
	h.Handle(OnVideoNote, func(c Context) error {
		assert.NotNil(t, c.Message().VideoNote)
		return nil
	})
	h.Handle(OnContact, func(c Context) error {
		assert.NotNil(t, c.Message().Contact)
		return nil
	})
	h.Handle(OnLocation, func(c Context) error {
		assert.NotNil(t, c.Message().Location)
		return nil
	})
	h.Handle(OnVenue, func(c Context) error {
		assert.NotNil(t, c.Message().Venue)
		return nil
	})
	h.Handle(OnDice, func(c Context) error {
		assert.NotNil(t, c.Message().Dice)
		return nil
	})
	h.Handle(OnInvoice, func(c Context) error {
		assert.NotNil(t, c.Message().Invoice)
		return nil
	})
	h.Handle(OnPayment, func(c Context) error {
		assert.NotNil(t, c.Message().Payment)
		return nil
	})
	h.Handle(OnRefund, func(c Context) error {
		assert.NotNil(t, c.Message().RefundedPayment)
		return nil
	})
	h.Handle(OnAddedToGroup, func(c Context) error {
		assert.NotNil(t, c.Message().GroupCreated)
		return nil
	})
	h.Handle(OnUserJoined, func(c Context) error {
		assert.NotNil(t, c.Message().UserJoined)
		return nil
	})
	h.Handle(OnUserLeft, func(c Context) error {
		assert.NotNil(t, c.Message().UserLeft)
		return nil
	})
	h.Handle(OnNewGroupTitle, func(c Context) error {
		assert.Equal(t, "title", c.Message().NewGroupTitle)
		return nil
	})
	h.Handle(OnNewGroupPhoto, func(c Context) error {
		assert.NotNil(t, c.Message().NewGroupPhoto)
		return nil
	})
	h.Handle(OnGroupPhotoDeleted, func(c Context) error {
		assert.True(t, c.Message().GroupPhotoDeleted)
		return nil
	})
	h.Handle(OnMigration, func(c Context) error {
		from, to := c.Migration()
		assert.Equal(t, int64(1), from)
		assert.Equal(t, int64(2), to)
		return nil
	})
	h.Handle(OnEdited, func(c Context) error {
		assert.Equal(t, "edited", c.Message().Text)
		return nil
	})
	h.Handle(OnChannelPost, func(c Context) error {
		assert.Equal(t, "post", c.Message().Text)
		return nil
	})
	h.Handle(OnEditedChannelPost, func(c Context) error {
		assert.Equal(t, "edited post", c.Message().Text)
		return nil
	})
	h.Handle(OnCallback, func(c Context) error {
		if data := c.Callback().Data; data[0] != '\f' {
			assert.Equal(t, "callback", data)
		}
		return nil
	})
	h.Handle("\funique", func(c Context) error {
		assert.Equal(t, "callback", c.Callback().Data)
		return nil
	})
	h.Handle(OnQuery, func(c Context) error {
		assert.Equal(t, "query", c.Data())
		return nil
	})
	h.Handle(OnInlineResult, func(c Context) error {
		assert.Equal(t, "result", c.InlineResult().ResultID)
		return nil
	})
	h.Handle(OnShipping, func(c Context) error {
		assert.Equal(t, "shipping", c.ShippingQuery().ID)
		return nil
	})
	h.Handle(OnCheckout, func(c Context) error {
		assert.Equal(t, "checkout", c.PreCheckoutQuery().ID)
		return nil
	})
	h.Handle(OnPoll, func(c Context) error {
		assert.Equal(t, "poll", c.Poll().ID)
		return nil
	})
	h.Handle(OnPollAnswer, func(c Context) error {
		assert.Equal(t, "poll", c.PollAnswer().PollID)
		return nil
	})

	h.Handle(OnWebApp, func(c Context) error {
		assert.Equal(t, "webapp", c.Message().WebAppData.Data)
		return nil
	})

	h.ProcessUpdate(b, Update{Message: &Message{Photo: &Photo{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Text: "/start"}})
	h.ProcessUpdate(b, Update{Message: &Message{Text: "/start@other_bot"}})
	h.ProcessUpdate(b, Update{Message: &Message{Text: "hello"}})
	h.ProcessUpdate(b, Update{Message: &Message{Text: "text"}})
	h.ProcessUpdate(b, Update{Message: &Message{PinnedMessage: &Message{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Photo: &Photo{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Voice: &Voice{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Audio: &Audio{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Animation: &Animation{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Document: &Document{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Sticker: &Sticker{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Video: &Video{}}})
	h.ProcessUpdate(b, Update{Message: &Message{VideoNote: &VideoNote{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Contact: &Contact{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Location: &Location{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Venue: &Venue{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Invoice: &Invoice{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Payment: &Payment{}}})
	h.ProcessUpdate(b, Update{Message: &Message{RefundedPayment: &RefundedPayment{}}})
	h.ProcessUpdate(b, Update{Message: &Message{Dice: &Dice{}}})
	h.ProcessUpdate(b, Update{Message: &Message{GroupCreated: true}})
	h.ProcessUpdate(b, Update{Message: &Message{UserJoined: &User{ID: 1}}})
	h.ProcessUpdate(b, Update{Message: &Message{UsersJoined: []User{{ID: 1}}}})
	h.ProcessUpdate(b, Update{Message: &Message{UserLeft: &User{}}})
	h.ProcessUpdate(b, Update{Message: &Message{NewGroupTitle: "title"}})
	h.ProcessUpdate(b, Update{Message: &Message{NewGroupPhoto: &Photo{}}})
	h.ProcessUpdate(b, Update{Message: &Message{GroupPhotoDeleted: true}})
	h.ProcessUpdate(b, Update{Message: &Message{Chat: &Chat{ID: 1}, MigrateTo: 2}})
	h.ProcessUpdate(b, Update{EditedMessage: &Message{Text: "edited"}})
	h.ProcessUpdate(b, Update{ChannelPost: &Message{Text: "post"}})
	h.ProcessUpdate(b, Update{ChannelPost: &Message{PinnedMessage: &Message{}}})
	h.ProcessUpdate(b, Update{EditedChannelPost: &Message{Text: "edited post"}})
	h.ProcessUpdate(b, Update{Callback: &Callback{MessageID: "inline", Data: "callback"}})
	h.ProcessUpdate(b, Update{Callback: &Callback{Data: "callback"}})
	h.ProcessUpdate(b, Update{Callback: &Callback{Data: "\funique|callback"}})
	h.ProcessUpdate(b, Update{Query: &Query{Text: "query"}})
	h.ProcessUpdate(b, Update{InlineResult: &InlineResult{ResultID: "result"}})
	h.ProcessUpdate(b, Update{ShippingQuery: &ShippingQuery{ID: "shipping"}})
	h.ProcessUpdate(b, Update{PreCheckoutQuery: &PreCheckoutQuery{ID: "checkout"}})
	h.ProcessUpdate(b, Update{Poll: &Poll{ID: "poll"}})
	h.ProcessUpdate(b, Update{PollAnswer: &PollAnswer{PollID: "poll"}})
	h.ProcessUpdate(b, Update{Message: &Message{WebAppData: &WebAppData{Data: "webapp"}}})
}

func TestHandler_OnError(t *testing.T) {
	h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})

	var ok bool
	h.onError = func(err error, c Context) {
		assert.Equal(t, b, c.(*nativeContext).b)
		assert.NotNil(t, err)
		ok = true
	}

	h.runHandler(func(c Context) error {
		return errors.New("not nil")
	}, &nativeContext{b: b})

	assert.True(t, ok)
}

func TestBotMiddleware(t *testing.T) {
	t.Run("calling order", func(t *testing.T) {
		var trace []string

		handler := func(name string) HandlerFunc {
			return func(c Context) error {
				trace = append(trace, name)
				return nil
			}
		}

		middleware := func(name string) MiddlewareFunc {
			return func(next HandlerFunc) HandlerFunc {
				return func(c Context) error {
					trace = append(trace, name+":in")
					err := next(c)
					trace = append(trace, name+":out")
					return err
				}
			}
		}

		h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})
		b, err := NewBot(Settings{Handler: h})
		if err != nil {
			t.Fatal(err)
		}

		h.Use(middleware("global1"), middleware("global2"))
		h.Handle("/a", handler("/a"), middleware("handler1a"), middleware("handler2a"))

		group := h.Group()
		group.Use(middleware("group1"), middleware("group2"))
		group.Handle("/b", handler("/b"), middleware("handler1b"))

		h.ProcessUpdate(b, Update{
			Message: &Message{Text: "/a"},
		})
		assert.Equal(t, []string{
			"global1:in", "global2:in",
			"handler1a:in", "handler2a:in",
			"/a",
			"handler2a:out", "handler1a:out",
			"global2:out", "global1:out",
		}, trace)

		trace = trace[:0]

		h.ProcessUpdate(b, Update{
			Message: &Message{Text: "/b"},
		})
		assert.Equal(t, []string{
			"global1:in", "global2:in",
			"group1:in", "group2:in",
			"handler1b:in",
			"/b",
			"handler1b:out",
			"group2:out", "group1:out",
			"global2:out", "global1:out",
		}, trace)
	})

	fatal := func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			t.Fatal("fatal middleware should not be called")
			return nil
		}
	}

	nop := func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			return next(c)
		}
	}

	t.Run("combining with global middleware", func(t *testing.T) {
		h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})
		b, err := NewBot(Settings{Handler: h})
		if err != nil {
			t.Fatal(err)
		}

		// Pre-allocate middleware slice to make sure
		// it has extra capacity after group-level middleware is added.
		h.Group().middleware = make([]MiddlewareFunc, 0, 2)
		h.Use(nop)

		h.Handle("/a", func(c Context) error { return nil }, nop)
		h.Handle("/b", func(c Context) error { return nil }, fatal)

		h.ProcessUpdate(b, Update{Message: &Message{Text: "/a"}})
	})

	t.Run("combining with group middleware", func(t *testing.T) {
		h := NewHandler(HandlerSettings{Synchronous: true, Offline: true})
		b, err := NewBot(Settings{Handler: h})
		if err != nil {
			t.Fatal(err)
		}

		g := h.Group()
		// Pre-allocate middleware slice to make sure
		// it has extra capacity after group-level middleware is added.
		g.middleware = make([]MiddlewareFunc, 0, 2)
		g.Use(nop)

		g.Handle("/a", func(c Context) error { return nil }, nop)
		g.Handle("/b", func(c Context) error { return nil }, fatal)

		h.ProcessUpdate(b, Update{Message: &Message{Text: "/a"}})
	})
}
