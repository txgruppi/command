package command_test

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/txgruppi/command"
)

func TestErrorGroup(t *testing.T) {
	Convey("ErrorGroup", t, func() {
		errs := []error{
			errors.New("Error 1"),
			errors.New("Error 2"),
			errors.New("Error 3"),
			errors.New("Error 4"),
		}
		var err error = &command.ErrorGroup{
			Errors: errs,
		}

		Convey("Errors array", func() {
			Convey("the Errors array should be accessible", func() {
				errGroup := err.(*command.ErrorGroup)
				for i, v := range errGroup.Errors {
					So(v, ShouldEqual, errs[i])
				}
			})
		})

		Convey(".Error()", func() {
			Convey("the .Error() method should return all error messages joined by a new line", func() {
				expected := "Error 1\nError 2\nError 3\nError 4"
				So(err.Error(), ShouldEqual, expected)
			})
		})
	})
}
