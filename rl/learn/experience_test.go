package learn

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExperiences(t *testing.T) {
	Convey("With  experiences", t, func() {
		exps := newExperiences()
		Convey("We should be able to add one", func() {
			exp := &experience{}
			exps.add(exp)
			Convey("And get it back", func() {
				So(exps.pick(), ShouldResemble, exp)
			})
		})
	})
}
