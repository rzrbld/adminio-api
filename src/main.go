package main

import (
	"fmt"
	"github.com/iris-contrib/middleware/cors"
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

	cnf "github.com/rzrbld/adminio-api/config"
	authh "github.com/rzrbld/adminio-api/handlers-auth"
	bcth "github.com/rzrbld/adminio-api/handlers-bucket"
	kvh "github.com/rzrbld/adminio-api/handlers-config"
	grph "github.com/rzrbld/adminio-api/handlers-groups"
	plch "github.com/rzrbld/adminio-api/handlers-policy"
	srvh "github.com/rzrbld/adminio-api/handlers-server"
	usrh "github.com/rzrbld/adminio-api/handlers-users"
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
	)

	fmt.Println("\033[31m\r\n ________   ________   _____ ______    ___   ________    ___   ________     \r\n|\\   __  \\ |\\   ___ \\ |\\   _ \\  _   \\ |\\  \\ |\\   ___  \\ |\\  \\ |\\   __  \\    \r\n\\ \\  \\|\\  \\\\ \\  \\_|\\ \\\\ \\  \\\\\\__\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\   \r\n \\ \\   __  \\\\ \\  \\ \\\\ \\\\ \\  \\\\|__| \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\  \r\n  \\ \\  \\ \\  \\\\ \\  \\_\\\\ \\\\ \\  \\    \\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\ \r\n   \\ \\__\\ \\__\\\\ \\_______\\\\ \\__\\    \\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\_______\\\r\n    \\|__|\\|__| \\|_______| \\|__|     \\|__| \\|__| \\|__| \\|__| \\|__| \\|_______|\r\n                                                                            \r\n                                                                            \r\n                                                                            \033[m")
	fmt.Println("\033[33mAdmin REST API for http://min.io (minio) s3 server")
	fmt.Println("version  : 0.9 ")
	fmt.Println("Author   : rzrbld")
	fmt.Println("License  : MIT")
	fmt.Println("Git-repo : https://github.com/rzrbld/adminio \033[m \r\n")

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{cnf.AdminioCORS}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	v1auth := app.Party("/auth/", crs).AllowMethods(iris.MethodOptions)
	{
		v1auth.Get("/logout/", authh.Logout)
		v1auth.Get("/", authh.Root)
		v1auth.Get("/check", authh.Check)
		v1auth.Get("/callback", authh.Callback)
	}

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{
		v1.Get("/list-buckets", bcth.List)
		v1.Post("/make-bucket", bcth.Make)
		v1.Get("/list-buckets-extended", bcth.ListExtended)
		v1.Post("/delete-bucket", bcth.Delete)
		v1.Post("/get-bucket-lifecycle", bcth.GetLifecycle)
		v1.Post("/set-bucket-lifecycle", bcth.SetLifecycle)
		v1.Post("/get-bucket-events", bcth.GetEvents)
		v1.Post("/set-bucket-events", bcth.SetEvents)
		v1.Post("/remove-bucket-events", bcth.RemoveEvents)

		v1.Get("/list-users", usrh.List)
		v1.Post("/set-status-user", usrh.SetStats)
		v1.Post("/delete-user", usrh.Delete)
		v1.Post("/add-user", usrh.Add)
		v1.Post("/create-user-extended", usrh.CreateExtended)
		v1.Post("/set-user", usrh.Set)

		v1.Get("/list-policies", plch.List)
		v1.Post("/add-policy", plch.Add)
		v1.Post("/delete-policy", plch.Delete)
		v1.Post("/set-policy", plch.Set)

		v1.Post("/set-status-group", grph.SetStatus)
		v1.Post("/get-description-group", grph.SetDescription)
		v1.Post("/update-members-group", grph.UpdateMembers)
		v1.Get("/list-groups", grph.List)

		v1.Get("/server-info", srvh.ServerInfo)
		v1.Get("/disk-info", srvh.DiskInfo)

		v1.Post("/get-kv", kvh.Get)

	}

	app.Run(iris.Addr(cnf.ServerHostPort))
}
