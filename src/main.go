package main

import (
	"errors"
	"fmt"
	"github.com/iris-contrib/middleware/cors"
	iris "github.com/kataras/iris/v12"
	minio "github.com/minio/minio-go/v6"
	madmin "github.com/minio/minio/pkg/madmin"
	log "log"
	"os"
	strconv "strconv"
	strings "strings"

	"github.com/kataras/iris/v12/sessions"

	"github.com/gorilla/securecookie"
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
)

var (
	sessionsManager *sessions.Sessions
	server          = getEnv("MINIO_HOST_PORT", "localhost:9000")
	maccess         = getEnv("MINIO_ACCESS", "test")
	msecret         = getEnv("MINIO_SECRET", "testtest123")
	region          = getEnv("MINIO_REGION", "us-east-1")
	ssl, _          = strconv.ParseBool(getEnv("MINIO_SSL", "false"))
	serverHostPort  = getEnv("ADMINIO_HOST_PORT", "localhost:8080")
	adminioCORS     = getEnv("ADMINIO_CORS_DOMAIN", "*")
	// AES only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	scHashKey  = getEnv("ADMINIO_COOKIE_HASH_KEY", "NRUeuq6AdskNPa7ewZuxG9TrDZC4xFat")
	scBlockKey = getEnv("ADMINIO_COOKIE_BLOCK_KEY", "bnfYuphzxPhJMR823YNezH83fuHuddFC")
	// ---------------
	scCookieName      = getEnv("ADMINIO_COOKIE_NAME", "adminiosessionid")
	oauthEnable, _    = strconv.ParseBool(getEnv("ADMINIO_OAUTH_ENABLE", "false"))
	auditLogEnable, _ = strconv.ParseBool(getEnv("ADMINIO_AUDIT_LOG_ENABLE", "false"))
	oauthProvider     = getEnv("ADMINIO_OAUTH_PROVIDER", "github")
	oauthClientId     = getEnv("ADMINIO_OAUTH_CLIENT_ID", "1111")
	oauthClientSecret = getEnv("ADMINIO_OAUTH_CLIENT_SECRET", "22222")
	oauthCallback     = getEnv("ADMINIO_OAUTH_CALLBACK", "http://"+serverHostPort+"/auth/callback")
	oauthCustomDomain = getEnv("ADMINIO_OAUTH_CUSTOM_DOMAIN", "")
)

func getEnv(key, fallback string) string {
	value, exist := os.LookupEnv(key)

	if !exist {
		return fallback
	}
	return value
}

func init() {
	cookieName := scCookieName
	hashKey := []byte(scHashKey)
	blockKey := []byte(scBlockKey)
	secureCookie := securecookie.New(hashKey, blockKey)

	sessionsManager = sessions.New(sessions.Config{
		Cookie: cookieName,
		Encode: secureCookie.Encode,
		Decode: secureCookie.Decode,
	})
}

type defaultRes struct {
	Success string
}

type User struct {
	accessKey string `json:"accessKey"`
	secretKey string `json:"secretKey"`
}

type policySet struct {
	policyName string `json:"policyName"`
	entityName string `json:"entityName"`
	isGroup    string `json:"isGroup"`
}

type Policy struct {
	policyName   string `json:"policyName"`
	policyString string `json:"policyString"`
}

type UserStatus struct {
	accessKey string               `json:"accessKey"`
	status    madmin.AccountStatus `json:"status"`
}

type bucketComplex struct {
	bucket string
	// bucketInfo minio.BucketInfo
	// bucketEvents minio.BucketNotification
}

type candidate struct {
	name       string
	interests  []string
	language   string
	experience bool
}

var GetProviderName = func(ctx iris.Context) (string, error) {
	return oauthProvider, nil
}

func BeginAuthHandler(ctx iris.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}

	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

