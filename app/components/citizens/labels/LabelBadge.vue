<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { Label } from '~~/gen/ts/resources/citizens/labels/labels';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    id?: number;
    label?: Label;
    expiresAt?: Timestamp;
    expired?: boolean;
    hideExpiresAt?: boolean;
}>();

const completorStore = useCompletorStore();

const {
    data: labels,
    pending,
    error,
} = useLazyAsyncData('citizens-labels', () => completorStore.completeCitizenLabels(''), {
    // Load labels when no label object but an id is given
    immediate: !props.label && !!props.id,
});

const label = computed(() => (props.label ? props.label : labels.value?.find((l) => l.id === props.id)));

const expiresAt = computed(() => props.expiresAt ?? label.value?.expiresAt);
</script>

<template>
    <USkeleton v-if="!props.label && pending" class="h-6 w-[185px]" />
    <UBadge v-else-if="error" :label="$t('common.error')" />
    <UBadge
        v-else-if="!pending && label"
        :class="[isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!', expired && 'line-through']"
        :style="{ backgroundColor: label.color }"
        :icon="label.icon && label.icon !== '' ? convertComponentIconNameToDynamic(label.icon) : undefined"
    >
        <div class="inline-flex flex-col gap-1">
            <span>{{ label.name }}</span>

            <div v-if="expiresAt && !hideExpiresAt">
                ({{ $t('common.expires_at') }} <GenericTime :value="expiresAt" type="short" />)
            </div>
        </div>
    </UBadge>
    <UBadge
        v-else
        icon="i-mdi-question-mark"
        variant="subtle"
        :label="`${$t('common.unknown')} ($t('common.id'): {{ id ?? label?.id }})`"
    />
</template>
