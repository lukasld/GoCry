package main

import (
	"log"
    "internal/cliCalls"
)

/*
1. Create a module that handles 1Password CLI commands
    1.1 Commands - Context with Timeout
    1.2 checks if 1Password and the 1PW-CLI is installed
    1.3 checks if vault exists
    - error handling (wrapping)
    - logging
*/


func main(){

    //res, err := clitools.NewOPCliCall([]string{"vault", "list"}, 1)
    err := cliCalls.NewOpLoginCall([]string{"signin"}, 1)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Printf("result: %v", res)

}
