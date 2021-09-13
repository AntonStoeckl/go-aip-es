package domain_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/AntonStoeckl/go-aip-es/domain"
)

func TestPersonName_Equals(t *testing.T) {
	Convey("Given a PersonName", t, func() {
		personName := domain.RebuildPersonName("Lib", "Gallagher")

		Convey("When it is compared with an identical PersonName", func() {
			identicalPersonName := domain.RebuildPersonName(personName.GivenName(), personName.FamilyName())
			isEqual := personName.Equals(identicalPersonName)

			Convey("Then it should be equal", func() {
				So(isEqual, ShouldBeTrue)
			})
		})

		Convey("When it is compared with another PersonName with different givenName", func() {
			differentPersonName := domain.RebuildPersonName("Phillip", personName.FamilyName())
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})

		Convey("When it is compared with another PersonName with different familyName", func() {
			differentPersonName := domain.RebuildPersonName(personName.GivenName(), "Jackson")
			isEqual := personName.Equals(differentPersonName)

			Convey("Then it should not be equal", func() {
				So(isEqual, ShouldBeFalse)
			})
		})
	})
}
