<script lang="ts" setup>
import type { Notification } from '~~/gen/ts/resources/notifications/notifications';
import GenericTime from '~/components/partials/elements/GenericTime.vue';

const props = defineProps<{
    notification: Notification;
    icon: string;
    full?: boolean;
}>();

const emit = defineEmits<{
    (e: 'navigate', notification: Notification): void;
    (e: 'state-changed', notification: Notification): void;
}>();

const notificationCenter = useNotificationCenterStore();
const { t } = useI18n();
const notification = ref<Notification>(props.notification);
watch(
    () => props.notification,
    (value) => (notification.value = value),
);

const actor = computed(() => notification.value.data?.causedBy);
const actorName = computed(() => [actor.value?.firstname, actor.value?.lastname].filter(Boolean).join(' '));

async function updateState(state: { unread?: boolean; starred?: boolean; archived?: boolean }): Promise<void> {
    await notificationCenter.updateState({ ids: [notification.value.id], ...state });
    notification.value = {
        ...notification.value,
        ...(state.unread !== undefined && { readAt: state.unread ? undefined : toTimestamp(new Date()) }),
        ...(state.starred !== undefined && { starred: state.starred }),
        ...(state.archived !== undefined && { archivedAt: state.archived ? toTimestamp(new Date()) : undefined }),
    };
    emit('state-changed', notification.value);
}

const actions = computed(() => [
    {
        label: notification.value.readAt ? t('components.notifications.mark_unread') : t('components.notifications.mark_read'),
        icon: notification.value.readAt ? 'i-mdi-email-mark-as-unread' : 'i-mdi-check',
        onClick: () => updateState({ unread: !!notification.value.readAt }),
    },
    {
        label: notification.value.starred ? t('components.notifications.unstar') : t('components.notifications.star'),
        icon: notification.value.starred ? 'i-mdi-star-off-outline' : 'i-mdi-star-outline',
        onClick: () => updateState({ starred: !notification.value.starred }),
    },
    {
        label: notification.value.archivedAt ? t('components.notifications.restore') : t('components.notifications.archive'),
        icon: notification.value.archivedAt ? 'i-mdi-archive-arrow-up-outline' : 'i-mdi-archive-arrow-down-outline',
        onClick: () => updateState({ archived: !notification.value.archivedAt }),
    },
]);
</script>

<template>
    <li class="group flex gap-3 hover:bg-elevated/50" :class="full ? 'px-4 py-4' : 'px-4 py-3'">
        <UIcon class="mt-0.5 shrink-0 text-muted" :class="full ? 'size-6' : 'size-5'" :name="icon" />
        <div class="min-w-0 flex-1">
            <UButton
                v-if="notification.data?.link"
                class="h-auto w-full justify-start p-0 text-left text-lg font-bold"
                variant="link"
                :color="notification.readAt ? 'neutral' : 'primary'"
                :to="notification.data.link.to"
                :external="notification.data.link.external"
                :label="$t(notification.title?.key ?? '', notification.title?.parameters ?? {})"
                @click="emit('navigate', notification)"
            />
            <p
                v-else
                class="truncate text-lg font-bold"
                :class="[full ? 'text-base' : 'text-sm', notification.readAt ? 'text-muted' : 'text-highlighted']"
            >
                {{ $t(notification.title?.key ?? '', notification.title?.parameters ?? {}) }}
            </p>

            <p class="mt-0.5 line-clamp-2 text-muted" :class="full ? 'text-sm' : 'text-xs'">
                {{ $t(notification.content?.key ?? '', notification.content?.parameters ?? {}) }}
            </p>

            <div class="mt-1 flex flex-wrap items-center gap-x-2 gap-y-1 text-xs text-dimmed">
                <span><GenericTime :value="notification.createdAt" ago /></span>
                <span v-if="actor?.userId" class="inline-flex items-center gap-1">
                    <UAvatar :src="actor.profilePicture" :alt="actorName" size="xs" />
                    <span class="truncate">{{ actorName || $t('common.na') }}</span>
                </span>
            </div>

            <UButton
                v-if="full && notification.data?.link"
                class="mt-2"
                color="neutral"
                variant="link"
                size="xs"
                trailing-icon="i-mdi-arrow-right"
                :to="notification.data.link.to"
                :external="notification.data.link.external"
                :label="$t('common.goto')"
                @click="emit('navigate', notification)"
            />
        </div>

        <UDropdownMenu :items="actions" :content="{ align: 'end' }">
            <UButton
                class="opacity-100 md:opacity-0 md:group-hover:opacity-100 md:focus:opacity-100"
                color="neutral"
                variant="ghost"
                size="xs"
                icon="i-mdi-dots-vertical"
                :aria-label="$t('common.more')"
            />
        </UDropdownMenu>
    </li>
</template>
