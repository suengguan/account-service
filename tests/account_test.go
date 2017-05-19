package test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	_ "app-service/account-service/routers"
	"model"
)

const (
	base_url = "http://localhost:8080/v1/account"
)

func Test_Create(t *testing.T) {
	var user model.User
	user.Id = 0
	user.Name = "user"
	user.Email = "user@xx.com"
	user.Company = "company"
	user.Active = true
	user.Role = 1
	user.EncryptedPassword = "user"
	var resource model.Resource
	resource.AlgorithmResource = "algorithm list"
	resource.CpuTotalResource = 1
	resource.CpuUsageResource = 0.5
	resource.CpuUnit = "core"
	resource.MemoryTotalResource = 1
	resource.MemoryUsageResource = 0.5
	resource.MemoryUnit = "Gi"
	user.Resource = &resource

	// post create action
	requestData, err := json.Marshal(&user)
	if err != nil {
		t.Log("erro : ", err)
		return
	}

	res, err := http.Post(base_url+"/register/", "application/x-www-form-urlencoded", bytes.NewBuffer(requestData))
	if err != nil {
		t.Log("erro : ", err)
		return
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Log("erro : ", err)
		return
	}

	t.Log(string(resBody))

	var response model.Response
	json.Unmarshal(resBody, &response)
	if err != nil {
		t.Log("erro : ", err)
		return
	}

	if response.Reason == "success" {
		t.Log("PASS OK")
	} else {
		t.Log("ERROR:", response.Reason)
		t.FailNow()
	}
}