func redirectOnCallback(ctx iris.Context) {
	url := GetState(ctx)
	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

func GetAuthURL(ctx iris.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session := sessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

var SetState = func(ctx iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}

	return "state"
}

var GetState = func(ctx iris.Context) string {
	return ctx.URLParam("state")
}

var CompleteUserAuth = func(ctx iris.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := sessionsManager.Start(ctx)
	value := session.GetString(providerName)

	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Set(providerName, sess.Marshal())
	return provider.FetchUser(sess)
}

func Logout(ctx iris.Context) error {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return err
	}
	session := sessionsManager.Start(ctx)
	session.Delete(providerName)
	return nil
}

func defaultResConstructor(ctx iris.Context, err error) iris.Map {
	var resp iris.Map
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	} else {
		resp = iris.Map{"Success": "OK"}
	}
	return resp
}

func bodyResConstructor(ctx iris.Context, err error, body interface{}) interface{} {
	var resp interface{}
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	} else {
		resp = body
	}
	return resp
}

func defaultResHandler(ctx iris.Context, err error) iris.Map {
	if oauthEnable {
		if gothUser, err := CompleteUserAuth(ctx); err == nil {
			defaultAuditLog(gothUser, ctx)
			return defaultResConstructor(ctx, err)
		} else {
			return iris.Map{"auth": false, "oauth": oauthEnable}
		}
	} else {
		return defaultResConstructor(ctx, err)
	}

	return nil
}

func bodyResHandler(ctx iris.Context, err error, body interface{}) interface{} {
	if oauthEnable {
		if gothUser, err := CompleteUserAuth(ctx); err == nil {
			defaultAuditLog(gothUser, ctx)
			return bodyResConstructor(ctx, err, body)
		} else {
			return iris.Map{"auth": false, "oauth": oauthEnable}
		}
	} else {
		return bodyResConstructor(ctx, err, body)
	}
	return nil
}

func defaultAuditLog(user goth.User, ctx iris.Context) {
	ctx.ViewData("", user)
	log.Print("user: ", user.NickName, "; method:", ctx.RouteName())
}

