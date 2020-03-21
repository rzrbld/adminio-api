package handlers

import (
	iris "github.com/kataras/iris/v12"
	audit "github.com/rzrbld/adminio-api/audit"
	cnf "github.com/rzrbld/adminio-api/config"
	auth "github.com/rzrbld/adminio-api/oauth"
	log "log"
)

func DefaultResHandler(ctx iris.Context, err error) iris.Map {
	if cnf.OauthEnable {
		if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
			audit.DefaultAuditLog(gothUser, ctx)
			return DefaultResConstructor(ctx, err)
		} else {
			return iris.Map{"auth": false, "oauth": cnf.OauthEnable}
		}
	} else {
		return DefaultResConstructor(ctx, err)
	}

	return nil
}

func BodyResHandler(ctx iris.Context, err error, body interface{}) interface{} {
	if cnf.OauthEnable {
		if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
			audit.DefaultAuditLog(gothUser, ctx)
			return BodyResConstructor(ctx, err, body)
		} else {
			return iris.Map{"auth": false, "oauth": cnf.OauthEnable}
		}
	} else {
		return BodyResConstructor(ctx, err, body)
	}
	return nil
}

func BodyResConstructor(ctx iris.Context, err error, body interface{}) interface{} {
	var resp interface{}
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	} else {
		resp = body
	}
	return resp
}

func DefaultResConstructor(ctx iris.Context, err error) iris.Map {
	var resp iris.Map
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	} else {
		resp = iris.Map{"Success": "OK"}
	}
	return resp
}
