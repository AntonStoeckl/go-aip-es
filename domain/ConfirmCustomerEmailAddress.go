package domain

import (
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

type ConfirmCustomerEmailAddress struct {
	customerID       CustomerID
	confirmationHash ConfirmationHash
	messageID        es.MessageID
}

func BuildConfirmCustomerEmailAddress(
	customerID CustomerID,
	confirmationHash ConfirmationHash,
) ConfirmCustomerEmailAddress {

	command := ConfirmCustomerEmailAddress{
		customerID:       customerID,
		confirmationHash: confirmationHash,
		messageID:        es.GenerateMessageID(),
	}

	return command
}

func (command ConfirmCustomerEmailAddress) CustomerID() CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() ConfirmationHash {
	return command.confirmationHash
}

func (command ConfirmCustomerEmailAddress) MessageID() es.MessageID {
	return command.messageID
}
