package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"
)

type Event struct {
	from    *Process
	to      *Process
	message string
	time    uint32
}

type Process struct {
	name   string
	time   uint32
	events chan *Event
}

func (p *Process) String() string {
	return p.name
}

func (p *Process) Step() {
	atomic.AddUint32(&p.time, 1)
}

func (p *Process) Send(to *Process, message string) {
	p.Step()
	if p.name != to.name {
		e := &Event{
			from:    p,
			to:      to,
			message: message,
			time:    p.time,
		}
		to.events <- e
	}
}

func (p *Process) receive() {
	for {
		select {
		case e := <-p.events:
			if e.time > p.time {
				p.time = e.time
			}
			p.Step()
			e.time = p.time
			log.Printf("%v: %v => %v: %v", e.time, e.from, e.to, e.message)
		}
	}
}

func NewProcess(name string) *Process {
	p := &Process{
		name:   name,
		events: make(chan *Event),
	}
	go p.receive()
	return p
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	p1 := NewProcess("Node_1")
	p2 := NewProcess("Node_2")
	p3 := NewProcess("Node_3")

	processes := [3]*Process{p1, p2, p3}
	for {
		time.Sleep(time.Second)
		a := processes[rand.Int()%len(processes)]
		b := processes[rand.Int()%len(processes)]
		a.Send(b, "Hello from random process")
	}
}
