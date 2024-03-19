package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
	"time"
)

func TestModels(t *testing.T) {
	now := time.Now()
	fNow := now.Format("2006-01-02 15:04:05")
	c := Commodity{
		ID:             0,
		DefaultName:    "",
		ProduceAt:      MyTime{now},
		ProduceAddress: "",
		Category:       "",
	}
	//fmt.Println(c.ProduceAt)
	assert.Equal(t, c.ProduceAt.Format("2006-01-02 15:04:05"), fNow)
	bs, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(bs))
	assert.Equal(t, string(bs), fmt.Sprintf(`{"id":0,"default_name":"","produce_at":"%s","produce_address":"","category":""}`, fNow))
	err = json.Unmarshal(bs, &c)
	if err != nil {
		t.Fatal(err.Error())
	}
	// 2023-12-01 17:06:19 +0000 UTC does not equal 2023-12-01 17:06:19.0372598 +0800 CST m=+0.016521901
	assert.Equal(t, c.ProduceAt.Format("2006-01-02 15:04:05"), fNow)
}
