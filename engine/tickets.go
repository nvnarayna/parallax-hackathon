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
