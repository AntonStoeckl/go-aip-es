package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
)

func Register(command domain.RegisterCustomer) domain.CustomerRegistered {
	event := domain.BuildCustomerRegistered(
		command.CustomerID(),
		command.EmailAddress(),
		command.PersonName(),
		command.MessageID(),
		1,
	)

	return event
}
