package res

import (
	"main/ToDoList_grpc/pkg/e"
)

// Response 基础序列化器
type Response struct {
	Status uint        `json:"Status"`
	Data   interface{} `json:"Data"`
	Msg    string      `json:"Msg"`
	Error  string      `json:"Error"`
}

// DataList 带有总数的Data结构
type DataList struct {
	Item  interface{} `json:"Item"`
	Total uint        `json:"Total"`
}

// TokenData 带有token的Data结构
type TokenData struct {
	User  interface{} `json:"User"`
	Token string      `json:"Token"`
}

// BulidListResponse 带有总数的列表构建器
func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Status: 200,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Msg: e.GetMsg(uint(200)),
	}
}
