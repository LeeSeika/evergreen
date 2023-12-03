package biz

type RespCode int64

const (
	CodeSuccess RespCode = 1000 + iota
	CodeInvalidParam
	CodeUserExists
	CodeUserNotFound
	CodeInvalidPassword
	CodeServerBusy
	CodeAuthNotFound
	CodeInvalidAuth
	CodeInvalidToken
	CodeUserNotLogin
	CodeCommentDeleted
)

var codeMsgMap = map[RespCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "invalid parameters",
	CodeUserExists:      "user exists",
	CodeUserNotFound:    "user not found",
	CodeInvalidPassword: "invalid password",
	CodeServerBusy:      "server busy",
	CodeAuthNotFound:    "authorization not found",
	CodeInvalidAuth:     "invalid format of authorization",
	CodeInvalidToken:    "invalid token",
	CodeUserNotLogin:    "user not login",
	CodeCommentDeleted:  "comment has been deleted",
}

func (c RespCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
