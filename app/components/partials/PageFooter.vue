<script lang="ts" setup>
defineProps<{
    disableGradient?: boolean;
}>();

const { t } = useI18n();

const { website } = useAppConfig();

const settingsStore = useSettingsStore();
const { eventsDisabled } = storeToRefs(settingsStore);

const items = computed(() =>
    [
        {
            label: t('common.privacy_policy'),
            to: website.links?.privacyPolicy,
        },
        {
            label: t('common.imprint'),
            to: website.links?.imprint,
        },
        {
            label: t('components.partials.footer.toggle_event_effect'),
            onClick: () => (eventsDisabled.value = !eventsDisabled.value),
        },
        {
            label: t('pages.about.title'),
            to: '/about',
        },
    ].filter((l) => l.to !== undefined || l.onClick !== undefined),
);

const year = new Date().getFullYear();
</script>

<template>
    <div class="relative pt-10">
        <div
            v-if="!disableGradient"
            class="pointer-events-none absolute inset-x-0 top-0 h-10"
            :style="{ background: 'linear-gradient(to bottom, transparent, var(--ui-bg))' }"
        />

        <USeparator class="h-px" />

        <UFooter class="bg-default/75">
            <template #left>
                <p class="text-sm text-muted">
                    {{ $t('copyright', { year: year }) }}
                </p>
            </template>

            <UNavigationMenu :items="items" variant="link" />

            <template #right>
                <UButton
                    icon="i-simple-icons-github"
                    color="neutral"
                    variant="ghost"
                    to="https://github.com/fivenet-app/fivenet"
                    target="_blank"
                />

                <UButton
                    icon="i-simple-icons-discord"
                    color="neutral"
                    variant="ghost"
                    to="https://discord.gg/ASRPPr8CeT"
                    target="_blank"
                />
            </template>
        </UFooter>
    </div>
</template>
