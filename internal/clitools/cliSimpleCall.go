package clitools

import (
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
