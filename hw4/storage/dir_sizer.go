package storage

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Result represents the Size function result
type Result struct {
	// Total Size of File objects
	Size int64
	// Count is a count of File objects processed
	Count int64
}

type DirSizer interface {
	// Size calculate a size of given Dir, receive a ctx and the root Dir instance
	// will return Result or error if happened
	Size(ctx context.Context, d Dir) (Result, error)
}

// sizer implement the DirSizer interface
type sizer struct {
	// maxWorkersCount number of workers for asynchronous run
	maxWorkersCount int
	countMadeDir    *int64
	countDir        *int64
	err             error
	dirCh           chan Dir
	fileCh          chan File
	endCh           chan struct{}
}

// NewSizer returns new DirSizer instance
func NewSizer() DirSizer {
	return &sizer{
		maxWorkersCount: 10,
		countMadeDir:    new(int64),
		countDir:        new(int64),
		dirCh:           make(chan Dir, 10),
		fileCh:          make(chan File),
		endCh:           make(chan struct{}, 10),
	}
}

func (a *sizer) Size(ctx context.Context, d Dir) (Result, error) {
	res := Result{Size: 0, Count: 0}
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	for i := 1; i <= a.maxWorkersCount; i++ {
		wg.Add(1)
		go a.dirProcessor(ctx, wg, mutex)
	}
	wg.Add(1)
	go a.fileProcessor(ctx, wg, mutex, &res)
	a.dirCh <- d
	atomic.AddInt64(a.countDir, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			a.endCh <- struct{}{}
		case <-a.endCh:
			a.endCh <- struct{}{}
			return
		}
	}()
	wg.Wait()
	a.close()
	if a.err != nil {
		return Result{}, a.err
	}
	return res, nil
}

func (a *sizer) close() {
	close(a.dirCh)
	close(a.fileCh)
	close(a.endCh)
}

func (a *sizer) dirProcessor(ctx context.Context, wg *sync.WaitGroup, mutex *sync.Mutex) {
	defer wg.Done()
	for {
		select {
		case dir, ok := <-a.dirCh:
			if !ok {
				return
			}
			dirSlice, fileSlice, err := dir.Ls(ctx)
			if err != nil {
				mutex.Lock()
				a.err = fmt.Errorf("dirProcessor : %w", err)
				mutex.Unlock()
				a.endCh <- struct{}{}
				return
			}
			atomic.AddInt64(a.countDir, int64(len(dirSlice)))
			for _, d := range dirSlice {
				a.dirCh <- d
			}
			for _, file := range fileSlice {
				a.fileCh <- file
			}
			atomic.AddInt64(a.countMadeDir, 1)
			if *a.countMadeDir == *a.countDir {
				a.endCh <- struct{}{}
				return
			}
		case <-a.endCh:
			a.endCh <- struct{}{}
			return
		}
	}
}

func (a *sizer) fileProcessor(ctx context.Context, wg *sync.WaitGroup, mutex *sync.Mutex, res *Result) {
	defer wg.Done()
	for {
		select {
		case file, ok := <-a.fileCh:
			if !ok {
				return
			}
			mutex.Lock()
			if a.err != nil {
				return
			}
			mutex.Unlock()
			size, err := file.Stat(ctx)
			if err != nil {
				mutex.Lock()
				a.err = fmt.Errorf("fileProcessor : %w", err)
				mutex.Unlock()
				a.endCh <- struct{}{}
				return
			}
			res.Count++
			res.Size += size
		case <-a.endCh:
			a.endCh <- struct{}{}
			return
		}
	}
}
