package notificationsstore

import (
	"context"

	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Store) ListPreferences(
	ctx context.Context,
	userID int32,
) ([]*resourcesnotifications.NotificationPreference, error) {
	stmt := tPreferences.
		SELECT(
			tPreferences.Category.AS("NotificationPreference.category"),
			tPreferences.Kind.AS("NotificationPreference.kind"),
			tPreferences.InboxEnabled.AS("NotificationPreference.inbox_enabled"),
			tPreferences.ToastEnabled.AS("NotificationPreference.toast_enabled"),
			tPreferences.SoundEnabled.AS("NotificationPreference.sound_enabled"),
		).
		FROM(tPreferences).
		WHERE(tPreferences.UserID.EQ(mysql.Int32(userID))).
		ORDER_BY(tPreferences.Category, tPreferences.Kind).
		LIMIT(64)

	var preferences []*resourcesnotifications.NotificationPreference
	if err := stmt.QueryContext(ctx, s.db, &preferences); err != nil {
		return nil, err
	}

	return preferences, nil
}

func (s *Store) UpsertPreference(
	ctx context.Context,
	userID int32,
	preference *resourcesnotifications.NotificationPreference,
) error {
	stmt := tPreferences.
		INSERT(
			tPreferences.UserID,
			tPreferences.Category,
			tPreferences.Kind,
			tPreferences.InboxEnabled,
			tPreferences.ToastEnabled,
			tPreferences.SoundEnabled,
		).
		VALUES(
			userID,
			int32(preference.GetCategory()),
			int32(preference.GetKind()),
			preference.InboxEnabled,
			preference.ToastEnabled,
			preference.SoundEnabled,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tPreferences.InboxEnabled.SET(mysql.BoolExp(mysql.RawString("VALUES(`inbox_enabled`)"))),
			tPreferences.ToastEnabled.SET(mysql.BoolExp(mysql.RawString("VALUES(`toast_enabled`)"))),
			tPreferences.SoundEnabled.SET(mysql.BoolExp(mysql.RawString("VALUES(`sound_enabled`)"))),
		)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) DeletePreference(
	ctx context.Context,
	userID int32,
	category resourcesnotifications.NotificationCategory,
	kind resourcesnotifications.NotificationKind,
) error {
	stmt := tPreferences.
		DELETE().
		WHERE(
			tPreferences.UserID.EQ(mysql.Int32(userID)).
				AND(tPreferences.Category.EQ(mysql.Int32(int32(category)))).
				AND(tPreferences.Kind.EQ(mysql.Int32(int32(kind)))),
		).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

// ResolveDelivery resolves global, category, and kind overrides in one query.
func (s *Store) ResolveDelivery(
	ctx context.Context,
	userID int32,
	category resourcesnotifications.NotificationCategory,
	kind resourcesnotifications.NotificationKind,
) (*resourcesnotifications.NotificationDelivery, error) {
	target := mysql.SELECT(mysql.Int32(userID).AS("user_id")).AsTable("target")
	targetUserID := mysql.IntegerColumn("user_id").From(target)
	global := tPreferences.AS("global_preference")
	categoryPreference := tPreferences.AS("category_preference")
	kindPreference := tPreferences.AS("kind_preference")

	stmt := target.
		SELECT(
			mysql.BoolExp(mysql.COALESCE(kindPreference.InboxEnabled, categoryPreference.InboxEnabled, global.InboxEnabled, mysql.Bool(true))).
				AS("NotificationDelivery.inbox_enabled"),
			mysql.BoolExp(mysql.COALESCE(kindPreference.ToastEnabled, categoryPreference.ToastEnabled, global.ToastEnabled, mysql.Bool(true))).
				AS("NotificationDelivery.toast_enabled"),
			mysql.BoolExp(mysql.COALESCE(kindPreference.SoundEnabled, categoryPreference.SoundEnabled, global.SoundEnabled, mysql.Bool(true))).
				AS("NotificationDelivery.sound_enabled"),
		).
		FROM(
			target.
				LEFT_JOIN(global, global.UserID.EQ(targetUserID).
					AND(global.Category.EQ(mysql.Int32(0))).
					AND(global.Kind.EQ(mysql.Int32(0)))).
				LEFT_JOIN(categoryPreference, categoryPreference.UserID.EQ(targetUserID).
					AND(categoryPreference.Category.EQ(mysql.Int32(int32(category)))).
					AND(categoryPreference.Kind.EQ(mysql.Int32(0)))).
				LEFT_JOIN(kindPreference, kindPreference.UserID.EQ(targetUserID).
					AND(kindPreference.Category.EQ(mysql.Int32(int32(category)))).
					AND(kindPreference.Kind.EQ(mysql.Int32(int32(kind))))),
		)

	delivery := &resourcesnotifications.NotificationDelivery{}
	if err := stmt.QueryContext(ctx, s.db, delivery); err != nil {
		return nil, err
	}

	return delivery, nil
}
