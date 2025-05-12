package metrics

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"time"
)

var (
	influxClient influxdb2.Client
	writeAPI     api.WriteAPIBlocking
)

func InitInfluxDB() {
	influxClient = influxdb2.NewClient(
		"http://localhost:8086",
		"HAEsVLUmN0vtdX6kztmZcxRFgfzPVVc6qvfZoC9pp_8z8HHoHuei5nmiRG5uWIf13C-utm95_k_SPXEgkqV-dg==",
	)

	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()
	//
	//health, err := influxClient.Health(ctx)
	//if err != nil {
	//	log.Printf("Ошибка подключения к InfluxDB: %v", err)
	//	return
	//}
	//
	//if health.Status != "pass" {
	//	log.Printf("InfluxDB не готов: %s", &health.Message)
	//} else {
	//	log.Println("InfluxDB успешно подключен")
	//}

	writeAPI = influxClient.WriteAPIBlocking("MPT", "metrics")
}

func CloseInfluxDB() {
	if influxClient != nil {
		influxClient.Close()
	}
}

func RecordRequestDuration(path, method, status string, duration float64) {
	point := write.NewPoint(
		"http_request_duration_seconds",
		map[string]string{
			"path":   path,
			"method": method,
			"status": status,
		},
		map[string]interface{}{
			"duration": duration,
		},
		time.Now(),
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Printf("Ошибка записи метрики длительности запроса: %v", err)
	}
}

func RecordActiveIPs(ip string, active bool) {
	value := 0.0
	if active {
		value = 1.0
	}

	point := write.NewPoint(
		"active_ips_total",
		map[string]string{
			"ip": ip,
		},
		map[string]interface{}{
			"value": value,
		},
		time.Now(),
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Printf("Ошибка записи метрики активных IP: %v", err)
	}
}

func RecordTotalRequests(path, method, userID string) {
	point := write.NewPoint(
		"total_requests",
		map[string]string{
			"path":    path,
			"method":  method,
			"user_id": userID,
		},
		map[string]interface{}{
			"count": 1.0,
		},
		time.Now(),
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Printf("Ошибка записи метрики общего количества запросов: %v", err)
	}
}

func RecordResponseStatus(status string) {
	point := write.NewPoint(
		"response_status",
		map[string]string{
			"status": status,
		},
		map[string]interface{}{
			"count": 1.0,
		},
		time.Now(),
	)

	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
		log.Printf("Ошибка записи метрики статуса ответа: %v", err)
	}
}
