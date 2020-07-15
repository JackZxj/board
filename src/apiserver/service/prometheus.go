package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type DashboardInfo struct {
	IsOverMaxLimit  bool          `json:"is_over_max_limit"`
	IsOverMinLimit  bool          `json:"is_over_min_limit"`
	NodeCount       int           `json:"node_count"`
	NodeListData    []NodeList    `json:"node_list_data"`
	ServiceCount    int           `json:"service_count"`
	ServiceListData []ServiceList `json:"service_list_data"`
	TimeUnit        string        `json:"time_unit"`
}

type NodeList struct {
	Name         string     `json:"name"`
	NodeLogsData []NodeLogs `json:"node_logs_data"`
}

type NodeLogs struct {
	TimeStamp    int64   `json:"timestamp"`
	CPUUsage     float64 `json:"cpu_usage"`
	MemoryUsage  float64 `json:"memory_usage"`
	StorageTotal int     `json:"storage_total"`
	StorageUsed  int     `json:"storage_used"`
}

type ServiceList struct {
	Name            string        `json:"name"`
	ServiceLogsData []ServiceLogs `json:"service_logs_data"`
}

type ServiceLogs struct {
	TimeStamp       int64 `json:"timestamp"`
	ContainerNumber int   `json:"container_number"`
	PodNumber       int   `json:"pod_number"`
}

type RequestPayload struct {
	TimeStamp int64  `json:"timestamp"`
	TimeCount int    `json:"time_count"`
	TimeUnit  string `json:"time_unit"`
}

func GetDashBoardData(request RequestPayload, nodename, servicename string) (DashboardInfo, error) {
	var para DashboardInfo

	client, err := api.NewClient(api.Config{
		Address: "http://prometheus:9090/",
	})
	if err != nil {
		return DashboardInfo{}, err
	}
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	para.TimeUnit = request.TimeUnit
	timeStampArray := make([]int64, request.TimeCount)
	var timeRange v1.Range
	switch request.TimeUnit {
	case "second":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*5)
		}
		timeRange.Step = time.Second * 5
	case "minute":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*60)
		}
		timeRange.Step = time.Minute
	case "hour":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*3600)
		}
		timeRange.Step = time.Hour
	case "day":
		for i := range timeStampArray {
			timeStampArray[i] = request.TimeStamp - int64((len(timeStampArray)-i)*3600*24)
		}
		timeRange.Step = time.Hour * 24
	default:
		return DashboardInfo{}, errors.New("wrong time unit")
	}
	timeRange.Start = time.Unix(timeStampArray[0], 0)
	timeRange.End = time.Unix(timeStampArray[request.TimeCount-1], 0)

	podInfo, err := para.GetServiceInfo(ctx, v1api, timeRange, timeStampArray)
	if err != nil {
		return DashboardInfo{}, err
	}
	containerInfo, err := para.CountPod(ctx, v1api, timeRange, timeStampArray, podInfo)
	if err != nil {
		return DashboardInfo{}, err
	}
	para.CountContainer(timeStampArray, containerInfo)

	//filtering
	for i := 0; i < len(para.ServiceListData); i++ {
		if para.ServiceListData[i].Name != servicename {
			para.ServiceListData[i].ServiceLogsData = []ServiceLogs{}
		}
	}

	//init node info
	kubeNodeQuery := `kube_node_info`
	kubeNodeResult, _, err := v1api.QueryRange(ctx, kubeNodeQuery, timeRange)
	if err != nil {
		return DashboardInfo{}, err
	}
	linesOfNode := strings.Split(kubeNodeResult.String(), kubeNodeQuery)
	para.NodeCount = len(linesOfNode) - 1
	para.NodeListData = make([]NodeList, len(linesOfNode))
	para.NodeListData[0].Name = "average"
	para.NodeListData[0].NodeLogsData = make([]NodeLogs, request.TimeCount)
	for i, v := range linesOfNode[1:] {
		para.NodeListData[i+1].Name = grepContent(v, `, node="([^"]+)"`)[0][1]
		para.NodeListData[i+1].NodeLogsData = make([]NodeLogs, request.TimeCount)
		for j := 0; j < request.TimeCount; j++ {
			para.NodeListData[i+1].NodeLogsData[j].TimeStamp = timeStampArray[j]
		}
	}

	nodeInfoSelector := []string{"storageCap", "storageUsed", "memory", "CPU"}
	nodeInfoQuery := []string{
		`kube_node_status_capacity{resource="ephemeral_storage"}`,
		`kube_node_status_capacity{resource="ephemeral_storage"} - kube_node_status_allocatable{resource="ephemeral_storage"}`,
		`(1 - kube_node_status_allocatable_memory_bytes / kube_node_status_capacity_memory_bytes) * 100`,
		`100 * (1 - sum by (instance)(node_cpu_seconds_total{mode="idle"}) / sum by (instance)(node_cpu_seconds_total))`}

	for i, q := range nodeInfoQuery {
		err = para.GetData(q, nodeInfoSelector[i], v1api, ctx, timeRange)
		if err != nil {
			return DashboardInfo{}, err
		}
	}

	//filtering
	for i := 0; i < len(para.NodeListData); i++ {
		if para.NodeListData[i].Name != nodename {
			para.NodeListData[i].NodeLogsData = []NodeLogs{}
		}
	}

	if nodename == "average" {
		for i, q := range nodeInfoQuery {
			avgQuery := fmt.Sprintf("avg(%s)", q)
			err = para.GetAvgData(avgQuery, nodeInfoSelector[i], v1api, ctx, timeRange, timeStampArray)
			if err != nil {
				return DashboardInfo{}, err
			}
		}
	}

	return para, nil
}

