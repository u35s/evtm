package evtm

import (
	"log"
	"reflect"
)

type EventPriority int
type EventType uint

const (
	EventPriorityHigh EventPriority = iota
	EventPriorityMiddle
	EventPriorityLow
	EventPriorityDefault
	EventPriorityNum
)

type eventCall struct {
	vfn  reflect.Value
	n    int // 期望的参数个数
	name string
}

type event struct {
	toCall [EventPriorityNum][]*eventCall
	num    uint
}

type EventManager struct {
	events []*event
	name   string
}

func (this *EventManager) Init(name string, num EventType) {
	this.events = make([]*event, int(num))
	this.name = name
}

func (this *EventManager) Register(fn interface{}, tp EventType, name string) {
	this.RegisterWithPriority(fn, tp, EventPriorityDefault, name)
}

func (this *EventManager) RegisterWithPriority(fn interface{}, tp EventType, p EventPriority, name string) {
	if length := len(this.events); int(tp) >= length {
		log.Printf("%v,%v,注册,类型%v超出容量%v", this.name, name, tp, length)
		return
	}
	if p >= EventPriorityNum {
		log.Printf("%v,%v,注册,类型%v优先级%v不存在", this.name, name, tp, p)
		return
	}
	evt := this.events[tp]
	if evt == nil {
		evt = &event{}
		this.events[tp] = evt
	}
	tfn := reflect.TypeOf(fn)
	call := &eventCall{reflect.ValueOf(fn), tfn.NumIn(), name}
	evt.toCall[p] = append(evt.toCall[p], call)
	evt.num++
	log.Printf("%v,%v,类型,%v,注册数量,%v", this.name, name, tp, evt.num)
}

func (this *EventManager) Exec(tp EventType, args ...interface{}) {
	event := this.events[tp]
	if event == nil {
		return
	}
	in := make([]reflect.Value, len(args))
	for i := range args {
		in[i] = reflect.ValueOf(args[i])
	}
	var num uint
	for i := range event.toCall {
		for _, call := range event.toCall[i] {
			call.vfn.Call(in)
			num++
		}
	}
}
