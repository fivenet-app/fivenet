<script lang="ts" setup>
// TODO
</script>

<template>
    <UPage>
        <UHeader :links="[]">
            <template #logo> Nuxt UI Pro <UBadge label="SaaS" variant="subtle" class="mb-0.5" /> </template>

            <template #right>
                <UButton label="Sign in" color="gray" to="/login" />
                <UButton
                    label="Sign up"
                    icon="i-heroicons-arrow-right-20-solid"
                    trailing
                    color="black"
                    to="/signup"
                    class="hidden lg:flex"
                />
            </template>
        </UHeader>

        <UMain>
            <slot />
        </UMain>

        <UFooter>
            <template #top>
                <UFooterColumns :links="[]"> </UFooterColumns>
            </template>

            <template #left>
                <p class="text-sm text-gray-500 dark:text-gray-400">
                    Copyright © {{ new Date().getFullYear() }}. All rights reserved.
                </p>
            </template>

            <template #right>
                <UColorModeButton size="sm" />

                <UButton
                    to="https://github.com/nuxt-ui-pro/saas"
                    target="_blank"
                    icon="i-simple-icons-github"
                    aria-label="GitHub"
                    color="gray"
                    variant="ghost"
                />
            </template>
        </UFooter>
    </UPage>
</template>
