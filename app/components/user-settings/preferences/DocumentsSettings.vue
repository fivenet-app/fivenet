<script lang="ts" setup>
import SignatureSettings from '~/components/user-settings/preferences/SignatureSettings.vue';

const settingsStore = useSettingsStore();
const { design } = storeToRefs(settingsStore);

const designDocumentListStyle = ref(design.value.documents.listStyle === 'double');

watch(designDocumentListStyle, (value) => {
    design.value.documents.listStyle = value ? 'double' : 'single';
});
</script>

<template>
    <div class="space-y-4">
        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="designDocumentListStyle"
            :label="$t('components.auth.user_settings.documents_lists_style.title')"
            :description="$t('components.auth.user_settings.documents_lists_style.description')"
        >
            <div class="flex items-center justify-end gap-2 text-sm">
                <span>{{ $t('components.auth.user_settings.documents_lists_style.single') }}</span>
                <USwitch v-model="designDocumentListStyle" />
                <span>{{ $t('components.auth.user_settings.documents_lists_style.double') }}</span>
            </div>
        </UFormField>

        <UFormField
            class="flex items-center gap-4"
            name="designDocumentViewCollapsedTitle"
            :label="$t('components.auth.user_settings.documents_view_collapsed_title.title')"
            :description="$t('components.auth.user_settings.documents_view_collapsed_title.description')"
            :ui="{ wrapper: 'min-w-0 flex-1', container: 'flex shrink-0 items-center justify-end' }"
        >
            <USwitch v-model="design.documents.viewCollapsedTitle" />
        </UFormField>

        <USeparator :label="$t('components.auth.user_settings.approval_signature.title')" />

        <p class="text-sm text-muted">
            {{ $t('components.auth.user_settings.approval_signature.description') }}
        </p>

        <SignatureSettings />
    </div>
</template>
