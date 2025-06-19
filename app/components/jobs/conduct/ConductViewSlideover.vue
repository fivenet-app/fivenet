<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';
import { ConductType, type ConductEntry } from '~~/gen/ts/resources/jobs/conduct';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { conductTypesToBadgeColor } from './helpers';

const props = defineProps<{
    entry: ConductEntry & { creator?: { value: Colleague } };
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const { isOpen } = useSlideover();

const notifications = useNotificationsStore();

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.JOBS_CONDUCT, () =>
    notifications.add({
        title: { key: 'notifications.jobs.conduct.client_view_update.title', parameters: {} },
        description: { key: 'notifications.jobs.conduct.client_view_update.content', parameters: {} },
        timeout: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                click: () => emits('refresh'),
            },
        ],
    }),
);

if (props.entry.id > 0) {
    sendClientView(props.entry.id);
}
</script>

<template>
    <USlideover class="flex flex-1" :ui="{ width: 'w-screen max-w-xl' }" :overlay="false">
        <UCard
            class="flex flex-1 flex-col"
            :ui="{
                body: {
                    base: 'flex-1 min-h-[calc(100dvh-(2*var(--header-height)))] max-h-[calc(100dvh-(2*var(--header-height)))] overflow-y-auto',
                    padding: 'px-1 py-2 sm:p-2',
                },
                ring: '',
                divide: 'divide-y divide-gray-100 dark:divide-gray-800',
            }"
        >
            <template #header>
                <div class="flex items-center justify-between">
                    <div class="inline-flex items-center">
                        <IDCopyBadge :id="entry.id" class="mx-2" prefix="CON" />
                        <h3 class="text-2xl font-semibold leading-6">
                            {{ $t('common.entry') }}
                        </h3>
                    </div>

                    <UButton class="-my-1" color="gray" variant="ghost" icon="i-mdi-window-close" @click="isOpen = false" />
                </div>
            </template>

            <div>
                <dl class="divide-neutral/10 divide-y">
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.created_at') }}
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <GenericTime :value="entry.createdAt" />
                        </dd>
                    </div>
                    <div v-if="entry.updatedAt" class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.updated_at') }}
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <GenericTime :value="entry.updatedAt" />
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.expires_at') }}
                        </dt>
                        <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                            <GenericTime :value="entry.expiresAt" />
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.type') }}
                        </dt>
                        <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                            <UBadge :color="conductTypesToBadgeColor(entry.type)">
                                {{ $t(`enums.jobs.ConductType.${ConductType[entry.type ?? 0]}`) }}
                            </UBadge>
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.message') }}
                        </dt>
                        <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                            <p class="max-h-14 overflow-y-scroll break-words">
                                {{ entry.message ?? $t('common.na') }}
                            </p>
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.target') }}
                        </dt>
                        <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                            <CitizenInfoPopover :user="entry.targetUser" />
                        </dd>
                    </div>
                    <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                        <dt class="text-sm font-medium leading-6">
                            {{ $t('common.creator') }}
                        </dt>
                        <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                            <CitizenInfoPopover :user="entry.creator?.value" />
                        </dd>
                    </div>
                </dl>
            </div>

            <template #footer>
                <UButton class="flex-1" color="black" block @click="isOpen = false">
                    {{ $t('common.close', 1) }}
                </UButton>
            </template>
        </UCard>
    </USlideover>
</template>
