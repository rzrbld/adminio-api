package handlers

import (
	iris "github.com/kataras/iris/v12"
	madmin "github.com/minio/minio/pkg/madmin"
	resph "github.com/rzrbld/adminio-api/response"
	log "log"
	strconv "strconv"
	"strings"
)

var GrSetStatus = func(ctx iris.Context) {
	var group = ctx.FormValue("group")
	var status = madmin.GroupStatus(ctx.FormValue("status"))

	err = madmClnt.SetGroupStatus(group, status)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var GrSetDescription = func(ctx iris.Context) {
	var group = ctx.FormValue("group")

	grp, err := madmClnt.GetGroupDescription(group)
	var res = resph.BodyResHandler(ctx, err, grp)
	ctx.JSON(res)
}

var GrUpdateMembers = func(ctx iris.Context) {
	gar := madmin.GroupAddRemove{}
	gar.Group = ctx.FormValue("group")
	gar.Members = strings.Split(ctx.FormValue("members"), ",")

	gar.IsRemove, err = strconv.ParseBool(ctx.FormValue("IsRemove"))
	if err != nil {
		log.Print(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	}

	err = madmClnt.UpdateGroupMembers(gar)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var GrList = func(ctx iris.Context) {
	lg, err := madmClnt.ListGroups()
	var res = resph.BodyResHandler(ctx, err, lg)
	ctx.JSON(res)
}
