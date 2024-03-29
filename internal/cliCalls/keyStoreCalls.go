package cliCalls

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

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type Item struct {
    Ts string   `json:"ts"`
    Id string   `json:"id"`
}
//type StoreItemMap map[string]StoreItem
//items   string  `json:"items"`

type KSData struct {
    Items map[string]Item `json:"items"`
}


func unmarshalKS(pJs *string) (*KSData, error){
    var kSD KSData
    err := json.Unmarshal([]byte(*pJs), &kSD)
    if err != nil {
        return nil, fmt.Errorf("KSUnmarshalErr", err)
    }
    return &kSD, nil
}

func cleanKsString(onePwStr *string) (error){
    // there seems to be a bug/inconsistency how 1pw-cli handles quotes - this func is a workaround:
    //https://1password.community/discussion/128123/cli-v2-returns-fields-surrounded-in-double-quotes-v1-didn-t
    *onePwStr = strings.ReplaceAll(*onePwStr, `""`, `"`)
    length := len(*onePwStr)
    if length > 2 {
        *onePwStr = (*onePwStr)[1:length-2]
    } else {
        return errors.New("KSCleanErr:")
    }
    return nil
}

func GetKeyStoreCall() (*KSData, error){
    if _, err := IsUserAuthDCall(); err != nil {
        return nil, err
    }
    args := oPReadArgs{
        flagsVals:  []string{"item", "get", "GoCry-KeyStore",
            "--vault", "GoCry-Vault",
            "--fields", "label=keystore",
        },
        tDMs:       800,
        numLn:      10,
    }
    //TODO: memguarded
    kSs, err := oPReadLines(args)
    if err != nil {
        return nil, fmt.Errorf("KSErr: %v", err)
    }
    // cleans the kSs inplace
    if err := cleanKsString(&kSs); err != nil {
        return nil, fmt.Errorf("GetKSCallErr: %w", err)
    }
    //TODO: memguarded
    kSD, err := unmarshalKS(&kSs)
    if err != nil {
        return nil, fmt.Errorf("KSErr: %v", err)
    }
    return kSD, nil
}


func updateKeyStoreCall(){
}
