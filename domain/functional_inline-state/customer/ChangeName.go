package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func ChangeName(eventStream es.EventStream, command domain.ChangeCustomerName) (es.RecordedEvents, error) {
	var customer struct {
		id                   domain.CustomerID
		personName           domain.PersonName
		currentStreamVersion uint
	}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			customer.id = actualEvent.CustomerID()
			customer.personName = actualEvent.PersonName()
		case domain.CustomerNameChanged:
			customer.personName = actualEvent.PersonName()
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	if customer.personName.Equals(command.PersonName()) {
		return nil, nil
	}

	event := domain.BuildCustomerNameChanged(
		command.CustomerID(),
		command.PersonName(),
		command.MessageID(),
		customer.currentStreamVersion+1,
	)

	return es.RecordedEvents{event}, nil
}
