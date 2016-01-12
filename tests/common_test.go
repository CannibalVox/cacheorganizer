package cacheorganizer_test

import "time"
import "github.com/cannibalvox/cacheorganizer"

func buildDataSource(getResult func() int, writeResult cacheorganizer.WriteOperation, pauseDuration time.Duration, timeout time.Duration) cacheorganizer.DataSource {
    retrieve := func() (resultChan cacheorganizer.ResultChan, cancelChan cacheorganizer.CancelChan) {
        resultChan = make(cacheorganizer.ResultChan)
        cancelChan = make(cacheorganizer.CancelChan)
        
        go func() {
            select {
                case <-cancelChan:
                    return
                case <-time.After(pauseDuration):
                    resultChan <- &cacheorganizer.Result{Value:getResult(), Error:nil}
            }
        }()
        
        return
    }
    
    return cacheorganizer.DataSource{Retrieve: retrieve, Timeout: timeout, Write: writeResult}
}

func nilDataSource(writeResult cacheorganizer.WriteOperation, pauseDuration time.Duration, timeout time.Duration) cacheorganizer.DataSource {
        retrieve := func() (resultChan cacheorganizer.ResultChan, cancelChan cacheorganizer.CancelChan) {
        resultChan = make(cacheorganizer.ResultChan)
        cancelChan = make(cacheorganizer.CancelChan)
        
        go func() {
            select {
                case <-cancelChan:
                    return
                case <-time.After(pauseDuration):
                    resultChan <- nil
            }
        }()
        
        return
    }
    
    return cacheorganizer.DataSource{Retrieve: retrieve, Timeout: timeout, Write: writeResult}
}

func writeTo(values []int, index int) cacheorganizer.WriteOperation {
    return func(result *cacheorganizer.Result) {
        values[index] = result.Value.(int)
    }
}

func readFrom(values []int, index int) func() int {
    return func() int {
        return values[index]
    }
}

func runBasicTest(values []int, nilValues int, timeouts []time.Duration, pauses []time.Duration) *cacheorganizer.Result {
    dataSources := make([]cacheorganizer.DataSource, len(timeouts))
    
    for i := 0; i < len(timeouts); i++ {
        if i < nilValues {
            dataSources[i] = nilDataSource(writeTo(values, i), pauses[i], timeouts[i])
        } else {
            dataSources[i] = buildDataSource(readFrom(values, i), writeTo(values, i), pauses[i], timeouts[i])
        }
    }
    
    resultChan, _ := cacheorganizer.Retrieve(dataSources...)
    time.Sleep(100 * time.Millisecond)
    return <-resultChan
}