package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type CustomerEmailAddressConfirmed struct {
	customerID   CustomerID
	emailAddress ConfirmedEmailAddress
	meta         es.EventMeta
}

func BuildCustomerEmailAddressConfirmed(
	customerID CustomerID,
	emailAddress ConfirmedEmailAddress,
	causationID es.MessageID,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func (event CustomerEmailAddressConfirmed) CustomerID() CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmed) EmailAddress() ConfirmedEmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmed) IsFailureEvent() bool {
	return false
}

func (event CustomerEmailAddressConfirmed) FailureReason() error {
	return nil
}
