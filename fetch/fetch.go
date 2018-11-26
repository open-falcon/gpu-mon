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

type uintmap map[uint]struct{}

//RawData 采集的单卡数据
type RawData struct {
	GpuID  uint
	Values MetricValues
}

func getValidValues(gpuID uint) MetricValues {
	dcgmSupported := int(1)
	common.Logger.WithFields(logrus.Fields{
		"gpuId": gpuID,
	}).Debug("Device support dcgm")

	values, err := fetchValues(gpuID)
	if err != nil {
		common.Logger.Error(err)
	}
	values.DcgmSupported = &dcgmSupported
	common.Logger.WithFields(logrus.Fields{
		"dcgmSupportedGPU": gpuID,
	}).Debug("Get DCGM supported Device Info")
	return values
}

func getEmptyValues(gpuID uint) MetricValues {
	dcgmSupported := -1
	common.Logger.WithFields(logrus.Fields{
		"gpuId": gpuID,
	}).Error("Device don't support dcgm")
	values := MetricValues{
		DcgmSupported: &dcgmSupported,
	}
	return values
}

func getRawdata(gpuID uint, dcgmGPUs uintmap) RawData {
	var values = MetricValues{}
	if _, ok := dcgmGPUs[gpuID]; ok {
		values = getValidValues(gpuID)
	} else {
		values = getEmptyValues(gpuID)
	}
	rawdata := RawData{
		GpuID:  gpuID,
		Values: values,
	}
	common.Logger.WithFields(logrus.Fields{
		"gpuId": gpuID,
	}).Info("successful fetch gpu info")
	return rawdata
}

func getDcgmGPUs() (uintmap, error) {
	// 获取支持DCGM的gpu设备，并构建DCGM map用于遍历
	gpus, err := dcgm.GetSupportedDevices()
	if err != nil {
		err = fmt.Errorf("failed to get GPU devices supporting DCGM | %s", err)
		return map[uint]struct{}{}, err
	}
	common.Logger.Debug("successful get GPU devices supporting dcgm")
	deviceSupportDCGM := make(map[uint]struct{})
	for _, gpuID := range gpus {
		deviceSupportDCGM[gpuID] = struct{}{}
	}
	return deviceSupportDCGM, err
}

func getRawdatalist() (rawdatas []RawData, err error) {
	// 获取所有GPU的数量
	count, err := dcgm.GetAllDeviceCount()
	if err != nil {
		return rawdatas, fmt.Errorf("unable to get GPU device number: %s", err)
	}
	common.Logger.Debug("Successful get GPU device number")

	dcgmGPUs, err := getDcgmGPUs()
	if err != nil {
		return rawdatas, err
	}
	rawdatas = []RawData{}
	for i := uint(0); i < count; i++ {
		rawdatas = append(rawdatas, getRawdata(i, dcgmGPUs))
	}
	return rawdatas, err
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
	common.Logger.Debug("Successful initialization of Dcgm")
	defer func() {
		if err := dcgm.Shutdown(); err != nil {
			common.Logger.Error(err) // Calls os.Exit(1) after loggings
		}
	}()
	datas, err := getRawdatalist()
	if err != nil {
		return datas, fmt.Errorf("unable to fetch dcgm status: %v", err)
	}
	common.Logger.Debug("Successful get Gpu Info")
	return datas, err
}
