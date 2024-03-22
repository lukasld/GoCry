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
    1. isSignedInCall   (listen to account is not signed in, or equivalent)
    2. logInCall        (longer timeout)

    For each call:

    0. is op installed?
    0. does doesGoCryVault exist
    0. is Signed in?

    3. createNewItem

*/


/*
1. Create a module that handles 1Password CLI commands
    1.1 Commands - Context with Timeout
    1.2 checks if 1Password and the 1PW-CLI is installed
    1.3 checks if vault exists
    - error handling (wrapping)
    - logging
*/

type oPCliCaller interface {
    invokeCommand()                         error
    handleLines(string, bool, *exec.Cmd)    (bool, error)
}
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



// Struct - to hold read-lines
type oPCliCallRl struct {
    cOPCall     oPCliCallerCommon
    rL          []string
}

func (rlCall *oPCliCallRl)invokeCommand() error{
    // pointer to call Result
    _, err := invokeCall(rlCall.cOPCall, rlCall)
    return err
}

func (rlCall *oPCliCallRl)handleLines(l string, ok bool, c *exec.Cmd) (bool, error){
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

func invokeCall(oPCliC oPCliCallerCommon, oPCli oPCliCaller) (bool, error) {
    /* invokes a command given by a string */

    simpleCall := oPCliC.getSc()
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
    out := make(chan string, oPCliC.getNumLn())

    // starting go-routing from buffer
    go readOutput(scanner, out)

    args := callResults{
        out: out,
        cmd: cmd,
        ctx: timeoutCtx,
    }
    // handles the reading of out
    _, err = handleCall(oPCli, args)
    return true, err
}


func handleCall(oPCli oPCliCaller, cR callResults) (bool, error){

    for {
        select {
        // closing the channel if timeout is hit
        case <- cR.ctx.Done():
            cR.cmd.Process.Kill()
            return false,
            fmt.Errorf("readLines: CLI - timed out: %w", cR.ctx.Err())
        case line, ok := <-cR.out:
            // handling of what happens with the text
            isDone, err := oPCli.handleLines(line, ok, cR.cmd)
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


