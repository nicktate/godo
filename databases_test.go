package godo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var db = Database{
	ID:          "da4e0206-d019-41d7-b51f-deadbeefbb8f",
	Name:        "dbtest",
	EngineSlug:  "pg",
	VersionSlug: "11",
	Connection: &DatabaseConnection{
		URI:      "postgres://doadmin:zt91mum075ofzyww@dbtest-do-user-3342561-0.db.ondigitalocean.com:25060/defaultdb?sslmode=require",
		Database: "",
		Host:     "dbtest-do-user-3342561-0.db.ondigitalocean.com",
		Port:     25060,
		User:     "doadmin",
		Password: "zt91mum075ofzyww",
		SSL:      true,
	},
	Users: []DatabaseUser{
		DatabaseUser{
			Name:     "doadmin",
			Role:     "primary",
			Password: "zt91mum075ofzyww",
		},
	},
	DBNames: []string{
		"defaultdb",
	},
	NumNodes:   3,
	RegionSlug: "sfo2",
	Status:     "online",
	CreatedAt:  time.Date(2019, 2, 26, 6, 12, 39, 0, time.UTC),
	MaintenanceWindow: &DatabaseMaintenanceWindow{
		Day:         "monday",
		Hour:        "13:51:14",
		Pending:     false,
		Description: nil,
	},
	SizeSlug: "db-s-2vcpu-4gb",
}

var dbJSON = `
{
	"id": "da4e0206-d019-41d7-b51f-deadbeefbb8f",
	"name": "dbtest",
	"engine": "pg",
	"version": "11",
	"connection": {
		"uri": "postgres://doadmin:zt91mum075ofzyww@dbtest-do-user-3342561-0.db.ondigitalocean.com:25060/defaultdb?sslmode=require",
		"database": "",
		"host": "dbtest-do-user-3342561-0.db.ondigitalocean.com",
		"port": 25060,
		"user": "doadmin",
		"password": "zt91mum075ofzyww",
		"ssl": true
	},
	"users": [
		{
			"name": "doadmin",
			"role": "primary",
			"password": "zt91mum075ofzyww"
		}
	],
	"db_names": [
		"defaultdb"
	],
	"num_nodes": 3,
	"region": "sfo2",
	"status": "online",
	"created_at": "2019-02-26T06:12:39Z",
	"maintenance_window": {
		"day": "monday",
		"hour": "13:51:14",
		"pending": false,
		"description": null
	},
	"size": "db-s-2vcpu-4gb"
}
`

var dbsJSON = fmt.Sprintf(`
{
  "databases": [
	%s
  ]
}
`, dbJSON)

func TestDatabases_List(t *testing.T) {
	setup()
	defer teardown()

	dbSvc := client.Databases

	want := []Database{db}

	mux.HandleFunc("/v2/databases", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, dbsJSON)
	})

	got, _, err := dbSvc.List(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestDatabases_Get(t *testing.T) {
	setup()
	defer teardown()

	body := fmt.Sprintf(`
{
  "database": %s
}
`, dbJSON)

	mux.HandleFunc("/v2/databases/da4e0206-d019-41d7-b51f-deadbeefbb8f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, body)
	})

	got, _, err := client.Databases.Get(ctx, "da4e0206-d019-41d7-b51f-deadbeefbb8f")
	require.NoError(t, err)
	require.Equal(t, &db, got)
}

