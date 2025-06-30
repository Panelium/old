package cookies

import (
	"net/http"
	"time"
)

// TODO: might have to url encode the names and tokens

func SetJWTCookie(header http.Header, name string, token string, expiration time.Time) {
	header.Add("Set-Cookie", name+"="+token+"; Path=/; HttpOnly; Secure; SameSite=Strict; Expires="+expiration.UTC().Format(http.TimeFormat))
}

func ClearJWTCookie(header http.Header, name string) {
	header.Add("Set-Cookie", name+"=; Path=/; HttpOnly; Secure; SameSite=Strict; Max-Age=0")
}
