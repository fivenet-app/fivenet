<script lang="ts" setup>
import PhoneNumber from '~/components/partials/users/PhoneNumber.vue';
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
            <PhoneNumber :number="user.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200">
            {{ user.dateofbirth }}
        </td>
    </tr>
</template>
