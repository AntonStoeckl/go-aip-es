package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type CustomerRegistered struct {
	customerID   CustomerID
	emailAddress UnconfirmedEmailAddress
	personName   PersonName
	meta         es.EventMeta
}

func BuildCustomerRegistered(
	customerID CustomerID,
	emailAddress UnconfirmedEmailAddress,
	personName PersonName,
	causationID es.MessageID,
	streamVersion uint,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:   customerID,
		emailAddress: emailAddress,
		personName:   personName,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerRegistered(
	customerID string,
	emailAddress string,
	confirmationHash string,
	givenName string,
	familyName string,
	meta es.EventMeta,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:   RebuildCustomerID(customerID),
		emailAddress: RebuildUnconfirmedEmailAddress(emailAddress, confirmationHash),
		personName:   RebuildPersonName(givenName, familyName),
		meta:         meta,
	}

	return event
}

func (event CustomerRegistered) CustomerID() CustomerID {
	return event.customerID
}

func (event CustomerRegistered) EmailAddress() UnconfirmedEmailAddress {
	return event.emailAddress
}

func (event CustomerRegistered) PersonName() PersonName {
	return event.personName
}

func (event CustomerRegistered) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerRegistered) IsFailureEvent() bool {
	return false
}

func (event CustomerRegistered) FailureReason() error {
	return nil
}
