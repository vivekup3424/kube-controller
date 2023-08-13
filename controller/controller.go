package controller

import (
	"context"
	"fmt"
	"sync"

	"kube-controller/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

var MetricsMap types.NameSpaceMetric

func StartController(clientSet *kubernetes.Clientset, metricClientSet *metricsv.Clientset){

	factory := informers.NewFilteredSharedInformerFactory(clientSet, 0, metav1.NamespaceAll, nil)
	informer := factory.Apps().V1().Deployments().Informer()


	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) { 
			fmt.Println("Added deployment")
			// time.Sleep(30 * time.Second)
			Compute(clientSet, metricClientSet)
		},
		UpdateFunc: func(oldObj, newObj interface{}) { 
			fmt.Println("Updated deployment")
			Compute(clientSet, metricClientSet)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Deleted deployment")
			Compute(clientSet, metricClientSet)
		},
	})

	ch := make(chan struct{})
	fmt.Println("Starting controller...")
	go informer.Run(ch)
}

func Compute(clientSet *kubernetes.Clientset, metricClientSet *metricsv.Clientset){
	var nsMap types.NameSpaceMetric = make(types.NameSpaceMetric)
	nsList, err := clientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		fmt.Println(err.Error())
	}
	mapWg := &sync.WaitGroup{}
	var nsMutex sync.Mutex
	for _, ns := range nsList.Items {
		go func(ns corev1.Namespace) {
			cpu := 0.0
			memory := 0.0
			var resourcesMutex sync.Mutex
			deploymentList, err := clientSet.AppsV1().Deployments(ns.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil{
				fmt.Println(err.Error())
			}
			for _, deployment := range deploymentList.Items {
				go func(deployment appsv1.Deployment) {
					podMetricList, err := metricClientSet.MetricsV1beta1().PodMetricses(ns.Name).List(context.TODO(), metav1.ListOptions{
						LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
					})
					if err != nil{
						fmt.Println(err.Error())
					}
		
					for _, podMetric := range podMetricList.Items {
						go func(podMetric v1beta1.PodMetrics) {
							for _, containerMetric := range podMetric.Containers {
								mapWg.Add(1)
								go func(containerMetric v1beta1.ContainerMetrics) {
									defer mapWg.Done()
									resourcesMutex.Lock()
									cpu += containerMetric.Usage.Cpu().AsApproximateFloat64()
									memory += containerMetric.Usage.Memory().AsApproximateFloat64()
									resourcesMutex.Unlock()
								}(containerMetric)
							}
						}(podMetric)
					}
				}(deployment)
			}
			nsMutex.Lock()
			nsMap[ns.Name] = types.Metric{
				CPU: cpu,
				Memory: memory,
			}
			nsMutex.Unlock()
		}(ns)
	}
	mapWg.Wait()
	MetricsMap = nsMap
}