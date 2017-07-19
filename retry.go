package httptool

import (
	"container/list"
	//"fmt"
	"sync"
	"time"
)

const (
	retrycount  = 5 //重试次数
	retrydetect = 5 //重试检测市场间隔
)

var retryList *list.List
var mutex *sync.Mutex

//retrydata 重试结构体
type retrydata struct {
	logid, url, param string
	count, flag       int
}

//init 初始化url重试请求
func init() {
	retryList = list.New()
	mutex = new(sync.Mutex)
	go func() {
		timer := time.NewTicker(time.Duration(retrydetect) * time.Second)
		defer timer.Stop()
		for {
			select {
			case _, ok := <-timer.C:
				if !ok {
					break
				}
				retry()
			}
		}
	}()
}

func retry() {
	mutex.Lock()
	defer mutex.Unlock()
	if retryList.Len() == 0 {
		return
	}
	//fmt.Println("retry len = ", retryList.Len())
	var err error
	for c := retryList.Front(); c != nil; {
		cd := c.Value.(*retrydata)
		isDelete := false
		if cd.flag == 1 {
			_, err = Post(cd.logid, cd.url, cd.param, true)
		} else {
			_, err = Get(cd.logid, cd.url, cd.param)
		}

		if err != nil {
			//fmt.Println(err.Error())
			cd.count++
			//超过重试次数删除
			if cd.count >= retrycount {
				isDelete = true
			}
		} else {
			//fmt.Println("res := ", res)
			isDelete = true
		}
		if isDelete {
			next := c.Next()
			retryList.Remove(c)
			c = next
		} else {
			c = c.Next()
		}
	}
}

//Put 添加重试 flag: 1.post,2.get
func Put(logid, url, param string, flag int) {
	if flag == 0 {
		flag = 1
	}
	mutex.Lock()
	defer mutex.Unlock()
	retryList.PushBack(&retrydata{logid, url, param, 0, flag})
}
