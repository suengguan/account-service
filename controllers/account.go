package controllers

import (
	"app-service/account-service/models"
	"app-service/account-service/service"
	"encoding/json"
	"fmt"
	"model"

	"github.com/astaxie/beego"
)

// Operations about Account
type AccountController struct {
	beego.Controller
}

// @Title Register
// @Description register user
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.Response
// @Failure 403 body is empty
// @router /register/ [post]
func (this *AccountController) Register() {
	var err error
	var user model.User
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	if err == nil {
		var svc service.AccountService
		var result []byte
		var newUser *model.User
		newUser, err = svc.Register(&user)
		if err == nil {
			result, err = json.Marshal(newUser)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug("Unmarshal data failed")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}

	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title Delete
// @Description delete users
// @Param	body		body 	models.User	true		"body for users content"
// @Success 200 {object} models.Response
// @Failure 403 body is empty
// @router / [delete]
func (this *AccountController) Delete() {
	var err error
	var users []*model.User
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &users)
	if err == nil {
		var svc service.AccountService
		err = svc.Delete(users)
		if err == nil {
			response.Status = model.MSG_RESULTCODE_SUCCESS
			response.Reason = "success"
			response.Result = ""
		}
	} else {
		beego.Debug("Unmarshal data failed")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}

	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title Update
// @Description update user
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.Response
// @Failure 403 body is empty
// @router / [put]
func (this *AccountController) Update() {
	var err error
	var user model.User
	var response models.Response

	// unmarshal
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &user)
	if err == nil {
		var svc service.AccountService
		var result []byte
		var newUser *model.User
		newUser, err = svc.Update(&user)
		if err == nil {
			result, err = json.Marshal(newUser)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug("Unmarshal data failed")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}

	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title GetAll
// @Description get all users exclude admin
// @Success 200 {object} models.Response
// @router / [get]
func (this *AccountController) GetAll() {
	var err error
	var response models.Response

	var svc service.AccountService
	var users []*model.User
	var result []byte
	users, err = svc.GetAll()
	if err == nil {
		result, err = json.Marshal(users)
		if err == nil {
			response.Status = model.MSG_RESULTCODE_SUCCESS
			response.Reason = "success"
			response.Result = string(result)
		}
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}
	this.Data["json"] = &response

	this.ServeJSON()
}

// @Title GetOneById
// @Description get user by id
// @Param	id		path 	int64	true		"The key for staticblock"
// @Success 200 {object} models.Response
// @Failure 403 :id is invalid
// @router /:id [get]
func (this *AccountController) GetOneById() {
	var err error
	var response models.Response

	var id int64
	id, err = this.GetInt64(":id")
	beego.Debug("GetOneById", id)
	if id > 0 && err == nil {
		var svc service.AccountService
		var user *model.User
		var result []byte
		user, err = svc.GetOneById(id)
		if err == nil {
			result, err = json.Marshal(user)
			if err == nil {
				response.Status = model.MSG_RESULTCODE_SUCCESS
				response.Reason = "success"
				response.Result = string(result)
			}
		}
	} else {
		beego.Debug(err)
		err = fmt.Errorf("%s", "user id is invalid")
	}

	if err != nil {
		response.Status = model.MSG_RESULTCODE_FAILED
		response.Reason = err.Error()
		response.RetryCount = 3
	}
	this.Data["json"] = &response

	this.ServeJSON()
}
