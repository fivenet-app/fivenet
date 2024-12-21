<script lang="ts" setup>
import SettingsModal from '~/components/livemap/controls/SettingsModal.vue';

const modal = useModal();
</script>

<template>
    <UTooltip :text="$t('common.setting', 2)">
        <UButton
            class="inset-0 inline-flex items-center justify-center rounded-md border border-black/20 bg-clip-padding text-black hover:bg-[#f4f4f4]"
            icon="i-mdi-cog"
            size="2xs"
            @click="modal.open(SettingsModal, {})"
        >
            {{ $t('common.setting', 2) }}
        </UButton>
    </UTooltip>
</template>
