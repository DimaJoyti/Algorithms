# 📊 Monitoring & Observability

## 🎯 The Three Pillars of Observability

### 1. Metrics - What is happening?
```go
// Prometheus metrics implementation
type MetricsCollector struct {
    requestCounter    *prometheus.CounterVec
    requestDuration   *prometheus.HistogramVec
    activeConnections prometheus.Gauge
    errorRate         *prometheus.CounterVec
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        requestCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status_code"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "endpoint"},
        ),
        activeConnections: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "active_connections",
                Help: "Number of active connections",
            },
        ),
        errorRate: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "errors_total",
                Help: "Total number of errors",
            },
            []string{"service", "error_type"},
        ),
    }
}

func (mc *MetricsCollector) RecordRequest(method, endpoint string, statusCode int, duration time.Duration) {
    mc.requestCounter.WithLabelValues(method, endpoint, strconv.Itoa(statusCode)).Inc()
    mc.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
    
    if statusCode >= 400 {
        errorType := "client_error"
        if statusCode >= 500 {
            errorType = "server_error"
        }
        mc.errorRate.WithLabelValues("user-service", errorType).Inc()
    }
}

// Metrics middleware
func (mc *MetricsCollector) MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
        
        next.ServeHTTP(wrapped, r)
        
        duration := time.Since(start)
        mc.RecordRequest(r.Method, r.URL.Path, wrapped.statusCode, duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

### 2. Logs - What happened and why?
```go
// Structured logging with context
type Logger struct {
    logger *logrus.Logger
}

func NewLogger() *Logger {
    logger := logrus.New()
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    logger.SetLevel(logrus.InfoLevel)
    
    return &Logger{logger: logger}
}

func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
    entry := l.logger.WithFields(logrus.Fields{})
    
    // Add trace information
    if traceID := getTraceID(ctx); traceID != "" {
        entry = entry.WithField("trace_id", traceID)
    }
    
    if spanID := getSpanID(ctx); spanID != "" {
        entry = entry.WithField("span_id", spanID)
    }
    
    // Add user information
    if userID := getUserID(ctx); userID != "" {
        entry = entry.WithField("user_id", userID)
    }
    
    // Add request information
    if requestID := getRequestID(ctx); requestID != "" {
        entry = entry.WithField("request_id", requestID)
    }
    
    return entry
}

func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
    l.WithContext(ctx).WithFields(logrus.Fields{
        "method":      method,
        "path":        path,
        "status_code": statusCode,
        "duration_ms": duration.Milliseconds(),
        "event_type":  "http_request",
    }).Info("HTTP request processed")
}

func (l *Logger) LogError(ctx context.Context, err error, message string) {
    l.WithContext(ctx).WithFields(logrus.Fields{
        "error":      err.Error(),
        "event_type": "error",
    }).Error(message)
}

func (l *Logger) LogBusinessEvent(ctx context.Context, eventType string, data map[string]interface{}) {
    fields := logrus.Fields{
        "event_type": eventType,
    }
    
    for k, v := range data {
        fields[k] = v
    }
    
    l.WithContext(ctx).WithFields(fields).Info("Business event")
}

// Log aggregation patterns
type LogAggregator struct {
    buffer    []LogEntry
    batchSize int
    interval  time.Duration
    output    LogOutput
    mutex     sync.Mutex
}

type LogEntry struct {
    Timestamp time.Time              `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields"`
    Service   string                 `json:"service"`
}

func (la *LogAggregator) Start() {
    ticker := time.NewTicker(la.interval)
    go func() {
        for range ticker.C {
            la.flush()
        }
    }()
}

func (la *LogAggregator) AddLog(entry LogEntry) {
    la.mutex.Lock()
    defer la.mutex.Unlock()
    
    la.buffer = append(la.buffer, entry)
    
    if len(la.buffer) >= la.batchSize {
        go la.flush()
    }
}

func (la *LogAggregator) flush() {
    la.mutex.Lock()
    if len(la.buffer) == 0 {
        la.mutex.Unlock()
        return
    }
    
    batch := make([]LogEntry, len(la.buffer))
    copy(batch, la.buffer)
    la.buffer = la.buffer[:0]
    la.mutex.Unlock()
    
    // Send batch to output (Elasticsearch, CloudWatch, etc.)
    if err := la.output.SendBatch(batch); err != nil {
        log.Printf("Failed to send log batch: %v", err)
    }
}
```

### 3. Traces - How did we get here?
```go
// OpenTelemetry distributed tracing (perfect for your background!)
type TracingService struct {
    tracer trace.Tracer
}

