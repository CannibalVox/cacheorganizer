/*
Package cacheorganizer provides an async multi-level cache coordinator
*/
package cacheorganizer

import "time"

// Value is an arbitrary obj type to be returned from/sent to operations
type Value interface{}
// Error is an arbitrary error type to be returned from operations if there was a problem
type Error error

// Result is a type that contains both a return value and a (probably-nil) error, returned from operations.
type Result struct {
    Value Value
    Error Error
}

const (
    cacheEventResult = iota
    cacheEventTimeout = iota
)

type event struct {
    EventID int
    SourceID int
    Result *Result
}

func newEvent(eventID int, sourceID int, result *Result) *event {
    return &event{EventID: eventID, SourceID: sourceID, Result:result}
}

// ResultChan is a channel of result objects, returned from methods that kick off async operations
type ResultChan chan *Result
// CancelChan is a channel of bools, returned from methods that kick off cancellable async operations
type CancelChan chan bool
type eventChan chan *event

// SyncOperation is a function that accepts a CancelChan to allow it to be cancelled in progress and returns a result.
// It's used to allow simple funcs from consumers to be converted to async.
// If you're implementing a SyncOperation, you should return nil for the result if you caught cancel to signal
// to the caller that cancel ocurred.
type SyncOperation func(cancelChan CancelChan) *Result
// AsyncOperation is a function that kicks off an async process and returns the resultchan & cancelchan used to control it.
type AsyncOperation func() (resultChan ResultChan, cancelChan CancelChan)
// WriteOperation is a function that writes values from deeper cache layers back to the upper layers
type WriteOperation func(result *Result)

// DataSource is a type that represents a single level of a cache.  An asynchronous operation and a timeout before
// the coordinator moves on to the next level of the cache.
type DataSource struct {
    Retrieve AsyncOperation
    Timeout time.Duration
    Write WriteOperation
}
