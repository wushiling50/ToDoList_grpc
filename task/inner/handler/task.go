package handler

import (
	"context"
	"main/ToDoList_grpc/pkg/e"
	"main/ToDoList_grpc/service"
	"main/ToDoList_grpc/task/inner/repository"
)

type TaskService struct {
	service.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (*TaskService) TaskCreate(ctx context.Context, req *service.TaskRequest) (resp *service.UsualResponse, err error) {
	var task repository.Task
	resp = new(service.UsualResponse)
	resp.Code = e.Success
	err = task.Create(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

func (*TaskService) TaskShow(ctx context.Context, req *service.TaskRequest) (resp *service.TasksDetailResponse, err error) {
	var task repository.Task
	resp = new(service.TasksDetailResponse)
	taskListRep, err := task.Show(req)
	resp.Code = e.Success
	if err != nil {
		resp.Code = e.Error
		return resp, err
	}
	resp.TaskDetail = repository.BuildTasks(taskListRep)
	return resp, nil
}

func (*TaskService) TaskUpdate(ctx context.Context, req *service.TaskRequest) (resp *service.UsualResponse, err error) {
	var task repository.Task
	resp = new(service.UsualResponse)
	resp.Code = e.Success
	err = task.Update(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}

func (*TaskService) TaskDelete(ctx context.Context, req *service.TaskRequest) (resp *service.UsualResponse, err error) {
	var task repository.Task
	resp = new(service.UsualResponse)
	resp.Code = e.Success
	err = task.Delete(req)
	if err != nil {
		resp.Code = e.Error
		resp.Msg = e.GetMsg(e.Error)
		resp.Data = err.Error()
		return resp, err
	}
	resp.Msg = e.GetMsg(uint(resp.Code))
	return resp, nil
}
