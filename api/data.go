package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"strings"
	"unsafe"

	"github.com/karalef/tgot/api/tg"
)

// NewData creates new Data object.
func NewData() *Data {
	return &Data{
		Params: make(map[string]string),
		Files:  make(map[string]*tg.InputFile),
	}
}

// Data contains query parameters and files data.
type Data struct {
	// contains query parameters and files that does not need to be uploaded.
	Params map[string]string

	// key is a multipart field name.
	Files map[string]*tg.InputFile
}

// Data encodes the values into “URL encoded” form or multipart/form-data.
func (d *Data) Data() (string, io.Reader) {
	if d == nil || len(d.Params) == 0 && len(d.Files) == 0 {
		return "", nil
	}
	if len(d.Files) > 0 {
		return d.writeMultipart()
	}

	vals := make(url.Values, len(d.Params))
	for k, v := range d.Params {
		vals.Set(k, v)
	}
	return "application/x-www-form-urlencoded", strings.NewReader(vals.Encode())
}

// Set sets the key to value.
func (d *Data) Set(k, v string, force ...bool) *Data {
	if v != "" || len(force) > 0 && force[0] {
		d.Params[k] = v
	}
	return d
}

// SetInt sets the key to int value.
func (d *Data) SetInt(key string, v int, force ...bool) *Data {
	return d.SetInt64(key, int64(v), force...)
}

// SetInt64 sets the key to int64 value.
func (d *Data) SetInt64(key string, v int64, force ...bool) *Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatInt(v, 10))
	}
	return d
}

// SetFloat sets the key to float value.
func (d *Data) SetFloat(key string, v float32, force ...bool) *Data {
	return d.SetFloat64(key, float64(v), force...)
}

// SetFloat64 sets the key to float64 value.
func (d *Data) SetFloat64(key string, v float64, force ...bool) *Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatFloat(v, 'f', 6, 64))
	}
	return d
}

// SetBool sets the key to bool value.
func (d *Data) SetBool(key string, v bool, force ...bool) *Data {
	if v || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatBool(v))
	}
	return d
}

// SetJSON sets the key to JSON value.
func (d *Data) SetJSON(key string, v any) *Data {
	if v != nil {
		b, _ := json.Marshal(v)
		d.Set(key, string(b))
	}
	return d
}

// SetFile sets file with thumbnail.
func (d *Data) SetFile(field string, file, thumb tg.Inputtable) {
	d.AddFile(field, file)
	d.AddFile("thumb", thumb)
}

// AddFile adds file.
func (d *Data) AddFile(field string, file tg.Inputtable) {
	if isNil(file) {
		return
	}

	if inp, ok := file.(*tg.InputFile); ok {
		d.Files[field] = inp
	} else {
		urlid, _ := file.FileData()
		d.Params[field] = urlid
	}
}

func (d *Data) writeMultipart() (string, io.Reader) {
	r, w := io.Pipe()
	mp := multipart.NewWriter(w)
	go func() {
		defer func() {
			w.CloseWithError(mp.Close())
		}()

		for field, v := range d.Params {
			if err := mp.WriteField(field, v); err != nil {
				w.CloseWithError(err)
				return
			}
		}

		for field, file := range d.Files {
			name, reader := file.FileData()
			part, err := mp.CreateFormFile(field, name)
			if err != nil {
				w.CloseWithError(err)
				return
			}
			if reader == nil {
				continue
			}
			_, err = io.Copy(part, reader)
			if err != nil {
				w.CloseWithError(err)
				return
			}
		}
	}()
	return mp.FormDataContentType(), r
}

func isNil(a any) bool {
	return (*[2]uintptr)(unsafe.Pointer(&a))[1] == 0
}
