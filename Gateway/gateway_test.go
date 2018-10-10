package main

import(
	//"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	)

func TestParseUser(t *testing.T) {
	Convey("Given some struct value", t, func(){
		var variable *http.Request

		Convey("When parsed", func(){
			var res UserExternal
			ParseUser(variable)

			Convey("The value should be entered", func(){
				So(variable, ShouldEqual, res)
				})
			})
		})
}