func TestDatabases_Create(t *testing.T) {
	setup()
	defer teardown()

	want := &Database{
		ID:          "8d91899c-0739-4a1a-acc5-deadbeefbb8f",
		Name:        "backend-test",
		EngineSlug:  "pg",
		VersionSlug: "10",
		Connection: &DatabaseConnection{
			URI:      "postgres://doadmin:zt91mum075ofzyww@dbtest-do-user-3342561-0.db.ondigitalocean.com:25060/defaultdb?sslmode=require",
			Database: "",
			Host:     "dbtest-do-user-3342561-0.db.ondigitalocean.com",
			Port:     25060,
			User:     "doadmin",
			Password: "zt91mum075ofzyww",
			SSL:      true,
		},
		Users:             nil,
		DBNames:           nil,
		NumNodes:          2,
		RegionSlug:        "nyc3",
		Status:            "creating",
		CreatedAt:         time.Date(2019, 2, 26, 6, 12, 39, 0, time.UTC),
		MaintenanceWindow: nil,
		SizeSlug:          "db-s-2vcpu-4gb",
	}
	createRequest := &DatabaseCreateRequest{
		Name:       "backend-test",
		EngineSlug: "pg",
		Version:    "10",
		Region:     "nyc3",
		SizeSlug:   "db-s-2vcpu-4gb",
		NumNodes:   2,
	}

	body := `
{
	"database": {
		"id": "8d91899c-0739-4a1a-acc5-deadbeefbb8f",
		"name": "backend-test",
		"engine": "pg",
		"version": "10",
		"connection": {
			"uri": "postgres://doadmin:zt91mum075ofzyww@dbtest-do-user-3342561-0.db.ondigitalocean.com:25060/defaultdb?sslmode=require",
			"database": "",
			"host": "dbtest-do-user-3342561-0.db.ondigitalocean.com",
			"port": 25060,
			"user": "doadmin",
			"password": "zt91mum075ofzyww",
			"ssl": true
		},
		"users": null,
		"db_names": null,
		"num_nodes": 2,
		"region": "nyc3",
		"status": "creating",
		"created_at": "2019-02-26T06:12:39Z",
		"maintenance_window": null,
		"size": "db-s-2vcpu-4gb"
	}
}`

	mux.HandleFunc("/v2/databases", func(w http.ResponseWriter, r *http.Request) {
		v := new(DatabaseCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, http.MethodPost)
		require.Equal(t, v, createRequest)
		fmt.Fprint(w, body)
	})

	got, _, err := client.Databases.Create(ctx, createRequest)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestDatabases_Delete(t *testing.T) {
	setup()
	defer teardown()

	path := "/v2/databases/deadbeef-dead-4aa5-beef-deadbeef347d"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Databases.Delete(ctx, "deadbeef-dead-4aa5-beef-deadbeef347d")
	require.NoError(t, err)
}

func TestDatabases_Resize(t *testing.T) {
	setup()
	defer teardown()

	resizeRequest := &DatabaseResizeRequest{
		SizeSlug: "db-s-16vcpu-64gb",
		NumNodes: 3,
	}

	path := "/v2/databases/deadbeef-dead-4aa5-beef-deadbeef347d/resize"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Databases.Resize(ctx, "deadbeef-dead-4aa5-beef-deadbeef347d", resizeRequest)
	require.NoError(t, err)
}

func TestDatabases_Migrate(t *testing.T) {
	setup()
	defer teardown()

	migrateRequest := &DatabaseMigrateRequest{
		Region: "lon1",
	}

	path := "/v2/databases/deadbeef-dead-4aa5-beef-deadbeef347d/migrate"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Databases.Migrate(ctx, "deadbeef-dead-4aa5-beef-deadbeef347d", migrateRequest)
	require.NoError(t, err)
}

func TestDatabases_UpdateMaintenance(t *testing.T) {
	setup()
	defer teardown()

	maintenanceRequest := &DatabaseUpdateMaintenanceRequest{
		Day:  "thursday",
		Hour: "16:00",
	}

	path := "/v2/databases/deadbeef-dead-4aa5-beef-deadbeef347d/maintenance"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.Databases.UpdateMaintenance(ctx, "deadbeef-dead-4aa5-beef-deadbeef347d", maintenanceRequest)
	require.NoError(t, err)
}

func TestDatabases_ListBackups(t *testing.T) {
	setup()
	defer teardown()

	want := []DatabaseBackup{
		DatabaseBackup{
			CreatedAt:     time.Date(2019, 1, 11, 18, 42, 27, 0, time.UTC),
			SizeGigabytes: 0.03357696,
		},
		DatabaseBackup{
			CreatedAt:     time.Date(2019, 1, 12, 18, 42, 29, 0, time.UTC),
			SizeGigabytes: 0.03364864,
		},
	}

	body := `
{
  "backups": [
    {
      "created_at": "2019-01-11T18:42:27Z",
      "size_gigabytes": 0.03357696
    },
    {
      "created_at": "2019-01-12T18:42:29Z",
      "size_gigabytes": 0.03364864
    }
  ]
}
`
	path := "/v2/databases/deadbeef-dead-4aa5-beef-deadbeef347d/backups"

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, body)
	})

	got, _, err := client.Databases.ListBackups(ctx, "deadbeef-dead-4aa5-beef-deadbeef347d", nil)
	require.NoError(t, err)
	require.Equal(t, want, got)
}
