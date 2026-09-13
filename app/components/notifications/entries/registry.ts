import type { Component } from 'vue';
import type { Notification } from '~~/gen/ts/resources/notifications/notifications';
import { notificationCategoryDefinition, notificationKindDefinition, type NotificationEntry } from '../definitions';
import CalendarEntry from './Calendar.vue';
import DefaultEntry from './Default.vue';
import DocumentEntry from './Document.vue';
import JobsGroupLeadershipEntry from './JobsGroupLeadership.vue';

const entryComponents: Record<NotificationEntry, Component> = {
    calendar: CalendarEntry,
    document: DocumentEntry,
    'jobs-group-leadership': JobsGroupLeadershipEntry,
};

/**
 * The kind registry takes precedence over category. Feature-specific rows can
 * add their own actions without modifying the inbox shell.
 */
export function notificationEntryComponent(notification: Notification): Component {
    const entry =
        notificationKindDefinition(notification.kind)?.entry ?? notificationCategoryDefinition(notification.category)?.entry;
    return entry ? entryComponents[entry] : DefaultEntry;
}
