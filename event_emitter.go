package playwright

import (
	"reflect"
	"sync"
)

type (
	eventRegister struct {
		once []interface{}
		on   []interface{}
	}
	EventEmitter struct {
		sync.Mutex
		events              map[string]*eventRegister
		addEventHandlers    []func(name string, handler interface{})
		removeEventHandlers []func(name string, handler interface{})
	}
)

func (e *EventEmitter) Emit(name string, payload ...interface{}) {
	e.Lock()
	defer e.Unlock()
	if _, ok := e.events[name]; !ok {
		return
	}

	payloadV := make([]reflect.Value, 0)

	for _, p := range payload {
		payloadV = append(payloadV, reflect.ValueOf(p))
	}

	for _, handler := range e.events[name].on {
		handlerV := reflect.ValueOf(handler)
		handlerV.Call(payloadV[:handlerV.Type().NumIn()])
	}
	for _, handler := range e.events[name].once {
		handlerV := reflect.ValueOf(handler)
		handlerV.Call(payloadV[:handlerV.Type().NumIn()])
	}
	e.events[name].once = make([]interface{}, 0)
}

func (e *EventEmitter) Once(name string, handler interface{}) {
	e.addEvent(name, handler, true)
}

func (e *EventEmitter) On(name string, handler interface{}) {
	e.addEvent(name, handler, false)
}

func (e *EventEmitter) addEventHandler(handler func(name string, handler interface{})) {
	e.addEventHandlers = append(e.addEventHandlers, handler)
}

func (e *EventEmitter) removeEventHandler(handler func(name string, handler interface{})) {
	e.removeEventHandlers = append(e.removeEventHandlers, handler)
}

func (e *EventEmitter) RemoveListener(name string, handler interface{}) {
	for _, mitm := range e.removeEventHandlers {
		mitm(name, handler)
	}
	e.Lock()
	defer e.Unlock()
	if _, ok := e.events[name]; !ok {
		return
	}
	handlerPtr := reflect.ValueOf(handler).Pointer()

	onHandlers := []interface{}{}
	for idx := range e.events[name].on {
		eventPtr := reflect.ValueOf(e.events[name].on[idx]).Pointer()
		if eventPtr != handlerPtr {
			onHandlers = append(onHandlers, e.events[name].on[idx])
		}
	}
	e.events[name].on = onHandlers

	onceHandlers := []interface{}{}
	for idx := range e.events[name].once {
		eventPtr := reflect.ValueOf(e.events[name].once[idx]).Pointer()
		if eventPtr != handlerPtr {
			onceHandlers = append(onceHandlers, e.events[name].once[idx])
		}
	}

	e.events[name].once = onceHandlers
}

func (e *EventEmitter) ListenerCount(name string) int {
	count := 0
	e.Lock()
	for key := range e.events {
		count += len(e.events[key].on) + len(e.events[key].once)
	}
	e.Unlock()
	return count
}

func (e *EventEmitter) addEvent(name string, handler interface{}, once bool) {
	for _, mitm := range e.addEventHandlers {
		mitm(name, handler)
	}
	e.Lock()
	if _, ok := e.events[name]; !ok {
		e.events[name] = &eventRegister{
			on:   make([]interface{}, 0),
			once: make([]interface{}, 0),
		}
	}
	if once {
		e.events[name].once = append(e.events[name].once, handler)
	} else {
		e.events[name].on = append(e.events[name].on, handler)
	}
	e.Unlock()
}

func (e *EventEmitter) initEventEmitter() {
	e.events = make(map[string]*eventRegister)
}
