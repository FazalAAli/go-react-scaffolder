package stripe

import (
	"context"
	"io"
	"net/http"

	appv1 "backend/gen/go/app/v1"
	"backend/gen/go/app/v1/appv1connect"
	"backend/internal/config"

	"connectrpc.com/connect"
	"github.com/labstack/echo/v5"
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/webhook"
)

type Client struct {
	cfg *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	stripe.Key = cfg.StripeSecretKey
	return &Client{cfg: cfg}, nil
}

func (c *Client) Mount(e *echo.Echo) {
	path, handler := appv1connect.NewStripeServiceHandler(c)
	e.Any(path+"*", echo.WrapHandler(handler))
	e.POST("/webhooks/stripe", c.handleWebhook)
}

func (c *Client) CreateCheckoutSession(
	ctx context.Context,
	req *connect.Request[appv1.CreateCheckoutSessionRequest],
) (*connect.Response[appv1.CreateCheckoutSessionResponse], error) {
	params := &stripe.CheckoutSessionParams{
		// TODO: set Mode, LineItems (req.Msg.PriceId), and any customer linkage
		// for your product. This is app-domain configuration.
		SuccessURL: stripe.String(c.cfg.FrontendBaseURL + "/checkout/success"),
		CancelURL:  stripe.String(c.cfg.FrontendBaseURL + "/checkout/cancel"),
	}
	s, err := session.New(params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&appv1.CreateCheckoutSessionResponse{Url: s.URL}), nil
}

func (c *Client) handleWebhook(ctx echo.Context) error {
	payload, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.NoContent(http.StatusServiceUnavailable)
	}
	event, err := webhook.ConstructEvent(
		payload,
		ctx.Request().Header.Get("Stripe-Signature"),
		c.cfg.StripeWebhookSecret,
	)
	if err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	switch event.Type {
	case "checkout.session.completed":
		// TODO: fulfil the order. Persist via your database feature (a.DB) by
		// threading a store into stripe.New; see the stripe feature README note.
	case "customer.subscription.updated", "customer.subscription.deleted":
		// TODO: update the customer's subscription state in your database.
	default:
	}
	return ctx.NoContent(http.StatusOK)
}
