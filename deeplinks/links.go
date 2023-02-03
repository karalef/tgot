package deeplinks

import "net/url"

// String creates a Telegram deeplink.
func String(path string, query url.Values) string {
	return New(path, query).String()
}

// TG returns a deeplink as a Telegram URL.
func TG(path string, query url.Values, slashes ...bool) string {
	return New(path, query).TG(len(slashes) > 0 && slashes[0])
}

// Resolve return 'tg://resolve?query' deeplink.
func Resolve(query url.Values) string {
	return TG("resolve", query, true)
}

// HTTP returns a deeplink as an HTTP URL.
func HTTP(path string, query url.Values, scheme ...bool) string {
	return New(path, query).HTTP(len(scheme) > 0 && scheme[0])
}

// HTTPS returns a deeplink as an HTTPS URL.
func HTTPS(path string, query url.Values) string {
	return New(path, query).HTTP(true)
}

// Username returns 'username.t.me' deeplink.
func Username(username string) string {
	return New("", nil).CustomHTTP(false, username+"."+HostTMe)
}
