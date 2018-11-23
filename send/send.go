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

	"github.com/open-falcon/gpu-mon/common"
	"github.com/open-falcon/gpu-mon/fetch"

	"github.com/sirupsen/logrus"
)

// 将metricValue值转化为meteData值
func buildMetaData(metricName string, metricValue interface{}, tags string, step int) (metaData MetaData) {
	globalConfig := cfg.Config()
	if globalConfig.Log.Level == "Debug" {
		cfg.CommonLogger.WithFields(logrus.Fields{
			"metricName":  metricName,
			"metricValue": metricValue,
			"tag":         tags}).Info("check value value")
	}

	metricValue, isNaN := checkMetricData(metricValue)
	if isNaN && globalConfig.Log.Level == "Debug" {
		cfg.CommonLogger.WithFields(logrus.Fields{
			"metricName":  metricName,
			"metricValue": metricValue,
			"tag":         tags}).Warn("metric value ")
	}

	endpoint := getEndPoint()
	cfg.CommonLogger.WithFields(logrus.Fields{
		"endpoint":   endpoint,
		"metricName": metricName,
	}).Info("send message to falcon")

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

// BuildMetaDatas 构建发送序列
func BuildMetaDatas(rawDataList []fetch.RawData) (metaDataList []MetaData) {
	for _, rawData := range rawDataList {
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
	}

	for sumMetricName, sumMetricValue := range sumMetric {
		metaData := buildMetaData(sumMetricName, &sumMetricValue, "type=sum", 60)
		metaDataList = append(metaDataList, metaData)
		aveValue := dvInt(sumMetricValue, len(rawDataList), 1)
		metaData = buildMetaData(sumMetricName, &aveValue, "type=ave", 60)
		metaDataList = append(metaDataList, metaData)
	}
	return metaDataList
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

func pushStdout(metaDataList []MetaData) error {
	js, err := json.Marshal(metaDataList)
	fmt.Println(string(js))
	return err
}

func pushAgent(metaDataList []MetaData) error {
	isError := false
	falconAgent := getFalconAgent()
	cfg.CommonLogger.WithFields(logrus.Fields{
		" falconAgent url": falconAgent,
	}).Info("send message to falconAgent")

	for _, metaData := range metaDataList {
		js, err := json.Marshal([]MetaData{metaData})
		if err != nil {
			isError = true
			cfg.CommonLogger.WithFields(logrus.Fields{
				"error info:": err,
				"metric name": metaData.Metric,
			}).Error("convert metaData from struct to json failed")
			continue
		}
		func([]byte) {
			res, err := http.Post(falconAgent,
				"Content-Type: application/json", bytes.NewBuffer(js))
			if err != nil {
				cfg.CommonLogger.WithFields(logrus.Fields{
					"error info:": err,
					"metric name": metaData.Metric,
				}).Error("send Data Error")
				isError = true
				return
			}
			if res.StatusCode != 200 {
				cfg.CommonLogger.WithFields(logrus.Fields{
					"StatusCode":  res.StatusCode,
					"metric name": metaData.Metric,
				}).Error("send Data failed")
				isError = true
				return
			}
			defer res.Body.Close()
		}(js)
	}
	if isError {
		return fmt.Errorf("something error when push to agent")
	}
	return nil
}

// Data 向falcon上报数据
func Data(metaDataList []MetaData) error {
	if cfg.Config().IsCrontab {
		return pushStdout(metaDataList)
	}
	return pushAgent(metaDataList)
}
