# Adminio-api
This is a simple admin "REST" API for [minio](https://min.io/) s3 server.
Here is a Web UI for this API - [adminio-ui](https://github.com/rzrbld/adminio-ui)

![Docker hub stats](https://img.shields.io/docker/pulls/rzrbld/adminio-api?style=flat-square) ![GitHub License](https://img.shields.io/github/license/rzrbld/adminio-api?style=flat-square)

## OpenAPI v3

see OpenAPI v3 specs at `openAPI/openapi_v3.yaml` or [html version](https://rzrbld.github.io/openapi/)

### Run full stack demo
obtain [docker-compose.yml](https://raw.githubusercontent.com/rzrbld/adminio-ui/master/docker-compose.yml) from [adminio-ui](https://github.com/rzrbld/adminio-ui) repository. And run it:
`docker-compose -f docker-compose.yml up`

it will bring up:

 - minio server on 9000 port
 - adminio API on 8080 port
 - adminio UI on 80 port

after that you can go to `http://localhost` and try out

### Run with docker
```
docker run -d \
      -p 8080:8080 \
      -e ADMINIO_HOST_PORT=":8080" \
      -e MINIO_HOST_PORT="localhost:9000" \
      -e MINIO_ACCESS="test" \
      -e MINIO_SECRET="testtest123" \
      rzrbld/adminio-api:latest

```

### Monitoring
Adminio-API expose metrics for [Prometheus](https://prometheus.io/) at `/metrics` if `ADMINIO_METRICS_ENABLE` is set to `true`.

### Run manually
 - [start](https://docs.min.io/) minio server
 - set env variables
 - go to `src` folder and compile with `go build main.go`, then run `./main` binary

### Config Env variables
| Variable   |      Description      |  Default |
|--------------|:-----------------------:|-----------:|
| `ADMINIO_HOST_PORT` | which host and port API should listening. This is Iris based API, so you will need to provide 0.0.0.0:8080 for listening on all interfaces | localhost:8080 |
| `MINIO_HOST_PORT` |  provide a minio server host and port  |  localhost:9000 |
| `MINIO_SSL` | enable or disable ssl |  false |
| `MINIO_REGION` | set minio region | us-east-1 |
| `MINIO_ACCESS` | set minio Access Key | test |
| `MINIO_SECRET` | set minio Secret Key | testtest123 |
| `ADMINIO_CORS_DOMAIN` | set adminio-api CORS policy domain  | * |
| `ADMINIO_OAUTH_ENABLE` | enable oauth over supported providers | false |
| `ADMINIO_OAUTH_PROVIDER` | oauth provider, for more information see the full list of supported providers | github |
| `ADMINIO_OAUTH_CLIENT_ID` | oauth app client id | my-github-oauth-app-client-id |
| `ADMINIO_OAUTH_CLIENT_SECRET` | oauth app secret | my-github-oauth-app-secret |
| `ADMINIO_OAUTH_CALLBACK` | oauth callback, default listener on /auth/callback | http://"+ADMINIO_HOST_PORT+"/auth/callback |
| `ADMINIO_OAUTH_CUSTOM_DOMAIN` | oauth custom domain, for supported providers (auth0\wso2) | - |
| `ADMINIO_COOKIE_HASH_KEY` | hash key for session cookies. AES only supports key sizes of 16, 24 or 32 bytes | NRUeuq6AdskNPa7ewZuxG9TrDZC4xFat |
| `ADMINIO_COOKIE_BLOCK_KEY` | block key for session cookies. AES only supports key sizes of 16, 24 or 32 bytes | bnfYuphzxPhJMR823YNezH83fuHuddFC |
| `ADMINIO_COOKIE_NAME` | name for the session cookie | adminiosessionid |
| `ADMINIO_AUDIT_LOG_ENABLE` | enable audit log, mae sense if oauth is enabled, othervise set to false | false |
| `ADMINIO_METRICS_ENABLE` | enable default iris/golang metrics and bucket sizes metric on `/metric/` uri path | false |
| `ADMINIO_PROBES_ENABLE` | enable liveness and readiness probes for k8s at `/ready/` and `/live/` uri path installations | false |

### List of supported oauth providers

 - amazon
 - auth0
 - bitbucket
 - box
 - digitalocean
 - dropbox
 - github
 - gitlab
 - heroku
 - onedrive
 - salesforce
 - slack
 - wso2


 ### Example config
 - prometheus config for adminio metrics: `examples/prometheus.yml`
 - bucket lifecycle: `examples/lifecycle.xml`
