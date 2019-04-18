# svc-html

This program implements a simple HTTP service which listens on port `8080` connects to remote service specified via `REMOTE_URL` environment variable to query it's `JSON` API and returns the result as static `HTML`

If `REMOTE_URL` the service fails to start.

# docker

You can build the service using the attached `Dockerfile` by running the command below:
```shell
docker build -t skillerwhale/svc-html:v1 .
```

You can run the service using the command below:
```shell
docker run --rm -p 8080:8080 -e REMOTE_URL="http://127.0.0.1" skillerwhale/svc-html:v1
```

If you would like to run it as a daemon:
```shell
docker run -d -p 8080:8080 skillerwhale/svc-html:v1
```

You can now query the service using simple `curl` command:
```shell
$ curl -i 127.0.0.1:8080
HTTP/1.1 200 OK
Date: Thu, 18 Apr 2019 00:41:19 GMT
Content-Length: 53
Content-Type: text/html; charset=utf-8

<b>Instance: e2a8a7a3-1b7a-4980-a296-4eeb8349c01f</b>
```

# deploy to k8s

You can deploy the service to kubernetes using the command below
```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: skillerwhale-svc-html
spec:
  containers:
    - name: skillerwhale-svc-html
      image: skillerwhale/svc-html:v1
      imagePullPolicy: Always
      ports:
        - containerPort: 8080
      env:
        - name: REMOTE_URL
          value: "http://123.45.67.80"
EOF
```

Or using the `pod.yaml` file in this repo:
```shell
kubectl apply -f pod.yaml
```
