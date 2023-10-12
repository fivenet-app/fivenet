<script lang="ts" setup>
import { useConfigStore } from '~/store/config';
import FiveNetLogo from './logos/FiveNetLogo.vue';

const configStore = useConfigStore();
const { appConfig } = storeToRefs(configStore);
</script>

<template>
    <footer class="bg-base-800 w-full">
        <div class="py-6 sm:py-8 flex flex-col items-center mx-auto w-fit gap-5 text-neutral">
            <div class="flex flex-row items-center gap-4">
                <FiveNetLogo class="w-16 h-auto" />
                <p class="text-center text-sm leading-5">
                    {{ $t('copyright') }}
                </p>
            </div>
            <p class="inline-flex flex-row underline gap-2">
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
