package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type CustomerNameChanged struct {
	customerID CustomerID
	personName PersonName
	meta       es.EventMeta
}

func BuildCustomerNameChanged(
	customerID CustomerID,
	personName PersonName,
	causationID es.MessageID,
	streamVersion uint,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: customerID,
		personName: personName,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerNameChanged(
	customerID string,
	givenName string,
	familyName string,
	meta es.EventMeta,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: RebuildCustomerID(customerID),
		personName: RebuildPersonName(givenName, familyName),
		meta:       meta,
	}

	return event
}

func (event CustomerNameChanged) CustomerID() CustomerID {
	return event.customerID
}

func (event CustomerNameChanged) PersonName() PersonName {
	return event.personName
}

func (event CustomerNameChanged) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerNameChanged) IsFailureEvent() bool {
	return false
}

func (event CustomerNameChanged) FailureReason() error {
	return nil
}
