package data

import (
	"github.com/derekpitt/weather_station/loop2packet"
	"time"
)

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
