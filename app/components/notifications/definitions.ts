import { NotificationCategory, NotificationKind } from '~~/gen/ts/resources/notifications/notifications';

export type NotificationEntry = 'calendar' | 'document' | 'jobs-group-leadership';

type NotificationCategoryDefinition = {
    category: NotificationCategory;
    icon: string;
    entry?: NotificationEntry;
};

type NotificationKindDefinition = {
    kind: NotificationKind;
    category: NotificationCategory;
    icon?: string;
    entry?: NotificationEntry;
};

// Keep these definitions in sync with
// services/notifications/notifications.go:validPreferenceScope. This is the
// single frontend source for notification categories, kinds, icons, and rows.
export const notificationCategoryDefinitions: readonly NotificationCategoryDefinition[] = [
    { category: NotificationCategory.GENERAL, icon: 'i-mdi-information-outline' },
    { category: NotificationCategory.DOCUMENT, icon: 'i-mdi-file-document-box-multiple-outline', entry: 'document' },
    { category: NotificationCategory.CALENDAR, icon: 'i-mdi-calendar-outline', entry: 'calendar' },
    { category: NotificationCategory.JOBS, icon: 'i-mdi-briefcase-outline' },
    { category: NotificationCategory.QUALIFICATIONS, icon: 'i-mdi-school-outline' },
    { category: NotificationCategory.MAILER, icon: 'i-mdi-email-outline' },
    { category: NotificationCategory.DISPATCH, icon: 'i-mdi-car-emergency' },
    { category: NotificationCategory.SYSTEM, icon: 'i-mdi-cog-outline' },
] as const;

export const notificationKindDefinitions: readonly NotificationKindDefinition[] = [
    {
        kind: NotificationKind.JOBS_GROUP_LEADERSHIP_ADDED,
        category: NotificationCategory.JOBS,
        icon: 'i-mdi-account-star-outline',
        entry: 'jobs-group-leadership',
    },
    {
        kind: NotificationKind.JOBS_GROUP_LEADERSHIP_REMOVED,
        category: NotificationCategory.JOBS,
        icon: 'i-mdi-account-star-off-outline',
        entry: 'jobs-group-leadership',
    },
    { kind: NotificationKind.JOBS_GROUP_MEMBER_ADDED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.JOBS_GROUP_MEMBER_REMOVED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.DOCUMENT_APPROVAL_ASSIGNED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_CREATED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_DECIDED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_CANCELLED, category: NotificationCategory.DOCUMENT },
] as const;

export function notificationCategoryDefinition(category: NotificationCategory) {
    return notificationCategoryDefinitions.find((definition) => definition.category === category);
}

export function notificationKindDefinition(kind: NotificationKind) {
    return notificationKindDefinitions.find((definition) => definition.kind === kind);
}
