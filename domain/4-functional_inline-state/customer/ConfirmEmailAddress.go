package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func ConfirmEmailAddress(eventStream es.EventStream, command domain.ConfirmCustomerEmailAddress) es.RecordedEvents {
	var customer struct {
		id                   domain.CustomerID
		emailAddress         domain.EmailAddress
		currentStreamVersion uint
	}

	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case domain.CustomerRegistered:
			customer.id = actualEvent.CustomerID()
			customer.emailAddress = actualEvent.EmailAddress()
		case domain.CustomerEmailAddressConfirmed:
			customer.emailAddress = actualEvent.EmailAddress()
		}

		customer.currentStreamVersion = event.Meta().StreamVersion()
	}

	switch actualEmailAddress := customer.emailAddress.(type) {
	case domain.ConfirmedEmailAddress:
		return nil
	case domain.UnconfirmedEmailAddress:
		confirmedEmailAddress, err := domain.ConfirmEmailAddressWithHash(actualEmailAddress, command.ConfirmationHash())

		if err != nil {
			return es.RecordedEvents{
				domain.BuildCustomerEmailAddressConfirmationFailed(
					command.CustomerID(),
					command.ConfirmationHash(),
					err,
					command.MessageID(),
					customer.currentStreamVersion+1,
				),
			}
		}

		return es.RecordedEvents{
			domain.BuildCustomerEmailAddressConfirmed(
				command.CustomerID(),
				confirmedEmailAddress,
				command.MessageID(),
				customer.currentStreamVersion+1,
			),
		}
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("ConfirmEmailAddress(): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}
}
