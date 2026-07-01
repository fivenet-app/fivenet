import mitt from 'mitt';
import type { Notification } from '~/types/notifications';

type NotificationToastEvents = {
    add: Notification;
};

export const notificationToastEvents = mitt<NotificationToastEvents>();
