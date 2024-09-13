package jwt

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	Convey("Given a valid user ID", t, func() {
		uid := int64(111)
		Convey("Generate a token", func() {
			token, err := GenerateToken(uid)
			So(err, ShouldBeNil)

			Convey("Check token", func() {
				id, err, valid := CheckToken(token)
				So(err, ShouldBeNil)
				So(id, ShouldEqual, uid)
				So(valid, ShouldBeTrue)
			})
		})
	})

}
