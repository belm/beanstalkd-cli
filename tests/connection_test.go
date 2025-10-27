package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/beanstalkd/go-beanstalk"
)

const (
	testHost = "127.0.0.1"
	testPort = "11300"
	testTube = "test-tube"
)

// TestConnection 测试连接功能
func TestConnection(t *testing.T) {
	conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", testHost, testPort))
	if err != nil {
		t.Skipf("无法连接到 beanstalkd 服务器: %v", err)
		return
	}
	defer conn.Close()

	t.Log("✓ 成功连接到 beanstalkd 服务器")
}

// TestConnectionTimeout 测试连接超时
func TestConnectionTimeout(t *testing.T) {
	timeout := 2 * time.Second
	conn, err := beanstalk.DialTimeout("tcp", fmt.Sprintf("%s:%s", testHost, testPort), timeout)
	if err != nil {
		t.Skipf("无法连接到 beanstalkd 服务器: %v", err)
		return
	}
	defer conn.Close()

	t.Log("✓ 带超时的连接成功")
}

// TestConnectionFailure 测试连接失败场景
func TestConnectionFailure(t *testing.T) {
	// 尝试连接到不存在的端口
	_, err := beanstalk.DialTimeout("tcp", "127.0.0.1:99999", 1*time.Second)
	if err == nil {
		t.Error("预期连接失败，但成功了")
	} else {
		t.Logf("✓ 正确处理连接失败: %v", err)
	}
}

// TestMultipleConnections 测试多个连接
func TestMultipleConnections(t *testing.T) {
	connections := make([]*beanstalk.Conn, 5)
	
	for i := 0; i < 5; i++ {
		conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", testHost, testPort))
		if err != nil {
			t.Skipf("无法连接到 beanstalkd 服务器: %v", err)
			return
		}
		connections[i] = conn
		defer conn.Close()
	}
	
	t.Logf("✓ 成功创建 %d 个并发连接", len(connections))
}

// BenchmarkConnection 连接性能基准测试
func BenchmarkConnection(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%s", testHost, testPort))
		if err != nil {
			b.Skip("无法连接到 beanstalkd 服务器")
			return
		}
		conn.Close()
	}
}
