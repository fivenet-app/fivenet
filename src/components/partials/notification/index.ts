import NotificationProvider from './NotificationProvider.vue';
import store from './store';

const { dispatchNotification } = store.actions;

export { NotificationProvider, dispatchNotification };
