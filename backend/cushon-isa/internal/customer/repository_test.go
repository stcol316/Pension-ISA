package customer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stcol316/cushon-isa/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateRetailCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	customer := &models.RetailCustomer{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	mock.ExpectExec("INSERT INTO retail_customers").
		WithArgs(customer.FirstName, customer.LastName, customer.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.createRetailCustomer(ctx, customer)
	assert.NoError(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRetailCustomerByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	email := "john.doe@example.com"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
		AddRow("1", "John", "Doe", email)

	mock.ExpectQuery("SELECT (.+) FROM retail_customers").
		WithArgs(email).
		WillReturnRows(rows)

	customer, err := repo.getRetailCustomerByEmail(ctx, email)
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, "John", customer.FirstName)
	assert.Equal(t, "Doe", customer.LastName)
	assert.Equal(t, email, customer.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetRetailCustomerByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	ctx := context.Background()
	id := "1"

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "email"}).
		AddRow(id, "John", "Doe", "john.doe@test.com")

	mock.ExpectQuery("SELECT (.+) FROM retail_customers").
		WithArgs(id).
		WillReturnRows(rows)

	customer, err := repo.getRetailCustomerByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, customer)
	assert.Equal(t, id, customer.ID)
	assert.Equal(t, "John", customer.FirstName)
	assert.Equal(t, "Doe", customer.LastName)
	assert.Equal(t, "john.doe@test.com", customer.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
