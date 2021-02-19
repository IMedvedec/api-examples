package builtin

import "net/http"

const (
	cookieMaxAge int = 1000
)

func domainCookieSetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		domainCookie := http.Cookie{
			Name:   "domain-cookie",
			Value:  "domain-cookie-value",
			Domain: "localhost",
			MaxAge: cookieMaxAge,
		}
		http.SetCookie(w, &domainCookie)

		next.ServeHTTP(w, r)
	})
}

func domainAndPathCookieSetHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:     "doman-path-cookie",
			Value:    "domain-path-cookie-value",
			Domain:   "localhost",
			Path:     "cookie/json",
			MaxAge:   cookieMaxAge,
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)

		next.ServeHTTP(w, r)
	})
}
