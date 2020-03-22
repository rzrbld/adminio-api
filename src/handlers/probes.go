package handlers

import (
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var Readiness = func(ctx iris.Context) {
	var res = resph.DefaultResConstructor(ctx, nil)
	ctx.JSON(res)
}

var Liveness = func(ctx iris.Context) {
	var res = resph.DefaultResConstructor(ctx, nil)
	ctx.JSON(res)
}
