package customer_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/domain/functional_shared-state/customer"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func TestConfirmEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := domain.GenerateCustomerID()
		emailAddress, err := domain.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		invalidConfirmationHash := domain.RebuildConfirmationHash("invalid_hash")
		personName, err := domain.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildConfirmCustomerEmailAddress(customerID, emailAddress.ConfirmationHash())
		commandWithInvalidHash := domain.BuildConfirmCustomerEmailAddress(customerID, invalidConfirmationHash)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			personName,
			es.GenerateMessageID(),
			1,
		)

		confirmedEmailAddress, err := domain.ConfirmEmailAddressWithHash(emailAddress, emailAddress.ConfirmationHash())
		So(err, ShouldBeNil)

		customerEmailAddressConfirmed := domain.BuildCustomerEmailAddressConfirmed(
			customerID,
			confirmedEmailAddress,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO 1: Confirm a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents = customer.ConfirmEmailAddress(eventStream, command)

					Convey("Then CustomerEmailAddressConfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmed)
						So(ok, ShouldBeTrue)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Confirm a Customer's emailAddress with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents = customer.ConfirmEmailAddress(eventStream, commandWithInvalidHash)

					Convey("Then CustomerEmailAddressConfirmationFailed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmationFailed)
						So(ok, ShouldBeTrue)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeTrue)
						So(event.FailureReason(), ShouldBeError)
						So(event.Meta().CausationID(), ShouldEqual, commandWithInvalidHash.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to confirm a Customer's emailAddress again with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents = customer.ConfirmEmailAddress(eventStream, command)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Try to confirm a Customer's emailAddress again with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents = customer.ConfirmEmailAddress(eventStream, commandWithInvalidHash)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
