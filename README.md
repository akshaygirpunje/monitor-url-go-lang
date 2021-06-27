# Monitor URLs

Monitor a list of URLs with Golang instrumented with Prometheus served on a Kubernetes Cluster

---

## Summary

-   A service written in Golang that queries two sample urls every 5 seconds:
    -   https://httpstat.us/503
    -   https://httpstat.us/200
-   The service checks:
    -   The external urls are up (based on http status code 200) return `1` if up, `0` if otherwise
    -   Response time in milliseconds
-   The service will run a simple http service that produces metrics (on `/metrics`) and output a Prometheus format when curling the service `/metrics` url

**Sample Response Format**:
```shell
sample_external_url_up{url="https://httpstat.us/200 "}  = 1
sample_external_url_response_ms{url="https://httpstat.us/200 "}  = [value]
sample_external_url_up{url="https://httpstat.us/503 "}  = 0
sample_external_url_response_ms{url="https://httpstat.us/503 "}  = [value]
```
___

## Technology Used

-   [Golang](https://golang.org/) - 1.15+
-   [Prometheus](https://github.com/prometheus/client_golang)
-   [Kubernetes](https://kubernetes.io/)
-   [Helm](https://helm.sh/)

## Project Configuration

This project uses [Go modules](https://blog.golang.org/using-go-modules) to manage dependencies

The `go` command resolves imports by using the specific dependency module versions listed in [go.mod](go.mod)

---

## Set-up

1. Configure [conf.json](conf.json) with URLs you wish to monitor. This is currently configured with two urls as an example.

```json
{
    "urls": ["https://httpstat.us/200", "https://httpstat.us/503"]
}
```

2. Build Docker image and push to repository of your choosing

```shell
docker build -t $USERNAME/monitor-urls .
docker push $USERNAME/monitor-urls
```

3. Create `monitoring` namespace

```shell
kubectl create namespace monitoring
```

4. Use `helm` to install Prometheus to namespace

```shell
helm install prometheus prometheus-community/prometheus --namespace monitoring

NAME: prometheus
LAST DEPLOYED: Thu Sep 10 20:15:45 2020
NAMESPACE: monitoring
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
The Prometheus server can be accessed via port 80 on the following DNS name from within your cluster:
prometheus-server.monitoring.svc.cluster.local


Get the Prometheus server URL by running these commands in the same shell:
  export POD_NAME=$(kubectl get pods --namespace monitoring -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
  kubectl --namespace monitoring port-forward $POD_NAME 9090
```

---

### Local Testing (Docker + MiniKube)

To locally test with [MiniKube](https://kubernetes.io/docs/tasks/tools/install-minikube/) ensure that it is properly installed according to your operating system and started through `minikube start`

1. Run `kubectl apply`

```shell
kubectl apply -f service.yml

service/monitor-urls created
deployment.apps/monitor-urls created
```

Note: In [service.yml](service.yml) change `image: akshaygirpunje/monitor-urls:latest` to newly built Docker image you done in the set-up

-   View the deployment

```shell
kubectl get deployments

NAME                            READY   UP-TO-DATE   AVAILABLE   AGE
monitor-urls                    1/1     1            1           11m
```

-   View the service

```shell
kubectl get services

NAME                            TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
monitor-urls                    LoadBalancer   10.104.174.222   <pending>     8080:32733/TCP   12m
```

2. Run service through `MiniKube`

```shell
minikube service monitor-urls

|-----------|--------------|-------------|-------------------------|
| NAMESPACE |     NAME     | TARGET PORT |           URL           |
|-----------|--------------|-------------|-------------------------|
| default   | monitor-urls |        8080 | http://172.17.0.3:32733 |
|-----------|--------------|-------------|-------------------------|
🏃  Starting tunnel for service monitor-urls.
|-----------|--------------|-------------|------------------------|
| NAMESPACE |     NAME     | TARGET PORT |          URL           |
|-----------|--------------|-------------|------------------------|
| default   | monitor-urls |             | http://127.0.0.1:58275 |
|-----------|--------------|-------------|------------------------|
🎉  Opening service default/monitor-urls in default browser...
❗  Because you are using a Docker driver on darwin, the terminal needs to be open to run it.
```

3. Test!

```shell
curl http://127.0.0.1:58275/metrics

# HELP sample_external_url_response_ms HTTP response in milliseconds
# TYPE sample_external_url_response_ms gauge
sample_external_url_response_ms{url="https://httpstat.us/200"} 129
sample_external_url_response_ms{url="https://httpstat.us/503"} 120
# HELP sample_external_url_up Boolean status of site up or down
# TYPE sample_external_url_up gauge
sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

-   Check Kubernetes Services Dashboard (Optional)

```shell
minikube dashboard

🔌  Enabling dashboard ...
🤔  Verifying dashboard health ...
🚀  Launching proxy ...
🤔  Verifying proxy health ...
🎉  Opening http://127.0.0.1:58952/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/ in your default browser...
```

---

### Local Testing (Golang Only)

1. Run `go run`

```shell
go run main.go helper.go
```

2. Curl `localhost`

```shell
curl localhost:8080/metrics

# HELP sample_external_url_response_ms HTTP response in milliseconds
# TYPE sample_external_url_response_ms gauge
sample_external_url_response_ms{url="https://httpstat.us/200"} 362
sample_external_url_response_ms{url="https://httpstat.us/503"} 453
# HELP sample_external_url_up Boolean status of site up or down
# TYPE sample_external_url_up gauge
sample_external_url_up{url="https://httpstat.us/200"} 1
sample_external_url_up{url="https://httpstat.us/503"} 0
```

---

## Tests

```shell
go test

PASS
ok      github.com/akshaygirpunje/monitor-urls-k8s  0.594s
```
