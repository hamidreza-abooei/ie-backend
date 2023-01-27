package monitor

import "errors"

type Scheduler struct {
	Mnt  *Monitor
	Quit chan struct{}
}

// NewScheduler creates a new scheduler instance with mnt as monitor
// it also creates a quit signal channel for emergency exits
func NewScheduler(mnt *Monitor) (*Scheduler, error) {
	sch := &Scheduler{Quit: make(chan struct{})}
	if mnt != nil {
		sch.Mnt = mnt
		return sch, nil
	}
	return nil, errors.New("cannot create a scheduler with nil monitor")
}
