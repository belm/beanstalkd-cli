package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

// TestPutJob 测试插入任务
func TestPutJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, testTube)
	
	testCases := []struct {
		name     string
		data     string
		priority uint32
		delay    time.Duration
		ttr      time.Duration
	}{
		{"普通任务", "test job data", 1024, 0, 60 * time.Second},
		{"高优先级任务", "urgent job", 10, 0, 60 * time.Second},
		{"延迟任务", "delayed job", 1024, 5 * time.Second, 60 * time.Second},
		{"中文任务", "测试中文数据", 1024, 0, 60 * time.Second},
		{"JSON任务", `{"name":"test","value":123}`, 1024, 0, 60 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tube.Put([]byte(tc.data), tc.priority, tc.delay, tc.ttr)
			if err != nil {
				t.Errorf("插入任务失败: %v", err)
				return
			}
			t.Logf("✓ 成功插入任务 ID: %d, 数据: %s", id, tc.data)
			
			// 清理：删除任务
			conn.Delete(id)
		})
	}
}

// TestReserveJob 测试预留任务
func TestReserveJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 先插入一个任务
	tube := beanstalk.NewTube(conn, testTube)
	testData := "reserve test data"
	jobID, err := tube.Put([]byte(testData), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}
	t.Logf("插入测试任务 ID: %d", jobID)

	// 预留任务
	tubeSet := beanstalk.NewTubeSet(conn, testTube)
	id, body, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}

	if string(body) != testData {
		t.Errorf("任务数据不匹配: 期望 %s, 得到 %s", testData, string(body))
	}

	t.Logf("✓ 成功预留任务 ID: %d, 数据: %s", id, string(body))

	// 清理：删除任务
	conn.Delete(id)
}

// TestDeleteJob 测试删除任务
func TestDeleteJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入任务
	tube := beanstalk.NewTube(conn, testTube)
	jobID, err := tube.Put([]byte("delete test"), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	// 删除任务
	err = conn.Delete(jobID)
	if err != nil {
		t.Errorf("删除任务失败: %v", err)
	} else {
		t.Logf("✓ 成功删除任务 ID: %d", jobID)
	}

	// 尝试删除不存在的任务
	err = conn.Delete(99999999)
	if err == nil {
		t.Error("删除不存在的任务应该失败")
	} else {
		t.Logf("✓ 正确处理删除不存在的任务错误")
	}
}

// TestReleaseJob 测试释放任务
func TestReleaseJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入并预留任务
	tube := beanstalk.NewTube(conn, testTube)
	jobID, err := tube.Put([]byte("release test"), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	tubeSet := beanstalk.NewTubeSet(conn, testTube)
	id, _, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}

	// 释放任务
	err = conn.Release(id, 2048, 1*time.Second)
	if err != nil {
		t.Errorf("释放任务失败: %v", err)
	} else {
		t.Logf("✓ 成功释放任务 ID: %d", id)
	}

	// 清理
	time.Sleep(2 * time.Second) // 等待延迟过期
	tubeSet.Reserve(5 * time.Second)
	conn.Delete(jobID)
}

// TestBuryJob 测试埋葬任务
func TestBuryJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入并预留任务
	tube := beanstalk.NewTube(conn, testTube)
	jobID, err := tube.Put([]byte("bury test"), 1024, 0, 60*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	tubeSet := beanstalk.NewTubeSet(conn, testTube)
	id, _, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}

	// 埋葬任务
	err = conn.Bury(id, 1024)
	if err != nil {
		t.Errorf("埋葬任务失败: %v", err)
	} else {
		t.Logf("✓ 成功埋葬任务 ID: %d", id)
	}

	// 清理：踢出并删除
	tube.Kick(1)
	tubeSet.Reserve(5 * time.Second)
	conn.Delete(jobID)
}

// TestTouchJob 测试触摸任务
func TestTouchJob(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	// 插入并预留任务
	tube := beanstalk.NewTube(conn, testTube)
	jobID, err := tube.Put([]byte("touch test"), 1024, 0, 10*time.Second)
	if err != nil {
		t.Fatalf("插入任务失败: %v", err)
	}

	tubeSet := beanstalk.NewTubeSet(conn, testTube)
	id, _, err := tubeSet.Reserve(5 * time.Second)
	if err != nil {
		t.Fatalf("预留任务失败: %v", err)
	}

	// 触摸任务
	err = conn.Touch(id)
	if err != nil {
		t.Errorf("触摸任务失败: %v", err)
	} else {
		t.Logf("✓ 成功触摸任务 ID: %d", id)
	}

	// 清理
	conn.Delete(jobID)
}

// TestKickJobs 测试踢出任务
func TestKickJobs(t *testing.T) {
	conn, err := connectToServer(t)
	if err != nil {
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, testTube)
	
	// 插入并埋葬几个任务
	var jobIDs []uint64
	for i := 0; i < 3; i++ {
		jobID, err := tube.Put([]byte(fmt.Sprintf("kick test %d", i)), 1024, 0, 60*time.Second)
		if err != nil {
			t.Fatalf("插入任务失败: %v", err)
		}
		jobIDs = append(jobIDs, jobID)
		
		tubeSet := beanstalk.NewTubeSet(conn, testTube)
		id, _, _ := tubeSet.Reserve(5 * time.Second)
		conn.Bury(id, 1024)
	}

	// 踢出任务
	kicked, err := tube.Kick(3)
	if err != nil {
		t.Errorf("踢出任务失败: %v", err)
	} else {
		t.Logf("✓ 成功踢出 %d 个任务", kicked)
	}

	// 清理
	for _, id := range jobIDs {
		tubeSet := beanstalk.NewTubeSet(conn, testTube)
		tubeSet.Reserve(5 * time.Second)
		conn.Delete(id)
	}
}

// BenchmarkPutJob 插入任务性能基准测试
func BenchmarkPutJob(b *testing.B) {
	conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", testHost, testPort))
	if err != nil {
		b.Skip("无法连接到 beanstalkd 服务器")
		return
	}
	defer conn.Close()

	tube := beanstalk.NewTube(conn, testTube)
	data := []byte("benchmark test data")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id, err := tube.Put(data, 1024, 0, 60*time.Second)
		if err != nil {
			b.Errorf("插入任务失败: %v", err)
		}
		conn.Delete(id)
	}
}

// connectToServer 辅助函数：连接到测试服务器
func connectToServer(t *testing.T) (*beanstalk.Conn, error) {
	conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", testHost, testPort))
	if err != nil {
		t.Skipf("无法连接到 beanstalkd 服务器: %v", err)
		return nil, err
	}
	return conn, nil
}
