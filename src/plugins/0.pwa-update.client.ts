import { useNotificationsStore } from '~/store/notifications';

export default async (_: any) => {
    //@ts-ignore Only available when in production
    const workbox = await window.$workbox;

    if (!workbox) {
        console.debug("Workbox couldn't be loaded.");
        return;
    }

    workbox.addEventListener('installed', (event: any) => {
        if (!event.isUpdate) {
            console.debug('The PWA is on the latest version.');
            return;
        }

        console.debug('There is an update for the PWA, reloading...');
        window.location.reload();
    });
};
