package handlers

import (
	"context"
	iris "github.com/kataras/iris/v12"
	madmin "github.com/minio/minio/pkg/madmin"
	resph "github.com/rzrbld/adminio-api/response"
	log "log"
	strconv "strconv"
	"strings"
)

var GrSetStatus = func(ctx iris.Context) {
	var group = ctx.FormValue("group")
	var status = ctx.FormValue("status")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		var status = madmin.GroupStatus(status)
		err = madmClnt.SetGroupStatus(context.Background(), group, status)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var GrSetDescription = func(ctx iris.Context) {
	var group = ctx.FormValue("group")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		grp, err := madmClnt.GetGroupDescription(context.Background(), group)
		var res = resph.BodyResHandler(ctx, err, grp)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var GrUpdateMembers = func(ctx iris.Context) {
	gar := madmin.GroupAddRemove{}
	gar.Group = ctx.FormValue("group")
	if ctx.FormValue("members") != "" {
		gar.Members = strings.Split(ctx.FormValue("members"), ",")
	}

	gar.IsRemove, err = strconv.ParseBool(ctx.FormValue("IsRemove"))
	if err != nil {
		log.Print(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.UpdateGroupMembers(context.Background(), gar)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}

}

var GrList = func(ctx iris.Context) {
	lg, err := madmClnt.ListGroups(context.Background())
	var res = resph.BodyResHandler(ctx, err, lg)
	ctx.JSON(res)
}
