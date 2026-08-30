<script setup lang="ts">
const leftTrim = defineModel<boolean>('leftTrim', { default: false });
const rightTrim = defineModel<boolean>('rightTrim', { default: false });
</script>

<template>
    <div class="flex flex-row gap-2">
        <UFormField
            class="justify-center"
            name="leftTrim"
            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_left')"
        >
            <USwitch v-model="leftTrim" />
        </UFormField>

        <UFormField
            class="justify-center"
            name="rightTrim"
            :label="$t('components.partials.tiptap_editor.extensions.template_var.trim_right')"
        >
            <USwitch v-model="rightTrim" />
        </UFormField>
    </div>
</template>
