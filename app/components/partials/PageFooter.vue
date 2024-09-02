<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const { t } = useI18n();

const { website } = useAppConfig();

const footerLinks = [
    {
        label: t('common.privacy_policy'),
        to: website.links?.privacyPolicy,
    },
    {
        label: t('common.imprint'),
        to: website.links?.imprint,
    },
].filter((l) => l.to !== undefined);

const year = new Date().getFullYear();
</script>

<template>
    <UFooter :links="footerLinks">
        <template #left>
            <FiveNetLogo class="mr-1 h-auto w-8" />

            <I18nT keypath="copyright">
                <template #year>
                    {{ year }}
                </template>
            </I18nT>
        </template>

        <template #right>
            <UButton
                icon="i-simple-icons-github"
                color="gray"
                variant="ghost"
                to="https://github.com/fivenet-app/fivenet"
                target="_blank"
            />
            <UButton
                icon="i-simple-icons-discord"
                color="gray"
                variant="ghost"
                to="https://discord.gg/sWvkHuVQA5"
                target="_blank"
            />
        </template>
    </UFooter>
</template>
