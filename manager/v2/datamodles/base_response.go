/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package datamodles

type ResultResponse struct {
    Status bool        `json:"status"`
    Msg    string      `json:"msg"`
    Data   interface{} `json:"data"`
}

func GenResultData(status bool, msg string, data interface{}) *ResultResponse {
    return &ResultResponse{status, msg, data}
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
