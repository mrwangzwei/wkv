package pool

import (
	"fmt"
	"log"
	"math/rand"
)

type (
	handleF  func(v interface{}) error
	pHandleF func(param interface{}, err string)
)

type task struct {
	h      handleF     //要执行的方法
	params interface{} //方法的参数
	ph     pHandleF    //异常或执行出错(error)后要执行的方法
	//taskId string      //自生成的任务id....暂时没什么用，看之后要不要扩展一个id:*task来做任务回源
}

//执行任务
func (t *task) execute(p *Pool) {
	defer func() {
		defer p.wg.Done()
		if r := recover(); r != nil {
			if t.ph != nil {
				msg := fmt.Sprintf("%s", r)
				t.ph(t.params, msg)
			} else {
				log.Printf("worker panic %s", r)
			}
		}
	}()

	err := t.h(t.params)
	if err != nil {
		panic(err.Error())
	}
}

func newTaskId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
