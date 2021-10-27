package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

var sse_conns = make(map[string]int)

var slock sync.Mutex

func getSSE(uri string, body *requestBody, ch chan string) error {
	slock.Lock()
	if _, ok := sse_conns[fmt.Sprintf("%s-%s", body.DeviceID, body.ApplicationID)]; ok {
		slock.Unlock()
		return nil
	}
	slock.Unlock()

	req, err := http.NewRequest("POST", uri, strings.NewReader(fmt.Sprintf(`{
		"identifiers":[{
			"device_ids":{
				"device_id":"%s",
				"application_ids":{"application_id":"%s"}
			}
		}]
	}`, body.DeviceID, body.ApplicationID)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", body.ApiKey))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	if err != nil {
		return fmt.Errorf("error NewRequest POST: %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error performing request for %s: %v", uri, err)
	}
	slock.Lock()
	sse_conns[fmt.Sprintf("%s-%s", body.DeviceID, body.ApplicationID)] = 1
	slock.Unlock()
	br := bufio.NewReader(res.Body)
	defer res.Body.Close()
	for {

		bs, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if len(bs) < 2 {
			continue
		}
		//fmt.Println(string(bs))
		go func() {
			ch <- string(bs)

		}()
		if err == io.EOF {
			fmt.Println("exit sse thread")
			break
		}
	}
	slock.Lock()
	delete(sse_conns, fmt.Sprintf("%s-%s", body.DeviceID, body.ApplicationID))
	slock.Unlock()
	return nil
}

func connSSE(body *requestBody, gid, uid string) {
	ch := make(chan string)
	go getSSE(fmt.Sprintf("%s/events", baseURL), body, ch)
	go func() {
		for {
			str := <-ch
			if msg, err := Decode(str); err == nil {
				if msg.Result.Name == "as.up.data.forward" && msg.Result.Data.Type == "type.googleapis.com/ttn.lorawan.v3.ApplicationUp" {
					bs, _ := json.Marshal(msg)
					go broadcast(gid, uid, string(bs))
					fmt.Println(string(bs))
				}
			} else {
				fmt.Println("decode error", err)
			}
		}
	}()
}
