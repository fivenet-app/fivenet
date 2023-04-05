<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import CitizenInfo from '../../components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@fivenet/gen/services/citizenstore/citizenstore_pb';
import { RpcError } from 'grpc-web';
import ClipboardButton from '../../components/clipboard/ClipboardButton.vue';
import { User } from '@fivenet/gen/resources/users/users_pb';
import { TypedRouteFromName } from '~~/.nuxt/typed-router/__router';

useHead({
    title: 'Citizen Profile',
});
definePageMeta({
    title: 'Citizen Profile',
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

const user = ref<User>();

async function getUser(): Promise<void> {
    const req = new GetUserRequest();
    req.setUserId(parseInt(route.params.id));

    try {
        const resp = await $grpc.getCitizenStoreClient().
            getUser(req, null);

        user.value = resp.getUser();
    } catch (e) {
        $grpc.handleRPCError(e as RpcError);
        return;
    }
}

onBeforeMount(async () => getUser());
</script>

<template>
    <ContentWrapper>
        <div v-if="user">
            <CitizenInfo :user="user" />
            <ClipboardButton />
        </div>
    </ContentWrapper>
</template>
