package cliSimpleCall

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
type oPCliCallWaitWProgress struct {
    cOPCall     oPCliCallerCommon
    chStep      <- chan time.Time
    msg         string
}

func (liCall *oPCliCallWaitWProgress)getCommonCall() oPCliCallerCommon {
    /* Returns the common Caller */
    return liCall.cOPCall
}

func (liCall *oPCliCallWaitWProgress)InvokeCommand() error{
    // pointer to call Result
    pM := &ProgressMessage{current: liCall.msg}
    go liCall.tickedProgressMsgRead(liCall.chStep, pM)

    err := invokeCall(liCall)
    return err
}

func (liCall *oPCliCallWaitWProgress)handleCallRes(cR callResults) (bool, error){
    // do we need a for loop or simply a single line?
    for {
        isDone, err := selectFunc(cR, liCall.handleLinesRes)
        if isDone || err != nil {
            return isDone, err
        }
    }
}

func (liCall *oPCliCallWaitWProgress)handleLinesRes(l string, ok bool, c *exec.Cmd) (bool, error){

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

func (liCall *oPCliCallWaitWProgress) tickedProgressMsgRead(c <-chan time.Time, pM *ProgressMessage){
    // counts down until a given time
    startT := int((liCall.cOPCall.getSc().tDMs*1e6).Seconds())
    for {
        select {
        case <- c:
            fmt.Printf("%v, exit t-%vs \n", pM.current, startT)
            startT --
        }
    }
}


type OpWaitArgs struct{
    TDMs        int
    TickTS      int
    Msg         string
}

func NewOpCliCallWaitProgress(flagsVals []string, wA OpWaitArgs) (error, oPCliCallWaitWProgress){
    // check that tickTS < tDMS
    wPCall := oPCliCallWaitWProgress {
                cOPCall: &commonOPCliCall {
                numLn: 1,
                sC: simpleCall {
                    cS: commandString{
                        command: "op",
                        flagsVals: flagsVals,
                    },
                    tDMs: time.Duration(wA.TDMs),
                },
            },
            chStep: time.Tick(1 * (time.Second * time.Duration(wA.TickTS))),
            msg: wA.Msg,
    }
    return nil, wPCall
}
