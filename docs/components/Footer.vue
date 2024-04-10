<script lang="ts" setup>
const year = new Date().getFullYear();
</script>

<template>
    <UFooter :links="[]">
        <template #left>
            <img src="/images/logo-200x200.png" alt="FiveNet Logo" title="FiveNet" class="mr-1 h-auto w-8" />

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
                to="https://github.com/galexrt/fivenet"
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
