package tests

import (
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// TestListTubes 测试列出所有 tubes
func TestListTubes(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tubes, err := conn.ListTubes()
	if err != nil {
		t.Errorf("列出 tubes 失败: %v", err)
		return
	}

	if len(tubes) == 0 {
		t.Error("tubes 列表为空")
		return
	}

	t.Logf("✓ 成功列出 %d 个 tubes: %v", len(tubes), tubes)
}

// TestMultipleTubes 测试多个 tubes 操作
func TestMultipleTubes(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	testTubes := []string{"test-tube-1", "test-tube-2", "test-tube-3"}
	
	// 在每个 tube 中插入任务
	for _, tubeName := range testTubes {
		tube := beanstalk.NewTube(conn, tubeName)
		jobID, err := tube.Put([]byte("test data for "+tubeName), 1024, 0, 60*time.Second)
		if err != nil {
			t.Errorf("在 tube %s 中插入任务失败: %v", tubeName, err)
			continue
		}
		t.Logf("✓ 在 tube %s 中插入任务 ID: %d", tubeName, jobID)
		
		// 清理
		defer func(id uint64) {
			conn.Delete(id)
		}(jobID)
	}

	// 验证 tubes 存在
	tubes, err := conn.ListTubes()
	if err != nil {
		t.Errorf("列出 tubes 失败: %v", err)
		return
	}

	for _, testTube := range testTubes {
		found := false
		for _, tube := range tubes {
			if tube == testTube {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("未找到 tube: %s", testTube)
		}
	}

	t.Logf("✓ 所有测试 tubes 都存在")
}

// TestTubeStats 测试获取 tube 统计信息
func TestTubeStats(t *testing.T) {
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

	// 获取统计信息
	stats, err := tube.Stats()
	if err != nil {
		t.Errorf("获取 tube 统计失败: %v", err)
		return
	}

	if len(stats) == 0 {
		t.Error("统计信息为空")
		return
	}

	t.Logf("✓ 成功获取 tube '%s' 的统计信息:", testTube)
	for key, value := range stats {
		t.Logf("  %s: %s", key, value)
	}
}

// TestTubeIsolation 测试 tube 隔离性
func TestTubeIsolation(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tube1 := beanstalk.NewTube(conn, "isolation-test-1")
	tube2 := beanstalk.NewTube(conn, "isolation-test-2")

	// 在两个不同的 tube 中插入任务
	data1 := "data for tube 1"
	data2 := "data for tube 2"
	
	jobID1, err := tube1.Put([]byte(data1), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("在 tube1 中插入任务失败: %v", err)
	}
	defer conn.Delete(jobID1)

	jobID2, err := tube2.Put([]byte(data2), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("在 tube2 中插入任务失败: %v", err)
	}
	defer conn.Delete(jobID2)

	// 从 tube1 预留任务
	tubeSet1 := beanstalk.NewTubeSet(conn, "isolation-test-1")
	id, body, err := tubeSet1.Reserve(5 * time.Second)
	if err != nil {
		t.Errorf("从 tube1 预留任务失败: %v", err)
		return
	}

	// 验证获取的是 tube1 的任务
	if string(body) != data1 {
		t.Errorf("从 tube1 获取到错误的数据: 期望 %s, 得到 %s", data1, string(body))
	}

	t.Logf("✓ Tube 隔离性测试通过")
	conn.Delete(id)
}

// TestConcurrentTubeAccess 测试并发访问不同 tubes
func TestConcurrentTubeAccess(t *testing.T) {
	testTubes := []string{"concurrent-1", "concurrent-2", "concurrent-3"}
	
	results := make(chan error, len(testTubes))
	
	for _, tubeName := range testTubes {
		go func(name string) {
			conn, err := connectToServer(t)
			if err != nil {
				results <- err
				return
			}
			defer conn.Close()

			tube := beanstalk.NewTube(conn, name)
			jobID, err := tube.Put([]byte("concurrent test"), 1024, 0, 60*time.Second)
			if err != nil {
				results <- err
				return
			}
			
			conn.Delete(jobID)
			results <- nil
		}(tubeName)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < len(testTubes); i++ {
		if err := <-results; err != nil {
			t.Errorf("并发 tube 访问失败: %v", err)
		}
	}

	t.Logf("✓ 并发访问 %d 个 tubes 成功", len(testTubes))
}
