package handlers

import (
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var Probes = func(ctx iris.Context) {
	var res = resph.DefaultResConstructor(nil)
	ctx.JSON(res)
}
