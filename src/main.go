package main

import (
	"fmt"
	"os"
	iris "github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	log     "log"
	madmin "github.com/minio/minio/pkg/madmin"
	minio "github.com/minio/minio-go/v6"
	strconv "strconv"
)

type defaultRes struct {
    Success string 
}

type User struct {
	accessKey  string `json:"accessKey"` 
	secretKey string `json:"secretKey"`
}

type policySet struct {
	policyName  string `json:"policyName"` 
	entityName string `json:"entityName"`
	isGroup string `json:"isGroup"`
}

type Policy struct {
	policyName  string `json:"policyName"` 
	policyString string `json:"policyString"`
}

type UserStatus struct {
	accessKey  string `json:"accessKey"` 
	status madmin.AccountStatus `json:"status"`
}

func defaultResHandler(ctx iris.Context,err error) iris.Map {
	var resp iris.Map
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	}else{
		resp = iris.Map{"Success": "OK"}
	}
	return resp
}

func bodyResHandler(ctx iris.Context,err error,body interface{}) interface{} {
	var resp interface{}
	if err != nil {
		log.Print(err)
		resp = iris.Map{"error": err.Error()}
	}else{
		resp = body
	}
	return resp
}


func main() {
	fmt.Println("\r\n ________   ________   _____ ______    ___   ________    ___   ________     \r\n|\\   __  \\ |\\   ___ \\ |\\   _ \\  _   \\ |\\  \\ |\\   ___  \\ |\\  \\ |\\   __  \\    \r\n\\ \\  \\|\\  \\\\ \\  \\_|\\ \\\\ \\  \\\\\\__\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\|\\  \\   \r\n \\ \\   __  \\\\ \\  \\ \\\\ \\\\ \\  \\\\|__| \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\  \r\n  \\ \\  \\ \\  \\\\ \\  \\_\\\\ \\\\ \\  \\    \\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\ \\  \\\\\\  \\ \r\n   \\ \\__\\ \\__\\\\ \\_______\\\\ \\__\\    \\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\__\\\\ \\_______\\\r\n    \\|__|\\|__| \\|_______| \\|__|     \\|__| \\|__| \\|__| \\|__| \\|__| \\|_______|\r\n                                                                            \r\n                                                                            \r\n                                                                            ")
	fmt.Println("simple admin API for min.io (minio) s3 server")
	fmt.Println("version  : 0.2 ")
	fmt.Println("Author   : rzrbld")
	fmt.Println("License  : MIT")
	fmt.Println("Git-repo : https://github.com/rzrbld/adminio \r\n")

	var ssl = false
	//config
	server, exists := os.LookupEnv("MINIO_HOST_PORT")
	if !exists {
		server = "localhost:9000"
	}

	maccess, exists := os.LookupEnv("MINIO_ACCESS")
	if !exists {
		maccess = "test"
	}

	msecret, exists := os.LookupEnv("MINIO_SECRET")
	if !exists {
		msecret = "testtest123"
	}

	region, exists := os.LookupEnv("MINIO_REGION")
	if !exists {
		region = "us-east-1"
	}	

	sslstr, exists := os.LookupEnv("MINIO_SSL")
	if exists {
		sslbool, err := strconv.ParseBool(sslstr)
		if err != nil {
			log.Print(err)
		}
		ssl = sslbool
	}

	serverHostPort, exists := os.LookupEnv("API_HOST_PORT")
	if !exists {
		serverHostPort = os.Getenv("API_HOST_PORT")
	}

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
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowCredentials: true,
	})
    
	
	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions) // <- important for the preflight.
	{

		v1.Post("/add-user", func(ctx iris.Context) {

			// debug body
			// rawBodyAsBytes, err := ioutil.ReadAll(ctx.Request().Body)
			// if err != nil { /* handle the error */ ctx.Writef("%v", err) }
			  
			// rawBodyAsString := string(rawBodyAsBytes)
			// println(rawBodyAsString) 

			user := User{}
			user.accessKey = ctx.FormValue("accessKey")
			user.secretKey = ctx.FormValue("secretKey")

    		err = madmClnt.AddUser(user.accessKey,user.secretKey)
    		var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
			
		})

		v1.Get("/list-users", func(ctx iris.Context) {
			st, err := madmClnt.ListUsers()
			var res = bodyResHandler(ctx,err,st)
			ctx.JSON(res)
		})

		v1.Post("/set-status-user", func(ctx iris.Context) {
			us := UserStatus{}
			us.accessKey = ctx.FormValue("accessKey")
			us.status = madmin.AccountStatus(ctx.FormValue("status"))

    		err = madmClnt.SetUserStatus(us.accessKey,us.status)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})

		v1.Post("/delete-user", func(ctx iris.Context) {
			user := User{}
			user.accessKey = ctx.FormValue("accessKey")

    		err = madmClnt.RemoveUser(user.accessKey)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})

		v1.Post("/make-bucket", func(ctx iris.Context) {
			var newBucket = ctx.FormValue("newBucket")
			var newBucketRegion = ctx.FormValue("newBucketRegion")
			if(newBucketRegion == ""){
				newBucketRegion = region
			}

			err = minioClnt.MakeBucket(newBucket, newBucketRegion)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})

		v1.Get("/list-buckets", func(ctx iris.Context) {
			lb, err := minioClnt.ListBuckets()
			var res = bodyResHandler(ctx,err,lb)
			ctx.JSON(res)
		})

		v1.Post("/delete-bucket", func(ctx iris.Context) {
			var bucketName = ctx.FormValue("bucketName")
			
			err := minioClnt.RemoveBucket(bucketName)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})


		v1.Get("/server-info", func(ctx iris.Context) {
			si, err := madmClnt.ServerInfo()
			var res = bodyResHandler(ctx,err,si)
			ctx.JSON(res)
		})

		v1.Get("/list-groups", func(ctx iris.Context) {
			lg, err := madmClnt.ListGroups()
			var res = bodyResHandler(ctx,err,lg)
			ctx.JSON(res)
		})

		v1.Post("/create-user-extended", func(ctx iris.Context) {

			p := policySet{}
			p.policyName = ctx.FormValue("policyName")
			p.entityName = ctx.FormValue("accessKey")

			u := User{}
			u.accessKey = ctx.FormValue("accessKey")
			u.secretKey = ctx.FormValue("secretKey")

    		err = madmClnt.AddUser(u.accessKey,u.secretKey)
			if err != nil {
				log.Print(err)
				ctx.JSON(iris.Map{"error": err.Error()})
			} else {
				err = madmClnt.SetPolicy(p.policyName,p.entityName,false)
				var res = defaultResHandler(ctx,err)
    			ctx.JSON(res)
			}
		})


		v1.Get("/list-policies", func(ctx iris.Context) {
			lp, err := madmClnt.ListCannedPolicies()
			var res = bodyResHandler(ctx,err,lp)
			ctx.JSON(res)
		})

		v1.Post("/add-policy", func(ctx iris.Context) {
			p := Policy{}
			p.policyName = ctx.FormValue("policyName")
			p.policyString = ctx.FormValue("policyString")

			err = madmClnt.AddCannedPolicy(p.policyName,p.policyString)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})
		
		v1.Post("/delete-policy", func(ctx iris.Context) {
			p := policySet{}
			p.policyName = ctx.FormValue("policyName")

			err = madmClnt.RemoveCannedPolicy(p.policyName)
			var res = defaultResHandler(ctx,err)
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

    		err = madmClnt.SetPolicy(p.policyName,p.entityName,isGroupBool)
			var res = defaultResHandler(ctx,err)
    		ctx.JSON(res)
		})
	}

	app.Run(iris.Addr(serverHostPort))
}


