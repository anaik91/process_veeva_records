package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var disable = "Objectlifecyclestateuseraction.study_person__clin.active_state__c.change_state_to_inactive_useraction3__c"
var enable = "Objectlifecyclestateuseraction.study_person__clin.inactive_state__c.change_state_to_active_useraction5__c"

func ModifyUserState(apiurl, apiversion, userid, action, sessionid string) (bool, string) {
	var state string
	if action == "enable" {
		state = enable
	} else {
		state = disable
	}
	url := apiurl + "/api/" + apiversion + "/vobjects/study_person__clin/" + userid + "/actions/" + state
	method := "POST"
	payload := strings.NewReader("")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", sessionid)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	defer res.Body.Close()

	out, err := ioutil.ReadAll(res.Body)
	//fmt.Println(string(out))
	if err != nil {
		fmt.Println(err)
		return false, string(out)
	} else {
		return true, string(out)
	}

}
