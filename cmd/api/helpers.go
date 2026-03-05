package main

import (
	"errors"

	"github.com/osamah22/order_service/internal/services"
)

func isNotFound(err error) bool {
	return errors.Is(err, services.ErrOrderNotFound) ||
		errors.Is(err, services.ErrProductNotFound)
}
