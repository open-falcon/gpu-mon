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
package fetch

import (
	"fmt"

	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/dcgm"
	"github.com/open-falcon/gpu-mon/common"
	"github.com/sirupsen/logrus"
)

func updateRawDataList(gpuNum uint, rawDataList []RawData, deviceSupportDCGM map[uint]struct{}) []RawData {
	var dcgmSupported = -1
	var values MetricValues //
	if _, ok := deviceSupportDCGM[gpuNum]; ok {
		gpuID := gpuNum
		//gpu支持DCGM
		dcgmSupported = 1
		common.Logger.WithFields(logrus.Fields{
			"gpuId": gpuID,
		}).Info("device support dcgm")
		values, err := fetchValues(gpuID)
		if err != nil {
			common.Logger.Error(err)
		}
		common.Logger.WithFields(logrus.Fields{
			"dcgmSupportedGPU": gpuID,
		}).Info("get DCGM supported Device Info")
		values.DcgmSupported = &dcgmSupported
	} else {
		common.Logger.WithFields(logrus.Fields{
			"gpuId": gpuNum,
		}).Error("device do not support dcgm")
		values.DcgmSupported = &dcgmSupported
	}

	rawData := RawData{
		GpuID:  gpuNum,
		Values: values,
	}
	rawDataList[gpuNum] = rawData
	common.Logger.WithFields(logrus.Fields{
		"gpuId": gpuNum,
	}).Info("successful fetch gpu info")
	return rawDataList
}

func buildDataList(count uint, deviceSupportDCGM map[uint]struct{}) []RawData {
	rawDataList := make([]RawData, count)
	for i := uint(0); i < count; i++ {
		rawDataList = updateRawDataList(i, rawDataList, deviceSupportDCGM)
	}
	return rawDataList
}
func getSupportedDevices() (map[uint]struct{}, error) {
	// 获取支持DCGM的gpu设备，并构建DCGM map用于遍历
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		err = fmt.Errorf("failed to get GPU devices supporting DCGM | %s", err)
		return map[uint]struct{}{}, err
	}

	common.Logger.Info("successful get GPU devices supporting dcgm")

	deviceSupportDCGM := make(map[uint]struct{})
	for _, gpuID := range gpus {
		deviceSupportDCGM[gpuID] = struct{}{}
	}
	return deviceSupportDCGM, err
}
func fetchData() ([]RawData, error) {
	var rawDataList []RawData
	// 获取所有GPU的数量
	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		return rawDataList, fmt.Errorf("unable to get GPU device number | %s", err)
	}
	common.Logger.Info("successful get all gpu device")

	deviceSupportDCGM, err := getSupportedDevices()
	if err != nil {
		return rawDataList, err
	}

	rawDataList = buildDataList(count, deviceSupportDCGM)
	return rawDataList, nil
}

//Data 获取监控数据
func Data() ([]RawData, error) {

	// 初始化 DCGM
	err := dcgm.Init(dcgm.Standalone, "localhost", "0")
	if err != nil {
		err = fmt.Errorf("initialization of DCGM failed | %s", err)
		common.Logger.Error(err)
		return nil, err
	}
	common.Logger.Info("Successful initialization of Dcgm")
	defer func() {
		if err := dcgm.Shutdown(); err != nil {
			common.Logger.Error(err) // Calls os.Exit(1) after loggings
		}
	}()

	rawDataList, err := fetchData()
	return rawDataList, err
}
