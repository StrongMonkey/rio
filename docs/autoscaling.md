# Autoscaling 

Rio deploys a simple autoscaler to watch metrics from workloads and scale applications based on current in-flight requests.

Note: Metrics are scraped from the linkerd-proxy sidecar, this requires your application to be injected with the linkerd sidecar.
This will happen by default when running new workloads.

To enable autoscaling:

```bash
$ rio run --scale 1-10 -p 8080 strongmonkey1992/autoscale:v0

# to give a higher concurrency
$ rio run --scale 1-10 --concurrency 20 -p 8080 strongmonkey1992/autoscale:v0 

# to scale to zero
$ rio run --scale 0-10 -p 8080 strongmonkey1992/autoscale:v0
```

To test putting load on the service, use [hey](https://github.com/rakyll/hey):

```bash
hey -z 3m -c 60 http://xxx-xx.xxxxxx.on-rio-io

# watch the scale of your service
$ watch rio ps
```

Note: `concurrency` means the maximum in-flight requests each pod can take. If your total in-flight request is 60 and concurrency 
is 10, Rio will scale workloads to 6 replicas.

Note: When scaling an application to zero, the first request will take longer.

##### Troubleshooting

Check autoscaler logs for more details on how metrics are collected and scale decision was made.

```bash
$ rio -s logs autoscaler
```
