<script lang="ts" setup>
import { Kind, type ContentDiffOp } from '~~/gen/ts/resources/common/content/diff_activity';

const props = defineProps<{
    ops: ContentDiffOp[];
}>();
</script>

<template>
    <div class="break-words whitespace-pre-wrap">
        <template v-for="(op, i) in props.ops" :key="i">
            <span v-if="op.kind === Kind.EQUAL">{{ op.text }}</span>

            <span v-else-if="op.kind === Kind.INSERT" class="rounded bg-green-200/30 px-0.5">
                {{ op.text }}
            </span>

            <span v-else class="rounded bg-red-200/30 px-0.5 text-red-700 line-through">
                {{ op.text }}
            </span>
        </template>
    </div>
</template>
