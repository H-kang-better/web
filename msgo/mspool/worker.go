package mspool

import (
	msgoLog "github.com/H-kang-better/msgo/log"
	"time"
)

type Worker struct {
	pool *Pool
	//task 任务队列
	task chan func()
	//lastTime 执行任务的最后的时间
	lastTime time.Time
}

func (w *Worker) run() {
	w.pool.incRunning()
	go w.running()
}

func (w *Worker) running() {
	defer func() {
		w.pool.workerCache.Put(w)
		w.pool.decRunning()
		if err := recover(); err != nil {
			if w.pool.PanicHandler != nil {
				w.pool.PanicHandler()
			} else {
				msgoLog.Default().Error(err)
			}
		}
		w.pool.cond.Signal()
	}()
	for f := range w.task {
		if f == nil {
			w.pool.workerCache.Put(w)
			return
		}
		f()
		// 任务运行完成，worker空闲
		w.pool.PutWorker(w)
	}
}
