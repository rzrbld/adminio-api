package handlers

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	iris "github.com/kataras/iris/v12"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/minio/minio-go/v7/pkg/notification"
	policy "github.com/minio/minio-go/v7/pkg/policy"
	"github.com/minio/minio-go/v7/pkg/sse"
	"github.com/minio/minio-go/v7/pkg/tags"

	madmin "github.com/minio/madmin-go"
	cnf "github.com/rzrbld/adminio-api/config"
	resph "github.com/rzrbld/adminio-api/response"
)

func getPolicyWithName(bucketName string) (string, string, error) {
	var p policy.BucketAccessPolicy

	bp, err := minioClnt.GetBucketPolicy(context.Background(), bucketName)
	policyShort := "none"
	if bp != "" {
		if err = json.Unmarshal([]byte(bp), &p); err != nil {
			log.Errorln("Error Unmarshal policy")
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
	lb, err := minioClnt.ListBuckets(context.Background())
	var res = resph.BodyResHandler(ctx, err, lb)
	ctx.JSON(res)
}

var BuckListExtended = func(ctx iris.Context) {

	var wg sync.WaitGroup

	lb, err := minioClnt.ListBuckets(context.Background())
	allBuckets := make([]iris.Map, len(lb))

	wg.Add(len(lb))
	for i := 0; i < len(lb); i++ {
		go func(i int) {
			bucket := lb[i]
			bn, err := minioClnt.GetBucketNotification(context.Background(), bucket.Name)
			if err != nil {
				log.Errorln("Error while getting bucket notification", err)
			}
			bq, _ := madmClnt.GetBucketQuota(context.Background(), bucket.Name)
			bt, bterr := minioClnt.GetBucketTagging(context.Background(), bucket.Name)
			be, _ := minioClnt.GetBucketEncryption(context.Background(), bucket.Name)

			pName, _, _ := getPolicyWithName(bucket.Name)

			btMap := map[string]string{}

			if bterr == nil {
				btMap = bt.ToMap()
			}

			br := iris.Map{"name": bucket.Name, "info": bucket, "events": bn, "quota": bq, "tags": btMap, "policy": pName, "encryption": be}
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
		log.Errorln("Error while parse bucket tags", err)
	}

	if resph.CheckAuthBeforeRequest(ctx) {
		err := minioClnt.SetBucketTagging(context.Background(), bucketName, bucketTags)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetTags = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
		bt, bterr := minioClnt.GetBucketTagging(context.Background(), bucketName)

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
	var newBucketObjectLockingStr = ctx.FormValue("newBucketObjectLocking")
	var newBucketObjectLocking bool

	if newBucketRegion == "" {
		newBucketRegion = cnf.Region
	}

	if newBucketObjectLockingStr == "" {
		newBucketObjectLocking = cnf.DefaultObjectLocking
	} else {
		newBucketObjectLocking, _ = strconv.ParseBool(newBucketObjectLockingStr)
	}

	if resph.CheckAuthBeforeRequest(ctx) {
		newBucketOpts := minio.MakeBucketOptions{
			Region:        newBucketRegion,
			ObjectLocking: newBucketObjectLocking,
		}
		err := minioClnt.MakeBucket(context.Background(), newBucket, newBucketOpts)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckDelete = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
		err := minioClnt.RemoveBucket(context.Background(), bucketName)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetLifecycle = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	lc, err := minioClnt.GetBucketLifecycle(context.Background(), bucketName)
	var res = resph.BodyResHandler(ctx, err, lc)
	ctx.JSON(res)
}

var BuckSetLifecycle = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")
	var lifecycleStr = ctx.FormValue("lifecycle")
	var lcc = lifecycle.NewConfiguration()

	if isJSON(lifecycleStr) {
		jdec := json.NewDecoder(strings.NewReader(lifecycleStr))

		if err := jdec.Decode(lcc); err != nil {
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		}
	} else {
		err := xml.Unmarshal([]byte(lifecycleStr), &lcc)

		if err != nil {
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		}
	}

	if resph.CheckAuthBeforeRequest(ctx) {
		err := minioClnt.SetBucketLifecycle(context.Background(), bucketName, lcc)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetEvents = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucket")

	bn, err := minioClnt.GetBucketNotification(context.Background(), bucket)
	var res = resph.BodyResHandler(ctx, err, bn)
	ctx.JSON(res)
}

var BuckSetEvents = func(ctx iris.Context) {
	var arrARN = strings.Split(ctx.FormValue("stsARN"), ":")

	var stsARN = notification.NewArn(arrARN[1], arrARN[2], arrARN[3], arrARN[4], arrARN[5])

	var bucket = ctx.FormValue("bucket")
	var eventTypes = strings.Split(ctx.FormValue("eventTypes"), ",")
	var filterPrefix = ctx.FormValue("filterPrefix")
	var filterSuffix = ctx.FormValue("filterSuffix")

	if resph.CheckAuthBeforeRequest(ctx) {
		bucketNotify, err := minioClnt.GetBucketNotification(context.Background(), bucket)

		var newNotification = notification.NewConfig(stsARN)
		for _, event := range eventTypes {
			switch event {
			case "put":
				newNotification.AddEvents(notification.ObjectCreatedAll)
			case "delete":
				newNotification.AddEvents(notification.ObjectRemovedAll)
			case "get":
				newNotification.AddEvents(notification.ObjectAccessedAll)
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
				log.Errorln("overlapping Topic configs")
			}
		case "sqs":
			if bucketNotify.AddQueue(newNotification) {
				log.Errorln("overlapping Queue configs")
			}
		case "lambda":
			if bucketNotify.AddLambda(newNotification) {
				log.Errorln("overlapping lambda configs")
			}
		}

		if err != nil {
			log.Errorln("Error:", err)
		}

		err = minioClnt.SetBucketNotification(context.Background(), bucket, bucketNotify)

		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckRemoveEvents = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucket")

	if resph.CheckAuthBeforeRequest(ctx) {
		err := minioClnt.RemoveAllBucketNotification(context.Background(), bucket)
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

	if resph.CheckAuthBeforeRequest(ctx) {
		err = madmClnt.SetBucketQuota(context.Background(), bucket, bucketQuota)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetQuota = func(ctx iris.Context) {
	var bucket = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
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

	if resph.CheckAuthBeforeRequest(ctx) {
		err = madmClnt.SetBucketQuota(context.Background(), bucket, bucketQuota)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetPolicy = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
		pName, bp, err := getPolicyWithName(bucketName)

		respBp := iris.Map{"policy": bp, "name": pName}
		var res = resph.BodyResHandler(ctx, err, respBp)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckSetPolicy = func(ctx iris.Context) {

	var bucket = ctx.FormValue("bucketName")
	var policyStr = ctx.FormValue("bucketPolicy")

	if !isJSON(policyStr) {
		if policyStr == "none" {
			policyStr = ""
		} else {
			bucketPolicy := stringToPolicy(policyStr)
			var p = policy.BucketAccessPolicy{Version: "2012-10-17"}
			p.Statements = policy.SetPolicy(p.Statements, policy.BucketPolicy(bucketPolicy), bucket, "")
			policyJSON, err := json.Marshal(p)
			if err != nil {
				log.Errorln("Error marshal json", err, string(policyJSON))
			}
			policyStr = string(policyJSON)
		}
	}

	if resph.CheckAuthBeforeRequest(ctx) {
		err = minioClnt.SetBucketPolicy(context.Background(), bucket, policyStr)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckGetEncryption = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
		ec, err := minioClnt.GetBucketEncryption(context.Background(), bucketName)
		var res = resph.BodyResHandler(ctx, err, ec)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}

var BuckSetEncryption = func(ctx iris.Context) {

	var bucketName = ctx.FormValue("bucketName")
	var bucketEncType = ctx.FormValue("bucketEncryptionType")
	var KmsMKey = ctx.FormValue("kmsMasterKey")
	var encErr error
	var sseConf *sse.Configuration

	switch strings.ToLower(bucketEncType) {
	case "sse-kms":
		sseConf = sse.NewConfigurationSSEKMS(KmsMKey)
	case "sse-s3":
		sseConf = sse.NewConfigurationSSES3()
	default:
		encErr = fmt.Errorf("invalid encryption algorithm %s", bucketEncType)
	}

	if encErr != nil {
		var res = resph.DefaultResHandler(ctx, encErr)
		ctx.JSON(res)
	} else {
		if resph.CheckAuthBeforeRequest(ctx) {
			err := minioClnt.SetBucketEncryption(context.Background(), bucketName, sseConf)
			var res = resph.DefaultResHandler(ctx, err)
			ctx.JSON(res)
		} else {
			ctx.JSON(resph.DefaultAuthError())
		}
	}
}

var BuckRemoveEncryption = func(ctx iris.Context) {
	var bucketName = ctx.FormValue("bucketName")

	if resph.CheckAuthBeforeRequest(ctx) {
		err := minioClnt.RemoveBucketEncryption(context.Background(), bucketName)
		var res = resph.DefaultResHandler(ctx, err)
		ctx.JSON(res)
	} else {
		ctx.JSON(resph.DefaultAuthError())
	}
}
