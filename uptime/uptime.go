// Package uptime provides utilities for tracking how long the current program
// has been running. It records the start time when the package is first
// initialized and offers functions to retrieve the start time and the elapsed
// uptime in various units.
package uptime

import "time"

// startTime indicates when the program has started.
var startTime = time.Now()

// GetStartTime returns the time when the program started.
func GetStartTime() time.Time {
	return startTime
}

// Duration returns the program uptime as [time.Duration].
func Duration() time.Duration {
	return time.Since(startTime)
}

// Seconds returns the program uptime in seconds.
func Seconds() float64 {
	return time.Since(startTime).Seconds()
}

// Minutes returns the program uptime in minutes.
func Minutes() float64 {
	return time.Since(startTime).Minutes()
}

// Hours returns the program uptime in hours.
func Hours() float64 {
	return time.Since(startTime).Hours()
}
