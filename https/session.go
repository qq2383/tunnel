package https

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	session    = map[string]*Session{}
	cookieName = "SESSIONID"
	// expires    = time.Hour
	path       = "/"
)

type Session struct {
	data   map[string]any
}

func NerSession(w http.ResponseWriter, r *http.Request) *Session {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		cookie = &http.Cookie{}
		cookie.Name = cookieName
		hosts := strings.Split(r.Host, ":")
		cookie.Domain = hosts[0]
		cookie.Path = path
		// cookie.MaxAge = 10
		// cookie.Expires = time.Now().UTC().Add(expires)

		u := uuid.New()
		cookie.Value = u.String()
		// w.Header().Set("cookie", cookie.String())
		http.SetCookie(w, cookie)
		// fmt.Println("setCookie")
	}
	s, ok := session[cookie.Value]
	if !ok {
		s = &Session{data: map[string]any{}}
		session[cookie.Value] = s
	}
	return s
}

func (s *Session) GetAttr(key string) (any, error) {
	if s.data == nil {
		return nil, fmt.Errorf("session failed")
	}
	v, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("session failed")
	}
	return v, nil
}

func (s *Session) SetAttr(key string, value any) {
	s.data[key] = value
}
