<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';

type WantedBadgeVariant = 'long' | 'short';

interface Props extends /* @vue-ignore */ Omit<BadgeProps, 'variant'> {
    wanted: boolean;
    variant?: WantedBadgeVariant;
}

const props = withDefaults(defineProps<Props>(), {
    variant: 'long',
    color: undefined,
    icon: undefined,
});

const color = computed(() => props.color ?? (props.wanted ? 'error' : 'success'));
const icon = computed(() => props.icon ?? (props.wanted ? 'i-mdi-account-alert' : 'i-mdi-account-check'));

const attrs = useAttrs();

const badgeProps = computed(() => {
    const { wanted: _, variant: __, color: ___, icon: ____, ...badgeProps } = props;

    return Object.assign({}, attrs, badgeProps, {
        color: color.value,
        icon: props.variant === 'long' ? icon.value : undefined,
    });
});

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UBadge v-bind="badgeProps" variant="subtle">
        <template v-if="variant === 'long'">
            <span class="font-medium">{{ $t('common.wanted') }}:</span>
            <span>{{ wanted ? $t('common.yes') : $t('common.no') }}</span>
        </template>
        <span v-else>{{ $t('common.wanted') }}</span>
    </UBadge>
</template>
