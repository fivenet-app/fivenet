<script lang="ts" setup>
import { ClipboardPlusIcon, EyeIcon } from 'mdi-vue3';
import PhoneNumber from '~/components/partials/citizens/PhoneNumber.vue';
import { attr } from '~/composables/can';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { User } from '~~/gen/ts/resources/users/users';

const clipboardStore = useClipboardStore();
const notifications = useNotificatorStore();

const props = defineProps<{
    user: User;
}>();

function addToClipboard(): void {
    clipboardStore.addUser(props.user);

    notifications.dispatchNotification({
        title: { key: 'notifications.clipboard.citizen_add.title', parameters: {} },
        content: { key: 'notifications.clipboard.citizen_add.content', parameters: {} },
        duration: 3250,
        type: 'info',
    });
}
</script>

<template>
    <tr :key="user.userId" class="transition-colors hover:bg-neutral/5 even:bg-base-800">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
            {{ user.firstname }}, {{ user.lastname }}
            <span
                v-if="user.props?.wanted"
                class="inline-flex items-center rounded-full bg-error-100 px-2.5 py-0.5 text-sm font-medium text-error-700 ml-1"
            >
                {{ $t('common.wanted').toUpperCase() }}
            </span>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm">
            {{ user.jobLabel }}
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm">
            {{ user.sex!.toUpperCase() }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
            class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm"
        >
            <PhoneNumber :number="user.phoneNumber" />
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm">
            {{ user.dateofbirth }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints')"
            class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm"
            :class="(user?.props?.trafficInfractionPoints ?? 0) >= 10 ? 'text-error-500' : ''"
        >
            {{ user.props?.trafficInfractionPoints ?? 0 }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')"
            class="whitespace-nowrap px-1 py-1 text-left text-error-500 text-sm"
        >
            <template v-if="(user.props?.openFines ?? 0n) > 0n">
                {{ $n(parseInt((user?.props?.openFines ?? 0n).toString(), 10), 'currency') }}
            </template>
        </td>
        <td class="whitespace-nowrap px-1 py-1 text-left text-base-200 text-sm">{{ user.height }}cm</td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
            <div v-if="can('CitizenStoreService.GetUser')" class="flex flex-row justify-end">
                <button class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard">
                    <ClipboardPlusIcon class="w-6 h-auto ml-auto mr-2.5" />
                </button>
                <NuxtLink
                    :to="{
                        name: 'citizens-id',
                        params: { id: user.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="w-6 h-auto ml-auto mr-2.5" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
