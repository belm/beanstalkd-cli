package tests

import (
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// TestPeekJob 测试查看任务
func TestPeekJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入任务
	tube := beanstalk.NewTube(conn, testTube)
	testData := "peek test data"
	jobID, err := tube.Put([]byte(testData), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}
	defer conn.Delete(jobID)

	// 查看任务
	body, err := conn.Peek(jobID)
	if err != nil {
		t.Errorf("查看任务失败: %v", err)
		return
	}

	if string(body) != testData {
		t.Errorf("任务数据不匹配: 期望 %s, 得到 %s", testData, string(body))
	}

	t.Logf("✓ 成功查看任务 ID: %d, 数据: %s", jobID, string(body))
}

// TestPeekReady 测试查看就绪任务
func TestPeekReady(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入任务
	tube := beanstalk.NewTube(conn, testTube)
	testData := "peek ready test"
	jobID, err := tube.Put([]byte(testData), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}
	defer conn.Delete(jobID)

	// 查看就绪任务
	id, body, err := tube.PeekReady()
	if err != nil {
		t.Errorf("查看就绪任务失败: %v", err)
		return
	}

	t.Logf("✓ 成功查看就绪任务 ID: %d, 数据: %s", id, string(body))
}

// TestPeekDelayed 测试查看延迟任务
func TestPeekDelayed(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入延迟任务
	tube := beanstalk.NewTube(conn, testTube)
	testData := "peek delayed test"
	jobID, err := tube.Put([]byte(testData), 1024, 10*time.Second, 60*time.Second)
	if err != nil {
		t.Fatalf("插入延迟任务失败: %v", err)
	}
	defer func() {
		// 等待延迟过期后删除
		time.Sleep(11 * time.Second)
		tubeSet := beanstalk.NewTubeSet(conn, testTube)
		tubeSet.Reserve(5 * time.Second)
		conn.Delete(jobID)
	}()

	// 查看延迟任务
	id, body, err := tube.PeekDelayed()
	if err != nil {
		t.Errorf("查看延迟任务失败: %v", err)
		return
	}

	if string(body) != testData {
		t.Errorf("任务数据不匹配: 期望 %s, 得到 %s", testData, string(body))
	}

	t.Logf("✓ 成功查看延迟任务 ID: %d, 数据: %s", id, string(body))
}

// TestPeekBuried 测试查看被埋葬的任务
func TestPeekBuried(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入并埋葬任务
	tube := beanstalk.NewTube(conn, testTube)
	testData := "peek buried test"
	jobID, err := tube.Put([]byte(testData), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	// 预留并埋葬
	tubeSet := beanstalk.NewTubeSet(conn, testTube)
	id, _, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}
	
	err = conn.Bury(id, 1024)
	if err != nil {
		t.Fatalf("埋葬任务失败: %v", err)
	}

	// 查看被埋葬的任务
	id, body, err := tube.PeekBuried()
	if err != nil {
		t.Errorf("查看被埋葬任务失败: %v", err)
	} else {
		t.Logf("✓ 成功查看被埋葬任务 ID: %d, 数据: %s", id, string(body))
	}

	// 清理
	tube.Kick(1)
	tubeSet.Reserve(5 * time.Second)
	conn.Delete(jobID)
}

// TestPeekNonExistentJob 测试查看不存在的任务
func TestPeekNonExistentJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 尝试查看不存在的任务
	_, err = conn.Peek(99999999)
	if err == nil {
		t.Error("查看不存在的任务应该失败")
	} else {
		t.Logf("✓ 正确处理查看不存在任务的错误: %v", err)
	}
}
