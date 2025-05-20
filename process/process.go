package process

import (
	"fmt"
	"sync"
	"time"
)

var (
	ErrNoExist = "the pid %s program does not exist"
	ErrStatus  = "unknown status"
	lock       sync.Mutex
)

type StatCode uint8

const (
	StatRunning StatCode = 0x01
	StatStop    StatCode = 0x00
	StatNone    StatCode = 0xFF
)

type IProcess interface {
	Close() error
	Start() error
	Stop() error
}

type Process struct {
	IProcess
	stat StatCode
}

var (
	list map[string]*Process
	wg   *sync.WaitGroup
)

func New(w *sync.WaitGroup) {
	list = make(map[string]*Process)
	wg = w
}

func Get(name string) (*Process, error) {
	defer lock.Unlock()
	lock.Lock()

	return get(name)
}

func get(name string) (*Process, error) {
	if p, ok := list[name]; ok {
		return p, nil
	}
	return nil, fmt.Errorf(ErrNoExist, name)
}

func Put(name string, ip IProcess) {
	defer lock.Unlock()
	lock.Lock()

	p := &Process{IProcess: ip}
	list[name] = p
}

func Remove(name string) error {
	defer lock.Unlock()
	lock.Lock()

	_, err := get(name)
	if err != nil {
		return err
	}
	delete(list, name)
	return nil
}

func Restart(name string, d time.Duration) error {
	defer lock.Unlock()
	lock.Lock()

	p, err := get(name)
	if err != nil {
		return err
	}

	p.Stop()
	time.Sleep(d)
	start(p)
	return err
}

func Start(name string) error {
	defer lock.Unlock()
	lock.Lock()

	p, err := get(name)
	if err != nil {
		return err
	}
	start(p)
	return err
}

func Starts() {
	defer lock.Unlock()
	lock.Lock()

	for _, p := range list {
		start(p)
	}
}

func start(p *Process) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		p.stat = StatRunning
		err := p.Start()
		if err != nil {
			p.stat = StatNone
		} 
	}(wg)
}

func Stop(name string) error {
	defer lock.Unlock()
	lock.Lock()

	p, err := get(name)
	if err != nil {
		return err
	}
	err = p.Close()
	if err == nil {
		p.stat = StatStop
	}
	return err
}

func Status(name string) StatCode {
	defer lock.Unlock()
	lock.Lock()

	p, err := get(name)
	if err != nil {
		return StatNone
	}
	return p.stat
}

func Statusall() map[string]StatCode {
	defer lock.Unlock()
	lock.Lock()

	ps := make(map[string]StatCode)
	for name, p := range list {
		ps[name] = p.stat
	}
	return ps
}
