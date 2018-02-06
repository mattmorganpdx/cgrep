package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
)

// ConsulKV a struct to hold consul data
type ConsulKV struct {
	LockIndex   int    `json:"LockIndex"`
	Key         string `json:"Key"`
	Flags       int    `json:"Flags"`
	Value       string `json:"Value"`
	CreateIndex int    `json:"CreateIndex"`
	ModifyIndex int    `json:"ModifyIndex"`
}

func (c ConsulKV) toString() string {
	return toJSON(c)
}

func toJSON(c interface{}) string {
    bytes, err := json.Marshal(c)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    return string(bytes)
}

func getKVs() []ConsulKV {
	raw, err := ioutil.ReadFile("consul.kv")
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c []ConsulKV
    json.Unmarshal(raw, &c)
    return c
}

func main() {
	kvs := getKVs()
	for _, kv := range kvs {
		d, err := base64.StdEncoding.DecodeString(kv.Value)
		if err == nil {
			fmt.Printf("Key: %s Value: %q\n", kv.Key, d)
		}
	}
	
	

}