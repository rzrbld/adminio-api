package clients

import (
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	log "github.com/sirupsen/logrus"

	madmin "github.com/minio/madmin-go"
	cnf "github.com/rzrbld/adminio-api/config"
)

var MadmClnt, MadmErr = madmin.New(cnf.Server, cnf.Maccess, cnf.Msecret, cnf.Ssl)

// var MinioClnt, MinioErr = minio.New(cnf.Server, cnf.Maccess, cnf.Msecret, cnf.Ssl)

var MinioClnt, MinioErr = minio.New(cnf.Server, &minio.Options{
	Creds:  credentials.NewStaticV4(cnf.Maccess, cnf.Msecret, ""),
	Secure: cnf.Ssl,
})

func main() {
	if MadmErr != nil {
		log.Fatalln("Error while connecting via admin client ", MadmErr)
	}

	if MinioErr != nil {
		log.Fatalln("Error while connecting via minio client ", MinioErr)
	}
}
