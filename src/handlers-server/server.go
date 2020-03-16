package server

import (
	iris "github.com/kataras/iris/v12"
	clients "github.com/rzrbld/adminio-api/clients"
	resph "github.com/rzrbld/adminio-api/handlers-response"
)

// clients
var madmClnt = clients.MadmClnt
var minioClnt = clients.MinioClnt
var err error

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
