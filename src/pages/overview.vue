<script lang="ts" setup>
import ContentCenterWrapper from '~/components/partials/ContentCenterWrapper.vue';
import { BriefcaseIcon, DocumentTextIcon, UsersIcon, MapIcon, TruckIcon } from '@heroicons/vue/24/outline';
import Cards from '~/components/partials/Cards.vue';
import { CardElements } from '~/utils/types';
import { CheckTokenRequest } from '@fivenet/gen/services/auth/auth_pb';
import { RpcError } from 'grpc-web';

useHead({
    title: 'common.overview',
});
definePageMeta({
    title: 'common.overview',
    requiresAuth: true,
});

const { t } = useI18n();

const features = [
    {
        title: `${t('common.citizen', 2)} ${t('common.search')}`,
        description: t('pages.overview.features.citizens'),
        href: { name: 'citizens' },
        permission: 'CitizenStoreService.FindUsers',
        icon: UsersIcon,
        iconForeground: 'text-purple-900',
        iconBackground: 'bg-purple-50',
    },
    {
        title: t('common.vehicle', 2),
        description: t('pages.overview.features.vehicles'),
        href: { name: 'vehicles' },
        permission: 'DMVService.FindVehicles',
        icon: TruckIcon,
        iconForeground: 'text-zinc-900',
        iconBackground: 'bg-zinc-50',
    },
    {
        title: t('common.document', 2),
        description: t('pages.overview.features.documents'),
        href: { name: 'documents' },
        permission: 'DocStoreService.FindDocuments',
        icon: DocumentTextIcon,
        iconForeground: 'text-sky-900',
        iconBackground: 'bg-sky-50',
    },
    {
        title: t('common.job', 2),
        description: t('pages.overview.features.jobs'),
        href: { name: 'jobs' },
        permission: 'Jobs.View',
        icon: BriefcaseIcon,
        iconForeground: 'text-yellow-900',
        iconBackground: 'bg-yellow-50',
    },
    {
        title: t('common.livemap'),
        description: t('pages.overview.features.livemap'),
        href: { name: 'livemap' },
        permission: 'LivemapperService.Stream',
        icon: MapIcon,
        iconForeground: 'text-teal-900',
        iconBackground: 'bg-teal-50',
    },
] as CardElements;

const { $grpc } = useNuxtApp();

async function checkToken(): Promise<void> {
    return new Promise(async (res, rej) => {
        const req = new CheckTokenRequest();

        try {
            const resp = await $grpc.getAuthClient().
                checkToken(req, null);

            console.log(toDateLocaleString(resp.getExpires()));

            return res();
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <ContentCenterWrapper>
        <button type="button" @click="checkToken"
            class="flex justify-center w-full px-3 py-2 text-sm font-semibold transition-colors rounded-md bg-secondary-600 text-neutral hover:bg-secondary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-base-300">
            CHECK TOKEN
        </button>
        <Cards :items="features" />
    </ContentCenterWrapper>
</template>
