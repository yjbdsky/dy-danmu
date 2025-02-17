package app

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400
	Unauthorized  = 401
	NotFound      = 404

	// 业务错误码从1000开始
	ErrInvalidRequest = 1000
	ErrDatabase       = 1001
)

var MsgFlags = map[int]string{
	SUCCESS:       "ok",
	ERROR:         "internal error",
	InvalidParams: "invalid parameters",
	Unauthorized:  "unauthorized",
	NotFound:      "not found",

	ErrInvalidRequest: "invalid request",
	ErrDatabase:       "database error",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