func NewTracingService(serviceName string) (*TracingService, error) {
    // Initialize OTLP exporter
    exporter, err := otlptrace.New(
        context.Background(),
        otlptracegrpc.NewClient(
            otlptracegrpc.WithEndpoint("http://jaeger:14268/api/traces"),
            otlptracegrpc.WithInsecure(),
        ),
    )
    if err != nil {
        return nil, err
    }
    
    // Create trace provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        )),
    )
    
    otel.SetTracerProvider(tp)
    
    return &TracingService{
        tracer: tp.Tracer(serviceName),
    }, nil
}

func (ts *TracingService) StartSpan(ctx context.Context, operationName string) (context.Context, trace.Span) {
    return ts.tracer.Start(ctx, operationName)
}

func (ts *TracingService) TraceHTTPHandler(operationName string, handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, span := ts.StartSpan(r.Context(), operationName)
        defer span.End()
        
        // Add HTTP attributes
        span.SetAttributes(
            semconv.HTTPMethodKey.String(r.Method),
            semconv.HTTPURLKey.String(r.URL.String()),
            semconv.HTTPUserAgentKey.String(r.UserAgent()),
        )
        
        // Wrap response writer to capture status code
        wrapped := &tracingResponseWriter{ResponseWriter: w, span: span}
        
        // Execute handler with traced context
        handler(wrapped, r.WithContext(ctx))
        
        // Set final attributes
        span.SetAttributes(semconv.HTTPStatusCodeKey.Int(wrapped.statusCode))
        
        if wrapped.statusCode >= 400 {
            span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", wrapped.statusCode))
        }
    }
}

type tracingResponseWriter struct {
    http.ResponseWriter
    span       trace.Span
    statusCode int
}

func (trw *tracingResponseWriter) WriteHeader(code int) {
    trw.statusCode = code
    trw.ResponseWriter.WriteHeader(code)
}

// Database tracing
func (ts *TracingService) TraceDBQuery(ctx context.Context, query string, args ...interface{}) (context.Context, trace.Span) {
    ctx, span := ts.StartSpan(ctx, "db.query")
    
    span.SetAttributes(
        semconv.DBStatementKey.String(query),
        semconv.DBSystemKey.String("postgresql"),
    )
    
    return ctx, span
}

// External service call tracing
func (ts *TracingService) TraceHTTPClient(ctx context.Context, method, url string) (context.Context, trace.Span) {
    ctx, span := ts.StartSpan(ctx, fmt.Sprintf("http.client.%s", strings.ToLower(method)))
    
    span.SetAttributes(
        semconv.HTTPMethodKey.String(method),
        semconv.HTTPURLKey.String(url),
    )
    
    return ctx, span
}

// Business operation tracing
func (ts *TracingService) TraceBusinessOperation(ctx context.Context, operation string, attributes map[string]interface{}) (context.Context, trace.Span) {
    ctx, span := ts.StartSpan(ctx, operation)
    
    for key, value := range attributes {
        switch v := value.(type) {
        case string:
            span.SetAttributes(attribute.String(key, v))
        case int:
            span.SetAttributes(attribute.Int(key, v))
        case int64:
            span.SetAttributes(attribute.Int64(key, v))
        case float64:
            span.SetAttributes(attribute.Float64(key, v))
        case bool:
            span.SetAttributes(attribute.Bool(key, v))
        }
    }
    
    return ctx, span
}
```

## 🚨 Alerting & SLA Monitoring

### SLI/SLO Implementation
```go
type SLICollector struct {
    availability    *prometheus.CounterVec
    latency         *prometheus.HistogramVec
    errorBudget     *prometheus.GaugeVec
    sloTargets      map[string]SLOTarget
}

type SLOTarget struct {
    Service           string
    AvailabilityTarget float64 // 99.9% = 0.999
    LatencyTarget     time.Duration // P95 < 200ms
    ErrorBudget       float64 // 0.1% = 0.001
}

func NewSLICollector() *SLICollector {
    return &SLICollector{
        availability: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "sli_availability_total",
                Help: "Total availability measurements",
            },
            []string{"service", "status"}, // status: "success" or "failure"
        ),
        latency: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "sli_latency_seconds",
                Help:    "Request latency for SLI calculation",
                Buckets: []float64{0.01, 0.05, 0.1, 0.2, 0.5, 1.0, 2.0, 5.0},
            },
            []string{"service"},
        ),
        errorBudget: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "slo_error_budget_remaining",
                Help: "Remaining error budget (0-1)",
            },
            []string{"service"},
        ),
        sloTargets: map[string]SLOTarget{
            "user-service": {
                Service:           "user-service",
                AvailabilityTarget: 0.999, // 99.9%
                LatencyTarget:     200 * time.Millisecond,
                ErrorBudget:       0.001, // 0.1%
            },
        },
    }
}

