package handler

import "errors"

func UserErrorExist(err error) {
	if err != nil {
		err = errors.New("用户服务错误：" + err.Error())
		panic(err)
	}
}

func TaskErrorExist(err error) {
	if err != nil {
		err = errors.New("任务服务错误：" + err.Error())
		panic(err)
	}
}
