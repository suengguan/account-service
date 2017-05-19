package service

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"config"
	"encoding/json"
	"model"

	daoApi "api/dao_service"
	systemApi "api/system_service"

	"github.com/astaxie/beego"
)

type AccountService struct {
}

func (this *AccountService) GetOneById(id int64) (*model.User, error) {
	var err error
	var user *model.User

	// get user
	beego.Debug("->get user")
	user, err = daoApi.UserDaoApi.GetById(id)
	if err != nil {
		return nil, err
	}

	// get resource
	beego.Debug("->get user's resource")
	var resource *model.Resource
	resource, err = daoApi.ResourceDaoApi.GetByUserId(id)
	if err != nil {
		return nil, err
	}
	user.Resource = resource

	beego.Debug("result:", *user)

	return user, err
}

func (this *AccountService) GetAll() ([]*model.User, error) {
	var err error
	var users []*model.User

	// get admin
	beego.Debug("->get admin")
	var admin *model.User
	admin, err = daoApi.UserDaoApi.GetByName("admin")
	if err != nil {
		return nil, err
	}

	// get users exclude admin
	beego.Debug("->get users exclude admin")
	users, err = daoApi.UserDaoApi.GetAllExcludeOneById(admin.Id)
	if err != nil {
		return nil, err
	}

	// get resource
	beego.Debug("->get users resource")
	var resource *model.Resource
	for _, u := range users {
		resource, err = daoApi.ResourceDaoApi.GetByUserId(u.Id)
		if err != nil {
			return nil, err
		}

		u.Resource = resource
	}

	beego.Debug("result:", users)

	return users, err
}

func (this *AccountService) Update(user *model.User) (*model.User, error) {
	var err error
	var newUser *model.User

	// get admin
	beego.Debug("->get admin")
	var admin *model.User
	admin, err = daoApi.UserDaoApi.GetByName("admin")
	if err != nil {
		return nil, err
	}

	// get admin resource
	beego.Debug("->get admin resource")
	var adminResource *model.Resource
	adminResource, err = daoApi.ResourceDaoApi.GetByUserId(admin.Id)
	if err != nil {
		return nil, err
	}

	// check resource
	beego.Debug("->check resource")
	err = this.checkResource(user.Resource, admin.Resource)
	if err != nil {
		beego.Debug("check resouce fialed:", err)
		return nil, err
	}

	// update user
	beego.Debug("->update user")
	newUser, err = daoApi.UserDaoApi.Update(user)
	if err != nil {
		beego.Debug(err)
		return nil, err
	}
	//beego.Debug(*newUser)
	//beego.Debug(*newUser.Resource)

	// get user resource
	beego.Debug("->get old resource")
	var oldResource *model.Resource
	oldResource, err = daoApi.ResourceDaoApi.GetByUserId(user.Id)
	if err != nil {
		return nil, err
	}

	// update resource
	beego.Debug("->update resource")
	var newResource *model.Resource
	var u model.User
	u.Id = newUser.Id
	newUser.Resource.User = &u
	newResource, err = this.updateResource(newUser.Resource)
	if err != nil {
		return nil, err
	}
	newUser.Resource = newResource

	// update admin resource
	beego.Debug("->update admin resource")
	adminResource.CpuUsageResource -= oldResource.CpuTotalResource
	adminResource.MemoryUsageResource -= oldResource.MemoryTotalResource
	adminResource.CpuUsageResource += newUser.Resource.CpuTotalResource
	adminResource.MemoryUsageResource += newUser.Resource.MemoryTotalResource
	_, err = this.updateResource(adminResource)
	if err != nil {
		return nil, err
	}

	beego.Debug("result:", *newUser)

	return newUser, err
}

