package tests

import (
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// TestProducerConsumerFlow 测试完整的生产者-消费者流程
func TestProducerConsumerFlow(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tubeName := "integration-test"
	tube := beanstalk.NewTube(conn, tubeName)

	// 生产者：插入任务
	testData := "integration test data"
	jobID, err := tube.Put([]byte(testData), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}
	t.Logf("生产者：插入任务 ID %d", jobID)

	// 消费者：预留任务
	tubeSet := beanstalk.NewTubeSet(conn, tubeName)
	id, body, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}
	t.Logf("消费者：预留任务 ID %d, 数据: %s", id, string(body))

	// 验证数据
	if string(body) != testData {
		t.Errorf("数据不匹配")
	}

	// 消费者：完成任务
	err = conn.Delete(id)
	if err != nil {
		t.Errorf("删除任务失败: %v", err)
	}
	t.Log("✓ 完整流程测试通过")
}

// TestRetryMechanism 测试任务重试机制
func TestRetryMechanism(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tubeName := "retry-test"
	tube := beanstalk.NewTube(conn, tubeName)
	tubeSet := beanstalk.NewTubeSet(conn, tubeName)

	// 插入任务
	_, err = tube.Put([]byte("retry test"), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	// 第一次尝试：失败并释放
	id, _, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("第一次预留失败: %v", err)
	}
	t.Log("第一次尝试处理任务...")

	err = conn.Release(id, 2048, 1*time.Second)
	if err != nil {
		t.Errorf("释放任务失败: %v", err)
	}
	t.Log("任务处理失败，已释放")

	// 等待延迟
	time.Sleep(2 * time.Second)

	// 第二次尝试：成功
	id, _, err = tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("第二次预留失败: %v", err)
	}
	t.Log("第二次尝试处理任务...")

	err = conn.Delete(id)
	if err != nil {
		t.Errorf("删除任务失败: %v", err)
	}
	t.Log("✓ 重试机制测试通过")
}

// TestPriorityQueue 测试优先级队列
func TestPriorityQueue(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tubeName := "priority-test"
	tube := beanstalk.NewTube(conn, tubeName)

	// 插入不同优先级的任务
	priorities := []struct {
		data     string
		priority uint32
	}{
		{"低优先级", 3000},
		{"高优先级", 10},
		{"中优先级", 1024},
	}

	for _, p := range priorities {
		_, err := tube.Put([]byte(p.data), p.priority, 0, 60*time.Second)
		if err != nil {
			t.Fatalf("插入任务失败: %v", err)
		}
	}

	// 按优先级顺序获取任务
	tubeSet := beanstalk.NewTubeSet(conn, tubeName)
	expectedOrder := []string{"高优先级", "中优先级", "低优先级"}

	for i, expected := range expectedOrder {
		id, body, err := tubeSet.Reserve(5 * time.Second)
		if err != nil {
			t.Errorf("预留任务 %d 失败: %v", i, err)
			continue
		}

		if string(body) != expected {
			t.Errorf("任务顺序错误: 期望 %s, 得到 %s", expected, string(body))
		}

		conn.Delete(id)
	}

	t.Log("✓ 优先级队列测试通过")
}
