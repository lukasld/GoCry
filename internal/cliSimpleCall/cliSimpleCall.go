package cliSimpleCall

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type oPCliCallerCommon interface {
    getSc()         simpleCall
    getNumLn()      int
}

type commonOPCliCall struct {
    // combines numLines and simpleCall
    // getSc and getNumln as methods
    sC          simpleCall
    numLn       int
}
type commandString struct {
    command         string
    flagsVals    []string
}

type simpleCall struct {
    cS              commandString
    tDMs            time.Duration
}
func (c *commonOPCliCall)getSc() simpleCall {
    return c.sC
}

func (c *commonOPCliCall)getNumLn() int {
    return c.numLn
}

/*
   - error handling (wrapping)
   - logging
*/

type oPCliCaller interface {
    InvokeCommand()                             error
    getCommonCall()                             oPCliCallerCommon
    handleCallRes(callResults)                  (bool, error)    // handle resulting call
    handleLinesRes(string, bool, *exec.Cmd)     (bool, error)    // handle resulting lines
}

/* args resulting from invokeCall */
type callResults struct {
    stdOut chan string
    stdErr chan string
    ctx context.Context
    cmd *exec.Cmd
}

func invokeCall(oPCli oPCliCaller) (error) {
    /* invokes a command given by a string */

    simpleCall := oPCli.getCommonCall().getSc()
    timeoutCtx, cancel := context.WithTimeout(
        context.Background(),
        time.Millisecond * simpleCall.tDMs)
    defer cancel()

    // the command
    //"op", "vault", "list")
    cmd := exec.CommandContext(timeoutCtx,
            simpleCall.cS.command, simpleCall.cS.flagsVals...)

    cmdReaderStdErr, _ := cmd.StderrPipe()
    cmdReaderStdOut, err := cmd.StdoutPipe()
    if err := cmd.Start(); err != nil{
        return fmt.Errorf("invokeCommand: Start cmd failed : %w", err)
    }

    scannerStdOut := bufio.NewScanner(cmdReaderStdOut)
    stdOut := make(chan string, oPCli.getCommonCall().getNumLn())

    scannerStdErr := bufio.NewScanner(cmdReaderStdErr)
    errOut := make(chan string, oPCli.getCommonCall().getNumLn())

    // starting go-routing from buffer
    go readOutput(scannerStdOut, stdOut)
    go readOutput(scannerStdErr, errOut)

    callRes := callResults{
        stdOut: stdOut,
        stdErr: errOut,
        cmd: cmd,
        ctx: timeoutCtx,
    }
    _, err = oPCli.handleCallRes(callRes)
    return err
}


type lineHandlefunc func(string, bool, *exec.Cmd ) (bool, error)
func selectFunc( cR callResults, lF lineHandlefunc ) (bool, error){
    /*
        picks between timeout and next-line
    */
    select {
    // closing the channel if timeout is hit
    case <- cR.ctx.Done():
        // TODO: POSIX only
        return false,
        fmt.Errorf("TimeoutErr: %v", cR.ctx.Err())
    case line, ok := <-cR.stdOut:
        // handling of what happens with the text
        isDone, err := lF(line, ok, cR.cmd)
        if isDone || err != nil {
            return isDone, err
        }
    case stdErrMsg, isErr := <-cR.stdErr:
        // handling of what happens with the text
        if isErr {
            return true, fmt.Errorf("1PW-StdErr: %v", stdErrMsg)
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

