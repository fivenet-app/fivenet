package utils

import "github.com/fivenet-app/fivenet/pkg/utils/http"

func ParseCookies(cookie string) ([]*http.Cookie, error) {
	return http.ParseCookie(cookie)
}
