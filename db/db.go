package db

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Location struct {
	id        int
	name      string
	lat, long float64
}

//If these locations do not exist in the database, they are automatically added in.
var defaultLocations = [5]struct {
	name      string
	lat, long float64
}{
	{"Tokio", 35.6584421, 139.7328635},
	{"Helsinki", 60.1697530, 24.9490830},
	{"New York", 40.7406905, -73.9938438},
	{"Amsterdam", 52.3650691, 4.9040238},
	{"Dubai", 25.092535, 55.1562243},
}

func getTableLength(db *sql.DB, table string) (int, error) {
	var len int
	row := db.QueryRow("select count(*) from " + table)
	err := row.Scan(&len)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return len, errors.WithStack(err)
}

func populateLocationsTable(db *sql.DB) error {
	for _, loc := range defaultLocations {
		_, err := db.Exec("insert into locations(name, lat, long) values (?, ?, ?)", loc.name, loc.lat, loc.long)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec("create table if not exists locations (id integer not null primary key autoincrement, name text not null, lat real not null, long real not null)")
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = db.Exec("create table if not exists logs (locationId integer not null, time int not null, temperature real not null)")
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func CreateDatabase(db *sql.DB) error {
	err := createTables(db)
	if err != nil {
		return err
	}

	len, err := getTableLength(db, "locations")
	if err != nil {
		return err
	}

	if len == 0 {
		err = populateLocationsTable(db)
		if err != nil {
			return err
		}
	}

	return nil
}
