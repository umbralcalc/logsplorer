package data

// DataStreamer
type DataStreamer interface {
	NextValue() []float64
}
