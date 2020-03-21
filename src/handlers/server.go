package handlers

import (
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var ServerInfo = func(ctx iris.Context) {
	si, err := madmClnt.ServerInfo()
	var res = resph.BodyResHandler(ctx, err, si)
	ctx.JSON(res)
}

var DiskInfo = func(ctx iris.Context) {
	du, err := madmClnt.DataUsageInfo()
	var res = resph.BodyResHandler(ctx, err, du)
	ctx.JSON(res)
}
