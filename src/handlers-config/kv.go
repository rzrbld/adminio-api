package kv

import (
	iris "github.com/kataras/iris/v12"
	clients "github.com/rzrbld/adminio-api/clients"
	resph "github.com/rzrbld/adminio-api/handlers-response"
)

// clients
var madmClnt = clients.MadmClnt
var err error

var Get = func(ctx iris.Context) {
	var keyString = ctx.FormValue("keyString")

	values, err := madmClnt.GetConfigKV(keyString)
	var res = resph.BodyResHandler(ctx, err, values)
	ctx.JSON(res)
}
