package datamodles

type ResultResponse struct {
    Status bool        `json:"status"`
    Msg    string      `json:"msg"`
    Data   interface{} `json:"data"`
}

func GenSuccessData(data interface{}) *ResultResponse {
    return &ResultResponse{true, "", data}
}

func GenSuccessMsg(msg string) *ResultResponse {
    return &ResultResponse{true, msg, ""}
}

func GenSuccess(msg string, data interface{}) *ResultResponse {
    return &ResultResponse{true, msg, data}
}

func GenFailedMsg(errMsg string) *ResultResponse {
    return &ResultResponse{false, errMsg, ""}
}
