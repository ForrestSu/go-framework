package framework

import "context"

type Event struct {
	Source  string
	Content string
}

//事件上报接口, 由framework实现
type EventReceiver interface {
	OnEvent(evt Event)
}
// framework 接口
type IFrameWork interface {
	RegisterService(name string, svr IService) error
	Start() error
	Stop() error
	Destroy() error
}

// IService
type IService interface {
	Init(event EventReceiver) error
	Start(ctx context.Context) error
	Stop() error
	Destroy() error
}
