<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const appConfig = useAppConfig();

const year = new Date().getFullYear();
</script>

<template>
    <footer class="w-full bg-base-800">
        <div class="mx-auto flex w-fit flex-col items-center gap-4 py-4 text-neutral sm:py-6">
            <div class="flex flex-row items-center gap-4">
                <FiveNetLogo class="h-auto w-16" />
                <p class="text-center text-sm leading-5">
                    <I18nT keypath="copyright">
                        <template #year>
                            {{ year }}
                        </template>
                    </I18nT>
                </p>
            </div>
            <p class="inline-flex flex-row gap-2 underline">
                <NuxtLink v-if="appConfig.links?.privacyPolicy" :external="true" :to="appConfig.links.privacyPolicy">
                    {{ $t('common.privacy_policy') }}
                </NuxtLink>
                <NuxtLink v-if="appConfig.links?.imprint" :external="true" :to="appConfig.links.imprint">
                    {{ $t('common.imprint') }}
                </NuxtLink>
            </p>
        </div>
    </footer>
</template>
