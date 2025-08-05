package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Định nghĩa Job, bao gồm channel để worker gửi kết quả
type Job struct {
	ID      string
	Payload interface{}
	Result  chan interface{}
}

var jobQueue chan *Job

func main() {
	// Tạo queue với buffer tuỳ bạn (ở đây 100)
	jobQueue = make(chan *Job, 100000)

	// Khởi 1 worker duy nhất, chạy tuần tự
	go worker()
	//TODO Authentication JWT
	r := gin.Default()
	r.POST("/process", handler)
	r.Run(":8081")
}

func handler(c *gin.Context) {
	// 1. Parse payload từ request (ví dụ JSON)
	var payload map[string]interface{}
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Tạo Job, kèm channel Result
	job := &Job{
		ID:      time.Now().Format("20060102T150405.000"),
		Payload: payload,
		Result:  make(chan interface{}),
	}

	// 3. Enqueue (có thể block hoặc trả lỗi nếu queue đầy)
	select {
	case jobQueue <- job:
		// 4. Block chờ worker xử lý xong
		//    Bạn nên để timeout tránh chờ quá lâu
		select {
		case res := <-job.Result:
			c.JSON(http.StatusOK, gin.H{"job_id": job.ID, "result": res})
		case <-time.After(30 * time.Second):
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "processing timeout"})
		}
	default:
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "server busy, try later"})
	}
}

func worker() {
	for job := range jobQueue {
		// --- XỬ LÝ TUẦN TỰ ---
		// Giả sử đây là công việc tốn thời gian
		// time.Sleep(10 * time)

		// Kết quả có thể là string, struct, map, v.v.
		result := map[string]interface{}{
			"processed_at": time.Now().Format(time.RFC3339),
			"echo":         job.Payload,
		}

		// Gửi kết quả về handler
		job.Result <- result
	}
}
