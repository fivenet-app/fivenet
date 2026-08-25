<script lang="ts" setup>
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';

const props = defineProps<{
    qualificationId: number;
    qualification?: QualificationShort;
}>();

const label = computed(() => {
    const abbreviation = props.qualification?.abbreviation?.trim() ?? '';
    const title = props.qualification?.title?.trim() ?? '';

    if (abbreviation && title) return `${abbreviation}: ${title}`;

    if (abbreviation) return abbreviation;

    if (title) return title;

    return `#${props.qualificationId}`;
});
</script>

<template>
    <UBadge color="neutral" size="sm" variant="soft" class="max-w-full">
        <span class="truncate">{{ label }}</span>
    </UBadge>
</template>
