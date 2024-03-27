package clitools

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"
)

/*
   - error handling (wrapping)
   - logging
*/


type oPCliCaller interface {
    invokeCommand()                             error
    getCommonCall()                             oPCliCallerCommon
    handleCallRes(callResults)                  (bool, error)    // handle resulting call
    handleLinesRes(string, bool, *exec.Cmd)     (bool, error)    // handle resulting lines
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
        picks between timeout and next-line
    */
    select {
    // closing the channel if timeout is hit
    case <- cR.ctx.Done():
        // TODO: POSIX only
        //cR.cmd.Process.Kill()
        //os.Exit(1)
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
