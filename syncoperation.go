package cacheorganizer

// RunAsync is a function that kicks off a goroutine that runs a SyncOperation and allows it to be cancelled
func (operation SyncOperation) RunAsync() (resultChan ResultChan, cancelChan CancelChan) {
    resultChan = make(ResultChan)
    cancelChan = make(CancelChan)
    innerResultChan := make(ResultChan)
    innerCancelChan := make(CancelChan)
    
    go func() {
        go func() {
            //Hopefully, the syncoperation actually makes use of the cancel chan- in which case, it'll return
            //nil if the cancel was received.  
            result := operation(innerCancelChan)
            
            if (result != nil) {
                select {
                    case innerResultChan <-result:
                        return;
                    case <-innerCancelChan:
                        return
                }
            }
        }()
        
        select {
            case result := <-innerResultChan:
                resultChan <- result;
            case <- cancelChan:
                innerCancelChan <- true
        }
    }()
    
    return
}

