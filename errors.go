package main

import (
	"fmt"
)

// ErrorType categorizes the kind of error that occurred
type ErrorType string

// ErrorSeverity indicates how severe the error is
type ErrorSeverity string

// Define common error types
const (
	ErrorTypeConfiguration ErrorType = "configuration" // Error in plugin configuration
	ErrorTypeProcessing    ErrorType = "processing"    // Error during data processing
	ErrorTypeParsing       ErrorType = "parsing"       // Error parsing input data
	ErrorTypeIO            ErrorType = "io"            // Error with input/output operations
	ErrorTypeSystem        ErrorType = "system"        // System-level error
)

// Define error severity levels
const (
	ErrorSeverityError   ErrorSeverity = "error"   // Serious error that should abort processing
	ErrorSeverityWarning ErrorSeverity = "warning" // Warning that doesn't necessarily stop processing
	ErrorSeverityInfo    ErrorSeverity = "info"    // Informational message about a potential issue
)

// ProcessorError is a rich error type with context for debugging
type ProcessorError struct {
	Err             error                  // Original error
	Type            ErrorType              // Category of error
	Severity        ErrorSeverity          // Severity level
	TransactionHash string                 // Transaction context
	LedgerSequence  uint32                 // Ledger context
	ContractID      string                 // Contract context
	Context         map[string]interface{} // Additional metadata
}

// Error implements the error interface
func (e *ProcessorError) Error() string {
	return fmt.Sprintf("[%s:%s] %v (tx: %s, ledger: %d, contract: %s, context: %v)",
		e.Type, e.Severity, e.Err, e.TransactionHash, e.LedgerSequence, e.ContractID, e.Context)
}

// Unwrap implements the errors unwrapping interface
func (e *ProcessorError) Unwrap() error {
	return e.Err
}

// NewProcessorError creates a new processor error with the given type and severity
func NewProcessorError(err error, errType ErrorType, severity ErrorSeverity) *ProcessorError {
	return &ProcessorError{
		Err:      err,
		Type:     errType,
		Severity: severity,
		Context:  make(map[string]interface{}),
	}
}

// WithTransaction adds transaction context to the error
func (e *ProcessorError) WithTransaction(hash string) *ProcessorError {
	e.TransactionHash = hash
	return e
}

// WithLedger adds ledger context to the error
func (e *ProcessorError) WithLedger(sequence uint32) *ProcessorError {
	e.LedgerSequence = sequence
	return e
}

// WithContract adds contract context to the error
func (e *ProcessorError) WithContract(id string) *ProcessorError {
	e.ContractID = id
	return e
}

// WithContext adds additional context to the error
func (e *ProcessorError) WithContext(key string, value interface{}) *ProcessorError {
	e.Context[key] = value
	return e
}
