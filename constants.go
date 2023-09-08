package bees

import "fmt"

var (
	ERROR_CREATING_POOL_NO_WORKER        = fmt.Errorf("Attempting to create worker pool with less than 1 worker!")
	ERROR_CREATING_POOL_NEGATIVE_CHANNEL = fmt.Errorf("Attempting to create worker pool with negative channel size.")
)
