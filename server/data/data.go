/*
  query for trends:


*/

package data

import (
	"database/sql"
	"github.com/derekpitt/weather_site/postdata"
	"github.com/derekpitt/weather_station/loop2packet"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var db *sql.DB = nil

func OpenDatabase(username, password, database string) (err error) {
	db, err = sql.Open("mysql", username+":"+password+"@/"+database+"?parseTime=true")
	return err
}

func CloseDatabase() {
	db.Close()
}

type databaseError struct {
	s string
}

func (e databaseError) Error() string {
	return e.s
}

func WriteSample(postData postdata.PostData) error {
	if db == nil {
		return &databaseError{"Database connection nil"}
	}
	_, err := db.Exec(`
    INSERT INTO weather 
      (TimeSampled, BarometerTrend, Barometer, 
       InsideTemperature, InsideHumidity, OutsideTemerature, 
       WindSpeed, WindDirection, AverageWindSpeed10Minute, 
       AverageWindSpeed2Minute, WindGust10Minute, WindDirectionGust10Minute, 
       DewPoint, OutsideHumidity, HeatIndex, WindChill, THSWIndex, 
       RainRate, UVIndex, SolarRadiation, StormRain, DailyRain, 
       Last15MinuteRain, LastHourRain, DailyET, Last24Rain) 
    VALUES 
      (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
  `,
		postData.Time,
		postData.Sample.BarometerTrend,
		postData.Sample.Barometer,

		postData.Sample.InsideTemperature,
		postData.Sample.InsideHumidity,
		postData.Sample.OutsideTemerature,

		postData.Sample.WindSpeed,
		postData.Sample.WindDirection,
		postData.Sample.AverageWindSpeed10Minute,

		postData.Sample.AverageWindSpeed2Minute,
		postData.Sample.WindGust10Minute,
		postData.Sample.WindDirectionGust10Minute,

		postData.Sample.DewPoint,
		postData.Sample.OutsideHumidity,
		postData.Sample.HeatIndex,
		postData.Sample.WindChill,
		postData.Sample.THSWIndex,

		postData.Sample.RainRate,
		postData.Sample.UVIndex,
		postData.Sample.SolarRadiation,
		postData.Sample.StormRain,
		postData.Sample.DailyRain,

		postData.Sample.Last15MinuteRain,
		postData.Sample.LastHourRain,
		postData.Sample.DailyET,
		postData.Sample.Last24Rain,
	)

	return err
}

type SampleFormat struct {
	Time   time.Time
	Sample loop2packet.Loop2Packet
}

func GetLatestSample() (data SampleFormat, err error) {
	if db == nil {
		err = &databaseError{"Database connection nil"}
		return
	}

	data = SampleFormat{}
	query := `
    SELECT
      TimeSampled, BarometerTrend, Barometer, 
      InsideTemperature, InsideHumidity, OutsideTemerature, 
      WindSpeed, WindDirection, AverageWindSpeed10Minute, 
      AverageWindSpeed2Minute, WindGust10Minute, WindDirectionGust10Minute, 
      DewPoint, OutsideHumidity, HeatIndex, WindChill, THSWIndex, 
      RainRate, UVIndex, SolarRadiation, StormRain, DailyRain, 
      Last15MinuteRain, LastHourRain, DailyET, Last24Rain
    FROM weather
    ORDER BY
      TimeSampled desc 
    LIMIT 1`

	err = db.QueryRow(query).Scan(
		&data.Time,
		&data.Sample.BarometerTrend,
		&data.Sample.Barometer,

		&data.Sample.InsideTemperature,
		&data.Sample.InsideHumidity,
		&data.Sample.OutsideTemerature,

		&data.Sample.WindSpeed,
		&data.Sample.WindDirection,
		&data.Sample.AverageWindSpeed10Minute,

		&data.Sample.AverageWindSpeed2Minute,
		&data.Sample.WindGust10Minute,
		&data.Sample.WindDirectionGust10Minute,

		&data.Sample.DewPoint,
		&data.Sample.OutsideHumidity,
		&data.Sample.HeatIndex,
		&data.Sample.WindChill,
		&data.Sample.THSWIndex,

		&data.Sample.RainRate,
		&data.Sample.UVIndex,
		&data.Sample.SolarRadiation,
		&data.Sample.StormRain,
		&data.Sample.DailyRain,

		&data.Sample.Last15MinuteRain,
		&data.Sample.LastHourRain,
		&data.Sample.DailyET,
		&data.Sample.Last24Rain,
	)

	return
}

/* Trends are just an average of some data point for an hour over like 3 hours */
type Trend struct {
	AverageData float64
	Year        int
	Month       int
	Day         int
	Hour        int
}

func Get3HourTrend(column string, latestTime time.Time) (trends []Trend, err error) {
	if db == nil {
		err = &databaseError{"Database connection nil"}
		return
	}

	trends = make([]Trend, 0)

	query := `
    select
      year(TimeSampled) as 'Year',
      month(TimeSampled) as 'Month',
      day(TimeSampled) as 'Day',
      hour(TimeSampled) as 'Hour',
      avg(` + column + `) as 'Data'
    from weather
    where
      TimeSampled >= ? and
      TimeSampled <= ?
     group by
      Year, Month, Day, Hour
  `

	// we subtract 2 since we are including the current hour
	threeHoursAgo := latestTime.Add(-2 * time.Hour)

	fromFormatted := threeHoursAgo.Format("2006-01-02 15:00:00")
	toFormatted := latestTime.Format("2006-01-02 15:59:59")

	rows, err := db.Query(query, fromFormatted, toFormatted)

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		t := Trend{}
		err = rows.Scan(&t.Year, &t.Month, &t.Day, &t.Hour, &t.AverageData)

		if err != nil {
			continue
		}

		trends = append(trends, t)
	}

	return
}

type HighLow struct {
	High  float64
	Low   float64
	Year  int
	Month int
	Day   int
}

func Get7DayHighLow(column string, latestTime time.Time) (highlows []HighLow, err error) {
	if db == nil {
		err = &databaseError{"Database connection nil"}
		return
	}

	highlows = make([]HighLow, 0)

	query := `
    select
      year(TimeSampled) as 'Year',
      month(TimeSampled) as 'Month',
      day(TimeSampled) as 'Day',
      max(` + column + `) as 'High',
      min(` + column + `) as 'Low'
    from weather
    where
      TimeSampled >= ?
     group by
      Year, Month, Day
  `

	thirtyDaysAgo := latestTime.Add(-7 * (time.Hour * 24))

	fromFormatted := thirtyDaysAgo.Format("2006-01-02")

	rows, err := db.Query(query, fromFormatted)

	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		hl := HighLow{}
		err = rows.Scan(&hl.Year, &hl.Month, &hl.Day, &hl.High, &hl.Low)

		if err != nil {
			continue
		}

		highlows = append(highlows, hl)
	}

	return
}