func main() {
	goth.UseProviders(
		github.New(oauthClientId, oauthClientSecret, oauthCallback),
		dropbox.New(oauthClientId, oauthClientSecret, oauthCallback),
		digitalocean.New(oauthClientId, oauthClientSecret, oauthCallback),
		bitbucket.New(oauthClientId, oauthClientSecret, oauthCallback),
		box.New(oauthClientId, oauthClientSecret, oauthCallback),
		salesforce.New(oauthClientId, oauthClientSecret, oauthCallback),
		amazon.New(oauthClientId, oauthClientSecret, oauthCallback),
		onedrive.New(oauthClientId, oauthClientSecret, oauthCallback),
		slack.New(oauthClientId, oauthClientSecret, oauthCallback),
		heroku.New(oauthClientId, oauthClientSecret, oauthCallback),
		gitlab.New(oauthClientId, oauthClientSecret, oauthCallback),
		auth0.New(oauthClientId, oauthClientSecret, oauthCallback, oauthCustomDomain),
	)

	fmt.Println("\033[31m\r\n ________   ________   _____ ______    ___   ________    ___   ________     \r\n|\\   __  \\ |\\   ___ \\ |\\   _ \\  _   \\ |\\  \\ |\\   ___  \\ |\\  \\ |\\   __  \\    \r\n\\ \\  \\|\\  \\\\ \\  \\_|\\ \\\\ \\  \\\\\\__\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\   \r\n \\ \\   __  \\\\ \\  \\ \\\\ \\\\ \\  \\\\|__| \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\  \r\n  \\ \\  \\ \\  \\\\ \\  \\_\\\\ \\\\ \\  \\    \\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\ \r\n   \\ \\__\\ \\__\\\\ \\_______\\\\ \\__\\    \\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\_______\\\r\n    \\|__|\\|__| \\|_______| \\|__|     \\|__| \\|__| \\|__| \\|__| \\|__| \\|_______|\r\n                                                                            \r\n                                                                            \r\n                                                                            \033[m")
	fmt.Println("\033[33mAdmin REST API for http://min.io (minio) s3 server")
	fmt.Println("version  : 0.9 ")
	fmt.Println("Author   : rzrbld")
	fmt.Println("License  : MIT")
	fmt.Println("Git-repo : https://github.com/rzrbld/adminio \033[m \r\n")

	// connect
	madmClnt, err := madmin.New(server, maccess, msecret, ssl)
	if err != nil {
		log.Print(err)
	}

	minioClnt, err := minio.New(server, maccess, msecret, ssl)
	if err != nil {
		log.Print(err)
	}

	app := iris.New()

	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{adminioCORS}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})

	v1auth := app.Party("/auth/", crs).AllowMethods(iris.MethodOptions)
	{
		v1auth.Get("/logout/", func(ctx iris.Context) {
			Logout(ctx)
			ctx.Redirect("/", iris.StatusTemporaryRedirect)
		})

		v1auth.Get("/", func(ctx iris.Context) {
			// try to get the user without re-authenticating
			if gothUser, err := CompleteUserAuth(ctx); err == nil {
				ctx.ViewData("", gothUser)
				ctx.JSON(iris.Map{"name": gothUser.NickName, "auth": true, "oauth": oauthEnable})
			} else {
				BeginAuthHandler(ctx)
			}
		})

		v1auth.Get("/check", func(ctx iris.Context) {
			// try to get the user without re-authenticating
			if gothUser, err := CompleteUserAuth(ctx); err == nil {
				ctx.ViewData("", gothUser)
				ctx.JSON(iris.Map{"name": gothUser.NickName, "auth": true, "oauth": oauthEnable})
			} else {
				ctx.JSON(iris.Map{"auth": false, "oauth": oauthEnable})
			}
		})

		v1auth.Get("/callback", func(ctx iris.Context) {
			_, err := CompleteUserAuth(ctx)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Writef("%v", err)
				return
			}
			// ctx.ViewData("", user)
			redirectOnCallback(ctx)
			// ctx.JSON(iris.Map{"name": user.NickName, "auth":true, "oauth":oauthEnable})
		})
	}

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions)
	{

		v1.Get("/list-buckets", func(ctx iris.Context) {
			lb, err := minioClnt.ListBuckets()
			var res = bodyResHandler(ctx, err, lb)
			ctx.JSON(res)
		})

		v1.Post("/set-status-group", func(ctx iris.Context) {
			var group = ctx.FormValue("group")
			var status = madmin.GroupStatus(ctx.FormValue("status"))

			err = madmClnt.SetGroupStatus(group, status)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/get-description-group", func(ctx iris.Context) {
			var group = ctx.FormValue("group")

			grp, err := madmClnt.GetGroupDescription(group)
			var res = bodyResHandler(ctx, err, grp)
			ctx.JSON(res)
		})

		v1.Post("/update-members-group", func(ctx iris.Context) {
			gar := madmin.GroupAddRemove{}
			gar.Group = ctx.FormValue("group")
			gar.Members = strings.Split(ctx.FormValue("members"), ",")

			gar.IsRemove, err = strconv.ParseBool(ctx.FormValue("IsRemove"))
			if err != nil {
				log.Print(err)
				ctx.JSON(iris.Map{"error": err.Error()})
			}

			err = madmClnt.UpdateGroupMembers(gar)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/add-user", func(ctx iris.Context) {

			// debug body
			// rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
			// if err != nil { /* handle the error */ ctx.Writef("%v", err) }

			// rawBodyAsString := string(rawBodyAsBytes)
			// println(rawBodyAsString)

			user := User{}
			user.accessKey = ctx.FormValue("accessKey")
			user.secretKey = ctx.FormValue("secretKey")

			err = madmClnt.AddUser(user.accessKey, user.secretKey)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)

		})

		v1.Post("/create-user-extended", func(ctx iris.Context) {

			p := policySet{}
			p.policyName = ctx.FormValue("policyName")
			p.entityName = ctx.FormValue("accessKey")

			u := User{}
			u.accessKey = ctx.FormValue("accessKey")
			u.secretKey = ctx.FormValue("secretKey")

			err = madmClnt.AddUser(u.accessKey, u.secretKey)
			if err != nil {
				log.Print(err)
				ctx.JSON(iris.Map{"error": err.Error()})
			} else {
				err = madmClnt.SetPolicy(p.policyName, p.entityName, false)
				var res = defaultResHandler(ctx, err)
				ctx.JSON(res)
			}
		})

		v1.Post("/set-user", func(ctx iris.Context) {
			u := User{}
			p := policySet{}
			us := UserStatus{}

			u.accessKey = ctx.FormValue("accessKey")
			u.secretKey = ctx.FormValue("secretKey")
			us.status = madmin.AccountStatus(ctx.FormValue("status"))
			p.policyName = ctx.FormValue("policyName")
			if u.secretKey == "" {
				err = madmClnt.SetUserStatus(u.accessKey, us.status)
			} else {
				err = madmClnt.SetUser(u.accessKey, u.secretKey, us.status)
			}
			if err != nil {
				log.Print(err)
				ctx.JSON(iris.Map{"error": err.Error()})
			} else {
				if p.policyName == "" {
					var res = defaultResHandler(ctx, err)
					ctx.JSON(res)
				} else {
					err = madmClnt.SetPolicy(p.policyName, u.accessKey, false)
					var res = defaultResHandler(ctx, err)
					ctx.JSON(res)
				}
			}
		})

		v1.Get("/list-users", func(ctx iris.Context) {
			st, err := madmClnt.ListUsers()
			var res = bodyResHandler(ctx, err, st)
			ctx.JSON(res)
		})

		v1.Post("/set-status-user", func(ctx iris.Context) {
			us := UserStatus{}
			us.accessKey = ctx.FormValue("accessKey")
			us.status = madmin.AccountStatus(ctx.FormValue("status"))

			err = madmClnt.SetUserStatus(us.accessKey, us.status)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/delete-user", func(ctx iris.Context) {
			user := User{}
			user.accessKey = ctx.FormValue("accessKey")

			err = madmClnt.RemoveUser(user.accessKey)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/make-bucket", func(ctx iris.Context) {
			var newBucket = ctx.FormValue("newBucket")
			var newBucketRegion = ctx.FormValue("newBucketRegion")
			if newBucketRegion == "" {
				newBucketRegion = region
			}

			err = minioClnt.MakeBucket(newBucket, newBucketRegion)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/get-bucket-events", func(ctx iris.Context) {
			var bucket = ctx.FormValue("bucket")
			bn, err := minioClnt.GetBucketNotification(bucket)

			var res = bodyResHandler(ctx, err, bn)
			ctx.JSON(res)
		})

		v1.Post("/remove-bucket-events", func(ctx iris.Context) {
			var bucket = ctx.FormValue("bucket")
			err := minioClnt.RemoveAllBucketNotification(bucket)

			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/set-bucket-events", func(ctx iris.Context) {
			var arrARN = strings.Split(ctx.FormValue("stsARN"), ":")

			var stsARN = minio.NewArn(arrARN[1], arrARN[2], arrARN[3], arrARN[4], arrARN[5])

			var bucket = ctx.FormValue("bucket")
			var eventTypes = strings.Split(ctx.FormValue("eventTypes"), ",")
			var filterPrefix = ctx.FormValue("filterPrefix")
			var filterSuffix = ctx.FormValue("filterSuffix")

			bucketNotify, err := minioClnt.GetBucketNotification(bucket)

			var newNotification = minio.NewNotificationConfig(stsARN)
			for _, event := range eventTypes {
				switch event {
				case "put":
					newNotification.AddEvents(minio.ObjectCreatedAll)
				case "delete":
					newNotification.AddEvents(minio.ObjectRemovedAll)
				case "get":
					newNotification.AddEvents(minio.ObjectAccessedAll)
				}
			}
			if filterPrefix != "" {
				newNotification.AddFilterPrefix(filterPrefix)
			}
			if filterSuffix != "" {
				newNotification.AddFilterSuffix(filterSuffix)
			}

			switch arrARN[2] {
			case "sns":
				if bucketNotify.AddTopic(newNotification) {
					err = fmt.Errorf("Overlapping Topic configs")
				}
			case "sqs":
				if bucketNotify.AddQueue(newNotification) {
					err = fmt.Errorf("Overlapping Queue configs")
				}
			case "lambda":
				if bucketNotify.AddLambda(newNotification) {
					err = fmt.Errorf("Overlapping lambda configs")
				}
			}

			err = minioClnt.SetBucketNotification(bucket, bucketNotify)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Get("/list-buckets-extended", func(ctx iris.Context) {
			lb, err := minioClnt.ListBuckets()
			allBuckets := []interface{}{}
			for _, bucket := range lb {
				bn, err := minioClnt.GetBucketNotification(bucket.Name)
				if err != nil {
					fmt.Errorf("Error while getting bucket notification")
				}
				br := iris.Map{"name": bucket.Name, "info": bucket, "events": bn}
				allBuckets = append(allBuckets, br)
			}

			var res = bodyResHandler(ctx, err, allBuckets)
			ctx.JSON(res)
		})

		v1.Post("/delete-bucket", func(ctx iris.Context) {
			var bucketName = ctx.FormValue("bucketName")

			err := minioClnt.RemoveBucket(bucketName)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/get-bucket-lifecycle", func(ctx iris.Context) {
			var bucketName = ctx.FormValue("bucketName")

			lc, err := minioClnt.GetBucketLifecycle(bucketName)
			var res = bodyResHandler(ctx, err, lc)
			ctx.JSON(res)
		})

		v1.Post("/set-bucket-lifecycle", func(ctx iris.Context) {
			var bucketName = ctx.FormValue("bucketName")
			var lifecycle = ctx.FormValue("lifecycle")

			err := minioClnt.SetBucketLifecycle(bucketName, lifecycle)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Get("/server-info", func(ctx iris.Context) {
			si, err := madmClnt.ServerInfo()
			var res = bodyResHandler(ctx, err, si)
			ctx.JSON(res)
		})

		v1.Get("/disk-info", func(ctx iris.Context) {
			du, err := madmClnt.DataUsageInfo()
			var res = bodyResHandler(ctx, err, du)
			ctx.JSON(res)
		})

		v1.Get("/list-groups", func(ctx iris.Context) {
			lg, err := madmClnt.ListGroups()
			var res = bodyResHandler(ctx, err, lg)
			ctx.JSON(res)
		})

		v1.Get("/list-policies", func(ctx iris.Context) {
			lp, err := madmClnt.ListCannedPolicies()
			var res = bodyResHandler(ctx, err, lp)
			ctx.JSON(res)
		})

		v1.Post("/add-policy", func(ctx iris.Context) {
			p := Policy{}
			p.policyName = ctx.FormValue("policyName")
			p.policyString = ctx.FormValue("policyString")

			err = madmClnt.AddCannedPolicy(p.policyName, p.policyString)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/delete-policy", func(ctx iris.Context) {
			p := policySet{}
			p.policyName = ctx.FormValue("policyName")

			err = madmClnt.RemoveCannedPolicy(p.policyName)
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/set-policy", func(ctx iris.Context) {
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
			var res = defaultResHandler(ctx, err)
			ctx.JSON(res)
		})

		v1.Post("/get-kv", func(ctx iris.Context) {
			var keyString = ctx.FormValue("keyString")

			values, err := madmClnt.GetConfigKV(keyString)
			var res = bodyResHandler(ctx, err, values)
			ctx.JSON(res)
		})

	}

	app.Run(iris.Addr(serverHostPort))
}
