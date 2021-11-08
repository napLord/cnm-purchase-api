package closer

import "sync"

type closeFunc = func()

//Closer is a holder for close functions.
type Closer struct {
	closeFuncs []closeFunc

	mx sync.Mutex
}

//NewCloser creates new Closer
func NewCloser() *Closer {
	return &Closer{}
}

//RunAll simultaneously runs all funcs that closeFuncs contains and wait them complete.
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

//Add adds closeFunc to closeFuncs
func (c *Closer) Add(f closeFunc) {
	c.mx.Lock()
	c.closeFuncs = append(c.closeFuncs, f)
	c.mx.Unlock()
}
