<script lang="ts" setup>
import type GrpcProvider from '~/composables/yjs/yjs';

const yjsProvider = inject<GrpcProvider | undefined>('yjsProvider', undefined);

const awareness = yjsProvider ? useAwarenessUsers(yjsProvider.awareness) : undefined;

const users = computed(() => {
    return awareness?.users?.value.filter((u) => u !== undefined && u !== null) || [];
});
</script>

<template>
    <UPopover v-if="awareness" :disabled="users.length === 0">
        <UButton
            :class="users.length === 0 && 'cursor-not-allowed'"
            color="neutral"
            variant="link"
            trailing-icon="i-heroicons-chevron-down-20-solid"
        >
            {{ users.length }} {{ $t('common.user', users.length) }}
        </UButton>

        <template #content>
            <div class="p-4">
                <ul class="grid grid-cols-2 gap-2">
                    <li
                        v-for="(user, idx) in users.filter((u) => u !== undefined && u !== null)"
                        :key="idx"
                        class="inline-flex items-center gap-1"
                    >
                        <UBadge class="shrink-0" :style="{ backgroundColor: user.color }" size="lg" />
                        <span>
                            {{ user.name }}
                        </span>
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
