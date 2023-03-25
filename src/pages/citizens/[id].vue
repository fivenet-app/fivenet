<script lang="ts">
import { defineLoader, RouteLocationNormalizedLoaded } from 'vue-router/auto';
import ClipboardButton from '../../components/ClipboardButton.vue';
import { useStore } from '../../store/store';

export const useUserData = defineLoader(async (route: RouteLocationNormalizedLoaded) => {
    route = route as RouteLocationNormalizedLoaded<'Citizens: Info'>;

    const req = new GetUserRequest();
    req.setUserId(parseInt(route.params.id));

    try {
        const resp = await getCitizenStoreClient()
            .getUser(req, null);
            return resp.getUser();
    } catch (e) {
        return;
    }
});
</script>

<script lang="ts" setup>
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import CitizenInfo from '../../components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';

const store = useStore();

const { data: user } = useUserData();

function addUserToClipboard() {
    store.commit('clipboard/addUser', user.value);
    console.log("ADDING USER " + user.value?.getUserId() + " TO CLIPBOARD");
}
</script>

<route lang="json">
{
    "name": "Citizens: Info",
    "meta": {
        "requiresAuth": true,
        "permission": "CitizenStoreService.FindUsers"
    }
}
</route>

<template>
    <ContentWrapper>
        <div v-if="user">
            <CitizenInfo :user="user" />
            <ClipboardButton :callback="addUserToClipboard" />
        </div>
    </ContentWrapper>
</template>
