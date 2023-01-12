package k8sLoggingMonitoring

import (
	"fmt"
	"log"
	"time"
)

type LoggingMonitoring struct {
	clusterName string
	logs        []string
	events      []string
	metrics     map[string]interface{}
}

func NewLoggingMonitoring(clusterName string) *LoggingMonitoring {
	return &LoggingMonitoring{
		clusterName: clusterName,
		logs:        []string{},
		events:      []string{},
		metrics:     make(map[string]interface{}),
	}
}

// CollectLogs to collect logs from given Kubernetes cluster
func (lm *LoggingMonitoring) CollectLogs() {
	// code to collect logs from cluster
	log.Printf("Collecting logs from %s cluster", lm.clusterName)
	lm.logs = append(lm.logs, "log1", "log2")
}

func (lm *LoggingMonitoring) StoreLogs() {
	log.Printf("Storing logs from %s cluster", lm.clusterName)
}

func (lm *LoggingMonitoring) AnalyzeLogs() {
	log.Printf("Analyzing logs from %s cluster", lm.clusterName)
	lm.events = append(lm.events, "event1", "event2")
	lm.metrics["timeStarted"] = time.Now()
}

func (lm *LoggingMonitoring) PrintData() {
	fmt.Printf("Logs:\n")
	for _, log := range lm.logs {
		fmt.Printf("- %s\n", log)
	}

	fmt.Printf("Events:\n")
	for _, event := range lm.events {
		fmt.Printf("- %s\n", event)
	}

	fmt.Printf("Metrics:\n")
	for key, value := range lm.metrics {
		fmt.Printf("- %s: %v\n", key, value)
	}
}

func (lm *LoggingMonitoring) CollectPodLogs() {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// list all Pods in the cluster
	pods, err := clientset.CoreV1().Pods(lm.clusterName).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing Pods: %s", err)
	}

	// collect logs from each Pod
	for _, pod := range pods.Items {
		req := clientset.CoreV1().Pods(lm.clusterName).GetLogs(pod.Name, &corev1.PodLogOptions{})
		logs, err := req.Do().Raw()
		if err != nil {
			log.Printf("Error getting logs for Pod %s: %s", pod.Name, err)
			continue
		}
		lm.logs = append(lm.logs, string(logs))
	}
}
