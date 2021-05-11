package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

func main() {
	apiurl := flag.String("apiurl", "default", "Example :  https://abcd.veevavault.com")
	apiversion := flag.String("apiversion", "default", "Example :  v20.3")
	username := flag.String("username", "default", "Example :   abcd@xyz.com ")
	password := flag.String("password", "default", "Example :  P###SSw0rd ")
	action := flag.String("action", "enable", "Example :  enable OR disable")
	flag.Parse()
	start := time.Now()
	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()

	session := SessionRespose{}
	getSessionId(*apiurl, *apiversion, *username, *password, &session)
	fmt.Println(session.SessionId)
	csv_data := readfile("users.csv")
	csv_data_list := strings.Split(strings.ReplaceAll(csv_data, "\r\n", "\n"), "\n")

	var user_ids []string
	for _, v := range csv_data_list[1:] {
		eachuserinfo := strings.Split(v, ",")
		if len(eachuserinfo) > 1 {
			user_ids = append(user_ids, strings.ReplaceAll(strings.TrimSpace(eachuserinfo[1]), "\"", ""))
		}
	}
	logFileName := "run.log"
	f, _ := initalizeLogFile(logFileName)
	defer f.Close()
	log.SetOutput(f)
	wg := sync.WaitGroup{}
	for _, user := range user_ids {
		wg.Add(1)
		go func(apiurl, apiversion, userid, action, sessionid string) {
			_, output := ModifyUserState(apiurl, apiversion, userid, action, sessionid)
			log.Println("Finished Processing User", userid, "With status :", output)
			wg.Done()
		}(*apiurl, *apiversion, user, *action, session.SessionId)
	}
	wg.Wait()
}
