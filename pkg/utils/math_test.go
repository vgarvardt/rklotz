package utils

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMin(t *testing.T) {
	Convey("Given two int numbers", t, func() {
		So(Min(1, 5), ShouldEqual, 1)
		So(Min(5, 1), ShouldEqual, 1)
		So(Min(5, 5), ShouldEqual, 5)
	})
}
