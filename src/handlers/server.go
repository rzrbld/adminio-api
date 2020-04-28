package handlers

import (
	"context"
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var ServerInfo = func(ctx iris.Context) {
	si, err := madmClnt.ServerInfo(context.Background())
	var res = resph.BodyResHandler(ctx, err, si)
	ctx.JSON(res)
}

var DiskInfo = func(ctx iris.Context) {
	du, err := madmClnt.DataUsageInfo(context.Background())
	var res = resph.BodyResHandler(ctx, err, du)
	ctx.JSON(res)
}
