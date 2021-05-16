package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func exec(cmds *list.List) {
	state := make(map[string]string)
	headers := make(http.Header)
	for e := cmds.Front(); e != nil; e = e.Next() {
		switch command := e.Value.(type) {
		case *request:
			switch command.Method {
			case "get":
				res, reqErr := http.Get("http://" + command.Link.URL)
				if reqErr != nil {
					fmt.Printf("error: %v\n", reqErr)
					return
				}
				data := make(map[string]interface{})
				dec := json.NewDecoder(res.Body)
				decErr := dec.Decode(&data)
				if decErr != nil {
					fmt.Printf("json decode error: %v\n", decErr)
					return
				}
				if command.Define != nil && command.Define.Variables != nil {
					v := ""
					if command.Define.Variables.From != "" {
						v = fmt.Sprint(data[command.Define.Variables.From])
					} else {
						v = fmt.Sprint(data[command.Define.Variables.Key])

					}
					state[command.Define.Variables.Key] = v
				}
			case "post":
				data := make(map[string]string)
				if command.Define != nil && command.Define.Variables != nil {
					v := ""
					if command.Define.Variables.From != "" {
						v = fmt.Sprint(state[command.Define.Variables.From])
					} else {
						v = fmt.Sprint(state[command.Define.Variables.Key])
					}
					data[command.Define.Variables.Key] = v
				}
				reqURL := "http://" + command.Link.URL
				req, reqErr := http.NewRequest("POST", reqURL, nil)
				if reqErr != nil {
					fmt.Printf("error: %v\n", reqErr)
					return
				}
				req.Header = headers
				switch command.Define.Keyword {
				case "json":
					body := bytes.NewBuffer([]byte{})
					enc := json.NewEncoder(body)
					encErr := enc.Encode(&data)
					if encErr != nil {
						fmt.Printf("json encode error: %v\n", encErr)
						return
					}
					req.Body = io.NopCloser(body)
					req.Header.Add("Content-Type", "application/json")
				case "params":
					form := url.Values{}
					for k, v := range data {
						form.Add(k, v)
					}
					req.Body = io.NopCloser(strings.NewReader(form.Encode()))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
				}
				_, resErr := http.DefaultClient.Do(req)
				if resErr != nil {
					fmt.Printf("server response error: %v\n", resErr)
					return
				}
			}
			break
		case *headerUpdate:
			headers.Add(command.Key, command.Value)
			break
		}
	}
}
