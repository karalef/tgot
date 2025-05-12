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
		Upload: make(map[string]*tg.InputFile),
	}
}

// NewDataFrom creates new Data object from v.
func NewDataFrom(v any, structTagName ...string) *Data {
	return NewData().AddObject(v, structTagName...)
}

// Data contains query parameters and files data.
type Data struct {
	// contains query parameters and files that does not need to be uploaded.
	Params map[string]string

	// contains files that will be uploaded, where key is a multipart field.
	Upload map[string]*tg.InputFile

	attachCounter int
}

// Copy copies Data's params.
func (d *Data) Copy() *Data {
	cp := NewData()
	for k, v := range d.Params {
		cp.Params[k] = v
	}
	return cp
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
	if d == nil || len(d.Params) == 0 && len(d.Upload) == 0 {
		return "", nil
	}
	if len(d.Upload) > 0 {
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

// SetFile sets the key to file.
func (d *Data) SetFile(key string, file tg.Inputtable) *Data {
	if isNil(file) {
		return d
	}

	if inp, ok := file.(*tg.InputFile); ok {
		return d.addFile(key, inp)
	}

	urlid, _ := file.FileData()
	d.Params[key] = urlid
	return d
}

func (d *Data) addFile(field string, file *tg.InputFile) *Data {
	d.Upload[field] = file
	return d
}

// SetInput sets the key to JSON value and adds the attachments if needed.
func (d *Data) SetInput(key string, v tg.Inputter) *Data {
	for _, inp := range v.GetInput() {
		d.AddAttach(inp)
	}
	return d.SetJSON(key, v)
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

		for field, file := range d.Upload {
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

// AddObject sets the struct fields as parameters.
// It uses the "tg" tag key with pattern `tg:"name,opt1,opt2,..."`.
// The name can be omitted (`tg:",opt1"`: field name will be automatically
// transformed to snake case) or replaced with '-' to ignore it.
// The 'force' option acts the same as argument in Data.Set* methods.
// The 'json' option allows to use structs with 'json' tags from tg package as
// anonymous fields, which will be embedded as the same struct with 'tg' tags.
// Supported struct field types: string, bool, int*, uint*, float*,
// struct/interface and array/slice of them all (automatically detects types
// which implement tg.Inputter or tg.Inputtable and adds attachments if needed).
func (d *Data) AddObject(o any, structTagName ...string) *Data {
	if c, ok := o.(Marshaler); ok {
		c.MarshalTg(d)
		return d
	}

	structTag := structTagTg
	if len(structTagName) > 0 {
		structTag = structTagName[0]
	}

	val := reflect.ValueOf(o)
	if !val.IsValid() {
		return d
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
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		tag, ignore := parseTag(field.Tag.Get(structTag))
		if ignore {
			continue
		}
		if field.Anonymous && tag.name == "" {
			anon := structTag
			if tag.has("json") {
				anon = "json"
			}
			d.AddObject(value.Interface(), anon)
			continue
		}
		if tag.name == "" {
			tag.name = camelToSnake(field.Name)
		}
		force := tag.has("force")
		typ := value.Type()
		switch typ.Kind() {
		case reflect.String:
			d.Set(tag.name, value.String(), force)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			d.SetInt64(tag.name, value.Int(), force)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			d.SetUint64(tag.name, value.Uint(), force)
		case reflect.Float32, reflect.Float64:
			d.SetFloat64(tag.name, value.Float(), force)
		case reflect.Bool:
			d.SetBool(tag.name, value.Bool(), force)
		case reflect.Struct, reflect.Interface:
			marshalType(d, typ, tag.name, value.Interface())
		case reflect.Slice, reflect.Array:
			marshalTypeArray(d, typ, tag.name, value)
		default:
			panic("unsupported type " + field.Type.String())
		}
	}
	return d
}

func marshalType(dst *Data, typ reflect.Type, name string, val any) {
	if typ.Implements(inputtableType) {
		dst.SetFile(name, val.(tg.Inputtable))
	} else if typ.Implements(inputterType) {
		dst.SetInput(name, val.(tg.Inputter))
	} else {
		dst.SetJSON(name, val)
	}
}

func marshalTypeArray(dst *Data, typ reflect.Type, name string, val reflect.Value) {
	if typ.Elem().Implements(inputtableType) {
		for i := 0; i < val.Len(); i++ {
			dst.AddAttach(val.Index(i).Interface().(tg.Inputtable))
		}
	} else if typ.Elem().Implements(inputterType) {
		for i := 0; i < val.Len(); i++ {
			inp := val.Index(i).Interface().(tg.Inputter).GetInput()
			for _, i := range inp {
				dst.AddAttach(i)
			}
		}
	}
	dst.SetJSON(name, val.Interface())
}

func camelToSnake(s string) string {
	result := make([]rune, 0, len(s)+3)
	for i, v := range s {
		if unicode.IsUpper(v) && i != 0 {
			result = append(result, '_', unicode.ToLower(v))
			continue
		}
		result = append(result, v)
	}
	return string(result)
}

var (
	inputterType   = reflect.TypeOf((*tg.Inputter)(nil)).Elem()
	inputtableType = reflect.TypeOf((*tg.Inputtable)(nil)).Elem()
)

func parseTag(t string) (s structtag, i bool) {
	tag := strings.Split(t, ",")
	if len(tag) > 0 {
		s.name = tag[0]
	}
	if s.name == "-" || s.name == "_" {
		s.name = ""
		i = true
	}
	if len(tag) > 0 {
		s.opts = tag[1:]
	}
	return
}

type structtag struct {
	name string
	opts []string
}

func (t structtag) has(opt string) bool {
	for i := range t.opts {
		if t.opts[i] == opt {
			return true
		}
	}
	return false
}