func (this *AccountService) Delete(users []*model.User) error {
	var err error

	// get users
	beego.Debug("->get users")
	var deleteUsers []*model.User
	var user *model.User
	for _, u := range users {
		user, err = daoApi.UserDaoApi.GetById(u.Id)
		if err != nil {
			beego.Debug("delete users failed!")
			return err
		}
		deleteUsers = append(deleteUsers, user)
	}

	// get admin
	beego.Debug("->get admin")
	var admin *model.User
	admin, err = daoApi.UserDaoApi.GetByName("admin")
	if err != nil {
		return err
	}

	// get admin resource
	beego.Debug("->get admin resource")
	var adminResource *model.Resource
	adminResource, err = daoApi.ResourceDaoApi.GetByUserId(admin.Id)
	if err != nil {
		beego.Debug("get admin resource failed!")
		return err
	}

	// delete users
	beego.Debug("->delete users")
	var resource *model.Resource
	for _, u := range deleteUsers {
		// get user resource
		beego.Debug("->get user resource")
		resource, err = daoApi.ResourceDaoApi.GetByUserId(u.Id)
		if err != nil {
			beego.Debug("get user resource failed!")
			return err
		}

		// delete user
		beego.Debug("->delete user")
		beego.Debug(*u)
		err = this.deleteUser(u)
		if err != nil {
			beego.Debug("delete user failed")
			return err
		}

		// delete resource
		// beego.Debug("delete resource")
		// err = daoApi.ApiDeleteResourceById(u.Resource.Id)
		// if err != nil {
		// 	beego.Debug("delete resource failed")
		// 	return err
		// }

		// update admin resource
		beego.Debug("->update admin resource")
		adminResource.CpuUsageResource -= resource.CpuTotalResource
		adminResource.MemoryUsageResource -= resource.MemoryTotalResource
		_, err = this.updateResource(adminResource)
		if err != nil {
			beego.Debug("update admin resource failed")
			return err
		}

		// delete project
		beego.Debug("->delete user's projects")
		err = daoApi.BussinessDaoApi.DeleteAllProjects(u.Id)
		if err != nil {
			beego.Debug("delete projects failed")
			return err
		}
	}

	if err == nil {
		beego.Debug("result:success")
	}

	return err
}

func (this *AccountService) CreateAdmin() error {
	var err error

	// get admin by name
	var u *model.User
	u, err = daoApi.UserDaoApi.GetByName("admin")
	if u != nil && err == nil {
		beego.Debug("admin is already existed")
		return nil
	}

	// construct admin
	var admin model.User
	admin.Id = 0
	admin.Name = "admin"
	admin.Header = ""
	admin.Email = ""
	admin.Phone = ""
	admin.Company = ""
	admin.EncryptedPassword = "admin"
	admin.CreatedAt = time.Now().Unix()
	admin.UpdatedAt = time.Now().Unix()
	admin.Active = true
	admin.Role = model.USER_AUTHORITY_ADMIN
	admin.Resource, err = this.getTotalResource()
	if err != nil {
		return err
	}

	// check register
	err = this.checkRegister(&admin)
	if err != nil {
		beego.Debug(err)
		return nil
	}

	// create admin
	var newAdmin *model.User
	newAdmin, err = this.createUser(&admin)
	if err != nil {
		return err
	}
	beego.Debug("create user", *newAdmin)

	// create resource
	var newResource *model.Resource
	newResource, err = this.createResource(newAdmin)
	if err != nil {
		return err
	}
	newAdmin.Resource = newResource

	// create workspace dir
	err = this.createWorkspace(newAdmin)
	if err != nil {
		return err
	}

	return err
}

func (this *AccountService) getTotalResource() (*model.Resource, error) {
	var err error
	var resource model.Resource

	var algorithms []*model.Algorithm
	algorithms, err = daoApi.AlgorithmDaoApi.GetAll()
	if err != nil {
		return nil, err
	}

	var resultBody []byte
	resultBody, err = json.Marshal(&algorithms)
	if err != nil {
		return nil, err
	}
	resource.AlgorithmResource = string(resultBody)

	var totalCpu float64
	var totalMemory float64
	totalCpu, totalMemory, err = systemApi.ApiGetTotalCpuAndMemory()
	if err != nil {
		return nil, err
	}

	resource.CpuTotalResource = totalCpu
	resource.CpuUsageResource = 0.0
	resource.CpuUnit = "core"
	resource.MemoryTotalResource = totalMemory
	resource.MemoryUsageResource = 0.0
	resource.MemoryUnit = "Gi"

	return &resource, err
}

