# adminio-api
This is a simple admin "REST" API for [minio](https://min.io/) s3 server. 
Here is a Web UI for this API - [adminio-ui](https://github.com/rzrbld/adminio-ui)

## run with docker
`` docker run rzrbld/adminio-api:0.2 ``

## how to run manually
1. [start](https://docs.min.io/) minio server
2. set env variables:

| Variable   |      Description      |  Default |
|--------------|:-----------------------:|-----------:|
| API_HOST_PORT | which host and port API should listening. This is Iris based API, so you will need to provide 0.0.0.0:8080 for listening on all interfaces | localhost:8080 |
| MINIO_HOST_PORT |  provide a minio server host and port  |  localhost:9000 |
| MINIO_SSL | enable or disable ssl |  false |
| MINIO_REGION | set minio region | us-east-1 |
| MINIO_ACCESS | set minio Access Key | test |
| MINIO_SECRET | set minio Secret Key | testtest123 |

3. run ./main form `dist` folder
