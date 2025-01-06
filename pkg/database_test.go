package pkg

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertComponentData(t *testing.T) {
	tests := []struct {
		name            string
		jobID           string
		commit          string
		createdAt       string
		totalSuccess    int
		totalFailures   int
		totalErrors     int
		totalSkips      int
		mockQueryResult func(mock sqlmock.Sqlmock)
		expectedError   bool
	}{
		{
			name:          "Successful insertion",
			jobID:         "job123",
			commit:        "abc123",
			createdAt:     "2024-11-26T12:00:00Z",
			totalSuccess:  10,
			totalFailures: 2,
			totalErrors:   1,
			totalSkips:    5,
			mockQueryResult: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT OR REPLACE INTO dci_components").
					WithArgs("job123", "abc123", "2024-11-26T12:00:00Z", 10, 2, 1, 5).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:          "Database error",
			jobID:         "job456",
			commit:        "def456",
			createdAt:     "2024-11-26T13:00:00Z",
			totalSuccess:  5,
			totalFailures: 1,
			totalErrors:   0,
			totalSkips:    2,
			mockQueryResult: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT OR REPLACE INTO dci_components").
					WithArgs("job456", "def456", "2024-11-26T13:00:00Z", 5, 1, 0, 2).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
		{
			name:          "Empty commit hash",
			jobID:         "job789",
			commit:        "",
			createdAt:     "2024-11-26T14:00:00Z",
			totalSuccess:  3,
			totalFailures: 0,
			totalErrors:   0,
			totalSkips:    1,
			mockQueryResult: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT OR REPLACE INTO dci_components").
					WithArgs("job789", "", "2024-11-26T14:00:00Z", 3, 0, 0, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a mock database connection
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Set up mock behavior
			tc.mockQueryResult(mock)

			// Call the function
			err = insertComponentData(db, tc.jobID, tc.commit, tc.createdAt, tc.totalSuccess, tc.totalFailures, tc.totalErrors, tc.totalSkips)

			// Validate the results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestInsertQuayData(t *testing.T) {
	// Define test cases
	tests := []struct {
		name          string
		datetime      string
		count         int
		kind          string
		mockSetup     func(mock sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:     "Successful Insert",
			datetime: "2024-11-26T12:00:00Z",
			count:    100,
			kind:     "image_pulls",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT OR REPLACE INTO aggregated_logs \(datetime, count, kind\)`).
					WithArgs("2024-11-26T12:00:00Z", 100, "image_pulls").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:     "Insert with Missing Kind",
			datetime: "2024-11-26T12:00:00Z",
			count:    50,
			kind:     "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT OR REPLACE INTO aggregated_logs \(datetime, count, kind\)`).
					WithArgs("2024-11-26T12:00:00Z", 50, "").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:     "Database Error",
			datetime: "2024-11-26T12:00:00Z",
			count:    200,
			kind:     "image_pulls",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT OR REPLACE INTO aggregated_logs \(datetime, count, kind\)`).
					WithArgs("2024-11-26T12:00:00Z", 200, "image_pulls").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	// Loop through each test case
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Set up the mock database
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Apply the test-specific mock setup
			tc.mockSetup(mock)

			// Call the function
			err = insertQuayData(db, tc.datetime, tc.count, tc.kind)

			// Validate the results
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
