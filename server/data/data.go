package data

import (
	"database/sql"
	"github.com/derekpitt/weather_site/postdata"
	"github.com/derekpitt/weather_station/loop2packet"
	_ "github.com/go-sql-driver/mysql"
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
	_, err := db.Query(`
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
