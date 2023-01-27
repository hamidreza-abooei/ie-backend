package monitor

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gammazero/workerpool"
	"github.com/hamidreza-abooei/ie-project/db"
	"github.com/hamidreza-abooei/ie-project/model"
)

type Monitor struct {
	store      *db.Store
	Urls       []model.Url
	wp         *workerpool.WorkerPool
	workerSize int
}

// NewMonitor creates a Monitor instance with 'store' and 'Url'
// it also creates a worker pool of size 'workerSize'
// if 'Urls' is set to nil it will be initialized with an empty slice
func NewMonitor(store *db.Store, Urls []model.Url, workerSize int) *Monitor {
	mnt := new(Monitor)
	if Urls == nil {
		mnt.Urls = make([]model.Url, 0)
	}
	mnt.Urls = Urls
	mnt.store = store
	mnt.workerSize = workerSize
	// max number of workers
	mnt.wp = workerpool.New(workerSize)
	return mnt
}

// LoadFromDatabase loads all Urls from database into monitor to start working on them
// this function will replace all of saved Urls with the ones from database
func (mnt *Monitor) LoadFromDatabase() error {
	Urls, err := mnt.store.GetAllUrls()
	if err != nil {
		return err
	}
	mnt.Urls = Urls
	return nil
}

// RemoveUrl removes a Url from current list of monitor's Urls
// returns error if the Url to be deleted was not found
func (mnt *Monitor) RemoveUrl(Url model.Url) error {
	var index = -1
	for i := range mnt.Urls {
		if mnt.Urls[i].ID == Url.ID {
			index = i
		}
	}
	if index == -1 {
		return errors.New("Url to be deleted was not found in the slice")
	}
	// deleting from list efficiently
	mnt.Urls[index], mnt.Urls[len(mnt.Urls)-1] = mnt.Urls[len(mnt.Urls)-1], mnt.Urls[index]
	mnt.Urls = mnt.Urls[:len(mnt.Urls)-1]
	return nil
}

// AddUrl appends a slice of Urls to the current list of Urls
func (mnt *Monitor) AddUrl(Urls []model.Url) {
	mnt.Urls = append(mnt.Urls, Urls...)
}

// Cancel stops all tasks of fetching Urls
// it will wait for current running jobs to finish
// note that if you call this method, for reusing the monitor
// you need to instantiate it again.
func (mnt *Monitor) Cancel() error {
	mnt.wp.Stop()
	if !mnt.wp.Stopped() {
		return errors.New("could not stop monitor")
	}
	return nil
}

// DoUrl checks a single Url's response and saves it's request into database
func (mnt *Monitor) DoUrl(Url model.Url) {
	var wg sync.WaitGroup
	wg.Add(1)
	mnt.wp.Submit(func() {
		defer wg.Done()
		mnt.monitorUrl(Url)
	})
	wg.Wait()
}

// Do ranges over Urls currently inside Monitor instance
// and save each one's request inside database
// this function does not block
func (mnt *Monitor) Do() {
	var wg sync.WaitGroup

	for UrlIndex := range mnt.Urls {
		Url := mnt.Urls[UrlIndex]
		wg.Add(1)
		mnt.wp.Submit(func() {
			defer wg.Done()
			mnt.monitorUrl(Url)
		})
	}
	wg.Wait()
}

func (mnt *Monitor) monitorUrl(Url model.Url) {
	// sending request
	req, err := Url.SendRequest()
	if err != nil {
		fmt.Println(err, "could not make request")
		req = new(model.Request)
		req.UrlId = Url.ID
		req.Result = http.StatusBadRequest
	}
	// add request to database
	if err = mnt.store.AddRequest(req); err != nil {
		fmt.Println(err, "could not save request to database")
	}
	// status code was other than 2XX
	if req.Result/100 != 2 {
		if err = mnt.store.IncrementFailed(&Url); err != nil {
			fmt.Println(err, "could not increment failed times for Url")
		}
	}
}
