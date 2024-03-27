package cliSimpleCall

import (
	"fmt"
	"os/exec"
)

// Struct - to hold read-lines
type oPCliCallRl struct {
    cOPCall     oPCliCallerCommon
    rL          []string
}

func (rlCall *oPCliCallRl)getCommonCall() oPCliCallerCommon {
    /* Returns the common Caller */
    return rlCall.cOPCall
}

func (rlCall *oPCliCallRl)InvokeCommand() error{
    // pointer to call Result
    _, err := invokeCall(rlCall)
    return err
}

func (rlCall *oPCliCallRl)handleCallRes(cR callResults) (bool, error){
    // do we need a for loop or simply a single line?
    for {
        isDone, err := selectFunc(cR, rlCall.handleLinesRes)
        if isDone || err != nil {
            return isDone, err
        }
    }
}

func (rlCall *oPCliCallRl)handleLinesRes(l string, ok bool, c *exec.Cmd) (bool, error){
    // handles the read lines
    isDone := !ok
    if ok {
        rlCall.rL = append(rlCall.rL, fmt.Sprintf("%v\n", l))
        return isDone, nil
    }
    // we reached the end of the command
    if err := c.Wait(); err != nil{
        return false, fmt.Errorf("handleLines: cmd.wait failed : %w", err)
    }
    return isDone, nil
}

func (rlCall *oPCliCallRl)GetReadLines() []string{
    return rlCall.rL
}

func NewOpCliCallRl(flagsVals []string, numLn int) (error, oPCliCallRl){
    rlCall := oPCliCallRl {
                cOPCall: &commonOPCliCall {
                numLn: numLn,
                sC: simpleCall {
                    cS: commandString{
                        command: "op",
                        flagsVals: flagsVals,
                    },
                    tDMs: 100,
                },
            },
            rL: []string{},
        }
    return nil, rlCall
}



