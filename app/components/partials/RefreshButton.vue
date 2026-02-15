<script lang="ts" setup>
import type { ButtonProps } from '@nuxt/ui';

interface Props extends /* @vue-ignore */ ButtonProps {
    labelClass?: string;
    iconOnly?: boolean;
}

defineProps<Props>();

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UTooltip :text="$t('common.refresh')">
        <UButton color="gray" variant="link" icon="i-mdi-refresh" v-bind="$attrs">
            <span v-if="!iconOnly" :class="labelClass ? labelClass : 'hidden truncate sm:block'">
                {{ $t('common.refresh') }}
            </span>
        </UButton>
    </UTooltip>
</template>
