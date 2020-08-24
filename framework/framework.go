package framework

import (
	"context"
	"errors"
	"fmt"
	"log"
)

const (
	StateInit = iota
	StateRunning
)

/*
 * 定义框架接口
type IFrameWork interface {
	RegisterService(name string, svr IService) error
	Start() error
	Stop() error
	Destroy() error
	//private
	//startServices() error
	//stopServices() error
	//destroyServices() error
}*/

type FrameWork struct {
	services map[string]IService
	evtBuf   chan Event
	cancel   context.CancelFunc
	ctx      context.Context
	state    int
}

// 实例化框架对象
func NewFrameWork(sizeEvtBuf int) *FrameWork {
	frame := FrameWork{
		services: map[string]IService{},
		evtBuf:   make(chan Event, sizeEvtBuf),
		state:    StateInit,
	}
	return &frame
}

/**
 * 以下为 Framework 的对外接口
 */
func (f *FrameWork) RegisterService(name string, svr IService) error {
	f.services[name] = svr
	return nil
}

func (f *FrameWork) Start() error {
	if f.state != StateInit {
		return WrongStateError
	}
	f.state = StateRunning
	f.ctx, f.cancel = context.WithCancel(context.Background())
	//TODO 先启动一个事件处理协程
	go f.EventProcessGoroutine()
	// 然后才启动服务
	return f.startServices()
}

func (f *FrameWork) Stop() error {
	if f.state != StateRunning {
		return WrongStateError
	}
	f.state = StateInit
	// 通知停止所有的子协程
	f.cancel()
	return f.stopServices()
}

func (f *FrameWork) Destroy() error {
	if f.state != StateInit {
		return WrongStateError
	}
	return f.destroyServices()
}

// 实现事件通知接口，用于接收多个子服务的状态
func (f *FrameWork) OnEvent(evt Event) {
	f.evtBuf <- evt
}


/**
 * 以下为 Framework 功能实现
 * internal func
 */
func (f *FrameWork) startServices() error {
	var err error
	//var mutex sync.Mutex
	//FIXME: 注意，这里的服务启动的错误，应该捕捉不到!
	fmt.Println("startServices() start...")
	for name, svr := range f.services {
		go func(name string, svr IService, ctx context.Context) {
			//var errs ServicesError
			err = svr.Init(f)
			if err != nil {
				log.Println("fail to init "+ name, err.Error())
				return
			}
			err = svr.Start(ctx)
			if err != nil {
				log.Println("fail to start "+ name, err.Error())
				return
			}
		}(name, svr, f.ctx)
	}
	fmt.Println("startServices() end.")
	return nil
}

func (f *FrameWork) stopServices() error {
	var err error
	var errs ServicesError
	for name, svr := range f.services {
		if err = svr.Stop(); err != nil {
			errs.errArr = append(errs.errArr,
				errors.New("stop: "+name+":"+err.Error()))
		}
	}
	if len(errs.errArr) == 0 {
		return nil
	}
	return errs
}

func (f *FrameWork) destroyServices() error {
	var err error
	var errs ServicesError
	for name, svr := range f.services {
		if err = svr.Destroy(); err != nil {
			errs.errArr = append(errs.errArr,
				errors.New("unInit: "+name+":"+err.Error()))
		}
	}
	if len(errs.errArr) == 0 {
		return nil
	}
	return errs
}

// 由主服务启动一个事件处理器，处理所有子服务上报的事件
func (f *FrameWork) EventProcessGoroutine() {
	var evtSeg [10]Event
	for {
		for i := 0; i < 10; i++ {
			select {
			case evtSeg[i] = <-f.evtBuf:
			case <-f.ctx.Done():
				return
			}
		}
		fmt.Println("reported:", evtSeg)
	}
}

