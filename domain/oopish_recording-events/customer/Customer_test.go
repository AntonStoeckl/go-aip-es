package customer_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/AntonStoeckl/go-aip-es/domain"
	. "github.com/AntonStoeckl/go-aip-es/domain/oopish_recording-events/customer"
	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		customerID := domain.GenerateCustomerID()
		emailAddress, err := domain.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		personName, err := domain.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildRegisterCustomer(
			customerID,
			emailAddress,
			personName,
		)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				customer := Register(command)
				recordedEvents := customer.RecordedEvents()

				Convey("Then CustomerRegistered", func() {
					So(recordedEvents, ShouldHaveLength, 1)
					event, ok := recordedEvents[0].(domain.CustomerRegistered)
					So(ok, ShouldBeTrue)
					So(event.CustomerID().Equals(customerID), ShouldBeTrue)
					So(event.EmailAddress().Equals(emailAddress), ShouldBeTrue)
					So(event.EmailAddress().ConfirmationHash().Equals(emailAddress.ConfirmationHash()), ShouldBeTrue)
					So(event.PersonName().Equals(personName), ShouldBeTrue)
					So(event.IsFailureEvent(), ShouldBeFalse)
					So(event.FailureReason(), ShouldBeNil)
					So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
					So(event.Meta().MessageID(), ShouldNotBeEmpty)
					So(event.Meta().StreamVersion(), ShouldEqual, uint(1))
				})
			})
		})
	})
}

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
					customer := ReconstituteCustomer(eventStream)
					customer.ConfirmEmailAddress(command)
					recordedEvents = customer.RecordedEvents()

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
					customer := ReconstituteCustomer(eventStream)
					customer.ConfirmEmailAddress(commandWithInvalidHash)
					recordedEvents = customer.RecordedEvents()

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
						customer := ReconstituteCustomer(eventStream)
						customer.ConfirmEmailAddress(command)
						recordedEvents = customer.RecordedEvents()

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
						customer := ReconstituteCustomer(eventStream)
						customer.ConfirmEmailAddress(command)
						recordedEvents = customer.RecordedEvents()

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}

func TestChangeName(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := domain.GenerateCustomerID()
		emailAddress, err := domain.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		personName, err := domain.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)
		changedPersonName, err := domain.BuildPersonName("Latoya", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildChangeCustomerName(customerID, changedPersonName)
		commandWithOriginalName := domain.BuildChangeCustomerName(customerID, personName)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			personName,
			es.GenerateMessageID(),
			1,
		)

		Convey("\nSCENARIO 1: Change a Customer's name", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerName", func() {
					customer := ReconstituteCustomer(eventStream)
					customer.ChangeName(command)
					recordedEvents = customer.RecordedEvents()

					Convey("Then CustomerNameChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerNameChanged)
						So(ok, ShouldBeTrue)
						So(event, ShouldNotBeNil)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.PersonName().Equals(changedPersonName), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeFalse)
						So(event.FailureReason(), ShouldBeNil)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's name to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ChangeCustomerName", func() {
					customer := ReconstituteCustomer(eventStream)
					customer.ChangeName(commandWithOriginalName)
					recordedEvents = customer.RecordedEvents()

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's name to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerNameChanged", func() {
					nameChanged := domain.BuildCustomerNameChanged(
						customerID,
						changedPersonName,
						es.GenerateMessageID(),
						2,
					)

					eventStream = append(eventStream, nameChanged)

					Convey("When ChangeCustomerName", func() {
						customer := ReconstituteCustomer(eventStream)
						customer.ChangeName(command)
						recordedEvents = customer.RecordedEvents()

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
