package main

import (
	"context"

	"libs/common/ctxconst"
	"libs/common/logger"

	"template/internal/app/config"
)

func main() {
	// Configure logger
	cfg := &config.Config{
		Logger: config.LoggerConfig{
			Level:  -1, // Debug level
			Format: "json",
		},
	}

	if err := logger.Configure(cfg.Logger.Level, cfg.Logger.Format); err != nil {
		panic(err)
	}

	// Basic printf style logging
	logger.Printf("Server started on port %d", 8080)
	// Output: {"level":"info","timestamp":"2024-11-01T10:00:00Z","message":"Server started on port 8080"}

	// Basic context without any values
	ctx := context.Background()
	logger.Infof(ctx, "Simple info message")
	// Output: {"level":"info","timestamp":"2024-11-01T10:00:00Z","message":"Simple info message"}

	// Context with request ID, user ID and phone
	ctx = ctxconst.SetRequestID(ctx, "req-123456")
	ctx = ctxconst.SetUserID(ctx, "user-789")
	ctx = ctxconst.SetUserPhoneNumber(ctx, "+77771234567")

	// Logging with context - automatically includes request_id, user_id, and phone
	logger.Infof(ctx, "Processing payment of %d tenge", 5000)
	// Output: {
	//   "level": "info",
	//   "timestamp": "2024-11-01T10:00:00Z",
	//   "message": "Processing payment of 5000 tenge",
	//   "request_id": "req-123456",
	//   "user_id": "user-789",
	//   "user_phone_number": "+77771234567"
	// }

	// Different log levels
	logger.Debugf(ctx, "Attempting database connection to %s", "postgres://localhost:5432")
	logger.Warnf(ctx, "High CPU usage: %d%%", 85)
	logger.Errorf(ctx, "Failed to process transaction: %s", "insufficient funds")

	// Using WithFields for structured logging
	orderLogger := logger.WithFields(
		logger.Field{Key: "order_id", Value: "ord-123"},
		logger.Field{Key: "amount", Value: 15000},
		logger.Field{Key: "currency", Value: "KZT"},
	)

	orderLogger.Infof(ctx, "Order created successfully")
	// Output: {
	//   "level": "info",
	//   "timestamp": "2024-11-01T10:00:00Z",
	//   "message": "Order created successfully",
	//   "request_id": "req-123456",
	//   "user_id": "user-789",
	//   "user_phone_number": "+77771234567",
	//   "order_id": "ord-123",
	//   "amount": 15000,
	//   "currency": "KZT"
	// }

	// Error logging with additional fields
	paymentLogger := logger.WithFields(
		logger.Field{Key: "payment_id", Value: "pay-456"},
		logger.Field{Key: "status", Value: "failed"},
		logger.Field{Key: "error_code", Value: "ERR-001"},
	)

	paymentLogger.Errorf(ctx, "Payment authorization failed: %s", "card declined")
	// Output: {
	//   "level": "error",
	//   "timestamp": "2024-11-01T10:00:00Z",
	//   "message": "Payment authorization failed: card declined",
	//   "request_id": "req-123456",
	//   "user_id": "user-789",
	//   "user_phone_number": "+77771234567",
	//   "payment_id": "pay-456",
	//   "status": "failed",
	//   "error_code": "ERR-001"
	// }

	/*
	   logger.Fatalf(ctx, "Failed to initialize database: %s", "connection timeout")
	   logger.Panicf(ctx, "Critical security breach detected: %s", "unauthorized root access")
	*/
}
