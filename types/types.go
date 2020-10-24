package types

// DBInstance represents a particular RDS instance
type DBInstance struct {
	Identifier       string  // Instance Identifier
	AllocatedStorage float64 // allocated storage
	Iops             float64 // iops
}
