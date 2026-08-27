package demo

import (
	"context"
	"fmt"

	livemap "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap"
	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

const (
	demoMarkerSuffix = " [DEMO]"
)

type demoMarkerGenerator struct{}

func (g demoMarkerGenerator) Name() string { return "demo_markers" }

func (g demoMarkerGenerator) Enabled(d *Demo) bool {
	return d.livemapStore != nil && d.cfg.Demo.Features.Markers
}

// seedDemoMarkers creates recognizable public landmarks and a few shape
// examples. Labels are generic so demo mode does not import server-specific
// staff, faction, or creator data from the source marker export.
func (d *Demo) seedDemoMarkers(ctx context.Context) error {
	if d.livemapStore == nil {
		return nil
	}

	job := d.cfg.Demo.TargetJob
	markers := demoMarkers(job)
	tMarkers := table.FivenetCentrumMarkers
	if _, err := tMarkers.DELETE().WHERE(mysql.OR(
		tMarkers.Name.LIKE(mysql.String("%"+demoMarkerSuffix)),
	)).ExecContext(ctx, d.db); err != nil {
		return fmt.Errorf("failed to clear demo markers. %w", err)
	}

	for _, marker := range markers {
		if _, err := d.livemapStore.CreateMarker(ctx, marker, nil, job); err != nil {
			return fmt.Errorf("failed to create demo marker %q. %w", marker.GetName(), err)
		}
	}

	d.logger.Info("completed demo marker seeding", zap.Int("count", len(markers)))
	return nil
}

func (g demoMarkerGenerator) Run(ctx context.Context, d *Demo) error {
	return d.seedDemoMarkers(ctx)
}

func demoMarkers(job string) []*livemapmarkers.MarkerMarker {
	public := true
	return []*livemapmarkers.MarkerMarker{
		demoIconMarker(
			"Central Police Station",
			425.1,
			-979.5,
			"Police station",
			"PoliceBadgeIcon",
			"#243EAC",
			job,
			public,
		),
		demoIconMarker(
			"City Hospital",
			1152.9,
			-1529.6,
			"Emergency hospital",
			"HospitalBuildingIcon",
			"#26B7DA",
			job,
			public,
		),
		demoIconMarker(
			"Paleto Bay Hospital",
			-246.4,
			6341.9,
			"Regional hospital",
			"HospitalBuildingIcon",
			"#26B7DA",
			job,
			public,
		),
		demoIconMarker(
			"Los Santos Courthouse",
			235.1,
			-421.3,
			"Courthouse and civic services",
			"ScaleBalanceIcon",
			"#A855F7",
			job,
			public,
		),
		demoIconMarker(
			"Mirror Park Police Department",
			1126.8,
			-482.3,
			"Police station",
			"PoliceBadgeIcon",
			"#243EAC",
			job,
			public,
		),
		demoIconMarker(
			"Performance Garage",
			-97.8,
			-1805.1,
			"Vehicle repair and tuning",
			"HammerWrenchIcon",
			"#F59E0B",
			job,
			public,
		),
		demoIconMarker(
			"Marina",
			94.9,
			-2999.8,
			"Boat services",
			"SailBoatIcon",
			"#14B8A6",
			job,
			public,
		),
		demoIconMarker(
			"News Station",
			-589.1,
			-933.5,
			"News and media",
			"CameraIcon",
			"#64748B",
			job,
			public,
		),
		demoCircleMarker(
			"Prison Perimeter",
			1691.3,
			2586.3,
			"Restricted area",
			150,
			"#E11D48",
			job,
			public,
		),
		demoRectangleMarker(
			"Power Facility",
			2673.9,
			1385.8,
			2830.7,
			1720.4,
			"Restricted area",
			"#EF4444",
			job,
			public,
		),
		demoPolygonMarker(
			"National Guard Base",
			[]*livemap.Coords{
				{
					X: -1807.7,
					Y: 2651.8,
				},
				{
					X: -1564.0,
					Y: 3212.8,
				},
				{
					X: -2334.1,
					Y: 3612.1,
				},
				{
					X: -2828.7,
					Y: 3459.7,
				},
				{
					X: -2932.8,
					Y: 3276.8,
				},
			},
			"Restricted area",
			"#1A5D0A",
			job,
			public,
		),
	}
}

func demoIconMarker(
	name string,
	x, y float64,
	description, icon, color, job string,
	public bool,
) *livemapmarkers.MarkerMarker {
	marker := baseDemoMarker(name, x, y, description, color, job, public)
	marker.Type = livemapmarkers.MarkerType_MARKER_TYPE_ICON
	marker.Data = &livemapmarkers.MarkerData{}
	marker.Data.SetIcon(&livemapmarkers.IconMarker{Icon: icon})
	return marker
}

func demoCircleMarker(
	name string,
	x, y float64,
	description string,
	radius int32,
	color, job string,
	public bool,
) *livemapmarkers.MarkerMarker {
	marker := baseDemoMarker(name, x, y, description, color, job, public)
	marker.Type = livemapmarkers.MarkerType_MARKER_TYPE_CIRCLE
	marker.Data = &livemapmarkers.MarkerData{}
	marker.Data.SetCircle(&livemapmarkers.CircleMarker{Radius: radius, Opacity: new(float32(15))})
	return marker
}

func demoRectangleMarker(
	name string,
	x, y, endX, endY float64,
	description, color, job string,
	public bool,
) *livemapmarkers.MarkerMarker {
	marker := baseDemoMarker(name, x, y, description, color, job, public)
	marker.Type = livemapmarkers.MarkerType_MARKER_TYPE_RECTANGLE
	marker.Data = &livemapmarkers.MarkerData{}
	marker.Data.SetRectangle(
		&livemapmarkers.RectangleMarker{EndX: endX, EndY: endY, Opacity: new(float32(15))},
	)
	return marker
}

func demoPolygonMarker(
	name string,
	points []*livemap.Coords,
	description, color, job string,
	public bool,
) *livemapmarkers.MarkerMarker {
	marker := baseDemoMarker(
		name,
		points[0].GetX(),
		points[0].GetY(),
		description,
		color,
		job,
		public,
	)
	marker.Type = livemapmarkers.MarkerType_MARKER_TYPE_POLYGON
	marker.Data = &livemapmarkers.MarkerData{}
	marker.Data.SetPolygon(&livemapmarkers.PolygonMarker{Points: points, Opacity: new(float32(15))})
	return marker
}

func baseDemoMarker(
	name string,
	x, y float64,
	description, color, job string,
	public bool,
) *livemapmarkers.MarkerMarker {
	return &livemapmarkers.MarkerMarker{
		Name:        name + demoMarkerSuffix,
		X:           x,
		Y:           y,
		Description: &description,
		Color:       &color,
		Job:         job,
		Public:      &public,
	}
}
