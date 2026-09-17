<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

interface Props extends /* @vue-ignore */ BadgeProps {
    deletedAt: Date | Timestamp;
    hideDate?: boolean;
}

defineOptions({
    inheritAttrs: false,
});

const props = withDefaults(defineProps<Props>(), {
    hideDate: false,
    color: 'warning',
    size: 'xs',
    icon: 'i-mdi-calendar-remove',
});

const badgeProps = computed(() => {
    const { deletedAt: _, hideDate: __, ...propsWithoutDeletedAt } = props;

    return propsWithoutDeletedAt;
});
</script>

<template>
    <UBadge class="inline-flex gap-1" v-bind="{ ...badgeProps, ...$attrs }">
        {{ $t('common.deleted') }}
        <GenericTime v-if="!hideDate" :value="deletedAt" type="long" />
    </UBadge>
</template>
