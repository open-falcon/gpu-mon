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

//RawData 采集的单卡数据
type RawData struct {
	GpuID  uint
	Values MetricValues
}

// MetricValues 监控项的值
type MetricValues struct {
	GPUUtils             *uint    // % DCGM_FI_DEV_GPU_UTIL  GPU Utilization
	MemUtils             *uint    // % DCGM_FI_DEV_MEM_COPY_UTIL  Memory Utilization
	Encoder              *uint    // % DCGM_FI_DEV_ENC_UTIL  Encoder Utilization
	Decoder              *uint    // % DCGM_FI_DEV_DEC_UTIL  Decoder Utilization
	SingleBitError       *uint    // DCGM_FI_DEV_ECC_SBE_AGG_TOTAL Total single bit aggregate (persistent) ECC errors
	DoubleBitError       *uint    // DCGM_FI_DEV_ECC_DBE_AGG_TOTAL Total double bit aggregate (persistent) ECC errors
	SmClock              *uint    // MHz DCGM_FI_DEV_SM_CLOCK
	MemClock             *uint    // MHz DCGM_FI_DEV_MEM_CLOCK
	BAR1Used             *uint    // MB DCGM_FI_DEV_BAR1_USED Used BAR1 of the GPU in MB
	Rx                   *uint64  // MB DCGM_FI_DEV_PCIE_RX_THROUGHPUT  PCIe Rx utilization information
	Tx                   *uint64  // MB DCGM_FI_DEV_PCIE_TX_THROUGHPUT  PCIe Tx utilization information
	Replays              *uint64  // DCGM_FI_DEV_PCIE_REPLAY_COUNTER  PCIe replay counter
	Performance          *uint    // DCGM_FI_DEV_PSTATE Performance state (P-State) 0-15. 0=highest
	FanSpeed             *uint    // % DCGM_FI_DEV_FAN_SPEED  Fan speed for the device in percent 0-100
	PowerUsed            *float64 // W DCGM_FI_DEV_POWER_USAGE Power usage for the device in Watts
	FBUsed               *uint    // MB DCGM_FI_DEV_FB_USED   Used Frame Buffer in MB
	DeviceTemperature    *uint    // °C DCGM_FI_DEV_GPU_TEMP  Current temperature readings for the device, in degrees C
	MemTemperature       *uint    // °C DCGM_FI_DEV_MEMORY_TEMP  Memory temperature for the device
	SlowdownTemperature  *uint    // °C DCGM_FI_DEV_SLOWDOWN_TEMP Slowdown temperature for the device
	ShutdownTemperature  *uint    // °C DCGM_FI_DEV_SHUTDOWN_TEMP Shutdown temperature for the device Modules
	PowerCurrentLimit    *float64 // W DCGM_FI_DEV_POWER_MGMT_LIMIT Current Power limit for the device
	PowerMinManLimit     *float64 // W DCGM_FI_DEV_POWER_MGMT_LIMIT_MIN Minimum power management limit for the device
	PowerMaxManLimit     *float64 // W DCGM_FI_DEV_POWER_MGMT_LIMIT_MAX Maximum power management limit for the device
	PowerDefaultManLimit *float64 // W DCGM_FI_DEV_POWER_MGMT_LIMIT_DEF Default power management limit for the device
	PowerEnforcedLimit   *float64 // W DCGM_FI_DEV_ENFORCED_POWER_LIMIT Effective power limit that the driver enforces after taking into account all limiters
	PowerViolationTime   *float64 // W DCGM_FI_DEV_POWER_VIOLATION Power Violation time in usec
	FBtotal              *uint    // MB DCGM_FI_DEV_FB_TOTAL Total Frame Buffer of the GPU in MB
	FBfree               *uint    // MB DCGM_FI_DEV_FB_FREE Free Frame Buffer in MB
	MemAppClock          *uint    // Mz DCGM_FI_DEV_APP_MEM_CLOCK Memory Application clocks
	SMAppClock           *uint    // Mz DCGM_FI_DEV_APP_SM_CLOCK SM Application clocks
	VideoEnClock         *uint    // Mz DCGM_FI_DEV_VIDEO_CLOCK Video encoder/decoder clock for the device
	RPSingleError        *uint    // DCGM_FI_DEV_RETIRED_SBE Number of retired pages because of single bit errors Note: monotonically increasing
	RPDoubleError        *uint    // DCGM_FI_DEV_RETIRED_DBE Number of retired pages because of double bit errors Note: monotonically increasing
	PackagePend          *uint    // DCGM_FI_DEV_RETIRED_PENDING Number of pages pending retirement
	SBErrors             *uint    // DCGM_FI_DEV_ECC_SBE_VOL_TOTAL Total single bit volatile ECC errors
	DBErrors             *uint    // DCGM_FI_DEV_ECC_DBE_VOL_TOTAL Total double bit volatile ECC errors
	MemSBAErrors         *uint    // DCGM_FI_DEV_ECC_SBE_AGG_DEV Device memory single bit aggregate (persistent) ECC errors Note: monotonically increasing
	MemDBAErrors         *uint    // DCGM_FI_DEV_ECC_DBE_AGG_DEV Device memory double bit aggregate (persistent) ECC errors Note: monotonically increasing
	DeviceMemSBErrors    *uint    // DCGM_FI_DEV_ECC_SBE_VOL_DEV Device memory single bit volatile ECC errors
	DeviceMemDBErrors    *uint    // DCGM_FI_DEV_ECC_DBE_VOL_DEV Device memory double bit volatile ECC errors
	RegisterSBErrors     *uint    // DCGM_FI_DEV_ECC_SBE_VOL_REG Register file single bit volatile ECC errors
	RegisterDBErrors     *uint    // DCGM_FI_DEV_ECC_DBE_VOL_REG Register file double bit volatile ECC errors
	DcgmSupported        *int     // supported 1, not running -1
}

// 监控项常量，定义用于构造deviceDields
const (
	gpuUtils int = iota
	memUtils
	encoder
	decoder
	singleBitError
	doubleBitError
	smClock
	memClock
	bAR1Used
	rx
	tx
	replays
	performance
	fanSpeed
	powerUsed
	fBUsed
	deviceTemperature
	memTemperature
	slowdownTemperature
	shutdownTemperature
	powerCurrentLimit
	powerMinManLimit
	powerMaxManLimit
	powerDefaultManLimit
	powerEnforcedLimit
	powerViolationTime
	fBtotal
	fBfree
	memAppClock
	sMAppClock
	videoEnClock
	rPSingleError
	rPDoubleError
	packagePend
	sBErrors
	dBErrors
	memSBAErrors
	memDBAErrors
	deviceMemSBErrors
	deviceMemDBErrors
	registerSBErrors
	registerDBErrors
	fieldsCount
)
