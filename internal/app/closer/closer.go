package closer

import "sync"

type closeFunc = func()

type Closer struct {
	closeFuncs []closeFunc

	mx sync.Mutex
}

func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) RunAll() {
	wg := sync.WaitGroup{}

	for _, f := range c.closeFuncs {
		wg.Add(1)
		go func(f closeFunc) {
			defer wg.Done()

			f()
		}(f)
	}

	wg.Wait()
}

func (c *Closer) Add(f closeFunc) {
	c.mx.Lock()
	c.closeFuncs = append(c.closeFuncs, f)
	c.mx.Unlock()
}
