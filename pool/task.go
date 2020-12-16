package pool

import (
	"fmt"
	"log"
	"math/rand"
)

type (
	HandleF  func(v interface{}) error
	PHandleF func(param interface{}, err interface{})
)

type task struct {
	handle  HandleF     //要执行的方法
	params  interface{} //方法的参数
	pHandle PHandleF    //异常或执行出错(error)后要执行的方法
	taskId  string      //自生成的任务id....暂时没什么用，看之后要不要扩展一个id:*task来做任务回源
}

//执行任务
func (t *task) execute() {
	defer func() {
		defer pw.Done()
		if r := recover(); r != nil {
			if t.pHandle != nil {
				t.pHandle(t.params, r)
			} else {
				log.Printf("worker panic %s %s", r, t.taskId)
			}
		}
	}()

	err := t.handle(t.params)
	if err != nil {
		panic(err.Error())
	}
}

func (t *task) getTaskId() string {
	return t.taskId
}

func newTaskId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
