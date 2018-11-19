# gpu-mon

## 简介

gpu-mon 是[open-falcon](http://open-falcon.com/)用于监控GPU状态的一个插件，该插件基于[(DCGM) NVIDIA Data Center GPU Manager](https://developer.nvidia.com/data-center-gpu-manager-dcgm)实现，可以对GPU的状态（显存使用、设备温度等）进行监控。

### 监控项

1. 详细的监控项说明可以参考[metric](https://github.com/open-falcon/gpu-mon/metric)文件，其中常用的一些监控项说明如下:

    ```plain
    
    GPUUtils             GPU 使用率 (%)
    MemUtils             GPU 显存使用率(%)
    FBUsed               GPU 的显存占用(MB)
    Performance          GPU 的性能状态(0-15, 其中0表示最高)
    DeviceTemperature    当前GPU设备温度(℃)
    PowerUsed            GPU的功率使用
    SingleBitError       全部累积的单精度ECC错误
    DoubleBitError       全部累积的双精度ECC错误
    ```

2. 在metric信息上报中，如果一些监控项采集到的数据异常，会上报 -1 值
3. 对于`GPUUtils`、`MemUtils`、`FBUsed`监控项的值，可以通过`type=sum`和`type=ave`tag来查看在整个设备上的全部使用情况和平均使用情况
4. 对于单个GPU卡的相关监控项查看，可以通过 `GpuID` tag来查看(如`GpuId = 0`)。

## 安装使用

### 1. 相关依赖

1. 安装DCGM并开启nv-hostengine进程，推荐使用DCGM 1.4.2版本
2. 目前能够支持全部 DCGM 1.4.2版本的GPU型号包括：
    - K80之后的Tesla系列GPU
    - Maxwell架构非Tesla系列GPU

    关于 Dcgm支持的GPU型号及DCGM安装可以参考[(DCGM) NVIDIA Data Center GPU Manager](https://developer.nvidia.com/data-center-gpu-manager-dcgm)
3. 目前插件已测试支持的GPU型号包括：v100、p4、p40。

### 2. 安装及使用

目前支持两种方式推送数据到open-falcon, 分别是设置crontab定时推送、作为falcon的插件推送

#### 2.1 通过crontab定时推送

1. 编译GPU监控文件

    ```bash
    go get -u github.com/open-falcon/gpu-mon
    cd $GOPATH/src/github.com/open-falcon/gpu-mon
    make
    ```

2. 编辑crontab 的配置文件，设置定时任务

    ``` bash
    # WORKPATH 为 gpu-mon 和 cfg.json 所在目录
    echo '* * * * * cd ${WORKPATH} && ./gpu-mon -c cfg.json ' >> /var/spool/cron/root
    ```

#### 2.2 作为falcon的插件推送

1. 编译GPU监控文件

    ```bash
    go get -u github.com/open-falcon/gpu-mon
    cd $GOPATH/src/github.com/open-falcon/gpu-mon
    make
    ```

2. 复制gpu_mon、cfg.json、60_gpuMonitor.sh 文件到 open-falcon 安装路径的plugin目录下

## 配置文件

配置文件参考cfg.example.json文件，相关配置项说明如下：

```json
{
    "falcon": {
        // Agent: 上报falcon客户端的地址
        "Agent": "http://127.0.0.1:1988/v1/agent" //todo 完整地址http:ip:port/v1/agent
    },
    "metric":{
        // ignoreMetrics: 不进行上报的GPU监控配置项
        "ignoreMetrics": [
            "RPSingleError",
            "RPDoubleError",
            "PackagePend",
            "Tx",
            "Rx"
        ],
        // endpoint值，默认为机器主机名
        "endpoint": ""
    },

    "log":{
        // logLevel: 日志级别，支持：Warn、Error和Debug，
        "level": "Warn",
        // logDir: 日志存储目录
        "dir": "./logs"
    }
}
```