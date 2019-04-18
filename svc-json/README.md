# svc-json

This program implements a simple HTTP service which listens on port `80` and serves randomly generated [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier) encoded in `JSON` format as its instance `ID`.

You can specify the instance `ID` via environment variable `ID`. If no variable has been exported, the server generates its own according to `RFC 4122`.

# docker

You can build the service using the attached `Dockerfile` by running the command below:
```shell
docker build -t skillerwhale/svc-json:v1 .
```

You can run the service using the command below:
```shell
docker run --rm -p 80:80 -e ID="foobar" skillerwhale/svc-json:v1
2019/04/18 00:08:34 Starting HTTP server ID: foobar Addr: :80...
```

If you would like to run it as a daemon:
```shell
docker run -d -p 80:80 -e ID="foobar" skillerwhale/svc-json:v1
```

You dont have to specify `ID` environment variable:
```shell
docker run --rm -p 80:80 skillerwhale/svc-json:v1
2019/04/18 00:09:38 Starting HTTP server ID: 41c55c62-e07b-4fd5-835f-8f28525e13ba Addr: :80...
```

You can now query the instance ID using simple `curl` command:
```shell
$ curl -i localhost:80
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 18 Apr 2019 00:12:14 GMT
Content-Length: 46

{"id":"41c55c62-e07b-4fd5-835f-8f28525e13ba"}
```

# deploy to k8s

You can deploy the service to kubernetes using the command below
```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: skillerwhale-svc-json
spec:
  containers:
    - name: skillerwhale-svc-json
      image: skillerwhale/svc-json:v1
      imagePullPolicy: Always
      ports:
        - containerPort: 80
      env:
        - name: ID
          value: "foo-bar-id"
EOF
```

Or using the `pod.yaml` file in this repo:
```shell
kubectl apply -f pod.yaml
```

# environment variables

Besides  `ID` environment variable, the program allows to specify many other environment variables that control the program's behaviour:
* `ADDR`: allows to specify HTTP service listen address `URL`
* `DELAY`: simulates HTTP response delay in `ms` (each response is delayed by `DELAY` ms)
* `VERSION`: specifies version of the `JSON` data to be returned
* `READY`: allow to specify number of `seconds` after which `/readyz` HTTP endpoint starts returning `200` HTTP response code
* `CACHE_PATH`: if specified, the service will return `500` unless the fie in the provided path exists (or is created)
