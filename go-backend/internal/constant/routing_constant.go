package constant

const (
	ModelTypeChat = "chat"

	RoutingStrategyAuto        = "auto"
	RoutingStrategyCostFirst   = "cost_first"
	RoutingStrategyLatencyFirst = "latency_first"
	RoutingStrategyFixed       = "fixed"

	ModelStatusActive = "active"

	HealthStatusHealthy   = "healthy"
	HealthStatusDegraded  = "degraded"
	HealthStatusUnknown   = "unknown"
	HealthStatusUnhealthy = "unhealthy"
)
