package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type RegisterCustomer struct {
	customerID   CustomerID
	emailAddress UnconfirmedEmailAddress
	personName   PersonName
	messageID    es.MessageID
}

func BuildRegisterCustomer(
	customerID CustomerID,
	emailAddress UnconfirmedEmailAddress,
	personName PersonName,
) RegisterCustomer {

	command := RegisterCustomer{
		customerID:   customerID,
		emailAddress: emailAddress,
		personName:   personName,
		messageID:    es.GenerateMessageID(),
	}

	return command
}

func (command RegisterCustomer) CustomerID() CustomerID {
	return command.customerID
}

func (command RegisterCustomer) EmailAddress() UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command RegisterCustomer) PersonName() PersonName {
	return command.personName
}

func (command RegisterCustomer) MessageID() es.MessageID {
	return command.messageID
}
