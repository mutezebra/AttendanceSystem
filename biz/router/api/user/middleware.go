// Code generated by hertz generator.

package user

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/mutezebra/ClassroomRandomRollCallSystem/biz/middleware"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _actverifyMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _getverifycodeMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _authMw() []app.HandlerFunc {
	return []app.HandlerFunc{middleware.JWT()}
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _changepasswordMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _userinfoMw() []app.HandlerFunc {
	// your code...
	return nil
}
