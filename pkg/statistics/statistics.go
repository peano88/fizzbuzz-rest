package statistics

// NoStatsAvailable indicates that no previous request is available
type NoStatsAvailable struct{}

// Error is the error interface implementation
func (n NoStatsAvailable) Error() string {
	return "no previous requests available"
}
