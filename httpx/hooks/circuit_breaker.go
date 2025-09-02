package hooks

// CircuitBreaker is implements circuit breaking pattern for improving system resiliency
// CircuitBreaker is only used as client
type CircuitBreaker struct{}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{}
}
