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
1. Create a module that handles 1Password CLI commands
    1.1 Commands - Context with Timeout
    1.2 checks if 1Password and the 1PW-CLI is installed
    1.3 checks if vault exists
    - error handling (wrapping)
    - logging
*/
type oPCliCaller interface {
    getSc()                                 simpleCall
    getNumLn()                              int
    invokeCommand()                         error
    handleLines(string, bool, *exec.Cmd)    (bool, error)
}

type commandString struct {
    command         string
    flagsVals    []string
}

type simpleCall struct {
    cS              commandString
    tDMs            time.Duration
}

// type return-Lines
// implements interface oPCliCaller
type oPCliCallRl struct {
    sC          simpleCall
    numLn       int
    rL          []string
}

func (rlCall *oPCliCallRl)invokeCommand() error{
    // TODO: error handle
    _, err := invokeAndHandle(rlCall)
    return err
}

func (rlCall *oPCliCallRl)handleLines(line string, ok bool, cmd *exec.Cmd) (bool, error){
    isDone := !ok
    if ok {
        rlCall.rL = append(rlCall.rL, fmt.Sprintf("%v\n", line))
        return isDone, nil
    }
    // we reached the end of the command
    if err := cmd.Wait(); err != nil{
        return false, fmt.Errorf("readLines: cmd.wait failed : %w", err)
    }
    return isDone, nil
}

func (rl *oPCliCallRl)getSc() simpleCall {
    return rl.sC
}

func (rl *oPCliCallRl)getNumLn() int {
    return rl.numLn
}



// typfunce waits for Login to succeed
// implements interface oPCliCaller

type oPCliCallLogin struct {
    sC          simpleCall
    numLn       int
    callLogin   func()
}

func (li *oPCliCallLogin)getSc() simpleCall {
    return li.sC
}

func (li *oPCliCallLogin)invokeCommand() {
    return
}

func (li *oPCliCallLogin)getNumLn() int {
    return li.numLn
}


func invokeAndHandle(oPCli oPCliCaller) (bool, error) {
    /* invokes a command given by a string */

    simpleCall := oPCli.getSc()

    timeoutCtx, cancel := context.WithTimeout(
        context.Background(),
        time.Millisecond * simpleCall.tDMs)
    defer cancel()

    // the command
    cmd := exec.CommandContext(timeoutCtx, simpleCall.cS.command, simpleCall.cS.flagsVals...)
    cmdReader, err := cmd.StdoutPipe()

    if err != nil {
        return false,
        fmt.Errorf("invokeCommand: Stdout Execution failed : %w", err)
    }

    if err := cmd.Start(); err != nil{
        return false,
        fmt.Errorf("invokeCommand: Start cmd failed : %w", err)
    }

    scanner := bufio.NewScanner(cmdReader)
    out := make(chan string, oPCli.getNumLn())

    // starting go-routing from buffer
    go readOutput(scanner, out)

    for {
        select {
        // closing the channel if timeout is hit
        case <- timeoutCtx.Done():
            return false,
            fmt.Errorf("readLines: CLI - timed out: %w", timeoutCtx.Err())
        case line, ok := <-out:
            // handling of what happens with the text
            //oPCli.invokeCommand(line, ok)
            isDone, err := oPCli.handleLines(line, ok, cmd)

            fmt.Printf("isDone: %v \n", isDone)
            if isDone || err != nil {
                return isDone, err
            }
        }
    }
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
        sC: simpleCall {
            cS: commandString{
                command: "op",
                flagsVals: flagsVals,
            },
            tDMs: 200,
        },
        numLn: numLn,
        rL: []string{},
    }
    err := call.invokeCommand()
    if err != nil {
        return "", fmt.Errorf("NewOPCliCall: Call Error : %w", err)
    }
    return strings.Join(call.rL, ""), nil
}
