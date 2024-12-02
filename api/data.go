package api

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"unicode"
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

// NewDataFrom creates new Data object from v.
func NewDataFrom(v any, structTagName ...string) *Data {
	d := NewData()
	MarshalTo(d, v, structTagName...)
	return d
}

// Data contains query parameters and files data.
type Data struct {
	// contains query parameters and files that does not need to be uploaded.
	Params map[string]string

	// key is a multipart field name.
	Files map[string]*tg.InputFile

	attachCounter int
}

// Copy copies Data's params.
func (d *Data) Copy() *Data {
	p := d.Params
	d = NewData()
	for k, v := range p {
		d.Params[k] = v
	}
	return d
}

// WriteTo copies Data's params to dst and returns dst.
func (d *Data) WriteTo(dst *Data) *Data {
	if d == nil || dst == nil {
		return dst
	}
	for k, v := range d.Params {
		dst.Params[k] = v
	}
	return dst
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

// SetUint sets the key to uint value.
func (d *Data) SetUint(key string, v uint, force ...bool) *Data {
	return d.SetUint64(key, uint64(v), force...)
}

// SetUint64 sets the key to uint64 value.
func (d *Data) SetUint64(key string, v uint64, force ...bool) *Data {
	if v != 0 || len(force) > 0 && force[0] {
		d.Set(key, strconv.FormatUint(v, 10))
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

// SetObject sets the any marshalable struct as key value.
func (d *Data) SetObject(v any, structTagName ...string) *Data {
	MarshalTo(d, v, structTagName...)
	return d
}

// AddFile adds file.
func (d *Data) AddFile(field string, file tg.Inputtable) *Data {
	if isNil(file) {
		return d
	}

	if inp, ok := file.(*tg.InputFile); ok {
		return d.addFile(field, inp)
	}

	urlid, _ := file.FileData()
	d.Params[field] = urlid
	return d
}

func (d *Data) addFile(field string, file *tg.InputFile) *Data {
	d.Files[field] = file
	return d
}

// AddAttach links a file to the multipart field and adds it.
// If the file is not *tg.InputFile, it does nothing.
func (d *Data) AddAttach(file tg.Inputtable) *Data {
	f, ok := file.(*tg.InputFile)
	if !ok || f == nil {
		return d
	}
	d.attachCounter++
	field := "file-" + strconv.Itoa(d.attachCounter)
	return d.addFile(field, f.AsAttachment(field))
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

const structTagTg = "tg"

// Marshaler interface represents object that marshals on his own.
type Marshaler interface {
	MarshalTg(dst *Data)
}

// MarshalTo marshals the object to dst.
func MarshalTo(dst *Data, o any, structTagName ...string) {
	if c, ok := o.(Marshaler); ok {
		c.MarshalTg(dst)
		return
	}

	structTag := structTagTg
	if len(structTagName) > 0 {
		structTag = structTagName[0]
	}

	val := reflect.ValueOf(o)
	if !val.IsValid() {
		return
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		panic("not a struct")
	}

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		if !value.IsValid() {
			continue
		}
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		isPointer := value.Kind() == reflect.Ptr
		if isPointer {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}
		tag := strings.Split(field.Tag.Get(structTag), ",")
		var name string
		if len(tag) > 0 {
			name = tag[0]
		}
		if name == "-" || name == "_" {
			continue
		}

		if field.Anonymous && name == "" {
			useJSONTags := len(tag) > 1 && tag[1] == "json"
			if useJSONTags {
				MarshalTo(dst, value.Interface(), "json")
			} else {
				MarshalTo(dst, value.Interface(), structTag)
			}
			continue
		}
		if name == "" {
			name = camelToSnake(field.Name)
		}
		force := isPointer || (len(tag) > 1 && tag[1] == "force")
		typ := value.Type()
		switch typ.Kind() {
		case reflect.String:
			dst.Set(name, value.String(), force)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst.SetInt64(name, value.Int(), force)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dst.SetInt64(name, int64(value.Uint()), force)
		case reflect.Float32, reflect.Float64:
			dst.SetFloat64(name, value.Float(), force)
		case reflect.Bool:
			dst.SetBool(name, value.Bool(), force)
		case reflect.Struct:
			dst.SetJSON(name, value.Interface())
		case reflect.Interface:
			if typ.Implements(inputtableType) {
				dst.AddFile(name, value.Interface().(tg.Inputtable))
				break
			} else if typ.Implements(inputterType) {
				prepareMedia(dst, value.Interface().(tg.Inputter))
			}
			dst.SetJSON(name, value.Interface())
		case reflect.Slice, reflect.Array:
			if typ.Elem().Implements(inputtableType) {
				for i := 0; i < value.Len(); i++ {
					dst.AddFile(name, value.Index(i).Interface().(tg.Inputtable))
				}
				break
			} else if typ.Elem().Implements(inputterType) {
				for i := 0; i < value.Len(); i++ {
					prepareMedia(dst, value.Index(i).Interface().(tg.Inputter))
				}
			}
			dst.SetJSON(name, value.Interface())
		default:
			panic("unsupported type " + field.Type.String())
		}
	}
}

func camelToSnake(s string) string {
	var result []rune
	for i, v := range s {
		if unicode.IsUpper(v) && i != 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(v))
	}
	return string(result)
}

var inputterType = reflect.TypeOf((*tg.Inputter)(nil)).Elem()
var inputtableType = reflect.TypeOf((*tg.Inputtable)(nil)).Elem()

func prepareMedia(d *Data, media tg.Inputter) {
	med, thumb := media.GetMedia()
	if med == nil {
		return
	}
	d.AddAttach(med)
	d.AddAttach(thumb)
}
