package repository

import (
	"errors"
	"main/ToDoList_grpc/service"
	"strconv"
)

type Task struct {
	TaskID    uint `gorm:"primarykey"` // id
	UserID    uint `gorm:"index"`      // 用户id
	Status    int  `gorm:"default:0"`  //0表示代办，1表示完成
	Title     string
	Content   string `gorm:"type:longtext"`
	StartTime int64
	EndTime   int64
}

// 增
func (*Task) Create(req *service.TaskRequest) error {
	if req.Status > 1 {
		return errors.New("不规范的操作")
	}
	task := Task{
		UserID:    uint(req.UserID),
		Title:     req.Title,
		Content:   req.Content,
		Status:    int(req.Status),
		StartTime: int64(req.StartTime),
		EndTime:   int64(req.EndTime),
	}

	return DB.Create(&task).Error
}

// 查
func (*Task) Show(req *service.TaskRequest) (taskList []Task, err error) {
	// req.Action
	action := req.Action
	switch action {
	case 1: //查询所有事项
		err = DB.Model(Task{}).Where("user_id=?", req.UserID).
			Limit(10).Offset(int((req.PageNum - 1) * 10)).Find(&taskList).Error
	case 2: //查询所有代办事项
		err = DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status=?", 0).
			Limit(10).Offset(int((req.PageNum - 1) * 10)).Find(&taskList).Error
	case 3: //查询所有完成事项
		err = DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status=?", 1).
			Limit(10).Offset(int((req.PageNum - 1) * 10)).Find(&taskList).Error
	case 4: //关键词查询
		searchMsg := req.SearchMsg
		numMsg, err := strconv.ParseInt(searchMsg, 10, 0)
		if err != nil {
			err = DB.Model(Task{}).Where("user_id=?", req.UserID).
				Where("title LIKE ? OR content LIKE ?", "%"+searchMsg+"%", "%"+searchMsg+"%").
				Limit(10).Offset(int((req.PageNum - 1) * 10)).Find(&taskList).Error
		} else {
			err = DB.Model(Task{}).Where("user_id=?", req.UserID).
				Where("task_id LIKE ? OR status LIKE ? OR title LIKE ? OR content LIKE ?", numMsg, numMsg, "%"+searchMsg+"%", "%"+searchMsg+"%").
				Limit(10).Offset(int((req.PageNum - 1) * 10)).Find(&taskList).Error

		}
	default:
		err = errors.New("不规范的操作")
		return taskList, err
	}
	if err != nil {
		return taskList, err
	}
	return taskList, nil
}

// 改
func (*Task) Update(req *service.TaskRequest) error {
	action := req.Action
	switch action {
	case 1: //更改一条内容
		t := Task{}
		err := DB.Where("task_id=?", req.TaskID).First(&t).Error
		if err != nil {
			return err
		}
		t.Title = req.Title
		t.Content = req.Content
		t.Status = int(req.Status)
		t.StartTime = int64(req.StartTime)
		t.EndTime = int64(req.EndTime)
		return DB.Save(&t).Error
	case 2: //一条代办改为完成
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("task_id=?", req.TaskID).
			Where("status = ?", 0).Update("status", 1).Error
	case 3: //一条完成改为代办
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("task_id=?", req.TaskID).
			Where("status = ?", 1).Update("status", 0).Error
	case 4: //所有代办改为完成
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status = ?", 0).Update("status", 1).Error
	case 5: //所有完成改为代办
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status = ?", 1).Update("status", 0).Error
	default:
		err := errors.New("不规范的操作")
		return err
	}

}

// 删
func (*Task) Delete(req *service.TaskRequest) error {
	action := req.Action
	switch action {
	case 1: //删除一条内容
		return DB.Where("user_id=?", req.UserID).Where("task_id=?", req.TaskID).Delete(Task{}).Error
	case 2: //删除所有已完成
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status = ?", 1).Delete(Task{}).Error
	case 3: //删除所有代办
		return DB.Model(Task{}).Where("user_id=?", req.UserID).Where("status = ?", 0).Delete(Task{}).Error
	case 4: //删除所有
		return DB.Where("user_id=?", req.UserID).Delete(Task{}).Error
	default:
		err := errors.New("不规范的操作")
		return err
	}

}

// 序列化的构造器（task）
func BuildTask(item Task) *service.TaskModel {
	return &service.TaskModel{
		TaskID:    uint32(item.TaskID),
		UserID:    uint32(item.UserID),
		Status:    uint32(item.Status),
		Title:     item.Title,
		Content:   item.Content,
		StartTime: uint32(item.StartTime),
		EndTime:   uint32(item.EndTime),
	}
}

func BuildTasks(item []Task) (tList []*service.TaskModel) {
	for _, v := range item {
		f := BuildTask(v)
		tList = append(tList, f)
	}
	return tList
}
