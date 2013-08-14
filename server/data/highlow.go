package data

import (
	"log"
	"time"
)

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
