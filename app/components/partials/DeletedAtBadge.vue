<script lang="ts" setup>
import type { BadgeProps } from '@nuxt/ui';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { I18nDateTimeFormats } from '~/utils/time';

interface Props extends /* @vue-ignore */ BadgeProps {
    deletedAt: Date | Timestamp;
    hideDate?: boolean;
    type?: I18nDateTimeFormats;
}

const props = withDefaults(defineProps<Props>(), {
    hideDate: false,
    type: 'short',
    color: 'warning',
    size: 'sm',
    icon: 'i-mdi-calendar-remove',
});

const attrs = useAttrs();

const badgeProps = computed(() => {
    const { deletedAt: _, hideDate: __, type: ___, ...propsWithoutDeletedAt } = props;

    return Object.assign(
        {
            color: 'warning',
            size: 'sm',
            icon: 'i-mdi-calendar-remove',
        },
        attrs,
        propsWithoutDeletedAt,
    );
});

defineOptions({
    inheritAttrs: false,
});
</script>

<template>
    <UBadge class="inline-flex gap-1 font-bold" v-bind="badgeProps">
        <span>{{ $t('common.deleted') }}</span>

        <GenericTime v-if="!hideDate" :value="deletedAt" :type="type" />
    </UBadge>
</template>
