package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       CustomerID
	confirmationHash ConfirmationHash
	reason           error
	meta             es.EventMeta
}

func BuildCustomerEmailAddressConfirmationFailed(
	customerID CustomerID,
	confirmationHash ConfirmationHash,
	reason error,
	causationID es.MessageID,
	streamVersion uint,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       customerID,
		confirmationHash: confirmationHash,
		reason:           reason,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func (event CustomerEmailAddressConfirmationFailed) CustomerID() CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmationFailed) ConfirmationHash() ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressConfirmationFailed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmationFailed) IsFailureEvent() bool {
	return true
}

func (event CustomerEmailAddressConfirmationFailed) FailureReason() error {
	return event.reason
}
