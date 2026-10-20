<script lang="ts" setup>
const { activeChar } = useAuth();

const settingsStore = useSettingsStore();
const { startpage } = storeToRefs(settingsStore);

const { startpageItems } = useAppFeatures();

const selectedHomepage = computed<(typeof startpageItems.value)[number] | undefined>({
    get: () => startpageItems.value.find((item) => item.path === startpage.value),
    set: (item) => {
        startpage.value = item?.path ?? '/overview';
    },
});
</script>

<template>
    <UFormField
        class="grid grid-cols-2 items-center gap-2"
        name="selectedHomepage"
        :label="$t('components.auth.user_settings.set_startpage.title')"
        :description="$t('components.auth.user_settings.set_startpage.description')"
    >
        <ClientOnly v-if="activeChar">
            <USelectMenu
                v-model="selectedHomepage"
                class="w-full"
                :items="startpageItems"
                :search-input="{ placeholder: $t('common.search_field') }"
            >
                <template #leading>
                    <UIcon v-if="selectedHomepage?.icon" class="size-4" :name="selectedHomepage.icon" />
                </template>
            </USelectMenu>
        </ClientOnly>
        <p v-else class="text-sm">
            {{ $t('components.auth.user_settings.set_startpage.no_char_selected') }}
        </p>
    </UFormField>
</template>
