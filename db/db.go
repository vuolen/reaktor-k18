package db

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type TemperatureDatabase struct {
	*sqlx.DB
}

type Location struct {
	Id   int64   `json:"id" db:"id"`
	Name string  `json:"name" db:"name"`
	Lat  float64 `json:"lat" db:"lat"`
	Long float64 `json:"long" db:"long"`
}

type TemperatureLog struct {
	LocationId  int64   `json:"locationId" db:"locationId"`
	Time        int64   `json:"time" db:"time"`
	Temperature float64 `json:"temperature" db:"temperature"`
}

func (tdb TemperatureDatabase) CreateTables() error {
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

func (tdb TemperatureDatabase) GetLocations() ([]Location, error) {
	locations := make([]Location, 0)
	tdb.Select(&locations, "select * from locations")
	return locations, nil
}

func (tdb TemperatureDatabase) GetLogs() ([]TemperatureLog, error) {
	logs := make([]TemperatureLog, 0)
	tdb.Select(&logs, "select * from logs")
	return logs, nil
}

func (tdb TemperatureDatabase) AddLog(tlog TemperatureLog) error {
	_, err := tdb.Exec("insert into logs(locationId, time, temperature) values (?, ?, ?)", tlog.LocationId, tlog.Time, tlog.Temperature)
	return errors.WithStack(err)
}

func (tdb TemperatureDatabase) PopulateWithDefaults() error {
	var defaultLocations = [5]Location{
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

func (tdb TemperatureDatabase) IsValidLocationId(id int64) (bool, error) {
	validLocation := false
	locations, err := tdb.GetLocations()
	if err != nil {
		return false, err
	}
	for _, loc := range locations {
		if loc.Id == id {
			validLocation = true
			break
		}
	}
	return validLocation, nil
}
