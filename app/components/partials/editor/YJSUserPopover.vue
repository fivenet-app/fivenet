<script lang="ts" setup>
import type GrpcProvider from '~/composables/yjs/yjs';

const yjsProvider = inject<GrpcProvider | undefined>('yjsProvider', undefined);

const awareness = yjsProvider ? useAwarenessUsers(yjsProvider.awareness) : undefined;
</script>

<template>
    <UPopover v-if="awareness" :popper="{ placement: 'top' }" :disabled="(awareness?.users?.value.length || 0) === 0">
        <UButton
            :class="(awareness?.users?.value.length || 0) === 0 && 'cursor-not-allowed'"
            color="white"
            variant="link"
            trailing-icon="i-heroicons-chevron-down-20-solid"
        >
            {{ awareness?.users?.value.length || 0 }} {{ $t('common.user', awareness?.users?.value.length || 0) }}
        </UButton>

        <template #panel>
            <div class="p-4">
                <ul class="grid grid-cols-2 gap-2">
                    <li
                        v-for="(user, idx) in awareness?.users?.value.filter((u) => u !== undefined && u !== null)"
                        :key="idx"
                        class="inline-flex items-center gap-1"
                    >
                        <UBadge
                            class="shrink-0"
                            :style="{ backgroundColor: user.color }"
                            :ui="{ rounded: 'rounded-full' }"
                            size="lg"
                        />
                        {{ user.name }}
                    </li>
                </ul>
            </div>
        </template>
    </UPopover>
</template>
