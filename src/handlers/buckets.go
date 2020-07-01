package handlers

import (
	"context"
	"fmt"
	log "log"
	"strconv"
	"strings"
	"encoding/json"
	"sync"

	iris "github.com/kataras/iris/v12"
	minio "github.com/minio/minio-go/v6"
	"github.com/minio/minio-go/v6/pkg/tags"
	madmin "github.com/minio/minio/pkg/madmin"
	cnf "github.com/rzrbld/adminio-api/config"
	resph "github.com/rzrbld/adminio-api/response"
	policy "github.com/minio/minio-go/v6/pkg/policy"
)

func getPolicyWithName(bucketName string) (string, string, error) {
	var p policy.BucketAccessPolicy

	bp, err := minioClnt.GetBucketPolicyWithContext(context.Background(), bucketName)
	policyShort := "none"
	if bp != "" {
		if err = json.Unmarshal([]byte(bp), &p); err != nil {
			fmt.Println("Error Unmarshal policy")
		}
		pName := string(policy.GetPolicy(p.Statements, bucketName, ""))
		if pName == string(policy.BucketPolicyNone) && bp != "" {
			pName = "custom"
		}
		policyShort = policyToString(pName)
	}

	return policyShort, bp, err
}

func policyToString(policyName string) string {
	name := ""
	switch policyName {
	case "none":
			name = "none"
		case "readonly":
			name = "download"
		case "writeonly":
			name = "upload"
		case "readwrite":
			name = "public"
		case "custom":
			name = "custom"
	}
	return name
}

func stringToPolicy(strPolicy string) string {
	policy := ""
	switch strPolicy {
	case "none":
			policy = "none"
		case "download":
			policy = "readonly"
		case "upload":
			policy = "writeonly"
		case "public":
			policy = "readwrite"
		case "custom":
			policy = "custom"
	}
	return policy
}

func isJSON(s string) bool {
    var js map[string]interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}


var BuckList = func(ctx iris.Context) {
	lb, err := minioClnt.ListBuckets()
	var res = resph.BodyResHandler(ctx, err, lb)
	ctx.JSON(res)
}

var BuckListExtended = func(ctx iris.Context)  {

	var wg sync.WaitGroup

		lb, err := minioClnt.ListBuckets()
		allBuckets := make([]iris.Map, len(lb))

    wg.Add(len(lb))
    for i := 0; i < len(lb); i++ {
			go func(i int) {
				bucket := lb[i]
				bn, err := minioClnt.GetBucketNotification(bucket.Name)
				if err != nil {
					log.Print("Error while getting bucket notification", err)
				}
				bq, _ := madmClnt.GetBucketQuota(context.Background(), bucket.Name)
				bt, bterr := minioClnt.GetBucketTaggingWithContext(context.Background(), bucket.Name)

				pName, _, _ := getPolicyWithName(bucket.Name)

				btMap := map[string]string{}

				if bterr == nil {
					btMap = bt.ToMap()
				}

				br := iris.Map{"name": bucket.Name, "info": bucket, "events": bn, "quota": bq, "tags": btMap, "policy": pName}
				allBuckets[i] = br
				wg.Done()
			}(i)
    }
    wg.Wait()
		var res = resph.BodyResHandler(ctx, err, allBuckets)
		ctx.JSON(res)
}

var BuckSetTags = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")
	var tagsString = ctx.FormValue("bucketTags")

	bucketTags, err := tags.Parse(tagsString, true)

	if err != nil {
		log.Print("Error while parse bucket tags", err)
	}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err := minioClnt.SetBucketTaggingWithContext(context.Background(), bucketName, bucketTags)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetTags = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		bt, bterr := minioClnt.GetBucketTaggingWithContext(context.Background(), bucketName)

		btMap := map[string]string{}

		if bterr == nil {
			btMap = bt.ToMap()
		}

		var res = resph.BodyResHandler(ctx, err, btMap)
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

var BuckGetPolicy = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) != false {
		pName, bp, err := getPolicyWithName(bucketName)

		respBp := iris.Map{"policy":bp, "name":pName}
		var res = resph.BodyResHandler(ctx, err, respBp)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckSetPolicy = func(ctx iris.Context) {

	var bucket = ctx.FormValue("bucketName")
	var policyStr = ctx.FormValue("bucketPolicy")

	if !isJSON(policyStr){
		if policyStr == "none" {
			policyStr = ""
		} else {
			bucketPolicy := stringToPolicy(policyStr)
			var p = policy.BucketAccessPolicy{Version: "2012-10-17"}
			p.Statements = policy.SetPolicy(p.Statements, policy.BucketPolicy(bucketPolicy), bucket, "")
			policyJSON, err := json.Marshal(p)
			if err != nil {
				fmt.Println("Error marshal json", err, string(policyJSON))
			}
			policyStr = string(policyJSON)
		}
	}

	if resph.CheckAuthBeforeRequest(ctx) != false {
		err = minioClnt.SetBucketPolicyWithContext(context.Background(), bucket, policyStr)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}
