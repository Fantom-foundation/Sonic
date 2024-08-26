package dbutil

type MeasurableStore interface {
	IoStats() (string, error)
	UsedDiskSpace() (string, error)
}
