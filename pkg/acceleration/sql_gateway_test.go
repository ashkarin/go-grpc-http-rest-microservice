package acceleration_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"gotest.tools/assert"

	pkg "github.com/ashkarin/go-grpc-http-rest-microservice/pkg/acceleration"
)

var (
	testHasDocker bool
	db            *sql.DB
	entries       []*pkg.Acceleration
)

const (
	createTableStatement = `CREATE TABLE IF NOT EXISTS ` +
		`acceleration("id" SERIAL PRIMARY KEY, "timestamp" TIMESTAMP, ` +
		`"x" DOUBLE PRECISION, "y" DOUBLE PRECISION, "z" DOUBLE PRECISION);`
	storeStatement      = "INSERT INTO acceleration(id, timestamp, x, y, z) VALUES($1,$2,$3,$4,$5)"
	selectByIdStatement = "SELECT id, timestamp, x, y, z FROM acceleration AS a WHERE a.id=$1"
)

func init() {
	if _, err := exec.LookPath("docker"); err == nil {
		testHasDocker = true
	}

	// Create a number of entries
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatalf("An error '%s' was not expected when loading the time location", err)
	}

	entries = []*pkg.Acceleration{
		&pkg.Acceleration{
			ID: 1,
			// Nanosecods is set to 0, since TIMESTAMP datatype does not support it
			Timestamp: time.Date(2019, 04, 11, 9, 00, 20, 0, loc),
			X:         10,
			Y:         11.11,
			Z:         20.20,
		},
		&pkg.Acceleration{
			ID:        2,
			Timestamp: time.Date(2019, 04, 11, 9, 00, 30, 0, loc),
			X:         20,
			Y:         21.21,
			Z:         30.30,
		},
	}
}

func TestMain(m *testing.M) {
	var docker = os.Getenv("DOCKER_URL")
	var pool *dockertest.Pool
	var resource *dockertest.Resource

	if testHasDocker {
		var err error
		pool, err = dockertest.NewPool(docker)
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
		resource, err = connectMongoDB(pool)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run tests
	code := m.Run()

	// Purge resources if required
	if testHasDocker {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	os.Exit(code)
}

func connectMongoDB(pool *dockertest.Pool) (*dockertest.Resource, error) {
	resource, err := pool.Run("postgres", "9.5", nil)
	if err != nil {
		return nil, fmt.Errorf("Could not start resource: %s", err)
	}

	if err = pool.Retry(func() error {
		var err error
		db, err = sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, fmt.Errorf("Could not connect to docker: %s", err)
	}

	return resource, nil
}

func TestSQLGateway_Store(t *testing.T) {
	// Get an empty context
	ctx := context.Background()

	// Open database stub
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Setup a number of expected queries
	mock.ExpectExec(regexp.QuoteMeta(createTableStatement)).WillReturnResult(sqlmock.NewErrorResult(nil))
	mock.ExpectPrepare(regexp.QuoteMeta(storeStatement))
	for _, entry := range entries {
		mock.ExpectExec(regexp.QuoteMeta(storeStatement)).
			WithArgs(entry.ID, entry.Timestamp, entry.X, entry.Y, entry.Z).
			WillReturnResult(sqlmock.NewResult(entry.ID, 1))
	}

	// Create the gateway
	gw, err := pkg.NewSQLStorageGateway(db)
	if err != nil {
		t.Errorf("An error '%s' was not expected when creating a storage gateway", err)
	}

	// Run queries to DB
	gw.Store(ctx, entries...)

	// Check that no expectations left
	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIntegrationSQLGateway_Store(t *testing.T) {
	if !testHasDocker {
		t.Skip("Docker is not available")
	}

	// Get an empty context
	ctx := context.Background()

	// Create the gateway
	gw, err := pkg.NewSQLStorageGateway(db)
	if err != nil {
		t.Errorf("An error '%s' was not expected when creating a storage gateway", err)
	}

	// Run queries to DB
	if err = gw.Store(ctx, entries...); err != nil {
		t.Error(err)
	}

	// Select data and compare with expected
	for _, entry := range entries {
		var ID int64
		var Timestamp time.Time
		var X, Y, Z float64

		err = db.QueryRowContext(ctx, selectByIdStatement, entry.ID).Scan(&ID, &Timestamp, &X, &Y, &Z)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, entry.ID, ID)
		assert.DeepEqual(t, entry.Timestamp, Timestamp)
		assert.Equal(t, entry.X, X)
		assert.Equal(t, entry.Y, Y)
		assert.Equal(t, entry.Z, Z)
	}
}
