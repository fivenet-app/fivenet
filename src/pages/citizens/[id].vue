<script lang="ts" setup>
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { type TypedRouteFromName } from '@typed-router';
import Info from '~/components/citizens/info/Info.vue';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import { User } from '~~/gen/ts/resources/users/users';

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
    try {
        const call = $grpc.getCitizenStoreClient().getUser({
            userId: parseInt(route.params.id),
        });
        const { response } = await call;

        return response.user!;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}
</script>

<template>
    <ContentWrapper>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 1)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.citizen', 1)])"
            :message="$t(error.message)"
            :retry="refresh"
        />
        <div v-else>
            <Info :user="user!" />
            <ClipboardButton />
        </div>
    </ContentWrapper>
</template>
