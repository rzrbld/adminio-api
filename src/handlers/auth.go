package handlers

import (
	iris "github.com/kataras/iris/v12"
	cnf "github.com/rzrbld/adminio-api/config"
	auth "github.com/rzrbld/adminio-api/oauth"
)

var AuthLogout = func(ctx iris.Context) {
	auth.Logout(ctx)
	ctx.Redirect("/", iris.StatusTemporaryRedirect)
}

var AuthRoot = func(ctx iris.Context) {
	// try to get the user without re-authenticating
	if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
		auth.Redirect(ctx)
		ctx.JSON(iris.Map{"name": gothUser.UserID, "auth": true, "oauth": cnf.OauthEnable})
	} else {
		auth.BeginAuthHandler(ctx)
	}
}

var AuthCheck = func(ctx iris.Context) {
	if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
		ctx.ViewData("", gothUser)
		ctx.JSON(iris.Map{"name": gothUser.UserID, "auth": true, "oauth": cnf.OauthEnable})
	} else {
		ctx.JSON(iris.Map{"auth": false, "oauth": cnf.OauthEnable})
	}
}

var AuthCallback = func(ctx iris.Context) {
	_, err := auth.CompleteUserAuth(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Writef("%v", err)
		return
	}
	auth.Redirect(ctx)
}
