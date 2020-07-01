package main

import (
	"fmt"

	"github.com/iris-contrib/middleware/cors"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	iris "github.com/kataras/iris/v12"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/box"
	"github.com/markbates/goth/providers/digitalocean"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/heroku"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/salesforce"
	"github.com/markbates/goth/providers/slack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rzrbld/goth-provider-wso2"

	cnf "github.com/rzrbld/adminio-api/config"
	hdl "github.com/rzrbld/adminio-api/handlers"
)

func main() {
	goth.UseProviders(
		github.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		dropbox.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		digitalocean.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		bitbucket.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		box.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		salesforce.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		amazon.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		onedrive.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		slack.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		heroku.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		gitlab.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback),
		auth0.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback, cnf.OauthCustomDomain),
		wso2.New(cnf.OauthClientId, cnf.OauthClientSecret, cnf.OauthCallback, cnf.OauthCustomDomain),
	)

	fmt.Println("\033[31m\r\n ________   ________   _____ ______    ___   ________    ___   ________     \r\n|\\   __  \\ |\\   ___ \\ |\\   _ \\  _   \\ |\\  \\ |\\   ___  \\ |\\  \\ |\\   __  \\    \r\n\\ \\  \\|\\  \\\\ \\  \\_|\\ \\\\ \\  \\\\\\__\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\   \r\n \\ \\   __  \\\\ \\  \\ \\\\ \\\\ \\  \\\\|__| \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\  \r\n  \\ \\  \\ \\  \\\\ \\  \\_\\\\ \\\\ \\  \\    \\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\ \r\n   \\ \\__\\ \\__\\\\ \\_______\\\\ \\__\\    \\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\_______\\\r\n    \\|__|\\|__| \\|_______| \\|__|     \\|__| \\|__| \\|__| \\|__| \\|__| \\|_______|\r\n                                                                            \r\n                                                                            \r\n                                                                            \033[m")
	fmt.Println("\033[33mAdmin REST API for http://min.io (minio) s3 server")
	fmt.Println("Version    : 1.1")
	fmt.Println("Authors    : rzrbld, 0x003e")
	fmt.Println("License    : MIT")
	fmt.Println("GitHub     : https://github.com/rzrbld/adminio-api")
	fmt.Println("Docker Hub : https://hub.docker.com/r/rzrbld/adminio-api \033[m \r\n")

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{cnf.AdminioCORS}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	// prometheus metrics
	if cnf.MetricsEnable {
		m := prometheusMiddleware.New("adminio", 0.3, 1.2, 5.0)
		hdl.RecordMetrics()
		app.Use(m.ServeHTTP)
		app.Get("/metrics", iris.FromStd(promhttp.Handler()))
	}

	if cnf.ProbesEnable {
		app.Get("/ready", hdl.Probes)
		app.Get("/live", hdl.Probes)
	}

	v1auth := app.Party("/auth/", crs).AllowMethods(iris.MethodOptions)
	{
		v1auth.Get("/logout/", hdl.AuthLogout)
		v1auth.Get("/", hdl.AuthRoot)
		v1auth.Get("/check", hdl.AuthCheck)
		v1auth.Get("/callback", hdl.AuthCallback)
	}

	v2 := app.Party("/api/v2", crs).AllowMethods(iris.MethodOptions)
	{
		v2.Get("/buckets/list", hdl.BuckList)
		v2.Post("/bucket/create", hdl.BuckMake)
		v2.Get("/buckets/list-extended", hdl.BuckListExtended)
		v2.Post("/bucket/delete", hdl.BuckDelete)
		v2.Post("/bucket/get-lifecycle", hdl.BuckGetLifecycle)
		v2.Post("/bucket/set-lifecycle", hdl.BuckSetLifecycle)
		v2.Post("/bucket/get-events", hdl.BuckGetEvents)
		v2.Post("/bucket/set-events", hdl.BuckSetEvents)
		v2.Post("/bucket/remove-events", hdl.BuckRemoveEvents)
		v2.Post("/bucket/set-quota", hdl.BuckSetQuota)
		v2.Post("/bucket/get-quota", hdl.BuckGetQuota)
		v2.Post("/bucket/remove-quota", hdl.BuckRemoveQuota)
		v2.Post("/bucket/set-tags", hdl.BuckSetTags)
		v2.Post("/bucket/get-tags", hdl.BuckGetTags)
		v2.Post("/bucket/get-policy", hdl.BuckGetPolicy)
		v2.Post("/bucket/set-policy", hdl.BuckSetPolicy)

		v2.Get("/users/list", hdl.UsrList)
		v2.Post("/user/set-status", hdl.UsrSetStats)
		v2.Post("/user/delete", hdl.UsrDelete)
		v2.Post("/user/create", hdl.UsrAdd)
		v2.Post("/user/create-extended", hdl.UsrCreateExtended)
		v2.Post("/user/update", hdl.UsrSet)

		v2.Get("/policies/list", hdl.PolList)
		v2.Post("/policy/create", hdl.PolAdd)
		v2.Post("/policy/delete", hdl.PolDelete)
		v2.Post("/policy/update", hdl.PolSet)

		v2.Post("/group/set-status", hdl.GrSetStatus)
		v2.Post("/group/get-description", hdl.GrSetDescription)
		v2.Post("/group/update-members", hdl.GrUpdateMembers)
		v2.Get("/groups/list", hdl.GrList)

		v2.Get("/server/common-info", hdl.ServerInfo)
		v2.Get("/server/disk-info", hdl.DiskInfo)

		v2.Post("/kv/get", hdl.KvGet)
	}

	app.Run(iris.Addr(cnf.ServerHostPort))
}
