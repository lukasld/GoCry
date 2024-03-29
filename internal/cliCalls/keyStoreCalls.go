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
    "fmt"
)

type storeItem struct {
    ts string   `json:"ts"`
    id int      `json:"id"`
}

type storeItemMap map[string]storeItem
type KSData struct {
    storeItem storeItemMap `json:"items"`
}

func unmarshalKS(pJs *string) *KSData{
    var kSD KSData
    err := json.Unmarshal([]byte(*pJs), &kSD)
    if err != nil {
        fmt.Println("Error:", err)
    }
    return &kSD
}



func getKeyStoreCall(){



}

func updateKeyStoreCall(){
}
