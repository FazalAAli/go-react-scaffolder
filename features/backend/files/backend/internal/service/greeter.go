package service

import (
	"context"

	appv1 "backend/gen/go/app/v1"
	"backend/gen/go/app/v1/appv1connect"
	"backend/internal/app"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v4"
)

type Greeter struct {
	app *app.App
}

func NewGreeter(a *app.App) *Greeter {
	return &Greeter{app: a}
}

func (g *Greeter) Greet(
	ctx context.Context,
	req *connect.Request[appv1.GreetRequest],
) (*connect.Response[appv1.GreetResponse], error) {
	return connect.NewResponse(&appv1.GreetResponse{
		Greeting: "Hello, " + req.Msg.Name + "!",
	}), nil
}

func (g *Greeter) Mount(e *echo.Echo) {
	path, handler := appv1connect.NewGreeterServiceHandler(g)
	e.Any(path+"*", echo.WrapHandler(handler))
}
