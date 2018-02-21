package db

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

type TemperatureDatabase struct {
	*sqlx.DB
}

type location struct {
	Id   int64   `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type temperatureLog struct {
	LocationId  int64   `json:"locationId"`
	Time        int64   `json:"time"`
	Temperature float64 `json:"temperature"`
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

func (tdb TemperatureDatabase) GetTableLength(table string) (int, error) {
	var len int
	row := tdb.QueryRow("select count(*) from " + table)
	err := row.Scan(&len)
	if err != nil {
		return -1, errors.WithStack(err)
	}
	return len, nil
}

func (tdb TemperatureDatabase) GetLocations() ([]location, error) {
	tableLength, err := tdb.GetTableLength("locations")
	if err != nil {
		return nil, err
	}
	locations := make([]location, 0, tableLength)
	tdb.Select(&locations, "select * from locations")
	return locations, nil
}

func (tdb TemperatureDatabase) GetLogs() ([]temperatureLog, error) {
	tableLength, err := tdb.GetTableLength("logs")
	if err != nil {
		return nil, err
	}
	logs := make([]temperatureLog, 0, tableLength)
	tdb.Select(&logs, "select * from logs")
	return logs, nil
}

func (tdb TemperatureDatabase) AddLog(tlog temperatureLog) error {
	if tlog.Time < 0 || tlog.Time > time.Now().UTC().Unix() {
		return errors.New("Time out of bounds")
	}
	// 373.15 kelvin equals 100 degrees Celcius
	if tlog.Temperature < 0 || tlog.Temperature > 373.15 {
		return errors.New("Temperature out of bounds")
	}
	locations, err := tdb.GetLocations()
	if err != nil {
		log.Print(err)
		return errors.WithMessage(err, "Internal database error")
	}
	validLocation := false
	for _, loc := range locations {
		if loc.Id == tlog.LocationId {
			validLocation = true
		}
	}
	if !validLocation {
		return errors.New("Invalid location")
	}
	_, err = tdb.Exec("insert into logs(locationId, time, temperature) values (?, ?, ?)", tlog.LocationId, tlog.Time, tlog.Temperature)
	return errors.WithStack(err)
}

func (tdb TemperatureDatabase) PopulateWithDefaults() error {
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

func OpenDatabase(file string) (TemperatureDatabase, error) {
	db, err := sqlx.Connect("sqlite3", file)
	if err != nil {
		return TemperatureDatabase{}, errors.WithStack(err)
	}
	tdb := TemperatureDatabase{db}
	err = tdb.CreateTables()
	if err != nil {
		return TemperatureDatabase{}, errors.WithStack(err)
	}
	err = tdb.PopulateWithDefaults()
	if err != nil {
		return TemperatureDatabase{}, errors.WithStack(err)
	}
	return tdb, nil
}
