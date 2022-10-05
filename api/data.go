package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/karalef/tgot/tg"
)

// Data contains query parameters and files data.
type Data struct {
	Params url.Values
	Files  []File
}

func (d Data) data() (string, io.Reader) {
	for i := range d.Files {
		if _, r := d.Files[i].FileData(); r != nil {
			return writeMultipart(d)
		}
	}
	for i := range d.Files {
		urlid, _ := d.Files[i].FileData()
		d.Set(d.Files[i].Field, urlid)
	}
	return "application/x-www-form-urlencoded", strings.NewReader(d.Params.Encode())
}

func (d Data) forEach(f func(k, v string) error) error {
	for k, v := range d.Params {
		if err := f(k, v[0]); err != nil {
			return err
		}
	}
	return nil
}

// Set sets the key to value.
func (d Data) Set(k, v string, force ...bool) Data {
	if v != "" {
		d.Params.Set(k, v)
	}
	return d
}

// SetInt sets the key to int value.
func (d Data) SetInt(key string, v int, force ...bool) Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.Itoa(v))
	}
	return d
}

// SetInt64 sets the key to int64 value.
func (d Data) SetInt64(key string, v int64, force ...bool) Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatInt(v, 10))
	}
	return d
}

// SetFloat sets the key to float value.
func (d Data) SetFloat(key string, v float32, force ...bool) Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatFloat(float64(v), 'f', 6, 32))
	}
	return d
}

// SetBool sets the key to bool value.
func (d Data) SetBool(key string, v bool) Data {
	if v {
		d.Set(key, strconv.FormatBool(v))
	}
	return d
}

// SetJSON sets the key to JSON value.
func (d Data) SetJSON(key string, v interface{}) Data {
	if v != nil && !reflect.ValueOf(v).IsZero() {
		b, _ := json.Marshal(v)
		d.Set(key, string(b))
	}
	return d
}

// File contains the file data with field.
type File struct {
	// contains query key or multipart field.
	Field string
	tg.FileSignature
}

func writeMultipart(d Data) (string, io.Reader) {
	r, w := io.Pipe()
	mp := multipart.NewWriter(w)
	go func() {
		defer func() {
			w.CloseWithError(mp.Close())
		}()

		err := d.forEach(mp.WriteField)
		if err != nil {
			w.CloseWithError(err)
			return
		}

		for _, file := range d.Files {
			if _, ok := file.FileSignature.(*tg.InputFile); !ok {
				urlid, _ := file.FileData()
				err := mp.WriteField(file.Field, urlid)
				if err != nil {
					w.CloseWithError(err)
					return
				}
				continue
			}

			name, reader := file.FileData()
			part, err := mp.CreateFormFile(file.Field, name)
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
