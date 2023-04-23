<script lang="ts" setup>
import ContentWrapper from '~/components/partials/ContentWrapper.vue';
import CitizenInfo from '~/components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@fivenet/gen/services/citizenstore/citizenstore_pb';
import { RpcError } from 'grpc-web';
import ClipboardButton from '~/components/clipboard/ClipboardButton.vue';
import { User } from '@fivenet/gen/resources/users/users_pb';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';
import DataPendingBlock from '~/components/partials/DataPendingBlock.vue';
import DataErrorBlock from '~/components/partials/DataErrorBlock.vue';

useHead({
    title: 'Citizen File',
});
definePageMeta({
    title: 'Citizen File',
    requiresAuth: true,
    permission: 'CitizenStoreService.FindUsers',
    validate: async (route) => {
        route = route as TypedRouteFromName<'citizens-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const { $grpc } = useNuxtApp();
const route = useRoute('citizens-id');

const { data: user, pending, refresh, error } = await useLazyAsyncData(`citizen-${route.params.id}`, () => getUser());

async function getUser(): Promise<User> {
    return new Promise(async (res, rej) => {
        const req = new GetUserRequest();
        req.setUserId(parseInt(route.params.id));

        try {
            const resp = await $grpc.getCitizenStoreClient().
                getUser(req, null);

            return res(resp.getUser()!);
        } catch (e) {
            $grpc.handleRPCError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}
</script>

<template>
    <ContentWrapper>
        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.citizen', 1)])" />
        <DataErrorBlock v-else-if="error" :title="$t('common.unable_to_load', [$t('common.citizen', 1)])" :retry="refresh" />
        <div v-else>
            <CitizenInfo :user="user!" />
            <ClipboardButton />
        </div>
    </ContentWrapper>
</template>
