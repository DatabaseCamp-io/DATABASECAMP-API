package application

type App interface {
	Get(path string, handles ...func(Context))
	Post(path string, handles ...func(Context))
	Put(path string, handles ...func(Context))
	Delete(path string, handle ...func(Context))
	Group(path string, handles ...func(Context)) Router
}

type Router interface {
	Get(path string, handles ...func(Context))
	Post(path string, handles ...func(Context))
	Put(path string, handles ...func(Context))
	Delete(path string, handle ...func(Context))
	Group(path string, handles ...func(Context)) Router
}

type Context interface {
	Bind(v interface{}) error
	JSON(statuscode int, v interface{})
	Error(err error) error
	Params(key string, defaultValue ...string) string
	Locals(key string, value ...interface{}) (val interface{})

	Next() error

	GetHeader(key string) string
}
