package customer

import (
	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func ConfirmEmailAddress(eventStream es.EventStream, command domain.ConfirmCustomerEmailAddress) (es.RecordedEvents, error) {
	customer := buildCurrentStateFrom(eventStream)

	switch actualEmailAddress := customer.emailAddress.(type) {
	case domain.ConfirmedEmailAddress:
		return nil, nil
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
			}, nil
		}

		return es.RecordedEvents{
			domain.BuildCustomerEmailAddressConfirmed(
				command.CustomerID(),
				confirmedEmailAddress,
				command.MessageID(),
				customer.currentStreamVersion+1,
			),
		}, nil
	default:
		// until Go has "union types" we need to use an interface and this case could exist - we don't want to hide it
		panic("ConfirmEmailAddress(): emailAddress is neither UnconfirmedEmailAddress nor ConfirmedEmailAddress")
	}
}
