package clients

import (
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	madmin "github.com/minio/minio/pkg/madmin"
	cnf "github.com/rzrbld/adminio-api/config"
	"log"
)

var MadmClnt, MadmErr = madmin.New(cnf.Server, cnf.Maccess, cnf.Msecret, cnf.Ssl)

// var MinioClnt, MinioErr = minio.New(cnf.Server, cnf.Maccess, cnf.Msecret, cnf.Ssl)

var MinioClnt, MinioErr = minio.New(cnf.Server, &minio.Options{
	Creds:  credentials.NewStaticV4(cnf.Maccess, cnf.Msecret, ""),
	Secure: cnf.Ssl,
})

func main() {
	if MadmErr != nil {
		log.Fatal("Error while connecting via admin client ", MadmErr)
	}

	if MinioErr != nil {
		log.Fatal("Error while connecting via minio client ", MinioErr)
	}
}
