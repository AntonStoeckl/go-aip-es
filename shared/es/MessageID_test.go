package es_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/AntonStoeckl/go-aip-es/shared/es"
)

func TestMessageID_Generate(t *testing.T) {
	Convey("When a MessageID is generated", t, func() {
		messageID := es.GenerateMessageID()

		Convey("It should not be empty", func() {
			So(messageID, ShouldNotBeEmpty)
		})
	})
}

func TestMessageID_Build(t *testing.T) {
	Convey("When a MessageID is built from another MessageID", t, func() {
		otherMessageID := es.GenerateMessageID()
		messageID := es.BuildMessageID(otherMessageID)

		Convey("It should not be empty", func() {
			So(messageID, ShouldNotBeEmpty)

			Convey("And it should equal the input messageID", func() {
				So(messageID.Equals(otherMessageID), ShouldBeTrue)
			})
		})
	})
}
