import type { Component } from 'vue';
import { NotificationCategory, NotificationKind, type Notification } from '~~/gen/ts/resources/notifications/notifications';
import CalendarEntry from './Calendar.vue';
import DefaultEntry from './Default.vue';
import DocumentEntry from './Document.vue';
import JobsGroupLeadershipEntry from './JobsGroupLeadership.vue';

const categoryEntries: Partial<Record<NotificationCategory, Component>> = {
    [NotificationCategory.DOCUMENT]: DocumentEntry,
    [NotificationCategory.CALENDAR]: CalendarEntry,
};

const kindEntries: Partial<Record<NotificationKind, Component>> = {
    [NotificationKind.JOBS_GROUP_LEADERSHIP_ADDED]: JobsGroupLeadershipEntry,
    [NotificationKind.JOBS_GROUP_LEADERSHIP_REMOVED]: JobsGroupLeadershipEntry,
};

/**
 * The kind registry takes precedence over category. Feature-specific rows can
 * add their own actions without modifying the inbox shell.
 */
export function notificationEntryComponent(notification: Notification): Component {
    return kindEntries[notification.kind] ?? categoryEntries[notification.category] ?? DefaultEntry;
}
