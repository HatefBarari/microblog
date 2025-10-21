package httputil

type SuccessResp struct {
	Data interface{} `json:"data"`
	Meta *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Message string `json:"message"`
}

func OK(data interface{}) SuccessResp {
	return SuccessResp{Data: data, Meta: &Meta{Message: "success"}}
}