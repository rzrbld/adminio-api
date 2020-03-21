package handlers

import (
	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
	log "log"
	strconv "strconv"
)

var PolList = func(ctx iris.Context) {
	lp, err := madmClnt.ListCannedPolicies()
	var res = resph.BodyResHandler(ctx, err, lp)
	ctx.JSON(res)
}

var PolAdd = func(ctx iris.Context) {
	p := Policy{}
	p.policyName = ctx.FormValue("policyName")
	p.policyString = ctx.FormValue("policyString")

	err = madmClnt.AddCannedPolicy(p.policyName, p.policyString)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var PolDelete = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")

	err = madmClnt.RemoveCannedPolicy(p.policyName)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var PolSet = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")
	p.entityName = ctx.FormValue("entityName")
	p.isGroup = ctx.FormValue("isGroup")

	isGroupBool, err := strconv.ParseBool(p.isGroup)

	if err != nil {
		log.Print(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	}

	err = madmClnt.SetPolicy(p.policyName, p.entityName, isGroupBool)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}
