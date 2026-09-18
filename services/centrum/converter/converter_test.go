package centrumconverter

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type dispatchDBMock struct {
	dispatches []*centrumdispatches.Dispatch
}

func (m *dispatchDBMock) Create(
	_ context.Context,
	dsp *centrumdispatches.Dispatch,
) (*centrumdispatches.Dispatch, error) {
	m.dispatches = append(m.dispatches, dsp)
	return dsp, nil
}

func TestConverterGKSPhoneQueryScansIntoDestination(t *testing.T) {
	for _, tc := range []struct {
		name          string
		jobm          string
		wantDispatch  bool
		wantJob       string
		wantMessage   string
		wantAnon      bool
		wantCreatorID int32
	}{
		{name: "empty job", jobm: `[""]`},
		{
			name: "police job", jobm: `["police"]`, wantDispatch: true,
			wantJob: "police", wantMessage: "caller message", wantAnon: true, wantCreatorID: 7,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })

			mock.ExpectQuery(regexp.QuoteMeta("FROM gksphone_job_message")).
				WillReturnRows(sqlmock.NewRows([]string{
					"gksphone_job_message.id", "gksphone_job_message.jobm",
					"gksphone_job_message.anon", "gksphone_job_message.gps",
					"gksphone_job_message.message", "userid",
				}).AddRow(int32(41), tc.jobm, "1", "GPS: 1, 2", tc.wantMessage, int32(7)))
			mock.ExpectExec(regexp.QuoteMeta("UPDATE gksphone_job_message")).
				WithArgs(int32(1), int32(41), int64(1)).
				WillReturnResult(sqlmock.NewResult(0, 1))

			dispatchDB := &dispatchDBMock{}
			c := &Converter{
				logger:           zap.NewNop(),
				db:               db,
				dispatchCreateFn: dispatchDB.Create,
				convertJobs:      []string{"police"},
			}

			require.NoError(t, c.convertGKSPhoneJobMsgToDispatch(t.Context()))
			require.Len(t, dispatchDB.dispatches, boolToInt(tc.wantDispatch))
			if tc.wantDispatch {
				dsp := dispatchDB.dispatches[0]
				require.Equal(t, tc.wantJob, dsp.GetJobs().GetJobs()[0].GetName())
				require.Equal(t, tc.wantMessage, dsp.GetMessage())
				require.Equal(t, tc.wantAnon, dsp.GetAnon())
				require.Equal(t, tc.wantCreatorID, dsp.GetCreatorId())
				require.InEpsilon(t, float64(1), dsp.GetX(), 0.0001)
				require.InEpsilon(t, float64(2), dsp.GetY(), 0.0001)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestConverterLBPhoneQueryScansIntoDestination(t *testing.T) {
	for _, tc := range []struct {
		name         string
		job          string
		wantDispatch bool
	}{
		{name: "empty job", job: " "},
		{name: "police job", job: "police", wantDispatch: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })

			mock.ExpectQuery(regexp.QuoteMeta("FROM phone_services_channels")).
				WillReturnRows(sqlmock.NewRows([]string{
					"id", "company", "phone_number", "message", "x_pos", "y_pos", "userid",
				}).AddRow(int32(42), tc.job, "555-0100", "caller message", int32(123), int32(456), int32(8)))
			mock.ExpectExec(regexp.QuoteMeta("DELETE FROM phone_services_channels")).
				WithArgs(int32(42), int64(1)).
				WillReturnResult(sqlmock.NewResult(0, 1))

			dispatchDB := &dispatchDBMock{}
			c := &Converter{
				logger:           zap.NewNop(),
				db:               db,
				dispatchCreateFn: dispatchDB.Create,
				convertJobs:      []string{"police"},
			}

			require.NoError(t, c.convertLBPhoneJobMsgToDispatch(t.Context()))
			require.Len(t, dispatchDB.dispatches, boolToInt(tc.wantDispatch))
			if tc.wantDispatch {
				dsp := dispatchDB.dispatches[0]
				require.Equal(t, "police", dsp.GetJobs().GetJobs()[0].GetName())
				require.Equal(t, "caller message", dsp.GetMessage())
				require.InEpsilon(t, float64(123), dsp.GetX(), 0.0001)
				require.InEpsilon(t, float64(456), dsp.GetY(), 0.0001)
				require.Equal(t, int32(8), dsp.GetCreatorId())
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