func grepString(src, reg string) []string {
	re, _ := regexp.Compile(reg)
	return re.FindAllString(src, -1)
}

func grepContent(src, reg string) [][]string {
	regex := regexp.MustCompile(reg)
	return regex.FindAllStringSubmatch(src, -1)
}

func sliceToPrometheusMetrics(src []string) string {
	return strings.Join(src, ", ")
}

func (d *DashboardInfo) GetAvgData(query, which string, v1api v1.API, ctx context.Context, timeRange v1.Range, timeStampArray []int64) error {
	result, _, err := v1api.QueryRange(ctx, query, timeRange)
	if err != nil {
		return err
	}
	data := grepContent(result.String(), "\n([0-9.]+)")
	for j, w := range data {
		var digit interface{}
		switch which {
		case "storageUsed", "storageCap":
			digit, err = strconv.Atoi(string(w[1]))
		case "CPU", "memory":
			digit, err = strconv.ParseFloat(string(w[1]), 64)
		}
		if err != nil {
			return err
		}
		switch which {
		case "CPU":
			d.NodeListData[0].NodeLogsData[j].TimeStamp = timeStampArray[j]
			d.NodeListData[0].NodeLogsData[j].CPUUsage = digit.(float64)
		case "memory":
			d.NodeListData[0].NodeLogsData[j].MemoryUsage = digit.(float64)
		case "storageCap":
			d.NodeListData[0].NodeLogsData[j].StorageTotal = digit.(int)
		case "storageUsed":
			d.NodeListData[0].NodeLogsData[j].StorageUsed = digit.(int)
		}
	}
	return nil
}

func (d *DashboardInfo) GetData(query, which string, v1api v1.API, ctx context.Context, timeRange v1.Range) error {
	result, _, err := v1api.QueryRange(ctx, query, timeRange)
	if err != nil {
		return err
	}
	var delimiter string
	switch which {
	case "storageCap":
		delimiter = "kube_node_status_capacity"
	default:
		delimiter = "=>"
	}
	lines := strings.Split(result.String(), delimiter)
	for i, v := range lines[1:] {
		data := grepContent(v, "\n([0-9.]+)")
		for j, w := range data {
			var digit interface{}
			switch which {
			case "storageUsed", "storageCap":
				digit, err = strconv.Atoi(w[1])
			case "CPU", "memory":
				digit, err = strconv.ParseFloat(w[1], 64)
			}
			if err != nil {
				return err
			}
			switch which {
			case "CPU":
				d.NodeListData[i+1].NodeLogsData[j].CPUUsage = digit.(float64)
			case "memory":
				d.NodeListData[i+1].NodeLogsData[j].MemoryUsage = digit.(float64)
			case "storageUsed":
				d.NodeListData[i+1].NodeLogsData[j].StorageUsed = digit.(int)
			case "storageCap":
				d.NodeListData[i+1].NodeLogsData[j].StorageTotal = digit.(int)
			}
		}
	}
	return nil
}

