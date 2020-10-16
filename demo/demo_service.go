package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	fw "github.com/ForrestSu/go-framework/framework"
)

type DemoService struct {
	evtReceiver fw.EventReceiver
	agtCtx      context.Context
	stopChan    chan struct{}
	name        string
	content     string
}

func NewDemoService(name string, content string) *DemoService {
	return &DemoService{
		stopChan: make(chan struct{}),
		name:     name,
		content:  content,
	}
}

func (d *DemoService) Init(evtReceiver fw.EventReceiver) error {
	fmt.Println("initialize demo", d.name)
	d.evtReceiver = evtReceiver
	return nil
}

func (d *DemoService) Start(ctx context.Context) error {
	fmt.Println("start demo", d.name)
	for {
		select {
		case <-ctx.Done():
			//这里响应框架的通知,一旦收到通知,
			//处理剩余的工作, 这里可以处理然后主动退出

			d.stopChan <- struct{}{}
			break
		default:
			time.Sleep(time.Millisecond * 50)
			// TODO 向 go_framework 上报数据
			d.evtReceiver.OnEvent(fw.Event{Source: d.name, Content: d.content})
		}
	}
}

func (d *DemoService) Stop() error {
	fmt.Println("stop demo", d.name)
	select {
	case <-d.stopChan:
		return nil
	case <-time.After(time.Second * 1):
		return errors.New("failed to stop for timeout")
	}
}

func (d *DemoService) Destroy() error {
	fmt.Println(d.name, "released resources.")
	return nil
}
