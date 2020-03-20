package envlibs

import (
	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFetchPkg(t *testing.T) {
	Convey("p1", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		// Exact URL match
		httpmock.RegisterResponder("GET", "http://dev-info/env",
			httpmock.NewStringResponder(200, `{"canary/*":"testing","canary/staging":"staging", "canary/hello":"hello"}`))
		rs, err := GetValues("Eureka", "go-server-with-multi-values", "0.0.1-66ff51ee864ccfdbe411d092b362854c34e279fb", "./output", "canary", "hello7")
		t.Logf("%+v, %+v",rs, err)
		So(err, ShouldBeNil)
		So(rs, ShouldHaveLength, 1)
		So(rs[0], ShouldEqual, "values.testing.yaml")
	})

}
