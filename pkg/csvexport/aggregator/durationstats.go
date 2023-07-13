package aggregator

type DurationStats struct {
	Duration         int64
	DurationMin      int64
	DurationMax      int64
	DurationMedian   int64
	DurationAvg      int64
	DurationRelative float32
	Count            int64
}
