<script lang="ts">
import { defineComponent } from 'vue';

import Navbar from '../../components/partials/Navbar.vue';
import Footer from '../../components/partials/Footer.vue';
import ContentWrapper from '../../components/partials/ContentWrapper.vue';
import NavPageHeader from '../../components/partials/NavPageHeader.vue';
import CitizenInfo from '../../components/citizens/CitizenInfo.vue';
import { GetUserRequest } from '@arpanet/gen/services/citizenstore/citizenstore_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';
import { RpcError } from 'grpc-web';
import { getCitizenStoreClient, handleGRPCError } from '../../grpc';

export default defineComponent({
    components: {
        Navbar,
        Footer,
        ContentWrapper,
        NavPageHeader,
        CitizenInfo,
    },
    data() {
        return {
            user: undefined as undefined | User,
        };
    },
    beforeMount() {
        const req = new GetUserRequest();
        req.setUserId(this.$route.params.id);

        getCitizenStoreClient()
            .getUser(req, null)
            .then((resp) => {
                this.user = resp.getUser();
            })
            .catch((err: RpcError) => {
                handleGRPCError(err, this.$route);
            });
    },
});
</script>

<route lang="json">
{
    "name": "Citizens Info",
    "meta": {
        "requiresAuth": true,
        "permission": "citizenstoreservice-findusers",
        "breadCrumbs": [
            { "name": "Citizens", "href": "/citizens" },
            { "name": "Citizen Info: ...", "href": "#" }
        ]
    }
}
</route>

<template>
    <Navbar />
    <NavPageHeader title="Citizens" />
    <ContentWrapper>
        <div v-if="user">
            <CitizenInfo :user="user" />
        </div>
    </ContentWrapper>
    <Footer />
</template>
