<script lang="ts" setup>
import { User } from '~~/gen/ts/resources/users/users';
import { ClipboardDocumentIcon, EyeIcon } from '@heroicons/vue/24/solid';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificationsStore } from '~/store/notifications';

const store = useClipboardStore();
const notifications = useNotificationsStore();

const { t } = useI18n();

const props = defineProps<{
    user: User,
}>();

function addToClipboard(): void {
    store.addUser(props.user);

    notifications.dispatchNotification({
        title: t('notifications.clipboard.citizen_add.title'),
        content: t('notifications.clipboard.citizen_add.content'),
        duration: 3500,
        type: 'info'
    });
}
</script>

<template>
    <tr :key="user.userId">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ user.firstname }}, {{ user.lastname }}
            <span v-if="user.props?.wanted"
                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700 ml-1">
                {{ $t('common.wanted').toUpperCase() }}
            </span>
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.jobLabel }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.sex!.toUpperCase() }}
        </td>
        <td v-can="'CitizenStoreService.ListCitizens.Fields.PhoneNumber'"
            class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.phoneNumber }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.dateofbirth }}
        </td>
        <td class="whitespace-nowrap px-2 py-2 text-sm text-base-200">
            {{ user.height }}cm
        </td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
            <div class="flex flex-row justify-end">
                <button class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard">
                    <ClipboardDocumentIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
                <NuxtLink v-can="'CitizenStoreService.ListCitizens'"
                    :to="{ name: 'citizens-id', params: { id: user.userId ?? 0 } }"
                    class="flex-initial text-primary-500 hover:text-primary-400">
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
