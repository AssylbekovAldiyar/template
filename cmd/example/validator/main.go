package main

import (
	"context"

	"libs/common/logger"
	"libs/common/validator"

	"template/internal/app/config"
)

type Book struct {
	ID   int64  `json:"id" validate:"required,book_id"`
	Name string `json:"name" validate:"required,book_name,name_format,name_no_special"`
}

func main() {

	ctx := context.Background()

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

	// Initialize valid
	valid, err := validator.New()
	if err != nil {
		logger.Fatalf(ctx, "Failed to initialize valid: %v", err)
	}

	// Test cases
	testCases := []struct {
		name string
		book Book
	}{
		{
			name: "Valid book",
			book: Book{
				ID:   1,
				Name: "The Great Gatsby",
			},
		},
		{
			name: "Invalid ID (zero)",
			book: Book{
				ID:   0,
				Name: "Valid Name",
			},
		},
		{
			name: "Invalid ID (negative)",
			book: Book{
				ID:   -1,
				Name: "Valid Name",
			},
		},
		{
			name: "Empty name",
			book: Book{
				ID:   1,
				Name: "",
			},
		},
		{
			name: "Short name",
			book: Book{
				ID:   1,
				Name: "A",
			},
		},
		{
			name: "Long name (101 chars)",
			book: Book{
				ID:   1,
				Name: "This is a very long book name that exceeds the maximum allowed length for a book name in our valid rules...",
			},
		},
		{
			name: "Name with special characters",
			book: Book{
				ID:   1,
				Name: "Book@#$%^&*",
			},
		},
		{
			name: "Name with allowed punctuation",
			book: Book{
				ID:   1,
				Name: "Book Name: Part 1!",
			},
		},
		{
			name: "Only spaces",
			book: Book{
				ID:   1,
				Name: "    ",
			},
		},
		{
			name: "Mixed valid characters",
			book: Book{
				ID:   1,
				Name: "Book 123 - Chapter 1",
			},
		},
	}

	// Validate each test case
	for _, tc := range testCases {
		// Create logger with test case fields
		testLogger := logger.WithFields(
			logger.Field{Key: "test_case", Value: tc.name},
			logger.Field{Key: "book", Value: tc.book},
		)

		testLogger.Infof(ctx, "Processing test case: %s", tc.name)

		if errors, err := valid.Validate(tc.book); err != nil {
			// Convert valid errors to JSON for structured logging
			testLogger.Errorf(ctx, "Validation failed for test case: %s\nerrors: %s", tc.name, logger.Field{Key: "errors", Value: errors})
		} else {
			testLogger.Infof(ctx, "Validation passed for test case: %s", tc.name)
		}

	}

}
