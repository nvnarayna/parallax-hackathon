package engine

import (
	"context"
	"errors"
)

var ErrTicketClientUnavailable = errors.New("ticket client unavailable")

type Ticket struct {
	Subject     string
	Description string
	Severity    string
	Diagnostic  any
}

type TicketClient interface {
	Create(context.Context, Ticket) error
}

func (e *Engine) CheckForProblem(
	ctx context.Context,
	source string,
	data any,
) error {
	/*
		add problem detection logic here later.

		return nil when everything is normal.

		return e.SubmitTicket(ctx, Ticket{
			Subject:     "...",
			Description: "...",
			Severity:    "...",
			Diagnostic:  data,
		})
	*/

	return nil
}
