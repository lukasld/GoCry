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
    invokeAndHandle(rlCall)
    return nil
}

func (rlCall *oPCliCallRl)handleLines(line string, ok bool, cmd *exec.Cmd) error{
    if !ok {
        // we reached the end of the line
        if err := cmd.Wait(); err != nil{
            return fmt.Errorf("readLines: cmd.wait failed : %w", err)
        }
        return nil
    }
    rlCall.rL = append(rlCall.rL, fmt.Sprintf("%v\n", line))
    fmt.Println(rlCall.rL)
    return nil
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

// TODO: Generics
func (li *oPCliCallLogin)getSc() simpleCall {
    return li.sC
}

func (li *oPCliCallLogin)invokeCommand() {
    return
}

func (li *oPCliCallLogin)getNumLn() int {
    return li.numLn
}




type oPCliCaller interface {
    getSc()                                 simpleCall
    getNumLn()                              int
    invokeCommand()                         error
    handleLines(string, bool, *exec.Cmd)    error
}

//resHandleFunc *func) (error) {

func invokeAndHandle(oPCli oPCliCaller) error {
    /* invokes a command given by a string */

    ctx := context.Background()
    simpleCall := oPCli.getSc()
    timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond * simpleCall.tDMs)
    defer cancel()

    // the command
    cmd := exec.Command(simpleCall.cS.command, simpleCall.cS.flagsVals...)
    cmdReader, err := cmd.StdoutPipe()

    if err != nil {
        return fmt.Errorf("invokeCommand: Stdout Execution failed : %w", err)
    }

    if err := cmd.Start(); err != nil{
        return fmt.Errorf("invokeCommand: Start cmd failed : %w", err)
    }

    scanner := bufio.NewScanner(cmdReader)
    out := make(chan string, oPCli.getNumLn())

    // starting go-routing from buffer
    go readOutput(scanner, out)

    // here we need to pass in the readline or not
    /*
    if err := readLines(oPCli, timeoutCtx, out, cmd); err != nil {
        return fmt.Errorf("invokeCommand: Start cmd failed : %w", err)
    }
    */

     for {
        select {
        // closing the channel if timeout is hit
        case <- timeoutCtx.Done():
            return fmt.Errorf("readLines: CLI - timed out: %w",
                timeoutCtx.Err())
        case line, ok := <-out:
            // handling of what happens with the text
            //oPCli.invokeCommand(line, ok)
            oPCli.handleLines(line, ok, cmd)
        }
    }

}

/*
func readLines(oPtOCtx context.Context, out <-chan string, cmd *exec.Cmd) (error){

    //var returnedLines []string

    for {
        select {
        // closing the channel if timeout is hit
        case <- tOCtx.Done():
            return fmt.Errorf("readLines: CLI - timed out: %w", tOCtx.Err())
        case line, ok := <-out:


            //sdtOutHandler(ok, ln *string, rL *[]string, cmd)

            /*
            if !ok {
                // we reached the end of the line
                if err := cmd.Wait(); err != nil{
                    return nil, fmt.Errorf("readLines: cmd.wait failed : %w", err)
                }

                // TODO: here we could inject a function that handles this differently

                return returnedLines, nil

            }

            returnedLines = append(returnedLines, fmt.Sprintf("%v\n", line))
        }
    }
}
*/

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
            tDMs: 1000,
        },
        numLn: numLn,
        rL: []string{},
    }
    err := call.invokeCommand()
    if err != nil {
        return "", fmt.Errorf("NewOPCliCall: Call Error : %w", err)
    }
    fmt.Println("here")
    fmt.Println(call.rL)
    return strings.Join(call.rL, ""), nil
}
