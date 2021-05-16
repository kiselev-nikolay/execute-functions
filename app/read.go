package main

import (
	"container/list"
	"encoding/json"
	"fmt"
)

type headerUpdate struct {
	Key   string
	Value string
}

type request struct {
	Method string
	Define *define
	Link   *link
}

type link struct {
	Protocol string
	URL      string
}

type define struct {
	Keyword   string
	Variables *variables
}

type variables struct {
	Key  string
	From string
}

func maybeAssign(x interface{}, dest *string) {
	protocol, ok := x.(string)
	if ok {
		*dest = protocol
	}
}

func readRequest(cmd map[string]interface{}) *request {
	r := &request{
		Method: cmd["action"].(string),
		Define: &define{},
		Link:   &link{},
	}

	cmdLink, ok := cmd["link"].(map[string]interface{})
	if ok {
		maybeAssign(cmdLink["protocol"], &r.Link.Protocol)
		maybeAssign(cmdLink["url"], &r.Link.URL)
	}

	cmdDefine, ok := cmd["define"].(map[string]interface{})
	if ok {
		maybeAssign(cmdDefine["keyword"], &r.Define.Keyword)
	}

	cmdDefineVariables, ok := cmdDefine["variables"].(map[string]interface{})
	if ok {
		r.Define.Variables = &variables{}
		maybeAssign(cmdDefineVariables["key"], &r.Define.Variables.Key)
		maybeAssign(cmdDefineVariables["from"], &r.Define.Variables.From)
	}

	return r
}

func readHeaderUpdate(cmd map[string]interface{}) *headerUpdate {
	su := &headerUpdate{}
	maybeAssign(cmd["key"], &su.Key)
	maybeAssign(cmd["value"], &su.Value)
	return su
}

func token(t map[string]interface{}) *list.List {
	cmds := list.New()
	cmd := t["cmd"].(map[string]interface{})
	switch cmd["action"] {
	case "get", "post":
		cmds.PushBack(readRequest(cmd))
	case "set header":
		cmds.PushBack(readHeaderUpdate(cmd))
	}
	next, ok := t["next"]
	if !ok {
		return cmds
	}
	cmds.PushBackList(token(next.(map[string]interface{})))
	return cmds
}

func read(data string) *list.List {
	r := make(map[string]interface{})
	unmarshalErr := json.Unmarshal([]byte(data), &r)
	if unmarshalErr != nil {
		fmt.Println(unmarshalErr)
	}
	cmds := token(r)
	return cmds
}
