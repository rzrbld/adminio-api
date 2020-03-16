package auth

import (
	iris "github.com/kataras/iris/v12"
	cnf "github.com/rzrbld/adminio-api/config"
	auth "github.com/rzrbld/adminio-api/oauth"
)

var Logout = func(ctx iris.Context) {
	auth.Logout(ctx)
	ctx.Redirect("/", iris.StatusTemporaryRedirect)
}

var Root = func(ctx iris.Context) {
	// try to get the user without re-authenticating
	if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
		ctx.ViewData("", gothUser)
		ctx.JSON(iris.Map{"name": gothUser.NickName, "auth": true, "oauth": cnf.OauthEnable})
	} else {
		auth.BeginAuthHandler(ctx)
	}
}

var Check = func(ctx iris.Context) {
	if gothUser, err := auth.CompleteUserAuth(ctx); err == nil {
		ctx.ViewData("", gothUser)
		ctx.JSON(iris.Map{"name": gothUser.NickName, "auth": true, "oauth": cnf.OauthEnable})
	} else {
		ctx.JSON(iris.Map{"auth": false, "oauth": cnf.OauthEnable})
	}
}

var Callback = func(ctx iris.Context) {
	_, err := auth.CompleteUserAuth(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.Writef("%v", err)
		return
	}
	auth.RedirectOnCallback(ctx)
}
