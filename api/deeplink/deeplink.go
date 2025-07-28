package deeplink

import "net/url"

// Scheme constants.
const (
	SchemeTG    = "tg"
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
	HostTMe     = "t.me"
)

// New returns a new Deeplink.
func New(path string, query url.Values) Deeplink {
	return Deeplink{
		Path:  path,
		Query: query,
	}
}

// Deeplink represents a Telegram deeplink.
type Deeplink struct {
	Path  string
	Query url.Values
}

// String is a alias for Deeplink.TG(true).
func (d Deeplink) String() string {
	return d.TG(true)
}

func (d Deeplink) build(scheme, host string, slashes bool) string {
	u := url.URL{
		Scheme:   scheme,
		RawQuery: d.Query.Encode(),
	}
	if slashes {
		u.Host = host
		u.Path = d.Path
	} else {
		if host != "" && d.Path != "" && d.Path[0] != '/' {
			host = host + "/"
		}
		u.Opaque = host + d.Path
	}
	return u.String()
}

// TG returns a deeplink as a Telegram URL.
//
// slashes=false: tg:path?query
//
// slashes=true: tg://path/?query
func (d Deeplink) TG(slashes bool) string {
	return d.build(SchemeTG, "", slashes)
}

// HTTP returns a deeplink as an HTTP URL.
//
// scheme=false: t.me/path?query
//
// scheme=true: https://t.me/path?query
//
// scheme=true, http=true: http://t.me/path?query
func (d Deeplink) HTTP(scheme bool, http ...bool) string {
	return d.CustomHTTP(scheme, HostTMe, http...)
}

// CustomHTTP returns a deeplink as an HTTP URL with custom domain.
func (d Deeplink) CustomHTTP(scheme bool, host string, http ...bool) string {
	if !scheme {
		return d.build("", host, false)
	}
	if len(http) > 0 && http[0] {
		return d.build(SchemeHTTP, host, true)
	}
	return d.build(SchemeHTTPS, host, true)
}
