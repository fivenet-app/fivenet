<script lang="ts" setup>
import type { Colleague } from '~~/gen/ts/resources/jobs/colleagues';

defineProps<{
    colleague: Colleague;
    birthday?: boolean;
}>();
</script>

<template>
    <span>
        {{ colleague.props?.namePrefix ? colleague.props?.namePrefix + ' ' : '' }}{{ colleague.firstname }}
        {{ colleague.lastname }} {{ colleague.props?.nameSuffix ? colleague.props?.nameSuffix + ' ' : '' }}
        <template v-if="birthday"> ({{ colleague.dateofbirth }}) </template>
    </span>
</template>
