package middleware

import (
	"github.com/AZRV17/zlib-backend/internal/server/metrics"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}

	ip = r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}

	return ip
}

func getUserID(r *http.Request) string {
	return r.Header.Get("X-User-ID")
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if isWebSocketRequest(r) {
				next.ServeHTTP(w, r)
				return
			}

			startTime := time.Now()
			clientIP := getIP(r)
			userID := getUserID(r)

			metrics.InitializeIPMetric(clientIP)
			metrics.RecordIPActivity(clientIP)
			metrics.RecordActiveIPs(clientIP, true)
			metrics.TotalRequests.WithLabelValues(r.URL.Path, r.Method).Inc()
			metrics.RecordTotalRequests(r.URL.Path, r.Method, userID)

			recorder := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(recorder, r)

			metrics.ResponseStatus.WithLabelValues(strconv.Itoa(recorder.statusCode)).Inc()
			metrics.RecordResponseStatus(strconv.Itoa(recorder.statusCode))

			duration := time.Since(startTime).Seconds()
			metrics.RequestDuration.WithLabelValues(
				r.URL.Path,
				r.Method,
				strconv.Itoa(recorder.statusCode),
			).Observe(duration)
			metrics.RecordRequestDuration(
				r.URL.Path,
				r.Method,
				strconv.Itoa(recorder.statusCode),
				duration,
			)
		},
	)
}

func isWebSocketRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Connection")) == "upgrade" &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket"
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
