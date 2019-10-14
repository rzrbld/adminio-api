# adminio
This is a simple admin API for min.io (minio) s3 server

## how to run
1. start minio server
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
