package handlers

import (
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var KvGet = func(ctx iris.Context) {
	var keyString = ctx.FormValue("keyString")

	values, err := madmClnt.GetConfigKV(keyString)
	var res = resph.BodyResHandler(ctx, err, values)
	ctx.JSON(res)
}
