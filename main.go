package main

import (
	"encoding/json"
	"fmt"
)

var (
	policy    map[string]string = make(map[string]string, 0)
	principal                   = "user:123"
	group                       = "group:123"
	actions                     = []string{"allow", "deny"}
	subject                     = "user:456"
	verb                        = []string{"read", "write", "delete", "update", "create", "list", "execute", "admin"}
)

func buildPolicy() {
	// [requesterID] = "subject:subjectid:subjectgroup:groupid:[verb]:action"
	policy[principal] = fmt.Sprintf("%s:%s", subject, verb)
}

func main() {
	buildPolicy()
	js, err := json.Marshal(policy)
	if err != nil {
		fmt.Println("Error marshaling policy:", err)
		return
	}
	fmt.Println(string(js))
}
