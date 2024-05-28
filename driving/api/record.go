package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"src/config"
	"src/schemas"
	"src/utils"
)

func AddChat(addChatRequest schemas.AddChatRequest) (chatID int) {
	return 0
}

//func AddChatAndRecord()

func AddRecord(addRecordRequest schemas.AddRecordRequest) {
	fmt.Println(addRecordRequest)
	var err error
	jsonData, err := json.Marshal(addRecordRequest)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	response, err := http.Post(viper.GetString(config.EnvUrl)+"/api/records", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if response.StatusCode == 200 {
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}(response.Body)

	requestBody, err := io.ReadAll(response.Body)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	utils.Logger.Error(string(requestBody))
	//fmt.Println(string(requestBody))
	//var addRecordResponse AddRecordResponse
	//err = json.Unmarshal(requestBody, &addRecordResponse)
	//if err != nil {
	//	Logger.Error(err.Error())
	//	return
	//}
}
