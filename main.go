package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"

	"kube-controller/controller"
	"kube-controller/types"
)

var globalClientSet *kubernetes.Clientset
var globalMetricClientSet *metricsv.Clientset

func main(){
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	globalClientSet = clientSet

	metricClientSet, err := metricsv.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	globalMetricClientSet = metricClientSet

	controller.StartController(clientSet, metricClientSet)
	
	http.HandleFunc("/metrics", GetMetrics)
	http.HandleFunc("/update", UpdateMetrics)

	fmt.Println("Listening on port 8000...")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(types.MetricResponse{
		Metrics: controller.MetricsMap,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Fprintf(w, "%s", string(res))
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	controller.Compute(globalClientSet, globalMetricClientSet)
	fmt.Fprintf(w, "%s", "Updated metrics")
}