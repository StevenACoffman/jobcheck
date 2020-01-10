# jobcheck
### Example kubernetes job with file health checking

Restricting what's in your runtime container to precisely what's necessary for your app is a best practice employed by Google and other tech giants that have used containers in production for many years. It improves the signal to noise of scanners (e.g. CVE) and reduces the burden of establishing provenance to just what you need.

Many applications running for long periods of time eventually transition to broken states, and cannot recover except by being restarted. Kubernetes provides liveness probes to detect and remedy such situations.

A sentinel file that has it's timestamp modified on a regular basis can provide proof that the application is still responsive.

This small program is able to create and regularly update the timestamp a sentinel file, and check that the last modification of the file does not exceed a threshold.

### Job Metrics

[Prometheus](https://prometheus.io/) is a very well known open source monitoring solution. It provides a framework using which users can collect metrics from various sources, store them and create charts or alert based on these metrics.

Prometheus by default prefers a pull based metrics collection. That means Prometheus will periodically collect metrics from your systems and store it in its database. But once in a while there comes a situation where you cannot modify your system to expose the metrics like Prometheus needs. Or sometimes it is possible that your system itself cannot hold the data till Prometheus polls. In these cases its better to push metrics to Prometheus.

For pushing metrics to Prometheus, you need run one more piece of software called [Prometheus Push Gateway](https://github.com/prometheus/pushgateway). Then you need to configure Push gateway as one of the targets that Prometheus needs to collect metrics from. Then from your application you push metrics to push gateway.

Push Gateway is available for download [here](https://github.com/prometheus/pushgateway/releases).

The packages available in Golang for Prometheus are,
```
"github.com/prometheus/client_golang/prometheus/push"
"github.com/prometheus/client_golang/prometheus"
```
Sample code to push to the Prometheus Gateway
```go

gatewayUrl:="http://my-pushgateway-vm1:9091/"

throughputGuage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "throughput",
		Help: "Throughput in Mbps",
})
throughputGuage.Set(800)

if err := push.Collectors(
		"throughput_job", push.HostnameGroupingKey(),
		gatewayUrl, throughputGuage
); err != nil {
	fmt.Println("Could not push completion time to Pushgateway:", err)
}
```