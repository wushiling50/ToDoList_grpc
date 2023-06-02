package handler

import (
	"context"
	"main/pkg/e"
	"main/pkg/res"
	"main/pkg/util"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	TaskErrorExist(ginCtx.Bind(&tReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.Uid)
	taskService := ginCtx.Keys["task"].(service.TaskServiceClient) //断言

	taskResp, err := taskService.TaskCreate(context.Background(), &tReq)

	TaskErrorExist(err)
	r := res.Response{
		Data:   taskResp,
		Status: uint(taskResp.Code),
		Msg:    e.GetMsg(uint(taskResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func ListTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	TaskErrorExist(ginCtx.Bind(&tReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))

	tReq.UserID = uint32(claim.Uid)
	TaskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	TaskResp, err := TaskService.TaskShow(context.Background(), &tReq)
	TaskErrorExist(err)
	r := res.BuildListResponse(TaskResp.TaskDetail, uint(len(TaskResp.TaskDetail)))

	ginCtx.JSON(http.StatusOK, r)
}

func UpdateTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	TaskErrorExist(ginCtx.Bind(&tReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.Uid)
	TaskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	TaskResp, err := TaskService.TaskUpdate(context.Background(), &tReq)
	TaskErrorExist(err)
	r := res.Response{
		Data:   TaskResp,
		Status: uint(TaskResp.Code),
		Msg:    e.GetMsg(uint(TaskResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}

func DeleteTask(ginCtx *gin.Context) {
	var tReq service.TaskRequest
	TaskErrorExist(ginCtx.Bind(&tReq))
	claim, _ := util.ParseToken(ginCtx.GetHeader("Authorization"))
	tReq.UserID = uint32(claim.Uid)
	TaskService := ginCtx.Keys["task"].(service.TaskServiceClient)
	TaskResp, err := TaskService.TaskDelete(context.Background(), &tReq)
	TaskErrorExist(err)
	r := res.Response{
		Data:   TaskResp,
		Status: uint(TaskResp.Code),
		Msg:    e.GetMsg(uint(TaskResp.Code)),
	}
	ginCtx.JSON(http.StatusOK, r)
}
