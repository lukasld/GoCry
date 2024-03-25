package clitools

import (
	"bufio"
	"context"
	"fmt"
    "time"
	"os/exec"
    "strings"
)


/*
    - error handling (wrapping)
    - logging
*/


// typfunc waits for Login to succeed
// implements interface oPCliCaller
type oPCliCallLogin struct {
    cOPCall     oPCliCallerCommon
}

func (li *oPCliCallLogin)invokeCommand() {
    return
}

func (li *oPCliCallLogin)handleLogin() (bool, error){
    /*  call OP signin
     *  call whoami repeatedly
     *      -> check if 'dismissed...' is not in the string anymore
     *      -> if URL:..., we are logged in,
     *          -> then we return true
     */
    return false, nil
}


/* args resulting from invokeCall */
type callResults struct {
    out chan string
    ctx context.Context
    cmd *exec.Cmd
}

func invokeCall(oPCli oPCliCaller) (bool, error) {
    /* invokes a command given by a string */

    simpleCall := oPCli.getCommonCall().getSc()
    timeoutCtx, cancel := context.WithTimeout(
        context.Background(),
        time.Millisecond * simpleCall.tDMs)
    defer cancel()

    // the command
    cmd := exec.CommandContext(timeoutCtx,
        simpleCall.cS.command, simpleCall.cS.flagsVals...)
    cmdReader, err := cmd.StdoutPipe()
    if err := cmd.Start(); err != nil{
        return false,
        fmt.Errorf("invokeCommand: Start cmd failed : %w", err)
    }

    scanner := bufio.NewScanner(cmdReader)
    out := make(chan string, oPCli.getCommonCall().getNumLn())

    // starting go-routing from buffer
    go readOutput(scanner, out)

    callRes := callResults{
        out: out,
        cmd: cmd,
        ctx: timeoutCtx,
    }
    _, err = oPCli.handleCallRes(callRes)
    return true, err
}



type lineHandlefunc func(string, bool, *exec.Cmd ) (bool, error)
func selectFunc( cR callResults, lF lineHandlefunc ) (bool, error){
    /*

    */
    select {
    // closing the channel if timeout is hit
    case <- cR.ctx.Done():
        return false,
        fmt.Errorf("readLines: CLI - timed out: %w", cR.ctx.Err())
    case line, ok := <-cR.out:
        // handling of what happens with the text
        isDone, err := lF(line, ok, cR.cmd)
        if isDone || err != nil {
            return isDone, err
        }
    }
    return false, nil
}

func readOutput(scanner *bufio.Scanner, out chan string){
    /* reads a maximum amount of lines into a channel */
    defer close(out)
    for scanner.Scan() {
        if len(out) == cap(out){
            break
        }
        out <- scanner.Text()
    }
}





func NewOPCliCall(flagsVals []string, numLn int) (string, error){
    /* takes in args and makes onepsw call */
    call := oPCliCallRl {
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

    err := call.invokeCommand()
    if err != nil {
        return "", fmt.Errorf("NewOPCliCall: Call Error : %w", err)
    }

    return strings.Join(call.rL, ""), nil
}
