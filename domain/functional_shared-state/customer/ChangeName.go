package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func ChangeName(eventStream es.EventStream, command domain.ChangeCustomerName) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

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
