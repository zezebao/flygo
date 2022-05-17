package eventdispatcher

import (
	"container/list"
	"unsafe"
)

type Event struct {
	Type string
	Data interface{}
}

type EventHandler func(e *Event)

//事件侦听器
type EventDispatcher struct {
	HandlerMap map[string]*list.List
}

func NewEventDispatcher() *EventDispatcher {
	ed := &EventDispatcher{}
	ed.HandlerMap = make(map[string]*list.List)
	return ed
}

//添加事件侦听
func (e *EventDispatcher) AddListener(eType string, eFunc *EventHandler) *EventDispatcher {
	if funcList, ok := e.HandlerMap[eType]; ok {
		funcList.PushBack(eFunc)
	} else {
		funcList := list.New()
		funcList.PushBack(eFunc)
		e.HandlerMap[eType] = funcList
	}
	return e
}

//移除事件侦听
func (e *EventDispatcher) RemoveListener(eType string, eFunc *EventHandler) *EventDispatcher {
	if funcList, ok := e.HandlerMap[eType]; ok {
		for ir := funcList.Front(); ir != nil; ir = ir.Next() {
			a := *(*int)(unsafe.Pointer(ir.Value.(*EventHandler)))
			b := *(*int)(unsafe.Pointer(eFunc))
			if a == b {
				funcList.Remove(ir)
				break
			}
		}
	}
	return e
}

//发送事件
func (e *EventDispatcher) Dispatch(eType string, data interface{}) *EventDispatcher {
	evt := &Event{}
	evt.Type = eType
	evt.Data = data

	if funcList, ok := e.HandlerMap[eType]; ok {
		for ir := funcList.Front(); ir != nil; ir = ir.Next() {
			handler := ir.Value.(*EventHandler)
			(*handler)(evt)
		}
	}

	return e
}
