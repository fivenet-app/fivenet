<script lang="ts" setup>
import { notificationCategoryToIcon } from '~/components/notifications/helpers';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import {
    NotificationCategory,
    NotificationKind,
    type NotificationPreference,
} from '~~/gen/ts/resources/notifications/notifications';
import type { NotificationDeliverySetting } from '~/components/user-settings/NotificationDeliverySettingsRow.vue';

const { activeChar } = useAuth();
const notificationsClient = await getNotificationsNotificationsClient();

const notificationCategories = [
    NotificationCategory.GENERAL,
    NotificationCategory.DOCUMENT,
    NotificationCategory.CALENDAR,
    NotificationCategory.JOBS,
    NotificationCategory.QUALIFICATIONS,
    NotificationCategory.MAILER,
    NotificationCategory.SYSTEM,
];
const notificationDeliveryScopes = [
    { category: NotificationCategory.UNSPECIFIED, labelKey: 'components.auth.user_settings.notification_delivery.all' },
    ...notificationCategories.map((category) => ({
        category,
        labelKey: `enums.notifications.NotificationCategory.${NotificationCategory[category]}`,
    })),
];
const notificationKindDefinitions = [
    { kind: NotificationKind.JOBS_GROUP_LEADERSHIP_ADDED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.JOBS_GROUP_LEADERSHIP_REMOVED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.JOBS_GROUP_MEMBER_ADDED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.JOBS_GROUP_MEMBER_REMOVED, category: NotificationCategory.JOBS },
    { kind: NotificationKind.DOCUMENT_APPROVAL_ASSIGNED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_CREATED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_DECIDED, category: NotificationCategory.DOCUMENT },
    { kind: NotificationKind.DOCUMENT_REQUEST_CANCELLED, category: NotificationCategory.DOCUMENT },
] as const;

const pendingPreferenceScopes = ref<Set<string>>(new Set());

const { data: notificationPreferences } = useAuthedLazyAsyncData(
    'userState',
    'notification-preferences',
    async ({ signal }) => (await notificationsClient.getNotificationPreferences({}, { abort: signal })).response.preferences,
    {
        default: () => [],
        enabled: computed(() => activeChar.value !== null),
    },
);

function notificationPreference(
    category: NotificationCategory,
    kind = NotificationKind.UNSPECIFIED,
): NotificationPreference | undefined {
    return notificationPreferences.value.find((preference) => preference.category === category && preference.kind === kind);
}

function preferenceScopeKey(category: NotificationCategory, kind: NotificationKind): string {
    return `${category}:${kind}`;
}

function notificationKindsForCategory(category: NotificationCategory) {
    return notificationKindDefinitions.filter((definition) => definition.category === category);
}

function deliverySetting(
    category: NotificationCategory,
    kind: NotificationKind,
    setting: NotificationDeliverySetting,
): boolean {
    return (
        notificationPreference(category, kind)?.[setting] ??
        notificationPreference(category)?.[setting] ??
        notificationPreference(NotificationCategory.UNSPECIFIED)?.[setting] ??
        true
    );
}

function isPreferencePending(category: NotificationCategory, kind = NotificationKind.UNSPECIFIED): boolean {
    return pendingPreferenceScopes.value.has(preferenceScopeKey(category, kind));
}

function deliverySettings(
    category: NotificationCategory,
    kind: NotificationKind,
): Record<NotificationDeliverySetting, boolean> {
    return {
        inboxEnabled: deliverySetting(category, kind, 'inboxEnabled'),
        toastEnabled: deliverySetting(category, kind, 'toastEnabled'),
        soundEnabled: deliverySetting(category, kind, 'soundEnabled'),
    };
}

async function updateDeliverySetting(
    category: NotificationCategory,
    kind: NotificationKind,
    setting: NotificationDeliverySetting,
    value: boolean,
): Promise<void> {
    if (isPreferencePending(category, kind)) return;

    pendingPreferenceScopes.value = new Set(pendingPreferenceScopes.value).add(preferenceScopeKey(category, kind));
    const current = notificationPreference(category, kind);
    const preference: NotificationPreference = { category, kind, ...current, [setting]: value };

    try {
        const { response } = await notificationsClient.updateNotificationPreference({ preference });
        notificationPreferences.value = response.preferences;
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        const pending = new Set(pendingPreferenceScopes.value);
        pending.delete(preferenceScopeKey(category, kind));
        pendingPreferenceScopes.value = pending;
    }
}

