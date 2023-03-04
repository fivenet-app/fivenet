<script lang="ts">
import { clientAuthOptions, handleGRPCError } from '../../grpc';
import { ChooseCharacterRequest } from '@arpanet/gen/auth/auth_pb';
import { defineComponent } from 'vue';
import { mapActions, mapState } from 'vuex';
import * as grpcWeb from 'grpc-web';
import { Character } from '@arpanet/gen/common/character_pb';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import config from '../../config';

export default defineComponent({
    data() {
        return {
            'client': new AccountServiceClient(config.apiProtoURL, null, clientAuthOptions),
        };
    },
    computed: {
        ...mapState([
            'activeCharIdentifier',
        ]),
    },
    props: {
        'char': {
            required: true,
            type: Character,
        },
    },
    methods: {
        ...mapActions([
            'updateAccessToken',
            'updateActiveChar',
            'updateActiveCharIdentifier',
            'updatePermissions',
        ]),
        chooseCharacter() {
            const req = new ChooseCharacterRequest();
            req.setIdentifier(this.char.getIdentifier());

            this.client.chooseCharacter(req, null).then((resp) => {
                this.updateAccessToken(resp.getToken());
                this.updateActiveChar(this.char);
                this.updateActiveCharIdentifier(this.char.getIdentifier());
                this.updatePermissions(resp.getPermissionsList());
                console.log(resp.getPermissionsList());

                const path = this.$route.query.redirect?.toString() || '/overview';
                this.$router.push({ path: path, query: {} });
            }).catch((err: grpcWeb.RpcError) => {
                handleGRPCError(err, this.$route);
            });
        },
        getTimeInHoursAndMins(timeInsSeconds: number): string {
            if (timeInsSeconds && timeInsSeconds > 0) {
                const minsTemp = timeInsSeconds / 60;
                let hours = Math.floor(minsTemp / 60);
                const mins = minsTemp % 60;
                const hoursText = 'hrs';
                const minsText = 'mins';

                if (hours !== 0 && mins !== 0) {
                    if (mins >= 59) {
                        hours += 1;
                        return `${hours} ${hoursText} `;
                    } else {
                        return `${hours} ${hoursText} ${mins.toFixed(0)} ${minsText}`;
                    }
                } else if (hours === 0 && mins !== 0) {
                    return `${mins.toFixed(0)} ${minsText}`;
                } else if (hours !== 0 && mins === 0) {
                    return `${hours} ${hoursText}`;
                }
            }
            return '-';
        },
    },
});
</script>

<template>
    <li :key="char.getIdentifier()"
        class="col-span-2 flex flex-col divide-y divide-gray-200 rounded-lg bg-white text-center shadow">
        <div class="flex flex-1 flex-col p-8">
            <h2 class="mt-6 text-2xl font-medium text-gray-900">
                {{ char.getFirstname() }}, {{ char.getLastname() }}
                <span
                    :class="[char.getSex() == 'f' ? 'bg-purple-100' : 'bg-blue-100', 'inline-flex items-center rounded-md px-2.5 py-0.5 text-sm font-medium text-yellow-800']">
                    <svg v-if="char.getSex() == 'f'" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                        fill="currentColor" class="bi bi-gender-female" viewBox="0 0 16 16">
                        <path fill-rule="evenodd"
                            d="M8 1a4 4 0 1 0 0 8 4 4 0 0 0 0-8zM3 5a5 5 0 1 1 5.5 4.975V12h2a.5.5 0 0 1 0 1h-2v2.5a.5.5 0 0 1-1 0V13h-2a.5.5 0 0 1 0-1h2V9.975A5 5 0 0 1 3 5z" />
                    </svg>
                    <svg v-else-if="char.getSex() == 'm'" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
                        fill="currentColor" class="bi bi-gender-male" viewBox="0 0 16 16">
                        <path fill-rule="evenodd"
                            d="M9.5 2a.5.5 0 0 1 0-1h5a.5.5 0 0 1 .5.5v5a.5.5 0 0 1-1 0V2.707L9.871 6.836a5 5 0 1 1-.707-.707L13.293 2H9.5zM6 6a4 4 0 1 0 0 8 4 4 0 0 0 0-8z" />
                    </svg>
                </span>
            </h2>
            <dl class="mt-1 flex flex-grow flex-col justify-between">
                <dd>
                    <span v-if="activeCharIdentifier == char.getIdentifier()"
                        class="inline-flex items-center rounded-full bg-green-100 px-3 py-0.5 text-sm font-medium text-green-800">
                        Last Used
                    </span>
                    <br v-else />
                </dd>
                <dd class="mt-3">
                    <span
                        class="inline-flex items-center rounded-md bg-gray-100 px-2.5 py-0.5 text-sm font-medium text-gray-800">{{
                            char.getJob() }} (Rank: {{ char.getJobgrade() }})</span>
                </dd>
                <dt class="text-sm text-black">Date of Birth</dt>
                <dd class="text-sm text-gray-500">{{ char.getDateofbirth() }}</dd>
                <dt class="text-sm text-black">Height</dt>
                <dd class="text-sm text-gray-500">{{ char.getHeight() }}cm</dd>
                <dt class="text-sm text-black">Visum</dt>
                <dd class="text-sm text-gray-500">{{ char.getVisum() }}</dd>
                <dt class="text-sm text-black">Playtime</dt>
                <dd class="text-sm text-gray-500">{{ getTimeInHoursAndMins(char.getPlaytime()) }}</dd>
            </dl>
        </div>
        <div>
            <div class="-mt-px flex divide-x divide-gray-200">
                <div class="flex w-0 flex-1">
                    <button @click="chooseCharacter()"
                        class="relative -mr-px inline-flex w-0 flex-1 items-center justify-center gap-x-3 rounded-bl-lg border border-transparent py-4 text-sm font-semibold text-gray-900">
                        Choose
                    </button>
                </div>
            </div>
        </div>
    </li>
</template>
