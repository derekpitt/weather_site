package data

import (
  "log"
  "time"
)


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

