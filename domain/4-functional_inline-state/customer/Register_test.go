package customer_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/AntonStoeckl/go-aip-es/domain"
	"github.com/AntonStoeckl/go-aip-es/domain/4-functional_inline-state/customer"
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
				event := customer.Register(command)

				Convey("Then CustomerRegistered", func() {
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
