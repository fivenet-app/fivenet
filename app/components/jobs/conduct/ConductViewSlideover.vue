<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import IDCopyBadge from '~/components/partials/IDCopyBadge.vue';
import { ConductType, type ConductEntry } from '~~/gen/ts/resources/jobs/conduct';
import { ObjectType } from '~~/gen/ts/resources/notifications/client_view';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { conductTypesToBadgeColor } from './helpers';

const props = defineProps<{
    entry: ConductEntry;
}>();

const emits = defineEmits<{
    (e: 'close', v: boolean): void;
    (e: 'refresh'): void;
}>();

const notifications = useNotificationsStore();

// Handle the client update event
const { sendClientView } = useClientUpdate(ObjectType.JOBS_CONDUCT, () =>
    notifications.add({
        title: { key: 'notifications.jobs.conduct.client_view_update.title', parameters: {} },
        description: { key: 'notifications.jobs.conduct.client_view_update.content', parameters: {} },
        duration: 7500,
        type: NotificationType.INFO,
        actions: [
            {
                label: { key: 'common.refresh', parameters: {} },
                icon: 'i-mdi-refresh',
                onClick: () => emits('refresh'),
            },
        ],
    }),
);

if (props.entry.id > 0) {
    sendClientView(props.entry.id);
}
</script>

<template>
    <USlideover :title="$t('common.entry')" :overlay="false">
        <template #actions>
            <IDCopyBadge :id="entry.id" class="mx-2" prefix="CON" />
        </template>

        <template #body>
            <dl class="divide-y divide-default">
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.created_at') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <GenericTime :value="entry.createdAt" />
                    </dd>
                </div>
                <div v-if="entry.updatedAt" class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.updated_at') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <GenericTime :value="entry.updatedAt" />
                    </dd>
                </div>
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.expires_at') }}
                    </dt>
                    <dd class="mt-1 text-sm leading-6 sm:col-span-2 sm:mt-0">
                        <GenericTime :value="entry.expiresAt" />
                    </dd>
                </div>
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.type') }}
                    </dt>
                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                        <UBadge :color="conductTypesToBadgeColor(entry.type)">
                            {{ $t(`enums.jobs.ConductType.${ConductType[entry.type ?? 0]}`) }}
                        </UBadge>
                    </dd>
                </div>
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.message') }}
                    </dt>
                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                        <!-- eslint-disable-next-line vue/no-v-html -->
                        <p class="max-h-14 overflow-y-scroll break-words" v-html="entry.message" />
                    </dd>
                </div>
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.target') }}
                    </dt>
                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                        <CitizenInfoPopover :user="entry.targetUser" />
                    </dd>
                </div>
                <div class="px-4 py-3 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-0">
                    <dt class="text-sm leading-6 font-medium">
                        {{ $t('common.creator') }}
                    </dt>
                    <dd class="mt-2 max-h-24 text-sm sm:col-span-2 sm:mt-0">
                        <CitizenInfoPopover :user="entry.creator" />
                    </dd>
                </div>
            </dl>
        </template>

        <template #footer>
            <UButton class="flex-1" color="neutral" block @click="$emit('close', false)">
                {{ $t('common.close', 1) }}
            </UButton>
        </template>
    </USlideover>
</template>
