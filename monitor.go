package main

import (
        "context"
        "log"
        "net/http"
        "os"
        "strconv"
        "time"

        "github.com/nervosnetwork/ckb-sdk-go/rpc"
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
        // 默认的 CKB 节点 RPC 地址
        DefaultCkbRpcUrl = "https://mainnet.ckb.dev"
        // 默认的更新时间间隔，单位为秒
        DefaultUpdateInterval = 10
        // 指标名称
        MetricName = "ckb_block_info"
)

var (
        blockNumber = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: MetricName + "_number",
                Help: "The current block number",
        })

        blockTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
                Name: MetricName + "_timestamp",
                Help: "The current block timestamp",
        })
)

func main() {
        // 读取 CKB 节点 RPC 地址和更新时间间隔的环境变量
        ckbRpcUrl := os.Getenv("CKB_RPC_URL")
        if ckbRpcUrl == "" {
                ckbRpcUrl = DefaultCkbRpcUrl
        }

        updateIntervalStr := os.Getenv("UPDATE_INTERVAL")
        updateInterval, err := strconv.Atoi(updateIntervalStr)
        if err != nil || updateInterval <= 0 {
                updateInterval = DefaultUpdateInterval
        }

        // 创建 CKB 节点的 RPC 客户端
        client, err := rpc.Dial(ckbRpcUrl)
        if err != nil {
                log.Fatalf("Failed to create CKB node RPC client: %v", err)
        }

        // 注册指标
        prometheus.MustRegister(blockNumber)
        prometheus.MustRegister(blockTimestamp)

        // 定时获取区块信息并更新指标
        go func() {
                for {
                        header, err := client.GetTipHeader(context.Background())
                        if err != nil {
                                log.Printf("Failed to get tip header: %v", err)
                                continue
                        }

                        blockNumber.Set(float64(header.Number))
                        blockTimestamp.Set(float64(header.Timestamp))

                        log.Printf("Updated block info: number=%d timestamp=%d", header.Number, header.Timestamp)

                        time.Sleep(time.Duration(updateInterval) * time.Second)
                }
        }()

        // 启动 HTTP 服务，暴露指标
        log.Printf("Starting HTTP server on port 8080")
        http.Handle("/metrics", promhttp.Handler())
        log.Fatal(http.ListenAndServe(":8080", nil))
}
