<script lang="ts" setup>
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';
import { User } from '~~/gen/ts/resources/users/users';

const clipboardStore = useClipboardStore();
const notifications = useNotificationsStore();

const props = defineProps<{
    user: User;
}>();

function addToClipboard(): void {
    clipboardStore.addUser(props.user);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: [] },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: [] },
        duration: 3500,
        type: 'info',
    });
}
</script>

<template>
    <tr :key="user.userId">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-base font-medium text-neutral sm:pl-0">
            {{ user.firstname }}, {{ user.lastname }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">{{ user.jobGradeLabel }} ({{ user.jobGrade }})</td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            <span v-for="part in (user?.phoneNumber ?? '').match(/.{1,3}/g)" class="mr-1">{{ part }}</span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ user.dateofbirth }}
        </td>
    </tr>
</template>
