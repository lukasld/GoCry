package cliCalls

import (
    "fmt"
    "strings"

    "internal/cliSimpleCall"
)

/* TODO:
    - implement funcs
    - User Errors
    - Messages
*/

type oPReadArgs struct {
    flagsVals   []string
    tDMs        int
    numLn       int
}

func oPReadLines(args oPReadArgs) (string, error){
    /* takes in args and makes onepsw call */
    _, call := cliSimpleCall.NewOpCliCallRl(
        args.flagsVals,
        args.tDMs,
        args.numLn,
    )
    err := call.InvokeCommand()
    if err != nil {
        return "", fmt.Errorf("InvokeErr: %v", err)
    }
    return strings.Join(call.GetReadLines(), ""), err
}


func LoginCall() (bool, error) {
    /* takes in args and makes onepsw call */
    /* timout is a minute -> 60 000 ms by  */
    flagsVals := []string{"signin"}
    args := cliSimpleCall.OpWaitArgs{
        TDMs:        60000,
        TickTS:      1,
        Msg:         "goCry: please authorize 1Pw Account",
    }
    _, call := cliSimpleCall.NewOpCliCallWaitProgress(flagsVals, args)
    err := call.InvokeCommand()
    if err != nil {
        return false, fmt.Errorf("InvokeErr: %v", err)
    }
    return true, nil
}

// Auth and Exist calls
func IsUserAuthDCall() (bool, error){
    /* is the user authorized ?*/
    /* TODO String - Command String Store or so */
    args := oPReadArgs{
        flagsVals:  []string{"whoami"},
        tDMs:       100,
        numLn:      5,
    }
    lns, err := oPReadLines(args)
    switch{
        case err != nil:
            return false, fmt.Errorf("AuthErr : %v", err)
        case strings.Contains(lns, "account is not signed in"):
            return false, nil
        default:
            return true, nil
    }
}

func DoesGoCryVaultExistCall() (bool, error){
    /* check if GoCry Vault exist */
    /* */
    if _, err := IsUserAuthDCall(); err != nil {
        return false, err
    }
    args := oPReadArgs{
        flagsVals:  []string{"vault", "get", "GoCry-Vault"},
        tDMs:       500,
        numLn:      1,
    }
    _, err := oPReadLines(args)
    switch{
        case err != nil:
            return false, fmt.Errorf("VaultExistErr: %v", err)
        default:
            return true, nil
    }
}


func IsOpCliInstalledCall() (bool, error){
    /*  */
    args := oPReadArgs{
        flagsVals:  []string{"--help"},
        tDMs:       100,
        numLn:      1,
    }
    _, err := oPReadLines(args)
    switch{
        case err != nil:
            return false, fmt.Errorf("1PWCliErr: %v", err)
        default:
            return true, nil
    }
}

// VaultCalls

func SetupVaultAndKeystoreCall(){
}

// Keystore
/*
{
    "items" : {
        <uuid>: {
            "ts": "YYYY-MM-DD HH:MM:SS",
            "id": <1pwId>
        },
        ...
    }
}
*/

func getKeyStoreCall(){
    // unmarshal json


}

func updateKeyStoreCall(){
}





func DoesItemExistCall() (bool, error){

    if _, err := IsUserAuthDCall(); err != nil {
        return false, err
    }

    return false, nil
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
