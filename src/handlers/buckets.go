package handlers

import (
	"fmt"
	log "log"
	"strings"

	iris "github.com/kataras/iris/v12"
	minio "github.com/minio/minio-go/v6"
	cnf "github.com/rzrbld/adminio-api/config"
	resph "github.com/rzrbld/adminio-api/response"
)

var BuckList = func(ctx iris.Context) {
	lb, err := minioClnt.ListBuckets()
	var res = resph.BodyResHandler(ctx, err, lb)
	ctx.JSON(res)
}

var BuckListExtended = func(ctx iris.Context) {
	lb, err := minioClnt.ListBuckets()
	allBuckets := []interface{}{}
	for _, bucket := range lb {
		bn, err := minioClnt.GetBucketNotification(bucket.Name)
		if err != nil {
			log.Print("Error while getting bucket notification")
		}
		br := iris.Map{"name": bucket.Name, "info": bucket, "events": bn}
		allBuckets = append(allBuckets, br)
	}

	var res = resph.BodyResHandler(ctx, err, allBuckets)
	ctx.JSON(res)
}

var BuckMake = func(ctx iris.Context) {
	var newBucket = ctx.FormValue("newBucket")
	var newBucketRegion = ctx.FormValue("newBucketRegion")
	if newBucketRegion == "" {
		newBucketRegion = cnf.Region
	}

	err := minioClnt.MakeBucket(newBucket, newBucketRegion)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var BuckDelete = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	err := minioClnt.RemoveBucket(bucketName)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var BuckGetLifecycle = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	lc, err := minioClnt.GetBucketLifecycle(bucketName)
	var res = resph.BodyResHandler(ctx, err, lc)
	ctx.JSON(res)
}

var BuckSetLifecycle = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")
	var lifecycle = ctx.FormValue("lifecycle")

	err := minioClnt.SetBucketLifecycle(bucketName, lifecycle)
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var BuckGetEvents = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucket")
	bn, err := minioClnt.GetBucketNotification(bucket)

	var res = resph.BodyResHandler(ctx, err, bn)
	ctx.JSON(res)
}

var BuckSetEvents = func(ctx iris.Context) {
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
	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}

var BuckRemoveEvents = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucket")
	err := minioClnt.RemoveAllBucketNotification(bucket)

	var res = resph.DefaultResHandler(ctx, err)
	ctx.JSON(res)
}
