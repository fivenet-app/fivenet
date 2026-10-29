package properties

import (
	"errors"
	"time"

	propertypb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
)

var ErrInvalidLocalDate = errors.New("invalid local date")

// LocalDateToTime converts a calendar date to UTC midnight solely for DATE
// database storage. It must never be interpreted as an instant for display.
func LocalDateToTime(date *propertypb.LocalDate) (time.Time, error) {
	if date == nil || date.GetYear() < 1 || date.GetYear() > 9999 ||
		date.GetMonth() < 1 || date.GetMonth() > 12 ||
		date.GetDay() < 1 || date.GetDay() > 31 {
		return time.Time{}, ErrInvalidLocalDate
	}

	value := time.Date(
		int(date.GetYear()),
		time.Month(date.GetMonth()),
		int(date.GetDay()),
		0,
		0,
		0,
		0,
		time.UTC,
	)
	if value.Year() != int(date.GetYear()) || int(value.Month()) != int(date.GetMonth()) ||
		value.Day() != int(date.GetDay()) {
		return time.Time{}, ErrInvalidLocalDate
	}
	return value, nil
}

func LocalDateFromTime(value time.Time) *propertypb.LocalDate {
	year, month, day := value.Date()
	yearValue, yearOK := utils.ToInt32Checked(year)
	monthValue, monthOK := utils.ToInt32Checked(month)
	dayValue, dayOK := utils.ToInt32Checked(day)
	if year < 1 || year > 9999 ||
		month < 1 || month > 12 ||
		day < 1 || day > 31 ||
		!yearOK || !monthOK || !dayOK {
		return nil
	}
	return &propertypb.LocalDate{
		Year:  yearValue,
		Month: monthValue,
		Day:   dayValue,
	}
}
