package wavx

import "sync"

type Activator struct {
	sync.RWMutex
	inactive bool
}

func (a *Activator) Activate() {
	a.Lock()
	defer a.Unlock()
	a.inactive = false
}

func (a *Activator) Deactivate() {
	a.Lock()
	defer a.Unlock()
	a.inactive = true
}

func (a *Activator) IsActive() bool {
	a.RLock()
	defer a.RUnlock()
	return !a.inactive
}
