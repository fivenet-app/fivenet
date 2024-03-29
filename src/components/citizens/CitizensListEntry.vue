<script lang="ts" setup>
import { ClipboardPlusIcon, EyeIcon } from 'mdi-vue3';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
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
    <tr :key="user.userId" class="transition-colors even:bg-base-800 hover:bg-neutral/5">
        <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
            <span>{{ user.firstname }} {{ user.lastname }}</span>
            <span class="lg:hidden"> ({{ user.dateofbirth }}) </span>

            <span v-if="user.props?.wanted" class="ml-1 rounded-md bg-error-100 px-2 py-0.5 text-sm font-medium text-error-700">
                {{ $t('common.wanted').toUpperCase() }}
            </span>

            <dl class="font-normal lg:hidden">
                <dt class="sr-only">{{ $t('common.sex') }} - {{ $t('common.job') }}</dt>
                <dd class="mt-1 truncate text-accent-200">{{ user.sex!.toUpperCase() }} - {{ user.jobLabel }}</dd>
            </dl>
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-sm text-accent-200 lg:table-cell">
            {{ user.jobLabel }}
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-sm text-accent-200 lg:table-cell">
            {{ user.sex!.toUpperCase() }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'PhoneNumber')"
            class="whitespace-nowrap p-1 text-left text-sm text-accent-200"
        >
            <PhoneNumberBlock :number="user.phoneNumber" />
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-sm text-accent-200 lg:table-cell">
            {{ user.dateofbirth }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.TrafficInfractionPoints')"
            class="whitespace-nowrap p-1 text-left text-sm text-accent-200"
            :class="(user?.props?.trafficInfractionPoints ?? 0) >= 10 ? 'text-error-500' : ''"
        >
            {{ user.props?.trafficInfractionPoints ?? 0 }}
        </td>
        <td
            v-if="attr('CitizenStoreService.ListCitizens', 'Fields', 'UserProps.OpenFines')"
            class="whitespace-nowrap p-1 text-left text-sm text-error-500"
        >
            <template v-if="(user.props?.openFines ?? 0n) > 0n">
                {{ $n(parseInt((user?.props?.openFines ?? 0n).toString()), 'currency') }}
            </template>
        </td>
        <td class="hidden whitespace-nowrap p-1 text-left text-sm text-accent-200 md:table-cell">{{ user.height }}cm</td>
        <td class="whitespace-nowrap py-2 pl-3 pr-4 text-sm font-medium sm:pr-0">
            <div v-if="can('CitizenStoreService.GetUser')" class="flex flex-col justify-end gap-1 md:flex-row">
                <button type="button" class="flex-initial text-primary-500 hover:text-primary-400" @click="addToClipboard">
                    <ClipboardPlusIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                </button>

                <NuxtLink
                    :to="{
                        name: 'citizens-id',
                        params: { id: user.userId ?? 0 },
                    }"
                    class="flex-initial text-primary-500 hover:text-primary-400"
                >
                    <EyeIcon class="ml-auto mr-2.5 h-auto w-5" aria-hidden="true" />
                </NuxtLink>
            </div>
        </td>
    </tr>
</template>
