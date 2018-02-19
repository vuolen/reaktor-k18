package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type temperatureDatabase struct {
	*sql.DB
}

type location struct {
	Id        int
	Name      string
	Lat, Long float64
}

type temperatureLog struct {
	LocationId  int
	Time        int
	Temperature float64
}

func (tdb temperatureDatabase) CreateTables() error {
	_, err := tdb.Exec("create table if not exists locations (id integer not null primary key autoincrement, name text not null, lat real not null, long real not null)")
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = tdb.Exec("create table if not exists logs (locationId integer not null, time int not null, temperature real not null)")
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (tdb temperatureDatabase) GetTableLength(table string) (int, error) {
	var len int
	row := tdb.QueryRow("select count(*) from " + table)
	err := row.Scan(&len)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return len, nil
}

func (tdb temperatureDatabase) GetJson() ([]byte, error) {
	return nil, nil
}

func (tdb temperatureDatabase) PopulateWithDefaults() error {
	var defaultLocations = [5]location{
		{-1, "Tokio", 35.6584421, 139.7328635},
		{-1, "Helsinki", 60.1697530, 24.9490830},
		{-1, "New York", 40.7406905, -73.9938438},
		{-1, "Amsterdam", 52.3650691, 4.9040238},
		{-1, "Dubai", 25.092535, 55.1562243},
	}
	for _, loc := range defaultLocations {
		_, err := tdb.Exec("insert into locations(name, lat, long)"+
			"select ?1,?2,?3 where not exists(select 1 from locations where name=?1 and lat=?2 and long=?3)",
			loc.Name, loc.Lat, loc.Long)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

func OpenDatabase(file string) (temperatureDatabase, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return temperatureDatabase{}, errors.WithStack(err)
	}
	tdb := temperatureDatabase{db}
	err = tdb.CreateTables()
	if err != nil {
		return temperatureDatabase{}, errors.WithStack(err)
	}
	err = tdb.PopulateWithDefaults()
	if err != nil {
		return temperatureDatabase{}, errors.WithStack(err)
	}
	return tdb, nil
}

//If these locations do not exist in the database, they are automatically added in.
// var defaultLocations = [5]struct {
// 	name      string
// 	lat, long float64
// }{
// 	{"Tokio", 35.6584421, 139.7328635},
// 	{"Helsinki", 60.1697530, 24.9490830},
// 	{"New York", 40.7406905, -73.9938438},
// 	{"Amsterdam", 52.3650691, 4.9040238},
// 	{"Dubai", 25.092535, 55.1562243},
// }
//
// func createLog(db *sql.DB, l temperatureLog) error {
// 	_, err := db.Exec("insert into logs(locationId, time, temperature) values (?, ?, ?)", l.LocationId, l.Time, l.Temperature)
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}
// 	return nil
// }
//
// func populateLocationsTable(db *sql.DB) error {
// 	for _, loc := range defaultLocations {
// 		_, err := db.Exec("insert into locations(name, lat, long) values (?, ?, ?)", loc.name, loc.lat, loc.long)
// 		if err != nil {
// 			return errors.WithStack(err)
// 		}
// 	}
// 	return nil
// }
