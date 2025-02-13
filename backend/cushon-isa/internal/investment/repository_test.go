package investment

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	isaerrors "github.com/stcol316/cushon-isa/internal/errors"
	"github.com/stcol316/cushon-isa/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateInvestment(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	tests := []struct {
		name        string
		investment  *models.Investment
		setupMock   func(sqlmock.Sqlmock)
		expectError error
	}{
		{
			name: "successful investment creation",
			investment: &models.Investment{
				CustomerID: "customer1",
				FundID:     "fund1",
				Amount:     float64(100),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				// Expect check for existing fund
				mock.ExpectQuery("SELECT DISTINCT fund_id").
					WithArgs("customer1").
					WillReturnRows(sqlmock.NewRows([]string{"fund_id"}))

				// Expect transaction begin
				mock.ExpectBegin()

				// Expect investment insert
				mock.ExpectExec("INSERT INTO investments").
					WithArgs("customer1", "fund1", float64(100)).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Expect materialized view refresh
				mock.ExpectExec("REFRESH MATERIALIZED VIEW customer_fund_totals").
					WillReturnResult(sqlmock.NewResult(0, 0))

				// Expect transaction commit
				mock.ExpectCommit()
			},
			expectError: nil,
		},
		{
			name: "different fund not allowed",
			investment: &models.Investment{
				CustomerID: "customer1",
				FundID:     "fund2",
				Amount:     float64(100),
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				existingFundID := "fund1"
				mock.ExpectQuery("SELECT DISTINCT fund_id").
					WithArgs("customer1").
					WillReturnRows(sqlmock.NewRows([]string{"fund_id"}).AddRow(existingFundID))
			},
			expectError: isaerrors.ErrDifferentFundNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setupMock(mock)
			err := repo.createInvestment(ctx, test.investment)
			assert.Equal(t, test.expectError, err)
		})
	}
}

func TestListInvestmentsByCustomerID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	mock_investments := []models.Investment{
		{
			ID:         "inv1",
			CustomerID: "customer1",
			FundID:     "fund1",
			Amount:     float64(100),
			CreatedAt:  time.Now(),
		},
		{
			ID:         "inv2",
			CustomerID: "customer1",
			FundID:     "fund1",
			Amount:     float64(200),
			CreatedAt:  time.Now(),
		},
	}
	t.Run("successful listing", func(t *testing.T) {
		customerID := "customer1"
		expectedTotal := 2
		expectedInvestments := mock_investments
		page := 1
		pageSize := 10
		offset := 0

		// Expect count query
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

		// Expect investments query
		rows := sqlmock.NewRows([]string{"id", "customer_id", "fund_id", "amount", "created_at"})
		for _, inv := range expectedInvestments {
			rows.AddRow(inv.ID, inv.CustomerID, inv.FundID, inv.Amount, inv.CreatedAt)
		}

		mock.ExpectQuery("SELECT (.+) FROM investments").
			WithArgs(customerID, pageSize, offset).
			WillReturnRows(rows)

		investments, total, err := repo.listInvestmentsByCustomerID(ctx, customerID, page, pageSize)
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
		assert.Len(t, investments, len(expectedInvestments))
	})

	t.Run("test pagination", func(t *testing.T) {
		customerID := "customer1"
		expectedTotal := 1
		expectedInvestments := mock_investments
		page := 1
		pageSize := 1
		offset := 0

		// Expect count query
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedTotal))

		// Expect investments query
		rows := sqlmock.NewRows([]string{"id", "customer_id", "fund_id", "amount", "created_at"})
		for _, inv := range expectedInvestments {
			rows.AddRow(inv.ID, inv.CustomerID, inv.FundID, inv.Amount, inv.CreatedAt)
		}

		mock.ExpectQuery("SELECT (.+) FROM investments").
			WithArgs(customerID, pageSize, offset).
			WillReturnRows(rows)

		investments, total, err := repo.listInvestmentsByCustomerID(ctx, customerID, page, pageSize)
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
		assert.Len(t, investments, len(expectedInvestments))
	})
}

func TestGetInvestmentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful get", func(t *testing.T) {
		expectedInvestment := &models.Investment{
			ID:         "inv1",
			CustomerID: "customer1",
			FundID:     "fund1",
			Amount:     float64(100),
		}

		mock.ExpectQuery("SELECT (.+) FROM investments").
			WithArgs(expectedInvestment.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "customer_id", "fund_id", "amount"}).
				AddRow(expectedInvestment.ID, expectedInvestment.CustomerID, expectedInvestment.FundID, expectedInvestment.Amount))

		investment, err := repo.getInvestmentByID(ctx, expectedInvestment.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedInvestment, investment)
	})

	t.Run("investment not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM investments").
			WithArgs("non-existent").
			WillReturnError(sql.ErrNoRows)

		investment, err := repo.getInvestmentByID(ctx, "non-existent")
		assert.Error(t, err)
		assert.Nil(t, investment)
	})
}

func TestGetCustomerFundTotal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful get total", func(t *testing.T) {
		expectedSummary := &models.InvestmentSummary{
			CustomerID:      "customer1",
			FirstName:       "John",
			LastName:        "Doe",
			Email:           "john@example.com",
			FundID:          "fund1",
			FundName:        "Test Fund",
			TotalInvestment: float64(300),
		}

		mock.ExpectQuery("SELECT (.+) FROM customer_fund_totals").
			WithArgs(expectedSummary.CustomerID, expectedSummary.FundID).
			WillReturnRows(sqlmock.NewRows([]string{
				"customer_id", "first_name", "last_name", "email",
				"fund_id", "fund_name", "total_investment",
			}).AddRow(
				expectedSummary.CustomerID, expectedSummary.FirstName,
				expectedSummary.LastName, expectedSummary.Email,
				expectedSummary.FundID, expectedSummary.FundName,
				expectedSummary.TotalInvestment,
			))

		summary, err := repo.getCustomerFundTotal(ctx, expectedSummary.CustomerID, expectedSummary.FundID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSummary, summary)
	})
}
