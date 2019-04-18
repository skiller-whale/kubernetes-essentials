docker build . -t skillerwhale/svc-json:v1         -f Dockerfile
docker build . -t skillerwhale/svc-json:v2         -f Dockerfile-v2
docker build . -t skillerwhale/svc-json:init       -f Dockerfile-init
docker build . -t skillerwhale/svc-json:slow       -f Dockerfile-slow
docker build . -t skillerwhale/svc-json:readiness  -f Dockerfile-readiness
docker build . -t skillerwhale/svc-json:with-cache -f Dockerfile-with-cache

docker push skillerwhale/svc-json:v1
docker push skillerwhale/svc-json:v2
docker push skillerwhale/svc-json:init
docker push skillerwhale/svc-json:slow
docker push skillerwhale/svc-json:readiness
docker push skillerwhale/svc-json:with-cache
