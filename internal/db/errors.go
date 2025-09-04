package db

import "errors"

var (
	ErrLoggerNotInitialized        = errors.New("logger is not initialized")
	ErrComputeModuleNotInitialized = errors.New("compute module is not initialized")
	ErrStorageModuleNotInitialized = errors.New("storage module is not initialized")
)
