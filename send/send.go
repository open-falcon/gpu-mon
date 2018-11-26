/*
 * Copyright 2018 Xiaomi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package send

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/open-falcon/gpu-mon/common"
	"github.com/open-falcon/gpu-mon/fetch"
)

var sumMetric = map[string]int{
	"FBUsed":   -1,
	"GPUUtils": -1,
	"MemUtils": -1,
}

// MetaData 发送到falcon的数据结构
type MetaData struct {
	Metric      string      `json:"metric"`
	Endpoint    string      `json:"endpoint"`
	Timestamp   int64       `json:"timestamp"`
	Step        int         `json:"step"`
	Value       interface{} `json:"value"`
	CounterType string      `json:"counterType"`
	TAGS        string      `json:"tags"`
}

var falconAgent string

// 将metricValue值转化为meteData值
func buildMetaData(metricName string, metricValue interface{}, tags string, step int) (metaData MetaData) {
	globalConfig := common.Config()
	common.Logger.WithFields(logrus.Fields{
		"metricName":  metricName,
		"metricValue": metricValue,
		"tag":         tags}).Debug("Check value")

	metricValue, isNaN := checkMetricData(metricValue)
	//debug 时会记录异常数据
	if isNaN && globalConfig.Log.Level == "Debug" {
		common.Logger.WithFields(logrus.Fields{
			"metricName":  metricName,
			"metricValue": metricValue,
			"tag":         tags}).Warn("Annormal metric value")
	}

	endpoint := getEndPoint()
	common.Logger.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"metricName": metricName,
	}).Debug("Send message to falcon")

	metaData = MetaData{
		Metric:      metricName,
		Endpoint:    endpoint,
		Timestamp:   time.Now().Unix(),
		Step:        step,
		Value:       metricValue,
		CounterType: "GAUGE",
		TAGS:        tags,
	}
	return metaData
}

func updateSumMetric(metricName string, metricValue interface{}) {
	if sumValue, ok := sumMetric[metricName]; ok {
		value, _ := checkMetricData(metricValue)
		//sumValue == -1： sumMetric为原始状态
		if sumValue == -1 && value != -1 {
			sumMetric[metricName] = int(value.(uint))
		}
		if sumValue != -1 && value != -1 {
			sumMetric[metricName] += int(value.(uint))
		}
	}
}

func updateMetaDataList(rawData fetch.RawData, metaDataList []MetaData) []MetaData {
	gpuIDTag := "GpuId=" + strconv.Itoa(int(rawData.GpuID)) // Tag: GpuId=0 ...
	metricValues := rawData.Values
	// 遍历 metricValue结构体
	structType := reflect.TypeOf(metricValues)
	structValue := reflect.ValueOf(metricValues)
	for i := 0; i < structValue.NumField(); i++ {
		metricName := structType.Field(i).Name
		metricValue := structValue.Field(i).Interface()
		if isIgnore(metricName) {
			continue
		}
		metaData := buildMetaData(metricName, metricValue, gpuIDTag, 60)
		metaDataList = append(metaDataList, metaData)
		updateSumMetric(metricName, metricValue)
	}
	return metaDataList
}

func addSumMetric(metaDataList []MetaData, rawDataList []fetch.RawData) []MetaData {
	for sumMetricName, sumMetricValue := range sumMetric {
		metaData := buildMetaData(sumMetricName, &sumMetricValue, "type=sum", 60)
		metaDataList = append(metaDataList, metaData)
		aveValue := dvInt(sumMetricValue, len(rawDataList), 1)
		metaData = buildMetaData(sumMetricName, &aveValue, "type=ave", 60)
		metaDataList = append(metaDataList, metaData)
	}
	return metaDataList
}

// BuildMetaDatas 构建发送序列
func BuildMetaDatas(rawDataList []fetch.RawData) (metaDataList []MetaData) {
	for _, rawData := range rawDataList {
		metaDataList = updateMetaDataList(rawData, metaDataList)
	}
	metaDataList = addSumMetric(metaDataList, rawDataList)
	return metaDataList
}

func pushStdout(metaDataList []MetaData) error {
	js, err := json.Marshal(metaDataList)
	fmt.Println(string(js))
	if err != nil {
		return err
	}
	common.Logger.Info("Successful send message to stdout")
	return err
}

func sendDatas(js []byte, url string) error {
	res, err := http.Post(url,
		"Content-Type: application/json", bytes.NewBuffer(js))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code is not 200")
	}
	defer res.Body.Close()
	return nil
}

func pushAgent(metaDataList []MetaData) error {
	falconAgent := getFalconAgent()
	js, err := json.Marshal(metaDataList)
	if err != nil {
		common.Logger.Errorf("Convert metaData from struct to json failed: %v", err)
		return err
	}
	err = sendDatas(js, falconAgent)
	if err != nil {
		err = fmt.Errorf("send data failed: %v", err)
		common.Logger.Errorln(err)
		return err
	}
	common.Logger.WithFields(logrus.Fields{
		" falconAgent url": falconAgent,
	}).Info("Successful send message to falconAgent")
	return err
}

// Data 向falcon上报数据
func Data(metaDataList []MetaData) error {
	if common.Config().IsCrontab {
		return pushStdout(metaDataList)
	}
	return pushAgent(metaDataList)
}
