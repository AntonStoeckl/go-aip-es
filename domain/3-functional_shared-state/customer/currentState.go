package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type currentState struct {
	id                   domain.CustomerID
	personName           domain.PersonName
	emailAddress         domain.EmailAddress
	currentStreamVersion uint
}

func buildCurrentStateFrom(eventStream es.EventStream) currentState {
	customer := currentState{}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			customer.id = actualEvent.CustomerID()
			customer.personName = actualEvent.PersonName()
			customer.emailAddress = actualEvent.EmailAddress()
		case domain.CustomerEmailAddressConfirmed:
			customer.emailAddress = actualEvent.EmailAddress()
		case domain.CustomerNameChanged:
			customer.personName = actualEvent.PersonName()
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	return customer
}
