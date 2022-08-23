package clients

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	madmin "github.com/minio/madmin-go"
	cnf "github.com/rzrbld/adminio-api/config"
)

var MadmErr error
var MinioErr error
var MinioClnt *minio.Client
var MadmClnt *madmin.AdminClient

func init() {
	tr, err := customTransport()
	if err != nil {
		MadmErr = err
		MinioErr = err
		return
	}

	MinioClnt, MinioErr = minio.New(cnf.Server, &minio.Options{
		Creds:     credentials.NewStaticV4(cnf.Maccess, cnf.Msecret, ""),
		Secure:    cnf.Ssl,
		Transport: tr,
	})

	MadmClnt, MadmErr = madmin.New(cnf.Server, cnf.Maccess, cnf.Msecret, cnf.Ssl)
	if err == nil {
		MadmClnt.SetCustomTransport(tr)
	}

	if MadmErr != nil {
		log.Fatalln("Error while connecting via admin client ", MadmErr)
	}

	if MinioErr != nil {
		log.Fatalln("Error while connecting via minio client ", MinioErr)
	}
}

func customTransport() (*http.Transport, error) {

	if !cnf.Ssl {
		return minio.DefaultTransport(cnf.Ssl)
	}

	tlsConfig := &tls.Config{
		// Can't use SSLv3 because of POODLE and BEAST
		// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
		// Can't use TLSv1.1 because of RC4 cipher usage
		MinVersion: tls.VersionTLS12,
	}

	tr, err := minio.DefaultTransport(cnf.Ssl)
	if err != nil {
		return nil, err
	}

	if cnf.SSLCACertFile != "" {
		minioCACert, err := os.ReadFile(cnf.SSLCACertFile)
		if err != nil {
			return nil, err
		}

		if !isValidCertificate(minioCACert) {
			return nil, fmt.Errorf("minio CA Cert is not a valid x509 certificate")
		}

		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			// In some systems (like Windows) system cert pool is
			// not supported or no certificates are present on the
			// system - so we create a new cert pool.
			rootCAs = x509.NewCertPool()
		}
		rootCAs.AppendCertsFromPEM(minioCACert)
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.RootCAs = rootCAs
	}

	if cnf.SSLSkipVerify {
		tlsConfig.InsecureSkipVerify = true
	}

	tr.TLSClientConfig = tlsConfig

	return tr, nil
}

func isValidCertificate(c []byte) bool {
	p, _ := pem.Decode(c)
	if p == nil {
		return false
	}
	_, err := x509.ParseCertificates(p.Bytes)
	return err == nil
}
