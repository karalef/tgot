package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/karalef/tgot/api/tg"
)

// NewData creates new Data object.
func NewData() *Data {
	return &Data{
		Params: make(map[string]string),
		Files:  make(map[string]tg.Inputtable),
	}
}

// Data contains query parameters and files data.
type Data struct {
	Params map[string]string

	// key must be params key or multipart field.
	Files map[string]tg.Inputtable
}

// Data encodes the values into “URL encoded” form or multipart/form-data.
func (d *Data) Data() (string, io.Reader) {
	if d == nil || len(d.Params) == 0 && len(d.Files) == 0 {
		return "", nil
	}
	for _, f := range d.Files {
		_, r := f.FileData()
		if r != nil {
			return d.writeMultipart()
		}
	}
	vals := make(url.Values, len(d.Params)+len(d.Files))
	for k, v := range d.Params {
		vals.Set(k, v)
	}
	for k, f := range d.Files {
		urlid, _ := f.FileData()
		vals.Set(k, urlid)
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
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.Itoa(v))
	}
	return d
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
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
	return d
}

// SetBool sets the key to bool value.
func (d *Data) SetBool(key string, v bool) *Data {
	if v {
		d.Set(key, strconv.FormatBool(v))
	}
	return d
}

// SetJSON sets the key to JSON value.
func (d *Data) SetJSON(key string, v interface{}) *Data {
	if v != nil && !reflect.ValueOf(v).IsZero() {
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
	if !isNil(file) {
		d.Files[field] = file
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
			if _, ok := file.(*tg.InputFile); !ok {
				urlid, _ := file.FileData()
				err := mp.WriteField(field, urlid)
				if err != nil {
					w.CloseWithError(err)
					return
				}
				continue
			}

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
