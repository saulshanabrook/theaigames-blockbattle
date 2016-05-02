package game

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRotation(t *testing.T) {
	Convey("invert", t, func() {
		Convey("should be an additive inverse for all rotations", func() {
			for _, rot := range allRotations {
				So(rot.add(rot.invert()), ShouldEqual, Rotation(0))
			}
		})
	})

	Convey("distance", t, func() {
		Convey("should wrap around properly", func() {
			So(RotatedUp.distance(RotatedDown), ShouldEqual, 2)
			So(RotatedDown.distance(RotatedUp), ShouldEqual, 2)

			So(RotatedUp.distance(RotatedLeft), ShouldEqual, 1)
			So(RotatedLeft.distance(RotatedUp), ShouldEqual, 1)

			So(RotatedLeft.distance(RotatedRight), ShouldEqual, 2)
			So(RotatedRight.distance(RotatedLeft), ShouldEqual, 2)

			So(RotatedLeft.distance(RotatedDown), ShouldEqual, 1)
			So(RotatedDown.distance(RotatedLeft), ShouldEqual, 1)
		})
	})
}
