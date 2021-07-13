package handlers

import (
	"context"
	strconv "strconv"

	log "github.com/sirupsen/logrus"

	iris "github.com/kataras/iris/v12"
	resph "github.com/rzrbld/adminio-api/response"
)

var PolList = func(ctx iris.Context) {
	lp, err := madmClnt.ListCannedPolicies(context.Background())
	var res = resph.BodyResHandler(ctx, err, lp)
	ctx.JSON(res)
}

var PolAdd = func(ctx iris.Context) {
	p := Policy{}
	p.policyName = ctx.FormValue("policyName")
	p.policyString = ctx.FormValue("policyString")

	if resph.CheckAuthBeforeRequest(ctx) {
		// policy, err := iampolicy.ParseConfig(strings.NewReader(p.policyString))
		policy := p.policyString
		if err == nil {
			err = madmClnt.AddCannedPolicy(context.Background(), p.policyName, []byte(policy))
		}
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var PolDelete = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")

	if resph.CheckAuthBeforeRequest(ctx) {
		err = madmClnt.RemoveCannedPolicy(context.Background(), p.policyName)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var PolSet = func(ctx iris.Context) {
	p := policySet{}
	p.policyName = ctx.FormValue("policyName")
	p.entityName = ctx.FormValue("entityName")
	p.isGroup = ctx.FormValue("isGroup")

	isGroupBool, err := strconv.ParseBool(p.isGroup)

	if err != nil {
		log.Errorln(err)
		ctx.JSON(iris.Map{"error": err.Error()})
	}

	if resph.CheckAuthBeforeRequest(ctx) {
		err = madmClnt.SetPolicy(context.Background(), p.policyName, p.entityName, isGroupBool)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}
