package users

import (
	iris "github.com/kataras/iris/v12"
	madmin "github.com/minio/minio/pkg/madmin"
	clients "github.com/rzrbld/adminio-api/clients"
	resph "github.com/rzrbld/adminio-api/handlers-response"
	log "log"
)

type UserStatus struct {
	accessKey string               `json:"accessKey"`
	status    madmin.AccountStatus `json:"status"`
}

type User struct {
	accessKey string `json:"accessKey"`
	secretKey string `json:"secretKey"`
}

type policySet struct {
	policyName string `json:"policyName"`
	entityName string `json:"entityName"`
	isGroup    string `json:"isGroup"`
}

// clients
var madmClnt = clients.MadmClnt
var minioClnt = clients.MinioClnt

var err error

var List = func(ctx iris.Context) {
	st, err := madmClnt.ListUsers()
	var res = resph.BodyResHandler(ctx, err, st)
	ctx.JSON(res)
}

var SetStats = func(ctx iris.Context) {
	us := UserStatus{}
	us.accessKey = ctx.FormValue("accessKey")
	us.status = madmin.AccountStatus(ctx.FormValue("status"))

	err = madmClnt.SetUserStatus(us.accessKey, us.status)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var Delete = func(ctx iris.Context) {
	user := User{}
	user.accessKey = ctx.FormValue("accessKey")

	err = madmClnt.RemoveUser(user.accessKey)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var Add = func(ctx iris.Context) {
	user := User{}
	user.accessKey = ctx.FormValue("accessKey")
	user.secretKey = ctx.FormValue("secretKey")

	err = madmClnt.AddUser(user.accessKey, user.secretKey)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var CreateExtended = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")
	p.entityName = ctx.FormValue("accessKey")

	u := User{}
	u.accessKey = ctx.FormValue("accessKey")
	u.secretKey = ctx.FormValue("secretKey")

	err = madmClnt.AddUser(u.accessKey, u.secretKey)
	if err != nil {
		log.Print(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	} else {
		err = madmClnt.SetPolicy(p.policyName, p.entityName, false)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	}
}

var Set = func(ctx iris.Context) {
	u := User{}
	p := policySet{}
	us := UserStatus{}

	u.accessKey = ctx.FormValue("accessKey")
	u.secretKey = ctx.FormValue("secretKey")
	us.status = madmin.AccountStatus(ctx.FormValue("status"))
	p.policyName = ctx.FormValue("policyName")
	if u.secretKey == "" {
		err = madmClnt.SetUserStatus(u.accessKey, us.status)
	} else {
		err = madmClnt.SetUser(u.accessKey, u.secretKey, us.status)
	}
	if err != nil {
		log.Print(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	} else {
		if p.policyName == "" {
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		} else {
			err = madmClnt.SetPolicy(p.policyName, u.accessKey, false)
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		}
	}
}
