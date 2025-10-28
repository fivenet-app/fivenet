import mitt from 'mitt';
import type { ObjectEvent, ObjectType } from '~~/gen/ts/resources/notifications/client_view';

type ObjectEvents = {
    ready: boolean;
    // Object events
    [ObjectType.UNSPECIFIED]: ObjectEvent;
    [ObjectType.CITIZEN]: ObjectEvent;
    [ObjectType.DOCUMENT]: ObjectEvent;
    [ObjectType.WIKI_PAGE]: ObjectEvent;
    [ObjectType.JOBS_COLLEAGUE]: ObjectEvent;
    [ObjectType.JOBS_CONDUCT]: ObjectEvent;
};

export const notificationsEvents = mitt<ObjectEvents>();

export function useClientUpdate(objType: ObjectType, callback: (event: ObjectEvent) => void) {
    const notifications = useNotificationsStore();
    const { ready } = storeToRefs(notifications);

    let clientViewSent = false;
    const sendClientView = async (objId: number) => {
        if (ready.value) {
            notifications.sendClientView(objType, objId);
        } else {
            const handler = (ready: boolean) => {
                if (!ready) return;

                notifications.sendClientView(objType, objId);
            };

            notificationsEvents.on('ready', handler);
            onBeforeUnmount(() => {
                notificationsEvents.off('ready', handler);
                notifications.sendClientView(objType, objId);
            });
        }

        clientViewSent = true;
    };

    onMounted(() => notifications.onClientUpdate(objType, callback));
    onUnmounted(() => {
        notifications.offClientUpdate(objType, callback);

        if (clientViewSent) {
            notifications.sendClientView(objType); // Reset the client view after unmounting
        }
    });

    return {
        sendClientView,
    };
}
