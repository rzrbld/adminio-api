# Adminio-api
This is a simple admin "REST" API for [minio](https://min.io/) s3 server. 
Here is a Web UI for this API - [adminio-ui](https://github.com/rzrbld/adminio-ui)

### Run full stack demo
obtain [docker-compose.yml](https://raw.githubusercontent.com/rzrbld/adminio-ui/master/docker-compose.yml) from [adminio-ui](https://github.com/rzrbld/adminio-ui) repository. And run it:
`` docker-compose -f docker-compose.yml up ``

it will bring up:

 -  minio server on 9000 port 
 - adminio API on 8080 port
 - adminio UI on 80 port

after that you can go to `` http://localhost `` and try out

### Run with docker
`` docker run rzrbld/adminio-api:0.2 ``

### Run manually
 - [start](https://docs.min.io/) minio server
 - set env variables
 - run ./main form `dist` folder

### Env variables
| Variable   |      Description      |  Default |
|--------------|:-----------------------:|-----------:|
| API_HOST_PORT | which host and port API should listening. This is Iris based API, so you will need to provide 0.0.0.0:8080 for listening on all interfaces | localhost:8080 |
| MINIO_HOST_PORT |  provide a minio server host and port  |  localhost:9000 |
| MINIO_SSL | enable or disable ssl |  false |
| MINIO_REGION | set minio region | us-east-1 |
| MINIO_ACCESS | set minio Access Key | test |
| MINIO_SECRET | set minio Secret Key | testtest123 |
