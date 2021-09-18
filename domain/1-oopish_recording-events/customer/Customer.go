package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type Customer struct {
	id                   domain.CustomerID
	personName           domain.PersonName
	emailAddress         domain.EmailAddress
	currentStreamVersion uint
	recordedEvents       es.RecordedEvents
}

func ReconstituteCustomer(eventStream es.EventStream) *Customer {
	customer := &Customer{}
	customer.when(eventStream...)

	return customer
}

func Register(command domain.RegisterCustomer) *Customer {
	customer := &Customer{}

	registered := domain.BuildCustomerRegistered(
		command.CustomerID(),
		command.EmailAddress(),
		command.PersonName(),
		command.MessageID(),
		1,
	)

	customer.recordThat(registered)

	return customer
}

func (customer *Customer) ConfirmEmailAddress(command domain.ConfirmCustomerEmailAddress) {
	switch actualEmailAddress := customer.emailAddress.(type) {
	case domain.ConfirmedEmailAddress:
		// already confirmed
	case domain.UnconfirmedEmailAddress:
		confirmedEmailAddress, err := domain.ConfirmEmailAddressWithHash(actualEmailAddress, command.ConfirmationHash())

		if err != nil {
			emailAddressConfirmationFailed := domain.BuildCustomerEmailAddressConfirmationFailed(
				command.CustomerID(),
				command.ConfirmationHash(),
				err,
				command.MessageID(),
				customer.currentStreamVersion+1,
			)

			customer.recordThat(emailAddressConfirmationFailed)
		} else {
			emailAddressConfirmed := domain.BuildCustomerEmailAddressConfirmed(
				command.CustomerID(),
				confirmedEmailAddress,
				command.MessageID(),
				customer.currentStreamVersion+1,
			)

			customer.recordThat(emailAddressConfirmed)
		}
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("ConfirmEmailAddress(): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}
}

func (customer *Customer) ChangeName(command domain.ChangeCustomerName) {
	if !customer.personName.Equals(command.PersonName()) {
		nameChanged := domain.BuildCustomerNameChanged(
			command.CustomerID(),
			command.PersonName(),
			command.MessageID(),
			customer.currentStreamVersion+1,
		)

		customer.recordThat(nameChanged)
	}
}

func (customer *Customer) RecordedEvents() es.RecordedEvents {
	var events es.RecordedEvents

	events = append(events, customer.recordedEvents...)

	customer.recordedEvents = es.RecordedEvents{}

	return events
}

func (customer *Customer) when(events ...es.DomainEvent) {
	for _, event := range events {
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
}

func (customer *Customer) recordThat(event es.DomainEvent) {
	customer.recordedEvents = append(customer.recordedEvents, event)
	customer.when(event)
}
