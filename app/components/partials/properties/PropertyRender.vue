<script lang="ts" setup>
import { PropertyRenderer, type Property } from '~~/gen/ts/resources/properties/properties';
import CharSexBadge from '~/components/partials/citizens/CharSexBadge.vue';
import PhoneNumberBlock from '~/components/partials/citizens/PhoneNumberBlock.vue';
import EmailInfoPopover from '~/components/mailer/EmailInfoPopover.vue';

const props = defineProps<{ property: Property }>();
const { propertyScalarValue, renderProperty } = usePropertyRenderer();

const scalarValue = computed(() => propertyScalarValue(props.property));
const isEmail = computed(
    () => props.property.definition?.renderer === PropertyRenderer.EMAIL || props.property.key === 'user.email',
);
const isPhone = computed(
    () => props.property.definition?.renderer === PropertyRenderer.PHONE || props.property.key === 'user.phone_number',
);
const isSex = computed(() => props.property.definition?.renderer === PropertyRenderer.SEX || props.property.key === 'user.sex');
</script>

<template>
    <EmailInfoPopover v-if="isEmail && typeof scalarValue === 'string' && scalarValue.length > 0" :email="scalarValue" />
    <PhoneNumberBlock v-else-if="isPhone && typeof scalarValue === 'string' && scalarValue.length > 0" :number="scalarValue" />
    <span v-else-if="isSex" class="inline-flex items-center gap-2">
        <span>{{ renderProperty(property) }}</span>
        <CharSexBadge :sex="typeof scalarValue === 'string' ? scalarValue : ''" />
    </span>
    <span v-else>{{ renderProperty(property) }}</span>
</template>
