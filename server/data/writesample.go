package data

import (
	"github.com/derekpitt/weather_site/postdata"
)

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

