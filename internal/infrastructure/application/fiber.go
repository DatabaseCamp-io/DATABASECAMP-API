package application

import (
	"database-camp/internal/errs"
	"database-camp/internal/logs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func toFiberHandle(handle func(Context)) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		handle(NewFiberCtx(c))
		return nil
	}
}

func toFiberHandles(handles []func(Context)) []func(c *fiber.Ctx) error {
	fiberHandles := make([]func(c *fiber.Ctx) error, len(handles), cap(handles))
	for index := range handles {
		fiberHandles[index] = toFiberHandle(handles[index])
	}
	return fiberHandles
}

type FiberApp struct {
	*fiber.App
}

func NewFiberApp() *FiberApp {
	r := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Database Camp",
	})

	r.Use(cors.New())
	r.Use(logger.New())
	r.Use(recover.New())

	return &FiberApp{r}
}

func (r *FiberApp) Post(path string, handles ...func(Context)) {
	r.App.Post(path, toFiberHandles(handles)...)
}

func (r *FiberApp) Get(path string, handles ...func(Context)) {
	r.App.Get(path, toFiberHandles(handles)...)
}

func (r *FiberApp) Put(path string, handles ...func(Context)) {
	r.App.Put(path, toFiberHandles(handles)...)
}

func (r *FiberApp) Delete(path string, handles ...func(Context)) {
	r.App.Delete(path, toFiberHandles(handles)...)
}

func (r *FiberApp) Group(path string, handles ...func(Context)) Router {
	router := r.App.Group(path, toFiberHandles(handles)...)
	return NewFiberRouter(router)
}

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{Ctx: c}
}

func (c *FiberCtx) Next() error {
	return c.Ctx.Next()
}

func (c *FiberCtx) GetHeader(key string) string {
	return c.Ctx.Get(key)
}

func (c *FiberCtx) Locals(key string, value ...interface{}) (val interface{}) {
	return c.Ctx.Locals(key, value...)
}

func (c *FiberCtx) Params(key string, defaultValue ...string) string {
	return c.Ctx.Params(key, defaultValue...)
}

func (c *FiberCtx) Bind(v interface{}) error {
	err := c.BodyParser(v)
	if err != nil {
		logs.GetInstance().Error(err)
		return errs.ErrBadRequestError
	}
	return nil
}

func (c *FiberCtx) JSON(statuscode int, v interface{}) {
	c.Ctx.Status(statuscode).JSON(v)
}

func (c *FiberCtx) Error(err error) error {

	type message struct {
		Th string `json:"th_message"`
		En string `json:"en_message"`
	}

	switch e := err.(type) {
	case errs.AppError:
		return c.Status(e.Code).JSON(message{Th: e.ThMessage, En: e.EnMessage})
	case error:
		return c.Status(fiber.StatusInternalServerError).JSON(message{
			Th: errs.INTERNAL_SERVER_ERROR_TH,
			En: errs.INTERNAL_SERVER_ERROR_EN,
		})
	}
	return nil
}

type FiberRouter struct {
	fiber.Router
}

func NewFiberRouter(r fiber.Router) FiberRouter {
	return FiberRouter{Router: r}
}

func (r FiberRouter) Post(path string, handles ...func(Context)) {
	r.Router.Post(path, toFiberHandles(handles)...)
}

func (r FiberRouter) Get(path string, handles ...func(Context)) {
	r.Router.Get(path, toFiberHandles(handles)...)
}

func (r FiberRouter) Put(path string, handles ...func(Context)) {
	r.Router.Put(path, toFiberHandles(handles)...)
}

func (r FiberRouter) Delete(path string, handles ...func(Context)) {
	r.Router.Delete(path, toFiberHandles(handles)...)
}

func (r FiberRouter) Group(path string, handles ...func(Context)) Router {
	r.Router = r.Router.Group(path, toFiberHandles(handles)...)
	return r
}
