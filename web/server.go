package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

var (
	beanstalkdHost string
	serverPort     string
)

type Response struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type StatsResponse struct {
	Stats map[string]string `json:"stats"`
	Error string            `json:"error,omitempty"`
}

type TubesResponse struct {
	Tubes []string `json:"tubes"`
	Error string   `json:"error,omitempty"`
}

type PutRequest struct {
	Tube     string `json:"tube"`
	Data     string `json:"data"`
	Priority uint32 `json:"priority"`
	Delay    int    `json:"delay"`
}

type PutResponse struct {
	JobID uint64 `json:"job_id"`
	Error string `json:"error,omitempty"`
}

type ReserveRequest struct {
	Tube    string `json:"tube"`
	Timeout int    `json:"timeout"`
}

type ReserveResponse struct {
	JobID uint64 `json:"job_id"`
	Data  string `json:"data"`
	Error string `json:"error,omitempty"`
}

type DeleteRequest struct {
	JobID uint64 `json:"job_id"`
}

type KickRequest struct {
	Tube  string `json:"tube"`
	Bound int    `json:"bound"`
}

type KickResponse struct {
	Kicked int    `json:"kicked"`
	Error  string `json:"error,omitempty"`
}

func main() {
	// 解析命令行参数
	flag.StringVar(&beanstalkdHost, "beanstalkd", "", "Beanstalkd 服务器地址 (默认: 127.0.0.1:11300)")
	flag.StringVar(&serverPort, "port", "", "Web 服务器端口 (默认: 8080)")
	flag.Parse()

	// 从环境变量读取配置（如果命令行未指定）
	if beanstalkdHost == "" {
		beanstalkdHost = os.Getenv("BEANSTALKD_HOST")
		if beanstalkdHost == "" {
			beanstalkdHost = "127.0.0.1:11300" // 默认值
		}
	}

	if serverPort == "" {
		serverPort = os.Getenv("WEB_PORT")
		if serverPort == "" {
			serverPort = "8080" // 默认值
		}
	}

	// 静态文件服务
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	// API 路由
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/tubes", handleTubes)
	http.HandleFunc("/api/tubes/", handleTubeStats)
	http.HandleFunc("/api/put", handlePut)
	http.HandleFunc("/api/reserve", handleReserve)
	http.HandleFunc("/api/delete", handleDelete)
	http.HandleFunc("/api/kick", handleKick)

	fmt.Printf("🚀 Beanstalkd Web UI 服务器启动成功!\n")
	fmt.Printf("📡 访问地址: http://localhost:%s\n", serverPort)
	fmt.Printf("🔗 Beanstalkd: %s\n\n", beanstalkdHost)
	fmt.Printf("💡 提示: 使用 -h 查看更多配置选项\n\n")

	log.Fatal(http.ListenAndServe(":"+serverPort, nil))
}

// 连接到 Beanstalkd
func connectToBeanstalkd() (*beanstalk.Conn, error) {
	conn, err := beanstalk.Dial("tcp", beanstalkdHost)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 Beanstalkd: %w", err)
	}
	return conn, nil
}

// 获取服务器统计
func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(StatsResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	stats, err := conn.Stats()
	if err != nil {
		json.NewEncoder(w).Encode(StatsResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(StatsResponse{Stats: stats})
}

// 获取所有 tubes
func handleTubes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(TubesResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	tubes, err := conn.ListTubes()
	if err != nil {
		json.NewEncoder(w).Encode(TubesResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(TubesResponse{Tubes: tubes})
}

// 获取 tube 统计
func handleTubeStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 从 URL 中提取 tube 名称
	tubeName := r.URL.Path[len("/api/tubes/"):]
	if len(tubeName) > 6 && tubeName[len(tubeName)-6:] == "/stats" {
		tubeName = tubeName[:len(tubeName)-6]
	}

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(StatsResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, tubeName)
	stats, err := tube.Stats()
	if err != nil {
		json.NewEncoder(w).Encode(StatsResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(StatsResponse{Stats: stats})
}

// 插入任务
func handlePut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(PutResponse{Error: "只支持 POST 方法"})
		return
	}

	var req PutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(PutResponse{Error: "无效的请求数据"})
		return
	}

	if req.Tube == "" {
		req.Tube = "default"
	}

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(PutResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, req.Tube)
	delay := time.Duration(req.Delay) * time.Second
	jobID, err := tube.Put([]byte(req.Data), req.Priority, delay, 60*time.Second)
	if err != nil {
		json.NewEncoder(w).Encode(PutResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(PutResponse{JobID: jobID})
}

// 预留任务
func handleReserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ReserveResponse{Error: "只支持 POST 方法"})
		return
	}

	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(ReserveResponse{Error: "无效的请求数据"})
		return
	}

	if req.Tube == "" {
		req.Tube = "default"
	}

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(ReserveResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	tubeSet := beanstalk.NewTubeSet(conn, req.Tube)
	timeout := time.Duration(req.Timeout) * time.Second
	jobID, body, err := tubeSet.Reserve(timeout)
	if err != nil {
		json.NewEncoder(w).Encode(ReserveResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(ReserveResponse{
		JobID: jobID,
		Data:  string(body),
	})
}

// 删除任务
func handleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(Response{Error: "只支持 POST 方法"})
		return
	}

	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{Error: "无效的请求数据"})
		return
	}

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}
	defer conn.Close()

	if err := conn.Delete(req.JobID); err != nil {
		json.NewEncoder(w).Encode(Response{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(Response{Data: "success"})
}

// 踢出任务
func handleKick(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(KickResponse{Error: "只支持 POST 方法"})
		return
	}

	var req KickRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(KickResponse{Error: "无效的请求数据"})
		return
	}

	if req.Tube == "" {
		req.Tube = "default"
	}

	conn, err := connectToBeanstalkd()
	if err != nil {
		json.NewEncoder(w).Encode(KickResponse{Error: err.Error()})
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, req.Tube)
	kicked, err := tube.Kick(req.Bound)
	if err != nil {
		json.NewEncoder(w).Encode(KickResponse{Error: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(KickResponse{Kicked: kicked})
}
