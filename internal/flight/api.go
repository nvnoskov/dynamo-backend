package flight

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/nvnoskov/dynamo-backend/internal/errors"
	"github.com/nvnoskov/dynamo-backend/pkg/log"
	"github.com/nvnoskov/dynamo-backend/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	// the following endpoints require a valid JWT
	r.Use(authHandler)
	r.Get("/flights/<id>", res.get)
	r.Get("/flights", res.query)
	r.Post("/flights", res.create)
	r.Put("/flights/<id>", res.update)
	r.Delete("/flights/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	flight, err := r.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(flight)
}

func (r resource) query(c *routing.Context) error {

	input := SearchFlightRequest{
		Name:          c.Query("name"),
		Departure:     c.Query("departure"),
		Destination:   c.Query("destination"),
		DepartureTime: c.Query("departure_time"),
	}

	ctx := c.Request.Context()
	count, err := r.service.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	flights, err := r.service.Query(ctx, input, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = flights
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateFlightRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	flight, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(flight, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateFlightRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	flight, err := r.service.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(flight)
}

func (r resource) delete(c *routing.Context) error {
	flight, err := r.service.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(flight)
}
