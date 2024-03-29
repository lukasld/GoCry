package main

import (
	"fmt"
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

    /*
    if isLogIn, err := cliCalls.LoginCall(); err != nil{
        fmt.Println(err)
    }else{
        fmt.Printf("Login suceeded %v", isLogIn)
    }

    isAuth, err := cliCalls.IsUserAuthDCall()
    if err != nil{
        fmt.Println(err)
    }else{
        fmt.Println(isAuth)
    }

    isVault, err := cliCalls.DoesGoCryVaultExistCall()
    if err != nil{
        fmt.Println(err)
    }else{
        fmt.Println(isVault)
    }
    */

    isOPInst, err := cliCalls.IsOpCliInstalledCall()
    if err != nil{
        fmt.Println(err)
    }else{
        fmt.Println(isOPInst)
    }

    //log.Fatal(err)

}
