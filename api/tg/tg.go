// package tg implements all Telegram Bot API types.
// Type names and fields could have been renamed to be shorter and more readable,
// json fields are equal to Telegram field names.
package tg

import "time"

// Identifier is a common interface for all Telegram identifiers.
type Identifier interface{ id() }

// ChatID represents destination chat identifier.
type ChatID Identifier

// ID represents a Telegram integer identifier (e.g. user id, chat id, message id).
type ID int64

var _ ChatID = ID(0)

func (ID) id() {}

// ChatUsername creates a ChatID from channel username.
func ChatUsername(username string) ChatID { return Username(username) }

// Username represents channel username.
type Username string

var _ ChatID = Username("")

func (Username) id() {}

// Date represents a Unix timestamp.
type Date int64

func (d Date) Time() time.Time {
	return time.Unix(int64(d), 0)
}

// Duration represents a time duration in seconds.
type Duration int64

func (d Duration) Duration() time.Duration {
	return time.Duration(d) * time.Second
}

// RGB represents a color in RGB format.
type RGB uint32

func (rgb RGB) Red() uint8   { return uint8(rgb >> 16) }
func (rgb RGB) Green() uint8 { return uint8(rgb >> 8) }
func (rgb RGB) Blue() uint8  { return uint8(rgb) }

// ARGB represents a color in ARGB format.
type ARGB uint32

func (argb ARGB) Alpha() uint8 { return uint8(argb >> 24) }
func (argb ARGB) Red() uint8   { return uint8(argb >> 16) }
func (argb ARGB) Green() uint8 { return uint8(argb >> 8) }
func (argb ARGB) Blue() uint8  { return uint8(argb) }
