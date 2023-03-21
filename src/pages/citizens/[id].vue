<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import CitizenInfo from '../../components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';

const { data: user } = useUserData();
</script>

<script lang="ts">
import { defineLoader } from 'vue-router/auto';

export const useUserData = defineLoader(async (route) => {
    const req = new GetUserRequest();
    req.setUserId(route.params.id);

    let user: User;
    await getCitizenStoreClient()
        .getUser(req, null)
        .then((resp) => {
            if (resp.hasUser()) {
                user = resp.getUser();
            }
        });

    return user;
});
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
    <NavPageHeader title="Citizens" />
    <ContentWrapper>
        <div v-if="user">
            <CitizenInfo :user="user" />
        </div>
    </ContentWrapper>
    <Footer />
</template>
