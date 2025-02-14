// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"

	"github.com/go-vela/server/database/postgres/dml"
	"github.com/go-vela/types/library"
)

func TestPostgres_Client_GetWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectWorker, "worker_0").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "last_checked_in", "build_limit"},
	).AddRow(1, "worker_0", "localhost", "{}", true, 0, 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Worker
	}{
		{
			failure: false,
			want:    _worker,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetWorker("worker_0")

		if test.failure {
			if err == nil {
				t.Errorf("GetWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorker returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorker is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_GetWorkerByAddress(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Raw(dml.SelectWorkerByAddress, "localhost").Statement

	// create expected return in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "hostname", "address", "routes", "active", "last_checked_in", "build_limit"},
	).AddRow(1, "worker_0", "localhost", "{}", true, 0, 0)

	// ensure the mock expects the query for test case 1
	_mock.ExpectQuery(_query.SQL.String()).WillReturnRows(_rows)
	// ensure the mock expects the error for test case 2
	_mock.ExpectQuery(_query.SQL.String()).WillReturnError(gorm.ErrRecordNotFound)

	// setup tests
	tests := []struct {
		failure bool
		want    *library.Worker
	}{
		{
			failure: false,
			want:    _worker,
		},
		{
			failure: true,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := _database.GetWorkerByAddress("localhost")

		if test.failure {
			if err == nil {
				t.Errorf("GetWorkerByAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("GetWorkerByAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("GetWorkerByAddress is %v, want %v", got, test.want)
		}
	}
}

func TestPostgres_Client_CreateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// create expected return in mock
	_rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

	// ensure the mock expects the query
	_mock.ExpectQuery(`INSERT INTO "workers" ("hostname","address","routes","active","last_checked_in","build_limit","id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`).
		WithArgs("worker_0", "localhost", "{}", true, nil, nil, 1).
		WillReturnRows(_rows)

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.CreateWorker(_worker)

		if test.failure {
			if err == nil {
				t.Errorf("CreateWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CreateWorker returned err: %v", err)
		}
	}
}

func TestPostgres_Client_UpdateWorker(t *testing.T) {
	// setup types
	_worker := testWorker()
	_worker.SetID(1)
	_worker.SetHostname("worker_0")
	_worker.SetAddress("localhost")
	_worker.SetActive(true)

	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// ensure the mock expects the query
	_mock.ExpectExec(`UPDATE "workers" SET "hostname"=$1,"address"=$2,"routes"=$3,"active"=$4,"last_checked_in"=$5,"build_limit"=$6 WHERE "id" = $7`).
		WithArgs("worker_0", "localhost", "{}", true, nil, nil, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.UpdateWorker(_worker)

		if test.failure {
			if err == nil {
				t.Errorf("UpdateWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("UpdateWorker returned err: %v", err)
		}
	}
}

func TestPostgres_Client_DeleteWorker(t *testing.T) {
	// setup types
	// setup the test database client
	_database, _mock, err := NewTest()
	if err != nil {
		t.Errorf("unable to create new postgres test database: %v", err)
	}

	defer func() { _sql, _ := _database.Postgres.DB(); _sql.Close() }()

	// capture the current expected SQL query
	//
	// https://gorm.io/docs/sql_builder.html#DryRun-Mode
	_query := _database.Postgres.Session(&gorm.Session{DryRun: true}).Exec(dml.DeleteWorker, 1).Statement

	// ensure the mock expects the query
	_mock.ExpectExec(_query.SQL.String()).WillReturnResult(sqlmock.NewResult(1, 1))

	// setup tests
	tests := []struct {
		failure bool
	}{
		{
			failure: false,
		},
	}

	// run tests
	for _, test := range tests {
		err := _database.DeleteWorker(1)

		if test.failure {
			if err == nil {
				t.Errorf("DeleteWorker should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("DeleteWorker returned err: %v", err)
		}
	}
}

// testWorker is a test helper function to create a
// library Worker type with all fields set to their
// zero values.
func testWorker() *library.Worker {
	i64 := int64(0)
	str := ""
	arr := []string{}
	b := false

	return &library.Worker{
		ID:            &i64,
		Hostname:      &str,
		Address:       &str,
		Routes:        &arr,
		Active:        &b,
		LastCheckedIn: &i64,
		BuildLimit:    &i64,
	}
}
