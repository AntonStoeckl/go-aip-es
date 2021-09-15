package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type ChangeCustomerName struct {
	customerID CustomerID
	personName PersonName
	messageID  es.MessageID
}

func BuildChangeCustomerName(
	customerID CustomerID,
	personName PersonName,
) ChangeCustomerName {

	command := ChangeCustomerName{
		customerID: customerID,
		personName: personName,
		messageID:  es.GenerateMessageID(),
	}

	return command
}

func (command ChangeCustomerName) CustomerID() CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() PersonName {
	return command.personName
}

func (command ChangeCustomerName) MessageID() es.MessageID {
	return command.messageID
}
