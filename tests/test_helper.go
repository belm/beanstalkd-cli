package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// Helper functions for tests

// GetTestConfig 获取测试配置
func GetTestConfig() (host, port string) {
	host = os.Getenv("BEANSTALKD_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	port = os.Getenv("BEANSTALKD_PORT")
	if port == "" {
		port = "11300"
	}

	return
}

// CleanupTube 清理指定 tube 中的所有任务
func CleanupTube(t *testing.T, conn *beanstalk.Conn, tubeName string) {
	tube := beanstalk.NewTube(conn, tubeName)
	tubeSet := beanstalk.NewTubeSet(conn, tubeName)

	// 尝试清理就绪任务
	for i := 0; i < 100; i++ {
		id, _, err := tubeSet.Reserve(1 * time.Second)
		if err != nil {
			break
		}
		conn.Delete(id)
	}

	// 踢出被埋葬的任务并清理
	tube.Kick(100)
	for i := 0; i < 100; i++ {
		id, _, err := tubeSet.Reserve(1 * time.Second)
		if err != nil {
			break
		}
		conn.Delete(id)
	}

	t.Logf("已清理 tube: %s", tubeName)
}

// AssertNoError 断言没有错误
func AssertNoError(t *testing.T, err error, message string) {
	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError 断言有错误
func AssertError(t *testing.T, err error, message string) {
	if err == nil {
		t.Fatalf("%s: 期望有错误但没有", message)
	}
}

// AssertEqual 断言相等
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Errorf("%s: 期望 %v, 得到 %v", message, expected, actual)
	}
}

// CreateTestJobs 创建测试任务
func CreateTestJobs(t *testing.T, conn *beanstalk.Conn, tubeName string, count int) []uint64 {
	tube := beanstalk.NewTube(conn, tubeName)
	jobIDs := make([]uint64, 0, count)

	for i := 0; i < count; i++ {
		data := fmt.Sprintf("test job %d", i)
		jobID, err := tube.Put([]byte(data), 1024, 0, 60*time.Second)
		if err != nil {
			t.Fatalf("创建测试任务失败: %v", err)
		}
		jobIDs = append(jobIDs, jobID)
	}

	t.Logf("已创建 %d 个测试任务", count)
	return jobIDs
}

// DeleteTestJobs 删除测试任务
func DeleteTestJobs(t *testing.T, conn *beanstalk.Conn, jobIDs []uint64) {
	for _, id := range jobIDs {
		conn.Delete(id)
	}
	t.Logf("已删除 %d 个测试任务", len(jobIDs))
}
