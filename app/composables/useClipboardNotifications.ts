import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

type ClipboardNotification = 'copied' | 'link_copied' | 'email_address_copied';

export function useClipboardNotifications() {
    const notifications = useNotificationsStore();

    function notify(type: ClipboardNotification, duration = 3250): void {
        notifications.add({
            title: { key: `notifications.clipboard.${type}.title`, parameters: {} },
            description: { key: `notifications.clipboard.${type}.content`, parameters: {} },
            duration,
            type: NotificationType.INFO,
        });
    }

    return {
        notifyCopied: (duration?: number) => notify('copied', duration),
        notifyLinkCopied: (duration?: number) => notify('link_copied', duration),
        notifyEmailAddressCopied: (duration?: number) => notify('email_address_copied', duration),
    };
}
