<script lang="ts" setup>
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import type { PropertyValue } from '~~/gen/ts/resources/properties/properties';
import type { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';

const props = defineProps<{
    propertyKey: string;
    oldValue?: PropertyValue;
    newValue?: PropertyValue;
    lifecycle?: {
        expiresAt?: Timestamp;
        expiryCleared?: boolean;
        stateReason?: string;
        stateReasonCleared?: boolean;
    };
    expired?: boolean;
}>();

const { renderProperty } = usePropertyRenderer();

function render(value: PropertyValue | undefined): string {
    if (!value) return '';
    return renderProperty({ key: props.propertyKey, value });
}

const isClear = computed(() => props.oldValue !== undefined && props.newValue === undefined);
const isCreate = computed(() => props.oldValue === undefined && props.newValue !== undefined);
const isStateOnly = computed(
    () =>
        props.lifecycle !== undefined &&
        props.oldValue !== undefined &&
        props.newValue !== undefined &&
        render(props.oldValue) === render(props.newValue),
);
</script>

<template>
    <span>
        <I18nT
            :keypath="
                expired
                    ? 'components.properties.activity.expired'
                    : isStateOnly
                      ? 'components.properties.activity.lifecycle_updated'
                      : isClear
                        ? 'components.properties.activity.cleared'
                        : isCreate
                          ? 'components.properties.activity.set'
                          : 'components.properties.activity.changed'
            "
        >
            <template #property
                ><code class="font-semibold">{{ propertyKey }}</code></template
            >
            <template #old
                ><span class="font-semibold line-through">{{ render(oldValue) }}</span></template
            >
            <template #new
                ><span class="font-semibold">{{ render(newValue) }}</span></template
            >
        </I18nT>
        <template v-if="lifecycle?.expiresAt">
            — {{ $t('common.expiration') }}: <GenericTime :value="lifecycle.expiresAt" type="short" />
        </template>
        <template v-else-if="lifecycle?.expiryCleared"> — {{ $t('components.properties.activity.expiry_cleared') }} </template>
        <template v-if="lifecycle?.stateReason"> — {{ $t('common.reason') }}: {{ lifecycle.stateReason }} </template>
        <template v-else-if="lifecycle?.stateReasonCleared">
            — {{ $t('components.properties.activity.reason_cleared') }}
        </template>
    </span>
</template>
