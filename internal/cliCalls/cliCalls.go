package cliCalls

import (
    "fmt"
    "strings"
    "internal/cliSimpleCall"
)


func NewOPCliCall(flagsVals []string, numLn int) (string, error){
    /* takes in args and makes onepsw call */
    _, call := cliSimpleCall.NewOpCliCallRl(flagsVals, numLn)
    err := call.InvokeCommand()
    if err != nil {
        return "", fmt.Errorf("NewOPCliCall: Call Error : %w", err)
    }
    return strings.Join(call.GetReadLines(), ""), nil
}


func NewOpLoginCall(flagsVals []string, numLn int) error{
    /* takes in args and makes onepsw call */
    /* timout is a minute -> 60 000 ms by  */

    _, call := cliSimpleCall.NewOpCliCallWaitProgress(flagsVals, numLn)
    err := call.InvokeCommand()
    if err != nil {
        return fmt.Errorf("NewOPCliCall: Call Error : %w", err)
    }
    return nil
}


// Auth and Exist calls
func NewOpIsAuthDCall(){
    /* is the user authorized ?*/
    /* call whoami */
}

func NewOpDoesVaultExistCall(){
    /* check if GoCry Vault exist */
    /* */
}

func NewOpDoesItemExistCall(){
    /* check if GoCry Vault exists */
    /* */
}

func NewIsOpInstalledCall(){
    /*  */
    /* */
}



// Create Items
func CreateGoCryVaultAndKeyStoreCall(){
    /* create a vault called Gocry and Keystore */
}

func CreateGoCryItemCall(){
    /*
        1. check if vault exists
        3. pull the KeyStore
            3.a. create uuidv4
            3.b. update KeyStore
        4. create Entry and Psw in 1Password
        5. retrieve password
        6. encrypt file w. password
        7. put Keystore
    */
}

// Delete Items
func DeleteGoCryItemCall(){
    /*  */
    /*  */
}



// Update Item KeyStore
func UpdateKeyStoreCall(){
    /* Memguard */
    /* add / delete uuidv4: id */
}

func GetKeyStoreItemCall(){
    /* Memguard */
    /* get uuid */
}
