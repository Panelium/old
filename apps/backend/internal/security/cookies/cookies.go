package cookies

import (
	"net/http"
	"time"
)

func SetJWTCookie(header http.Header, name string, token string, expiration time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiration,
	}

	cookieString := cookie.String()
	if cookieString == "" {
		return
	}

	header.Add("Set-Cookie", cookieString)
}

func ClearJWTCookie(header http.Header, name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	}

	cookieString := cookie.String()
	if cookieString == "" {
		return
	}

	header.Add("Set-Cookie", cookieString)
}
