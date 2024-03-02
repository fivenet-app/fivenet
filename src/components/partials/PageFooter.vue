<script lang="ts" setup>
import FiveNetLogo from '~/components/partials/logos/FiveNetLogo.vue';

const { links } = useAppConfig();

const year = new Date().getFullYear();
</script>

<template>
    <footer class="w-full bg-base-800">
        <div class="mx-auto flex w-fit flex-col items-center gap-4 py-2 text-neutral sm:py-4">
            <div class="flex flex-row items-center gap-4">
                <FiveNetLogo class="h-auto w-10" />
                <div class="">
                    <p class="text-center text-base leading-5">
                        <I18nT keypath="copyright">
                            <template #year>
                                {{ year }}
                            </template>
                        </I18nT>
                    </p>

                    <p class="inline-flex flex-row gap-2 underline">
                        <NuxtLink v-if="links?.privacyPolicy" :external="true" :to="links.privacyPolicy">
                            {{ $t('common.privacy_policy') }}
                        </NuxtLink>
                        <NuxtLink v-if="links?.imprint" :external="true" :to="links.imprint">
                            {{ $t('common.imprint') }}
                        </NuxtLink>
                    </p>
                </div>
            </div>
        </div>
    </footer>
</template>
