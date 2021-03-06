package main

import (
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"regexp"
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

type simpleKV map[string]string

func (s simpleKV) match(r string) simpleKV {
	var searchValue = regexp.MustCompile(r)
	mkv := make(simpleKV)
	
	for k, v := range s {
		if searchValue.MatchString(k) || searchValue.MatchString(v) {
			mkv[k] = v
		}
	}
	return mkv
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

func getKVsFromServer(url string) []ConsulKV {
	_ = url
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	resp, err := http.Get("http://localhost:8500/v1/kv/?recurse=true")
	if err != nil {
		os.Exit(1)
	}

	defer resp.Body.Close()

	raw, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }

    var c []ConsulKV
    json.Unmarshal(raw, &c)
    return c
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

func decode(c []ConsulKV) simpleKV {
	skv := make(simpleKV)
	for _, kv := range c {
		d, err := base64.StdEncoding.DecodeString(kv.Value)
		if err == nil {
			skv[kv.Key] = fmt.Sprintf("%q", d)
		} else {
			skv[kv.Key] = kv.Value
		}
	}
	return skv
} 

func main() {
		
	if len(os.Args) >= 2 {		
		for k, v := range decode(getKVsFromServer("blah")).match(os.Args[1]) {
			fmt.Println("Key:", k)
			fmt.Println("Value:", v)
		}
	} else {
		fmt.Println("please enter a search string")
		os.Exit(1)
	}
}