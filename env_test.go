package envlibs

import (
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_t1(t *testing.T) {
	Convey("t1", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", "http://dev-info/env",
			httpmock.NewStringResponder(200, `{"canary/*":"testing","canary/staging":"staging", "canary/hello":"hello"}`))

		r, err := GetEnvs()
		So(err, ShouldBeNil)
		So(r["canary/hello"], ShouldEqual, "hello")
		t.Logf("%+v", r)

		v, ok := GetMatch("canary", "hello", r)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "hello")

		v, ok = GetMatch("canary", "qqq", r)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "testing")

		v, ok = GetMatch("canary", "staging", r)
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "staging")

		v, ok = GetMatch("canary1", "staging", r)
		So(ok, ShouldBeFalse)
		So(v, ShouldEqual, "")

		v, ok, err = GetOriginMatch("canary", "hello")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "hello")
		So(err, ShouldBeNil)

		v, ok, err = GetOriginMatch("canary", "qqq")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "testing")
		So(err, ShouldBeNil)


		v, ok, err = GetOriginMatch("canary", "staging")
		So(ok, ShouldBeTrue)
		So(v, ShouldEqual, "staging")
		So(err, ShouldBeNil)

		v, ok, err = GetOriginMatch("canary1", "staging")
		So(ok, ShouldBeFalse)
		So(v, ShouldEqual, "")
		So(err, ShouldBeNil)



	})

	Convey("iter", t, func() {
		m := GetAllExtraValuesFiles("./test/helm1")
		t.Logf("%+v", m)
		So(len(m), ShouldEqual, 4)
	})
}