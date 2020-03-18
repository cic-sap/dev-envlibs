package envlibs

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_t1(t *testing.T) {
	Convey("t1", t, func() {
		So(1, ShouldEqual, 1)
	})
}