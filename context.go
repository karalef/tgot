package tgot

import (
	"bytes"
	"io"
	"runtime"
	"strconv"

	"github.com/karalef/tgot/api"
	"github.com/karalef/tgot/api/tg"
	"github.com/karalef/tgot/logger"
)

// MakeContext creates new context.
//
// The context must be used in a separate goroutine because in case of fatal
// errors the bot will be closed and the goroutine will be terminated.
func (b *Bot) MakeContext(name string) Context {
	return Context{bot: b, name: name}
}

// Context type.
type Context struct {
	bot  *Bot
	name string
}

// Ctx returns Context.
func (c Context) Ctx() Context { return c }

// Child creates sub context.
func (c Context) Child(name string) Context {
	c.name += "::" + name
	return c
}

// Logger returns context logger.
func (c Context) Logger() *logger.Logger {
	return c.bot.log.Child(c.name)
}

// GetMe returns basic information about the bot.
func (c Context) GetMe() (*tg.User, error) {
	return method[*tg.User](c, "getMe")
}

// GetUserPhotos returns a list of profile pictures for a user.
func (c Context) GetUserPhotos(userID int64) (*tg.UserProfilePhotos, error) {
	d := api.NewData().SetInt64("user_id", userID)
	return method[*tg.UserProfilePhotos](c, "getUserProfilePhotos", d)
}

// GetFile returns basic information about a file
// and prepares it for downloading.
func (c Context) GetFile(fileID string) (*tg.File, error) {
	d := api.NewData().Set("file_id", fileID)
	return method[*tg.File](c, "getFile", d)
}

// DownloadReader downloads file as io.ReadCloser from Telegram servers.
func (c Context) DownloadReader(f *tg.File) (io.ReadCloser, error) {
	return c.bot.api.DownloadFile(f.FilePath)
}

// Download downloads file from Telegram servers.
func (c Context) Download(f *tg.File) ([]byte, error) {
	rc, err := c.DownloadReader(f)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DownloadReaderFile downloads file as io.ReadCloser from Telegram servers
// by file id.
func (c Context) DownloadReaderFile(fid string) (io.ReadCloser, error) {
	f, err := c.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return c.DownloadReader(f)
}

// DownloadFile downloads file from Telegram servers by file id.
func (c Context) DownloadFile(fid string) ([]byte, error) {
	f, err := c.GetFile(fid)
	if err != nil {
		return nil, err
	}
	return c.Download(f)
}

func (c Context) method(meth string, d ...*api.Data) error {
	_, err := method[api.Empty](c, meth, d...)
	return err
}

func method[T any](c Context, method string, d ...*api.Data) (T, error) {
	var data *api.Data
	if len(d) > 0 {
		data = d[0]
	}
	result, err := api.Request[T](c.bot.api, method, data)
	if err == nil {
		return result, nil
	}
	if e, ok := err.(*api.Error); ok {
		switch e.Err.Code {
		case 404: // Not Found
			panic("telegram method is not found but is present in the current implementation (" + method + ")")
		case 401: // Not Authorized
			c.bot.cancel(err)
			runtime.Goexit()
		case 400, 500: // Bad Request, Internal Server Error
			c.bot.log.Error("%s\n%s\n%s", err.Error(), c.name, traceback(2))
		}
	}

	return result, err
}

func traceback(skip int) string {
	pc := make([]uintptr, 0, 16)
	for {
		n := runtime.Callers(2+skip+len(pc), pc[len(pc):cap(pc)])
		pc = pc[:len(pc)+n]
		if len(pc) < cap(pc) {
			break
		}

		newpc := make([]uintptr, len(pc), len(pc)*2)
		copy(newpc, pc)
		pc = newpc
	}

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	frames := runtime.CallersFrames(pc)

	for {
		f, more := frames.Next()
		buf.WriteString(f.Function + "\n\t" + f.File)
		buf.WriteString(":" + strconv.Itoa(f.Line))
		if !more {
			break
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}
