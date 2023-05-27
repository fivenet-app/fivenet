<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import CitizenInfo from '~/components/citizens/CitizenInfo.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import { User } from '~~/gen/ts/resources/users/users';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';
import { RpcError } from 'grpc-web';

useHead({
    title: 'pages.citizens.id.title',
});
definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'CitizenStoreService.ListCitizens',
    validate: async (route) => {
        route = route as TypedRouteFromName<'citizens-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const { $grpc } = useNuxtApp();
const route = useRoute('citizens-id');

const { data: user, pending, refresh, error } = useLazyAsyncData(`citizen-${route.params.id}`, () => getUser());

async function getUser(): Promise<User> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getCitizenStoreClient().
                getUser({
                    userId: parseInt(route.params.id),
                });
            const { response } = await call;

            return res(response.user!);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <ContentWrapper>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 1)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
            :retry="refresh" />
        <div v-else>
            <CitizenInfo :user="user!" />
            <ClipboardButton />
        </div>
    </ContentWrapper>
</template>
