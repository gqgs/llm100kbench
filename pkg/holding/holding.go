package holding

import (
	"strings"
	"time"

	"github.com/gqgs/llminvestbench/pkg/order"
)

type Holdings []*Holding

type Holding struct {
	Ticket    string    `json:"ticket"`
	Sum       float64   `json:"sum"`
	Quantity  int       `json:"quantity"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (h *Holding) Update(update *order.Update) {
	switch strings.ToUpper(update.Action) {
	case "BUY":
		h.Quantity += update.Quantity
		h.Sum += update.Price * float64(update.Quantity)
	case "SELL":
		h.Quantity -= update.Quantity
		h.Sum -= update.Price * float64(update.Quantity)
	}
}

func (h *Holdings) ProcessOrder(order *order.Order) error {
Outer:
	for _, update := range order.Updates {
		for _, holding := range *h {
			if holding.Ticket == update.Ticket {
				holding.Update(update)
				continue Outer
			}

		}
		(*h) = append((*h), &Holding{
			Ticket:    update.Ticket,
			Quantity:  update.Quantity,
			Sum:       update.Price * float64(update.Quantity),
			UpdatedAt: time.Now(),
			CreatedAt: time.Now(),
		})
	}

	return nil
}
