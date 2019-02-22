package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

type (
	Store interface {
		sessions.Store
		Options(options Options)
	}

	Session interface {
		// Get returns the session value associated to the given key.
		Get(key interface{}) interface{}
		// GetAll
		GetAll() map[interface{}]interface{}
		// Set sets the session value associated to the given key.
		Set(key interface{}, val interface{})
		// Delete removes the session value associated to the given key.
		Delete(key interface{})
		// Clear deletes all values in the session.
		Clear()
		// AddFlash adds a flash message to the session.
		// A single variadic argument is accepted, and it is optional: it defines the flash key.
		// If not defined "_flash" is used by default.
		AddFlash(value interface{}, vars ...string)
		// Flashes returns a slice of flash messages from the session.
		// A single variadic argument is accepted, and it is optional: it defines the flash key.
		// If not defined "_flash" is used by default.
		Flashes(vars ...string) []interface{}
		// Options sets confuguration for a session.
		Options(Options)
		// Save saves all sessions used during the current request.
		Save() error
	}

	session struct {
		name    string
		request *http.Request
		store   Store
		session *sessions.Session
		written bool
		writer  http.ResponseWriter
	}

	Options struct {
		Path   string
		Domain string
		// MaxAge=0 means no 'Max-Age' attribute specified.
		// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
		// MaxAge>0 means Max-Age attribute present and given in seconds.
		MaxAge   int
		Secure   bool
		HttpOnly bool
	}
)

var (
	SESSION_STORE_KEY = "_session_store"
	errorFormat       = "[session] ERROR! %s\n"
)

func NewSession(name string, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess := &session{
				name,
				c.Request(),
				store,
				nil,
				false,
				c.Response().Writer,
			}
			c.Set(SESSION_STORE_KEY, sess)
			return next(c)
		}
	}
}

func Get(c echo.Context) Session {
	return c.Get(SESSION_STORE_KEY).(Session)
}

func (this *session) Get(key interface{}) interface{} {
	return this.Session().Values[key]
}

func (this *session) GetAll() map[interface{}]interface{} {
	return this.Session().Values
}

func (this *session) Set(key interface{}, val interface{}) {
	this.Session().Values[key] = val
	this.written = true
}

func (this *session) Delete(key interface{}) {
	delete(this.Session().Values, key)
	this.written = true
}

func (this *session) Clear() {
	for key := range this.Session().Values {
		this.Delete(key)
	}
}

func (this *session) AddFlash(value interface{}, vars ...string) {
	this.Session().AddFlash(value, vars...)
	this.written = true
}

func (this *session) Flashes(vars ...string) []interface{} {
	this.written = true
	return this.Session().Flashes(vars...)
}

func (this *session) Options(options Options) {
	this.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}

func (this *session) Save() error {
	if this.Written() {
		e := this.Session().Save(this.request, this.writer)
		if e == nil {
			this.written = false
		}
		return e
	}
	return nil
}

func (this *session) Session() *sessions.Session {
	if this.session == nil {
		var err error
		this.session, err = this.store.Get(this.request, this.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return this.session
}

func (this *session) Written() bool {
	return this.written
}
