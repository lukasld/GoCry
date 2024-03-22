package main

import (
	"fmt"
	"log"
    "internal/clitools"
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
    res, err := clitools.NewOPCliCall([]string{"--help"}, 5)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("result: %v", res)

}
