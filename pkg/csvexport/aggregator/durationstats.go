package aggregator

type DurationStats struct {
	Duration         uint32
	DurationMin      uint32
	DurationMax      uint32
	DurationMedian   uint32
	DurationAvg      uint32
	DurationRelative float32
	Count            uint32
}