func (sli *SLICollector) RecordRequest(service string, success bool, duration time.Duration) {
    status := "success"
    if !success {
        status = "failure"
    }
    
    sli.availability.WithLabelValues(service, status).Inc()
    sli.latency.WithLabelValues(service).Observe(duration.Seconds())
}

// Calculate SLO compliance
func (sli *SLICollector) CalculateErrorBudget(service string, windowDuration time.Duration) float64 {
    target, exists := sli.sloTargets[service]
    if !exists {
        return 0
    }
    
    // Query Prometheus for success/failure counts in the time window
    successCount := sli.getMetricValue(fmt.Sprintf(`sli_availability_total{service="%s",status="success"}[%s]`, service, windowDuration))
    failureCount := sli.getMetricValue(fmt.Sprintf(`sli_availability_total{service="%s",status="failure"}[%s]`, service, windowDuration))
    
    totalRequests := successCount + failureCount
    if totalRequests == 0 {
        return 1.0 // Full budget if no requests
    }
    
    actualAvailability := successCount / totalRequests
    allowedFailures := totalRequests * target.ErrorBudget
    actualFailures := failureCount
    
    remainingBudget := (allowedFailures - actualFailures) / allowedFailures
    if remainingBudget < 0 {
        remainingBudget = 0
    }
    
    sli.errorBudget.WithLabelValues(service).Set(remainingBudget)
    
    return remainingBudget
}
```

### Alert Manager Configuration
```go
type AlertManager struct {
    rules       []AlertRule
    channels    map[string]NotificationChannel
    silences    map[string]time.Time
    escalations map[string]EscalationPolicy
}

type AlertRule struct {
    Name        string
    Query       string
    Threshold   float64
    Duration    time.Duration
    Severity    AlertSeverity
    Labels      map[string]string
    Annotations map[string]string
}

type AlertSeverity string

const (
    SeverityCritical AlertSeverity = "critical"
    SeverityWarning  AlertSeverity = "warning"
    SeverityInfo     AlertSeverity = "info"
)

type EscalationPolicy struct {
    Steps []EscalationStep
}

type EscalationStep struct {
    Duration time.Duration
    Channels []string
}

func (am *AlertManager) EvaluateRules() {
    for _, rule := range am.rules {
        value := am.queryMetric(rule.Query)
        
        if value > rule.Threshold {
            alert := Alert{
                Name:        rule.Name,
                Value:       value,
                Threshold:   rule.Threshold,
                Severity:    rule.Severity,
                Labels:      rule.Labels,
                Annotations: rule.Annotations,
                StartsAt:    time.Now(),
            }
            
            am.fireAlert(alert)
        }
    }
}

func (am *AlertManager) fireAlert(alert Alert) {
    // Check if alert is silenced
    silenceKey := am.generateSilenceKey(alert)
    if silenceUntil, exists := am.silences[silenceKey]; exists && time.Now().Before(silenceUntil) {
        return
    }
    
    // Get escalation policy
    escalation, exists := am.escalations[alert.Name]
    if !exists {
        escalation = am.getDefaultEscalation(alert.Severity)
    }
    
    // Start escalation
    go am.escalateAlert(alert, escalation)
}

func (am *AlertManager) escalateAlert(alert Alert, policy EscalationPolicy) {
    for _, step := range policy.Steps {
        // Wait for step duration
        time.Sleep(step.Duration)
        
        // Check if alert is resolved
        if am.isAlertResolved(alert) {
            return
        }
        
        // Send notifications
        for _, channelName := range step.Channels {
            if channel, exists := am.channels[channelName]; exists {
                channel.Send(alert)
            }
        }
    }
}

// Notification channels
type SlackNotificationChannel struct {
    webhookURL string
    channel    string
}

func (snc *SlackNotificationChannel) Send(alert Alert) error {
    message := SlackMessage{
        Channel: snc.channel,
        Text:    fmt.Sprintf("🚨 Alert: %s", alert.Name),
        Attachments: []SlackAttachment{
            {
                Color: snc.getColorForSeverity(alert.Severity),
                Fields: []SlackField{
                    {Title: "Severity", Value: string(alert.Severity), Short: true},
                    {Title: "Value", Value: fmt.Sprintf("%.2f", alert.Value), Short: true},
                    {Title: "Threshold", Value: fmt.Sprintf("%.2f", alert.Threshold), Short: true},
                    {Title: "Time", Value: alert.StartsAt.Format(time.RFC3339), Short: true},
                },
            },
        },
    }
    
    return snc.sendToSlack(message)
}

type PagerDutyNotificationChannel struct {
    integrationKey string
}

