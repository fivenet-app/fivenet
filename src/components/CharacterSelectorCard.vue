<script lang="ts">
import authInterceptor from '../grpcauth';
import { AccountServiceClient } from '@arpanet/gen/auth/AuthServiceClientPb';
import { ChooseCharacterRequest } from '@arpanet/gen/auth/auth_pb';
import { defineComponent } from 'vue';
import { mapActions, mapState } from 'vuex';
import * as grpcWeb from 'grpc-web';
import config from '../config';

export default defineComponent({
    data() {
        return {
            'client': new AccountServiceClient(config.apiProtoURL, null, {
                unaryInterceptors: [authInterceptor],
                streamInterceptors: [authInterceptor],
            }),
        };
    },
    computed: {
        ...mapState([
            'activeCharIdentifier',
        ]),
    },
    props: [
        'char',
        'identifier',
    ],
    methods: {
        ...mapActions([
            'updateAccessToken',
            'updateActiveChar',
            'updateActiveCharIdentifier',
        ]),
        chooseCharacter() {
            const req = new ChooseCharacterRequest();
            req.setIdentifier(this.identifier);
            this.client.chooseCharacter(req, null).then((resp) => {
                this.updateAccessToken(resp.getToken());
                this.updateActiveChar(this.char);
                this.updateActiveCharIdentifier(this.identifier);

                const path = this.$route.query.redirect?.toString() || '/overview';
                this.$router.push({ path: path, query: {} });
            }).catch((err: grpcWeb.RpcError) => {
                console.log(err);
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
    <div :key="char.getIdentifier()" class="grid flex-grow card w-96 bg-base-100 shadow-xl rounded-box place-items-center">
        <div class="card-body items-center text-center">
            <h2 class="card-title text-xl">
                {{ char.getFirstname() }}, {{ char.getLastname() }}
                <div class="badge badge-accent">
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
                </div>
                <div v-if="activeCharIdentifier == char.getIdentifier()" class="badge badge-primary badge-outline">Last Used
                </div>
            </h2>
            <div class="badge">Job: {{ char.getJob() }} (Rank: {{ char.getJobgrade() }})</div>
            <div class="stats">
                <div class="stat place-items-center">
                    <div class="stat-figure">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                            class="bi bi-person-lines-fill inline-block w-8 h-8 stroke-current" viewBox="0 0 16 16">
                            <path
                                d="M6 8a3 3 0 1 0 0-6 3 3 0 0 0 0 6zm-5 6s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1H1zM11 3.5a.5.5 0 0 1 .5-.5h4a.5.5 0 0 1 0 1h-4a.5.5 0 0 1-.5-.5zm.5 2.5a.5.5 0 0 0 0 1h4a.5.5 0 0 0 0-1h-4zm2 3a.5.5 0 0 0 0 1h2a.5.5 0 0 0 0-1h-2zm0 3a.5.5 0 0 0 0 1h2a.5.5 0 0 0 0-1h-2z" />
                        </svg>
                    </div>
                    <div class="stat-title">Height</div>
                    <div class="stat-value">{{ char.getHeight() }}</div>
                </div>
                <div class="stat place-items-center">
                    <div class="stat-figure">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                            class="bi bi-balloon-lines-fill inline-block w-8 h-8 stroke-current" viewBox="0 0 16 16">
                            <path fill-rule="evenodd"
                                d="M8.48 10.901C11.211 10.227 13 7.837 13 5A5 5 0 0 0 3 5c0 2.837 1.789 5.227 4.52 5.901l-.244.487a.25.25 0 1 0 .448.224l.04-.08c.009.17.024.315.051.45.068.344.208.622.448 1.102l.013.028c.212.422.182.85.05 1.246-.135.402-.366.751-.534 1.003a.25.25 0 0 0 .416.278l.004-.007c.166-.248.431-.646.588-1.115.16-.479.212-1.051-.076-1.629-.258-.515-.365-.732-.419-1.004a2.376 2.376 0 0 1-.037-.289l.008.017a.25.25 0 1 0 .448-.224l-.244-.487ZM4.352 3.356a4.004 4.004 0 0 1 3.15-2.325C7.774.997 8 1.224 8 1.5c0 .276-.226.496-.498.542-.95.162-1.749.78-2.173 1.617a.595.595 0 0 1-.52.341c-.346 0-.599-.329-.457-.644Z" />
                        </svg>
                    </div>
                    <div class="stat-title">Birthdate</div>
                    <div class="stat-value">{{ char.getDateofbirth() }}</div>
                </div>
            </div>
            <div class="stats">
                <div class="stat place-items-center">
                    <div class="stat-figure">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                            class="bi bi-pass-fill inline-block w-8 h-8 stroke-current" viewBox="0 0 16 16">
                            <path
                                d="M10 0a2 2 0 1 1-4 0H3.5A1.5 1.5 0 0 0 2 1.5v13A1.5 1.5 0 0 0 3.5 16h9a1.5 1.5 0 0 0 1.5-1.5v-13A1.5 1.5 0 0 0 12.5 0H10ZM4.5 5h7a.5.5 0 0 1 0 1h-7a.5.5 0 0 1 0-1Zm0 2h4a.5.5 0 0 1 0 1h-4a.5.5 0 0 1 0-1Z" />
                        </svg>
                    </div>
                    <div class="stat-title">Visum</div>
                    <div class="stat-value">{{ char.getVisum() }}</div>
                </div>
                <div class="stat place-items-center">
                    <div class="stat-figure">
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                            class="bi bi-calendar-fill inline-block w-8 h-8 stroke-current" viewBox="0 0 16 16">
                            <path
                                d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V5h16V4H0V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5z" />
                        </svg>
                    </div>
                    <div class="stat-title">Playtime</div>
                    <div class="stat-value">{{ getTimeInHoursAndMins(char.getPlaytime()) }}</div>
                </div>
            </div>
            <div class="card-actions">
                <button @click="chooseCharacter()" class="btn btn-primary">Choose</button>
            </div>
        </div>
    </div>
</template>
