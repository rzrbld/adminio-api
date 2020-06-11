package handlers

import (
	"context"
	"fmt"
	log "log"
	"strconv"
	"strings"

	iris "github.com/kataras/iris/v12"
	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/tags"
	madmin "github.com/minio/minio/pkg/madmin"
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
			log.Print("Error while getting bucket notification", err)
		}
		bq, _ := madmClnt.GetBucketQuota(context.Background(), bucket.Name)
		bt, bterr := minioClnt.GetBucketTaggingWithContext(context.Background(), bucket.Name)

		btMap := map[string]string{}

		if bterr == nil {
			btMap = bt.ToMap()
		}

		br := iris.Map{"name": bucket.Name, "info": bucket, "events": bn, "quota": bq, "tags": btMap}
		allBuckets = append(allBuckets, br)
	}

	var res = resph.BodyResHandler(ctx, err, allBuckets)
	ctx.JSON(res)
}

var BuckSetTags = func(ctx iris.Context)  {
	var bucketName = ctx.FormValue("bucketName")
	var tagsString = ctx.FormValue("bucketTags")

	bucketTags, err := tags.Parse(tagsString, true)

	if err != nil {
		log.Print("Error while getting bucket notification", err)
	}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.SetBucketTaggingWithContext(context.Background(), bucketName, bucketTags)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckMake = func(ctx iris.Context) {
	var newBucket = ctx.FormValue("newBucket")
	var newBucketRegion = ctx.FormValue("newBucketRegion")

	if newBucketRegion == "" {
		newBucketRegion = cnf.Region
	}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.MakeBucket(newBucket, newBucketRegion)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckDelete = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.RemoveBucket(bucketName)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
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

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.SetBucketLifecycle(bucketName, lifecycle)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
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

	if resph.CheckAuthBeforeRequest(ctx) != false {
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
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckRemoveEvents = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucket")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.RemoveAllBucketNotification(bucket)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckSetQuota = func(ctx iris.Context) {

	var bucket = ctx.FormValue("bucketName")
	var quotaType = madmin.QuotaType(strings.ToLower(ctx.FormValue("quotaType")))
	var quotaStr = ctx.FormValue("quotaValue")
	var quota, _ = strconv.ParseUint(quotaStr, 10, 64)
	bucketQuota := &madmin.BucketQuota{Quota: quota, Type: quotaType}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.SetBucketQuota(context.Background(), bucket, bucketQuota)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetQuota = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		bq, err := madmClnt.GetBucketQuota(context.Background(), bucket)
		var res = resph.BodyResHandler(ctx, err, bq)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckRemoveQuota = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucketName")
	var quota, _ = strconv.ParseUint("0", 10, 64)
	bucketQuota := &madmin.BucketQuota{Quota: quota}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = madmClnt.SetBucketQuota(context.Background(), bucket, bucketQuota)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}