func (pdnc *PagerDutyNotificationChannel) Send(alert Alert) error {
    event := PagerDutyEvent{
        RoutingKey:  pdnc.integrationKey,
        EventAction: "trigger",
        Payload: PagerDutyPayload{
            Summary:  alert.Name,
            Source:   "monitoring-system",
            Severity: string(alert.Severity),
            CustomDetails: map[string]interface{}{
                "value":     alert.Value,
                "threshold": alert.Threshold,
                "labels":    alert.Labels,
            },
        },
    }
    
    return pdnc.sendToPagerDuty(event)
}
```

## 📈 Health Checks & Circuit Breakers

### Health Check Implementation
```go
type HealthChecker struct {
    checks map[string]HealthCheck
    mutex  sync.RWMutex
}

type HealthCheck interface {
    Name() string
    Check(ctx context.Context) HealthStatus
}

type HealthStatus struct {
    Status  string                 `json:"status"` // "healthy", "unhealthy", "degraded"
    Message string                 `json:"message,omitempty"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// Database health check
type DatabaseHealthCheck struct {
    db      *sql.DB
    timeout time.Duration
}

func (dhc *DatabaseHealthCheck) Name() string {
    return "database"
}

func (dhc *DatabaseHealthCheck) Check(ctx context.Context) HealthStatus {
    ctx, cancel := context.WithTimeout(ctx, dhc.timeout)
    defer cancel()
    
    if err := dhc.db.PingContext(ctx); err != nil {
        return HealthStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("Database ping failed: %v", err),
        }
    }
    
    // Check if we can perform a simple query
    var result int
    err := dhc.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
    if err != nil {
        return HealthStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("Database query failed: %v", err),
        }
    }
    
    return HealthStatus{
        Status: "healthy",
        Details: map[string]interface{}{
            "max_open_connections": dhc.db.Stats().MaxOpenConnections,
            "open_connections":     dhc.db.Stats().OpenConnections,
            "in_use":              dhc.db.Stats().InUse,
            "idle":                dhc.db.Stats().Idle,
        },
    }
}

// External service health check
type ExternalServiceHealthCheck struct {
    name       string
    url        string
    httpClient *http.Client
    timeout    time.Duration
}

func (eshc *ExternalServiceHealthCheck) Name() string {
    return eshc.name
}

func (eshc *ExternalServiceHealthCheck) Check(ctx context.Context) HealthStatus {
    ctx, cancel := context.WithTimeout(ctx, eshc.timeout)
    defer cancel()
    
    req, err := http.NewRequestWithContext(ctx, "GET", eshc.url, nil)
    if err != nil {
        return HealthStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("Failed to create request: %v", err),
        }
    }
    
    start := time.Now()
    resp, err := eshc.httpClient.Do(req)
    duration := time.Since(start)
    
    if err != nil {
        return HealthStatus{
            Status:  "unhealthy",
            Message: fmt.Sprintf("Request failed: %v", err),
            Details: map[string]interface{}{
                "response_time_ms": duration.Milliseconds(),
            },
        }
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return HealthStatus{
            Status: "healthy",
            Details: map[string]interface{}{
                "status_code":      resp.StatusCode,
                "response_time_ms": duration.Milliseconds(),
            },
        }
    }
    
    return HealthStatus{
        Status:  "unhealthy",
        Message: fmt.Sprintf("Unhealthy status code: %d", resp.StatusCode),
        Details: map[string]interface{}{
            "status_code":      resp.StatusCode,
            "response_time_ms": duration.Milliseconds(),
        },
    }
}

// Aggregate health endpoint
func (hc *HealthChecker) HealthHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    hc.mutex.RLock()
    checks := make(map[string]HealthCheck, len(hc.checks))
    for name, check := range hc.checks {
        checks[name] = check
    }
    hc.mutex.RUnlock()
    
    results := make(map[string]HealthStatus)
    overallStatus := "healthy"
    
    // Run all health checks concurrently
    var wg sync.WaitGroup
    var mutex sync.Mutex
    
    for name, check := range checks {
        wg.Add(1)
        go func(name string, check HealthCheck) {
            defer wg.Done()
            
            status := check.Check(ctx)
            
            mutex.Lock()
            results[name] = status
            if status.Status != "healthy" {
                overallStatus = "unhealthy"
            }
            mutex.Unlock()
        }(name, check)
    }
    
    wg.Wait()
    
    response := map[string]interface{}{
        "status": overallStatus,
        "checks": results,
        "timestamp": time.Now().UTC(),
    }
    
    statusCode := http.StatusOK
    if overallStatus != "healthy" {
        statusCode = http.StatusServiceUnavailable
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}
```

## 🔗 Next Steps
- Practice with real-world system design problems
- Learn about advanced monitoring patterns (chaos engineering, canary deployments)
- Study compliance and regulatory monitoring requirements
