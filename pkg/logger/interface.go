package logger

// Logger defines the interface for application logging
// This allows for easy testing and switching between different logging implementations
type Logger interface {
	// Info logs an informational message
	Info(msg string, fields ...Field)

	// Debug logs a debug message (only shown in development)
	Debug(msg string, fields ...Field)

	// Warn logs a warning message
	Warn(msg string, fields ...Field)

	// Error logs an error message
	Error(msg string, fields ...Field)

	// Fatal logs a fatal message and exits the application
	Fatal(msg string, fields ...Field)
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value any
}

// F is a helper function to create a Field
func F(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