func (this *AccountService) Register(user *model.User) (*model.User, error) {
	var err error

	// check register data
	beego.Debug("->check register data")
	err = this.checkRegister(user)
	if err != nil {
		beego.Debug("check register data fialed:", err)
		return nil, err
	}

	// get admin
	beego.Debug("->get admin")
	var admin *model.User
	admin, err = daoApi.UserDaoApi.GetByName("admin")
	if err != nil {
		return nil, err
	}

	// get admin resource
	beego.Debug("->get admin resource")
	var adminResource *model.Resource
	adminResource, err = daoApi.ResourceDaoApi.GetByUserId(admin.Id)
	if err != nil {
		return nil, err
	}

	// check resource
	beego.Debug("->check resource")
	err = this.checkResource(user.Resource, admin.Resource)
	if err != nil {
		beego.Debug("check resouce fialed:", err)
		return nil, err
	}

	// create user
	beego.Debug("->create user")
	user.Id = 0
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()
	user.Role = model.USER_AUTHORITY_USER
	user.Active = true
	var newUser *model.User
	newUser, err = this.createUser(user)
	if err != nil {
		return nil, err
	}

	// create resource
	beego.Debug("->create resource")
	var newResource *model.Resource
	newResource, err = this.createResource(newUser)
	if err != nil {
		return nil, err
	}
	newUser.Resource = newResource

	// update admin resource
	beego.Debug("->update admin resource")
	adminResource.CpuUsageResource += newUser.Resource.CpuTotalResource
	adminResource.MemoryUsageResource += newUser.Resource.MemoryTotalResource
	_, err = this.updateResource(adminResource)
	if err != nil {
		return nil, err
	}

	// create workspace
	beego.Debug("->create workspace")
	err = this.createWorkspace(newUser)
	if err != nil {
		return nil, err
	}

	beego.Debug("result:", *newUser)

	return newUser, err
}

func (this *AccountService) checkRegister(user *model.User) error {
	var err error

	if user.Name == "" {
		err = fmt.Errorf("%s", "please input your name")
		return err
	}

	if user.EncryptedPassword == "" {
		err = fmt.Errorf("%s", "please input your password")
		return err
	}

	// get user by name
	var u *model.User
	u, err = daoApi.UserDaoApi.GetByName(user.Name)
	if u != nil && err == nil {
		err = fmt.Errorf("%s", user.Name+" is already existed")
		return err
	} else {
		beego.Debug(err)
		err = nil
	}

	return err
}

func (this *AccountService) checkResource(applied *model.Resource, existed *model.Resource) error {
	return nil
}

func (this *AccountService) createUser(user *model.User) (*model.User, error) {
	var err error

	// create user
	var newUser *model.User
	newUser, err = daoApi.UserDaoApi.Create(user)
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create user failed")
		return nil, err
	}

	// create namespace(user name)
	err = systemApi.ApiCreateNamespace(newUser.Name)
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create namespace failed")
		return nil, err
	}

	return newUser, err
}

func (this *AccountService) deleteUser(user *model.User) error {
	var err error

	// delete in mysql
	err = daoApi.UserDaoApi.DeleteById(user.Id)
	if err != nil {
		return err
	}

	// delete namespace in kuberntes
	err = systemApi.ApiDeleteNamespace(user.Name)
	if err != nil {
		return err
	}

	return err
}

func (this *AccountService) createResource(user *model.User) (*model.Resource, error) {
	var err error
	var resource *model.Resource

	resource = user.Resource
	resource.CpuUnit = "core"
	resource.MemoryUnit = "Gi"
	var u model.User
	u.Id = user.Id
	resource.User = &u

	var newResource *model.Resource
	newResource, err = daoApi.ResourceDaoApi.Create(resource)
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create resource failed")
		return nil, err
	}

	// create resource quota in namespace(user name)
	cpu := strconv.FormatFloat(newResource.CpuTotalResource, 'f', 3, 32)
	memory := strconv.FormatFloat(newResource.MemoryTotalResource*1024.0, 'f', 0, 32) + "Mi"
	err = systemApi.ApiCreateQuota(user.Name, user.Name, cpu, memory)
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create quota failed")
		return nil, err
	}

	return newResource, err
}

func (this *AccountService) updateResource(resource *model.Resource) (*model.Resource, error) {
	var err error

	// update in mysql
	var newResource *model.Resource
	newResource, err = daoApi.ResourceDaoApi.Update(resource)
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "udpate resource failed")
		return nil, err
	}

	// update in kubernetes
	beego.Debug(newResource)

	return newResource, err
}

func (this *AccountService) createWorkspace(user *model.User) error {
	var err error

	inputPath := config.WORKSPACE_PATH + "/" + user.Name + "/data/input"
	outputPath := config.WORKSPACE_PATH + "/" + user.Name + "/data/output"

	//beego.Debug("input  path :", inputPath)
	//beego.Debug("output path :", outputPath)

	err = os.MkdirAll(inputPath, os.ModePerm) //生成多级目录
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create input path failed")
		return err
	}

	err = os.MkdirAll(outputPath, os.ModePerm) //生成多级目录
	if err != nil {
		beego.Debug(err)
		err = fmt.Errorf("%s", "create output path failed")
		return err
	}

	return err
}
