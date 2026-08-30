<script setup lang="ts">
import TemplateTrimControls from './TemplateTrimControls.vue';

const props = withDefaults(defineProps<{ leftTrim?: boolean; rightTrim?: boolean }>(), {
    leftTrim: false,
    rightTrim: false,
});

const emit = defineEmits<{
    (e: 'submit', value: { value: string; leftTrim: boolean; rightTrim: boolean }): void;
    (e: 'delete'): void;
}>();

const leftTrim = ref(props.leftTrim);
const rightTrim = ref(props.rightTrim);

function submit(): void {
    emit('submit', { value: 'end', leftTrim: leftTrim.value, rightTrim: rightTrim.value });
}
</script>

<template>
    <div class="flex flex-col gap-3">
        <TemplateTrimControls v-model:left-trim="leftTrim" v-model:right-trim="rightTrim" />

        <UFieldGroup>
            <UButton block :label="$t('common.save')" @click="submit" />
            <UTooltip :text="$t('common.delete')">
                <UButton color="error" icon="i-mdi-trash-can" @click="$emit('delete')" />
            </UTooltip>
        </UFieldGroup>
    </div>
</template>