func (d *DashboardInfo) GetServiceInfo(ctx context.Context, v1api v1.API, timeRange v1.Range, timeStampArray []int64) ([][]string, error) {
	serviceQuery := `kube_service_spec_selector{namespace!~"cadvisor|istio-system|kube-public|kube-system", service!="kubernetes"}`
	result, warnings, err := v1api.QueryRange(ctx, serviceQuery, timeRange)
	if err != nil {
		return [][]string{}, err
	}
	if len(warnings) > 0 {
		logs.Info("Warnings: %v\n", warnings)
	}
	lines := grepString(result.String(), "{[^}]*}")
	d.ServiceListData = make([]ServiceList, len(lines)+1)
	d.ServiceListData[0].Name = "total"
	d.ServiceListData[0].ServiceLogsData = make([]ServiceLogs, len(timeStampArray))
	for s, time := range timeStampArray {
		d.ServiceListData[0].ServiceLogsData[s].TimeStamp = time
	}
	podResults := make([][]string, len(lines))
	for i, line := range lines {
		d.ServiceCount++
		d.ServiceListData[i+1].Name = grepContent(line, `service="([^"]+)"`)[0][1]
		serviceSelectorLabels := grepString(line, "label_[^,}]+")
		if len(serviceSelectorLabels) == 0 {
			continue
		}
		d.ServiceListData[i+1].ServiceLogsData = make([]ServiceLogs, len(timeStampArray))
		for s, time := range timeStampArray {
			d.ServiceListData[i+1].ServiceLogsData[s].TimeStamp = time
		}
		grepName := "kube_pod_labels{" + sliceToPrometheusMetrics(serviceSelectorLabels) + "}"
		pods, _, err := v1api.QueryRange(ctx, grepName, timeRange)
		if err != nil {
			return [][]string{}, err
		}
		podlines := strings.Split(pods.String(), "kube_pod_labels")
		podResults[i] = podlines[1:]
	}
	return podResults, nil
}

func (d *DashboardInfo) CountPod(ctx context.Context, v1api v1.API, timeRange v1.Range, timeStampArray []int64, podResults [][]string) ([][]string, error) {
	containerResults := make([][]string, len(podResults))
	for i, svc := range podResults {
		if len(svc) == 0 {
			continue
		}
		for _, podline := range svc {
			for s, time := range timeStampArray {
				if strings.Contains(podline, strconv.FormatInt(time, 10)) {
					d.ServiceListData[0].ServiceLogsData[s].PodNumber++
					d.ServiceListData[i+1].ServiceLogsData[s].PodNumber++
				}
			}
			podName := grepString(podline, "pod=[^,}]+")
			containerGrepName := "kube_pod_container_info{" + sliceToPrometheusMetrics(podName) + "}"
			containers, _, err := v1api.QueryRange(ctx, containerGrepName, timeRange)
			if err != nil {
				return [][]string{}, err
			}
			containerlines := strings.Split(containers.String(), "kube_pod_container_info")
			containerResults[i] = append(containerResults[i], containerlines[1:]...)
		}
	}
	return containerResults, nil
}

func (d *DashboardInfo) CountContainer(timeStampArray []int64, containerResults [][]string) {
	for i, svc := range containerResults {
		if len(svc) == 0 {
			continue
		}
		for _, containerline := range svc {
			for s, time := range timeStampArray {
				if strings.Contains(containerline, strconv.FormatInt(time, 10)) {
					d.ServiceListData[0].ServiceLogsData[s].ContainerNumber++
					d.ServiceListData[i+1].ServiceLogsData[s].ContainerNumber++
				}
			}
		}
	}
}