package clitools

import (
	"fmt"
	"os/exec"
    "time"
    "sync"
)


// Progress Time
type ProgressMessage struct {
    current      string
    rwlock       sync.RWMutex
}

func (p *ProgressMessage) Set(v string) {
    p.rwlock.Lock()
    defer p.rwlock.Unlock()
    p.current = v

}

func (p *ProgressMessage) Get() string {
   p.rwlock.RLock()
   defer p.rwlock.RUnlock()
   return p.current
}


// Struct - to hold read-lines
type oPCliCallLogin struct {
    cOPCall     oPCliCallerCommon
    chStep           <- chan time.Time
}

func (liCall *oPCliCallLogin)getCommonCall() oPCliCallerCommon {
    /* Returns the common Caller */
    return liCall.cOPCall
}

func (liCall *oPCliCallLogin)invokeCommand() error{
    // pointer to call Result

    pM := &ProgressMessage{current: "goCry: please authorize 1Pw Account"}
    go liCall.tickedProgressMsgRead(liCall.chStep, pM)

    isDone, err := invokeCall(liCall)
    pM.Set(fmt.Sprintf(`isDone: %v`, isDone))

    return err
}

func (liCall *oPCliCallLogin)handleCallRes(cR callResults) (bool, error){
    // do we need a for loop or simply a single line?
    for {
        isDone, err := selectFunc(cR, liCall.handleLinesRes)
        if isDone || err != nil {
            return isDone, err
        }
    }
}

func (liCall *oPCliCallLogin)handleLinesRes(l string, ok bool, c *exec.Cmd) (bool, error){

    // handles the read lines
    isDone := !ok
    if ok {
        /* what to do in the end  for login*/
        return isDone, nil
    }
    // we reached the end of the command
    if err := c.Wait(); err != nil{
        return false, fmt.Errorf("handleLines: cmd.wait failed : %w", err)
    }
    return isDone, nil
}

func (liCall *oPCliCallLogin) tickedProgressMsgRead(c <-chan time.Time, pM *ProgressMessage){
    // counts down until a given time
    startT := int((liCall.cOPCall.getSc().tDMs*1e6).Seconds())
    for {
        select {
        case <- c:
            fmt.Printf("%v, -%v \n", pM.current, startT)
            startT--
        }
    }
}
