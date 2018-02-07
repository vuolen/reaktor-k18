package db

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type cache struct {
	Locations []location
	Logs      []temperatureLog
}

type cachedDatabase struct {
	Database *sql.DB
	Cache    cache
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

func openDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
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

func getTableLength(db *sql.DB, table string) (int, error) {
	var len int
	row := db.QueryRow("select count(*) from " + table)
	err := row.Scan(&len)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return len, nil
}

func CachedDatabase() (cachedDatabase, error) {
	database, err := openDatabase()
	if err != nil {
		return cachedDatabase{}, errors.WithStack(err)
	}
	err = createTables(database)
	if err != nil {
		return cachedDatabase{}, errors.WithStack(err)
	}
	cdb := cachedDatabase{
		database,
		cache{
			make([]location, 0),
			make([]temperatureLog, 0),
		},
	}
	cdb.load()
	return cdb, nil
}

func (cdb *cachedDatabase) Close() {
	cdb.Database.Close()
}

func (cdb *cachedDatabase) load() error {
	cdb.Cache.Locations = make([]location, 0)
	rows, err := cdb.Database.Query("select * from locations")
	if err != nil {
		return errors.WithStack(err)
	}
	for i := 0; rows.Next(); i++ {
		var loc location
		rows.Scan(&loc.Id, &loc.Name, &loc.Lat, &loc.Long)
		cdb.Cache.Locations = append(cdb.Cache.Locations, loc)
	}

	cdb.Cache.Logs = make([]temperatureLog, 0)
	rows, err = cdb.Database.Query("select * from logs")
	if err != nil {
		return errors.WithStack(err)
	}
	for i := 0; rows.Next(); i++ {
		var log temperatureLog
		rows.Scan(&log.LocationId, &log.Time, &log.Temperature)
		cdb.Cache.Logs = append(cdb.Cache.Logs, log)
	}

	return nil
}

func (cdb cachedDatabase) GetJson() ([]byte, error) {
	b, err := json.Marshal(cdb.Cache)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
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
