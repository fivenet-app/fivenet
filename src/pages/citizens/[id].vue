<script setup lang="ts">
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import CitizenInfo from '../../components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { getCitizenStoreClient } from '../../grpc/grpc';
import { ref, Ref, onBeforeMount } from 'vue';
import { useRoute } from 'vue-router/auto';

const user = ref() as Ref<undefined | User>;

const route = useRoute();

onBeforeMount(() => {
    const req = new GetUserRequest();
    req.setUserId(route.params.id);

    getCitizenStoreClient()
        .getUser(req, null)
        .then((resp) => {
            user.value = resp.getUser();
        });
});
</script>

<route lang="json">
{
    "name": "Citizens Info",
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
