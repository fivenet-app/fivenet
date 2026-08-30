<script setup lang="ts">
import { isTemplateBlockActionValid } from '~/composables/tiptap/extensions/TemplateBlockValidation';
import TemplateTrimControls from './TemplateTrimControls.vue';

const props = withDefaults(defineProps<{ editing?: boolean; value?: string; leftTrim?: boolean; rightTrim?: boolean }>(), {
    editing: false,
    value: '',
    leftTrim: false,
    rightTrim: false,
});

const emit = defineEmits<{
    (e: 'submit', value: { value: string; leftTrim: boolean; rightTrim: boolean }): void;
    (e: 'delete'): void;
}>();

const value = ref(props.value);
const leftTrim = ref(props.leftTrim);
const rightTrim = ref(props.rightTrim);

const isValid = computed(() => isTemplateBlockActionValid(value.value));

watch(
    () => [props.value, props.leftTrim, props.rightTrim] as const,
    ([v, l, r]) => {
        value.value = v;
        leftTrim.value = l;
        rightTrim.value = r;
    },
    { immediate: true },
);

function submit(): void {
    if (!isValid.value) return;
    emit('submit', { value: value.value.trim(), leftTrim: leftTrim.value, rightTrim: rightTrim.value });
}
</script>

<template>
    <div class="flex flex-col gap-3">
        <UTextarea v-model="value" :rows="3" autofocus />
        <TemplateTrimControls v-model:left-trim="leftTrim" v-model:right-trim="rightTrim" />

        <UFieldGroup>
            <UButton block :disabled="!isValid" :label="$t('common.save')" @click="submit" />
            <UTooltip :text="$t('common.delete')">
                <UButton color="error" icon="i-mdi-trash-can" @click="$emit('delete')" />
            </UTooltip>
        </UFieldGroup>
    </div>
</template>
