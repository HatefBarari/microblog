package httputil

type ErrResp struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func NewError(code int, msg string) ErrResp {
	return ErrResp{Error: msg, Code: code}
}