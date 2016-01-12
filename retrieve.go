package cacheorganizer

import "time"

func createOpProxy(eventChan eventChan, sourceID int, dataSource DataSource) CancelChan {
    cancel := make(CancelChan)
    
    go func() {
        resultChan, cancelChan := dataSource.Retrieve()
        
        for {
            timeout := time.After(dataSource.Timeout)
            select {
                case result := <-resultChan:
                    eventChan <- newEvent(cacheEventResult, sourceID, result)                    
                    return
                case <-cancel:
                    cancelChan <- true
                    return
                case <-timeout: {
                    eventChan <- newEvent(cacheEventTimeout, sourceID, nil)
                }
            }
        }
    }()
    
    return cancel
}

func cancelAll(cancelChans []CancelChan) {
    for i := 0; i < len(cancelChans); i++ {
        if cancelChans[i] != nil {
            cancelChans[i] <- true
        }
    }
}

func cancelAllBut(cancelChans []CancelChan, exception int) {
    for i := 0; i < len(cancelChans); i++ {
        if i != exception && cancelChans[i] != nil {
            cancelChans[i] <- true
        }
    }
}

func writeValues(dataSources []DataSource, event *event) {
    if event.Result == nil || event.Result.Error != nil {
        return
    }
    
    for i := 0; i < event.SourceID; i++ {
        go func(dataSourceIndex int) {
            dataSources[dataSourceIndex].Write(event.Result)
        }(i);        
    }
}

func acceptEvent(resultChan ResultChan, cancelChans []CancelChan, dataSources []DataSource, event *event) {
    cancelAllBut(cancelChans, event.SourceID)
    writeValues(dataSources, event)
    resultChan <- event.Result
}

func startNextSource(resultChan ResultChan, eventChan eventChan, cancelChans []CancelChan, dataSources []DataSource, event *event) bool {
    if event.SourceID + 1 < len(dataSources) && cancelChans[event.SourceID + 1] == nil {
        cancelChans[event.SourceID + 1] = createOpProxy(eventChan, event.SourceID + 1, dataSources[event.SourceID + 1])
    } else if event.SourceID + 1 >= len(dataSources) {
        acceptEvent(resultChan, cancelChans, dataSources, event)
        return true
    }
    
    return false
}

// Retrieve is a function that takes in some number of data sources, each one representing a different cache layer.
// This method kicks off the first data source- if the first data source times out or returns a nil result or error
// result, the next data source will be kicked off.  The non-nil/non-error result will be sent to the result channel.
// Any result (including nil/error) from the final data source will be sent to the result channel as soon as it is received.
// If the final data source times out, nil will be sent to the result channel.  A data source can supply the eventual response
// even if it has timed out, if it returns the first valid result before the final layer times out. 
// Once a valid value is retrieved, it will be written to all prior data sources (if any).
func Retrieve(dataSources ...DataSource) (resultChan ResultChan, cancelChan CancelChan) {
    resultChan = make(ResultChan)
    cancelChan = make(CancelChan)
    
    go func() {
        sourceCount := len(dataSources)
        innerCancelChans := make([]CancelChan, sourceCount)
        eventChan := make(eventChan)
        
        innerCancelChans[0] = createOpProxy(eventChan, 0, dataSources[0])
        
        for {
            select {
                case event := <-eventChan:
                    switch event.EventID {
                        case cacheEventResult:
                            if event.Result != nil && (event.Result.Error == nil || event.SourceID + 1 >= len(dataSources)) {
                                acceptEvent(resultChan, innerCancelChans, dataSources, event)
                                return
                            }
                            
                            innerCancelChans[event.SourceID] = nil
                            startNextSource(resultChan, eventChan, innerCancelChans, dataSources, event)
                        case cacheEventTimeout:
                            startNextSource(resultChan, eventChan, innerCancelChans, dataSources, event)
                    }
                case <-cancelChan:
                    cancelAll(innerCancelChans)
                    return
            }
        }
    }()
    
    return
}