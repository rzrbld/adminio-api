package handlers

import (
	"context"
	iris "github.com/kataras/iris/v12"
	madmin "github.com/minio/minio/pkg/madmin"
	resph "github.com/rzrbld/adminio-api/response"
	log "log"
)

var UsrList = func(ctx iris.Context) {
	st, err := madmClnt.ListUsers(context.Background())
	var res = resph.BodyResHandler(ctx, err, st)
	ctx.JSON(res)
}

var UsrSetStats = func(ctx iris.Context) {
	us := UserStatus{}
	us.accessKey = ctx.FormValue("accessKey")
	us.status = madmin.AccountStatus(ctx.FormValue("status"))

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.SetUserStatus(context.Background(), us.accessKey, us.status)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var UsrDelete = func(ctx iris.Context) {
	user := User{}
	user.accessKey = ctx.FormValue("accessKey")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.RemoveUser(context.Background(), user.accessKey)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var UsrAdd = func(ctx iris.Context) {
	user := User{}
	user.accessKey = ctx.FormValue("accessKey")
	user.secretKey = ctx.FormValue("secretKey")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.AddUser(context.Background(), user.accessKey, user.secretKey)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var UsrCreateExtended = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")
	p.entityName = ctx.FormValue("accessKey")

	u := User{}
	u.accessKey = ctx.FormValue("accessKey")
	u.secretKey = ctx.FormValue("secretKey")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.AddUser(context.Background(), u.accessKey, u.secretKey)
		if err != nil {
			log.Print(err)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else {
			err = madmClnt.SetPolicy(context.Background(), p.policyName, p.entityName, false)
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		}
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var UsrSet = func(ctx iris.Context) {
	u := User{}
	p := policySet{}
	us := UserStatus{}

	u.accessKey = ctx.FormValue("accessKey")
	u.secretKey = ctx.FormValue("secretKey")
	us.status = madmin.AccountStatus(ctx.FormValue("status"))
	p.policyName = ctx.FormValue("policyName")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		if u.secretKey == "" {
			err = madmClnt.SetUserStatus(context.Background(), u.accessKey, us.status)
		} else {
			err = madmClnt.SetUser(context.Background(), u.accessKey, u.secretKey, us.status)
		}
		if err != nil {
			log.Print(err)
			ctx.JSON(iris.Map{"error": err.Error()})
		} else {
			if p.policyName == "" {
				var res = resph.DefaultResHandler(ctx, err)
				ctx.JSON(res)
			} else {
				err = madmClnt.SetPolicy(context.Background(), p.policyName, u.accessKey, false)
				var res = resph.DefaultResHandler(ctx, err)
				ctx.JSON(res)
			}
		}
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}
