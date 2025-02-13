package fund

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_ListFunds(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful listing with pagination", func(t *testing.T) {
		// Mock count query
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery("SELECT COUNT.*FROM funds").
			WillReturnRows(countRows)

		// Mock data query
		rows := sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow("1", "Fund A", "Description A").
			AddRow("2", "Fund B", "Description B")

		mock.ExpectQuery("SELECT id, name, description.*FROM funds.*ORDER BY name.*LIMIT.*OFFSET.*").
			WithArgs(2, 0). // pageSize=2, offset=0
			WillReturnRows(rows)

		funds, total, err := repo.listFunds(ctx, 1, 2)

		assert.NoError(t, err)
		assert.Equal(t, 3, total)
		assert.Len(t, funds, 2)
		assert.Equal(t, "Fund A", funds[0].Name)
		assert.Equal(t, "Fund B", funds[1].Name)
	})

	t.Run("database error on count", func(t *testing.T) {
		mock.ExpectQuery("SELECT COUNT.*FROM funds").
			WillReturnError(sql.ErrConnDone)

		funds, total, err := repo.listFunds(ctx, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, funds)
		assert.Zero(t, total)
	})

	t.Run("database error on query", func(t *testing.T) {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(3)
		mock.ExpectQuery("SELECT COUNT.*FROM funds").
			WillReturnRows(countRows)

		mock.ExpectQuery("SELECT id, name, description.*FROM funds").
			WillReturnError(sql.ErrConnDone)

		funds, total, err := repo.listFunds(ctx, 1, 10)

		assert.Error(t, err)
		assert.Nil(t, funds)
		assert.Zero(t, total)
	})
}

func TestRepository_GetFundByID(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)
	ctx := context.Background()

	t.Run("successful fund retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "_id"}).
			AddRow("1", "Fund A", "Description A", "Low")

		mock.ExpectQuery("SELECT id, name, description, risk_level_id.*FROM funds.*WHERE id = .*").
			WithArgs("1").
			WillReturnRows(rows)

		fund, err := repo.getFundByID(ctx, "1")

		assert.NoError(t, err)
		assert.NotNil(t, fund)
		assert.Equal(t, "1", fund.ID)
		assert.Equal(t, "Fund A", fund.Name)
		assert.Equal(t, "Description A", fund.Description)
		assert.Equal(t, "Low", fund.RiskLevel)
	})

	t.Run("fund not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, description, risk_level_id.*FROM funds.*WHERE id = .*").
			WithArgs("999").
			WillReturnError(sql.ErrNoRows)

		fund, err := repo.getFundByID(ctx, "999")

		assert.Error(t, err)
		assert.Nil(t, fund)
		assert.Contains(t, err.Error(), "fund not found")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, description, risk_level_id.*FROM funds.*WHERE id = .*").
			WithArgs("1").
			WillReturnError(sql.ErrConnDone)

		fund, err := repo.getFundByID(ctx, "1")

		assert.Error(t, err)
		assert.Nil(t, fund)
		assert.Contains(t, err.Error(), "failed to get fund")
	})
}
