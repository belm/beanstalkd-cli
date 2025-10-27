package tests

import (
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// TestServerStats 测试服务器统计信息
func TestServerStats(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	stats, err := conn.Stats()
	if err != nil {
		t.Errorf("获取服务器统计失败: %v", err)
		return
	}

	if len(stats) == 0 {
		t.Error("统计信息为空")
		return
	}

	// 检查关键指标
	requiredKeys := []string{"current-jobs-ready", "current-jobs-reserved", "current-tubes", "version"}
	for _, key := range requiredKeys {
		if _, exists := stats[key]; !exists {
			t.Errorf("统计信息缺少关键指标: %s", key)
		}
	}

	t.Logf("✓ 成功获取服务器统计信息 (%d 个指标)", len(stats))
}

// TestJobStats 测试任务统计信息
func TestJobStats(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入任务
	tube := beanstalk.NewTube(conn, testTube)
	jobID, err := tube.Put([]byte("stats test"), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}
	defer conn.Delete(jobID)

	// 获取任务统计
	stats, err := conn.StatsJob(jobID)
	if err != nil {
		t.Errorf("获取任务统计失败: %v", err)
		return
	}

	t.Logf("✓ 成功获取任务 %d 的统计信息 (%d 个指标)", jobID, len(stats))
}
