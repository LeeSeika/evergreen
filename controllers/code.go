package controllers

type RespCode int64

const (
	CodeSuccess RespCode = 1000 + iota
	CodeInvalidParam
	CodeUserExists
	CodeUserNotFound
	CodeInvalidPassword
	CodeServerBusy
)

var codeMsgMap = map[RespCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "invalid parameters",
	CodeUserExists:      "user exists",
	CodeUserNotFound:    "user not found",
	CodeInvalidPassword: "invalid password",
	CodeServerBusy:      "server busy",
}

func (c RespCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