async function resetDeliveryPreference(category: NotificationCategory, kind = NotificationKind.UNSPECIFIED): Promise<void> {
    if (isPreferencePending(category, kind)) return;

    pendingPreferenceScopes.value = new Set(pendingPreferenceScopes.value).add(preferenceScopeKey(category, kind));
    try {
        const { response } = await notificationsClient.updateNotificationPreference({
            preference: { category, kind },
            reset: true,
        });
        notificationPreferences.value = response.preferences;
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        const pending = new Set(pendingPreferenceScopes.value);
        pending.delete(preferenceScopeKey(category, kind));
        pendingPreferenceScopes.value = pending;
    }
}
</script>

<template>
    <UPageCard
        v-if="activeChar"
        :title="$t('components.auth.user_settings.notification_delivery.title')"
        :description="$t('components.auth.user_settings.notification_delivery.description')"
    >
        <div class="space-y-2">
            <div v-for="scope in notificationDeliveryScopes" :key="scope.category" class="rounded-lg border border-default p-2">
                <div class="grid gap-2 lg:grid-cols-[minmax(12rem,1fr)_auto_auto_auto_auto] lg:items-center">
                    <UCollapsible
                        v-if="notificationKindsForCategory(scope.category).length > 0"
                        class="contents"
                        :unmount-on-hide="false"
                    >
                        <div class="flex items-center gap-2">
                            <UIcon :name="notificationCategoryToIcon(scope.category)" class="size-5 text-muted" />
                            <UButton
                                class="group min-w-0 flex-1 justify-start px-0 py-0"
                                color="neutral"
                                variant="link"
                                trailing-icon="i-mdi-chevron-down"
                                :label="$t(scope.labelKey)"
                                :ui="{
                                    label: 'font-medium truncate text-highlighted',
                                    trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                                }"
                            />
                        </div>
                        <template #content>
                            <div class="col-span-full mt-2 space-y-2 border-t border-default pt-2">
                                <NotificationDeliverySettingsRow
                                    v-for="definition in notificationKindsForCategory(scope.category)"
                                    :key="definition.kind"
                                    class="grid gap-2 lg:grid-cols-[minmax(12rem,1fr)_auto_auto_auto_auto] lg:items-center"
                                    :label="$t(`enums.notifications.NotificationKind.${NotificationKind[definition.kind]}`)"
                                    :settings="deliverySettings(scope.category, definition.kind)"
                                    :pending="isPreferencePending(scope.category, definition.kind)"
                                    :can-reset="notificationPreference(scope.category, definition.kind) !== undefined"
                                    @update="
                                        (setting, value) =>
                                            updateDeliverySetting(scope.category, definition.kind, setting, value)
                                    "
                                    @reset="resetDeliveryPreference(scope.category, definition.kind)"
                                />
                            </div>
                        </template>
                    </UCollapsible>
                    <div v-else class="flex items-center gap-2">
                        <UIcon :name="notificationCategoryToIcon(scope.category)" class="size-5 text-muted" />
                        <span class="font-medium">{{ $t(scope.labelKey) }}</span>
                    </div>
                    <NotificationDeliverySettingsRow
                        class="lg:col-span-4"
                        :settings="deliverySettings(scope.category, NotificationKind.UNSPECIFIED)"
                        :pending="isPreferencePending(scope.category)"
                        :can-reset="notificationPreference(scope.category) !== undefined"
                        @update="
                            (setting, value) =>
                                updateDeliverySetting(scope.category, NotificationKind.UNSPECIFIED, setting, value)
                        "
                        @reset="resetDeliveryPreference(scope.category)"
                    />
                </div>
            </div>
        </div>
    </UPageCard>
</template>
