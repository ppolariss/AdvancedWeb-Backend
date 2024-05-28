package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

func AddChat(addChatRequest AddChatRequest) (chatID int) {
	return 0
}

//func AddChatAndRecord()

func AddRecord(addRecordRequest AddRecordRequest) {
	fmt.Println(addRecordRequest)
	var err error
	jsonData, err := json.Marshal(addRecordRequest)
	if err != nil {
		Logger.Error(err.Error())
		return
	}

	response, err := http.Post(viper.GetString(EnvUrl)+"/api/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	if response.StatusCode == 200 {
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			Logger.Error(err.Error())
		}
	}(response.Body)

	requestBody, err := io.ReadAll(response.Body)
	if err != nil {
		Logger.Error(err.Error())
		return
	}
	Logger.Error(string(requestBody))
	//fmt.Println(string(requestBody))
	//var addRecordResponse AddRecordResponse
	//err = json.Unmarshal(requestBody, &addRecordResponse)
	//if err != nil {
	//	Logger.Error(err.Error())
	//	return
	//}
}
