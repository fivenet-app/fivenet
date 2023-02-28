<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { OrderBy } from '@arpanet/gen/common/database_pb';
import { defineComponent } from 'vue';

import authInterceptor from '../grpcauth';
import * as grpcWeb from 'grpc-web';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { FindUsersRequest } from '@arpanet/gen/users/users_pb';

export default defineComponent({
    data: function () {
        return {
            'searchFirstname': '',
            'searchLastname': '',
            'orderBys': [] as Array<OrderBy>,
            'users': [] as Array<Character>,
            'offset': 0,
            'totalCount': 0,
            'end': 0,
        };
    },
    methods: {
        findUsers: function (offset: number) {
            if (offset < 0) {
                return;
            }

            const req = new FindUsersRequest();
            req.setCurrent(offset);
            req.setFirstname(this.searchFirstname);
            req.setLastname(this.searchLastname);
            req.setOrderbyList(this.orderBys);

            client.
                findUsers(req, null).
                then((resp) => {
                    this.users = resp.getUsersList();
                    this.totalCount = resp.getTotalcount();
                    this.offset = resp.getCurrent();
                    this.end = resp.getEnd();
                }).catch((err: grpcWeb.RpcError) => {
                    console.log(err);
                });
        },
        toggleOrderBy: function (column: string) {
            // Check if the first one is the default entry, if so, remove if another column has been toggled
            if (this.orderBys.at(0)?.getColumn() != column) {
                this.orderBys.pop();
            }

            const index = this.orderBys.findIndex((o) => {
                return o.getColumn() == column;
            });

            let orderBy: OrderBy;
            if (index > -1) {
                //@ts-ignore I just checked if it exists, so it should exist
                orderBy = this.orderBys.at(index);

                if (orderBy.getDesc()) {
                    this.orderBys.splice(index);
                    if (this.orderBys.length == 0) {
                        this.orderBys.push(this.getDefaultOrderBy());
                    }
                } else {
                    orderBy.setDesc(true);
                }
            } else {
                orderBy = new OrderBy();
                orderBy.setColumn(column);
                orderBy.setDesc(false);
                this.orderBys.push(orderBy);
            }

            console.log(this.orderBys);
            this.findUsers(this.offset);
        },
        getDefaultOrderBy(): OrderBy {
            const defaultOrderBy = new OrderBy();
            defaultOrderBy.setColumn('firstname');
            defaultOrderBy.setDesc(false);
            return defaultOrderBy;
        },
        toTitleCase(value: string) {
            return value.replace(/(?:^|\s|-)\S/g, x => x.toUpperCase());
        },
    },
    mounted: function () {
        this.orderBys.push(this.getDefaultOrderBy());

        this.findUsers(this.offset);
    },
});

const client = new UsersServiceClient('https://localhost:8181', null, {
    unaryInterceptors: [authInterceptor],
    streamInterceptors: [authInterceptor],
});
</script>

<template>
    <form @submit.prevent="findUsers(offset)">
        <div class="grid grid-cols-2 gap-4">
            <div class="form-control">
                <label class="input-group input-group-vertical">
                    <span>First Name</span>
                    <input v-model="searchFirstname" v-on:keyup.enter="findUsers" type="text" placeholder="First Name"
                        class="input input-bordered" />
                </label>
            </div>
            <div class="form-control">
                <label class="input-group input-group-vertical">
                    <span>Last Name</span>
                    <input v-model="searchLastname" v-on:keyup.enter="findUsers" type="text" placeholder="Last Name"
                        class="input input-bordered" />
                </label>
            </div>
        </div>
    </form>
    <div class="overflow-x-auto">
        <table class="table table-compact w-full">
            <thead>
                <tr>
                    <th></th>
                    <th v-on:click="toggleOrderBy('firstname')">Name</th>
                    <th v-on:click="toggleOrderBy('job')">Job</th>
                    <th v-on:click="toggleOrderBy('sex')">Sex</th>
                    <th v-on:click="toggleOrderBy('dateofbirth')">Date Of Birth</th>
                    <th v-on:click="toggleOrderBy('height')">Height</th>
                    <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="(user, index) in users">
                    <th>{{ index + offset + 1 }}</th>
                    <td>{{ user.getFirstname() }}, {{ user.getLastname() }}</td>
                    <td>{{ toTitleCase(user.getJob()) }}</td>
                    <td>{{ user.getSex().toUpperCase() }}</td>
                    <td>{{ user.getDateofbirth() }}</td>
                    <td>{{ user.getHeight() }}cm</td>
                    <td><router-link
                            :to="{ name: '/citizens/[id]', params: { id: user.getIdentifier() } }">VIEW</router-link></td>
                </tr>
            </tbody>
            <tfoot>
                <tr>
                    <th></th>
                    <th>Name</th>
                    <th>Job</th>
                    <th>Sex</th>
                    <th>Date Of Birth</th>
                    <th>Height</th>
                    <th></th>
                </tr>
            </tfoot>
        </table>
    </div>

    <!-- Pagination -->
    <div class="flex flex-col items-center">
        <!-- Help text -->
        <span class="text-sm text-gray-700 dark:text-gray-400">
            Showing <span class="font-semibold text-gray-900 dark:text-white">{{ offset + 1 }}</span> to <span
                class="font-semibold text-gray-900 dark:text-white">{{ end }}</span> of <span
                class="font-semibold text-gray-900 dark:text-white">{{ totalCount }}</span> Entries
        </span>
        <div class="inline-flex mt-2 xs:mt-0">
            <!-- Buttons -->
            <button :class="[offset <= 0 ? 'disabled' : '' ]" :disabled="offset <= 0" v-on:click="findUsers(offset-users.length)"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-800 rounded-l hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                <svg aria-hidden="true" class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd"
                        d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
                        clip-rule="evenodd"></path>
                </svg>
                Prev
            </button>
            <button :class="[offset >= totalCount ? 'disabled' : '' ]" :disabled="offset >= totalCount" v-on:click="findUsers(end)"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-800 border-0 border-l border-gray-700 rounded-r hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
                Next
                <svg aria-hidden="true" class="w-5 h-5 ml-2" fill="currentColor" viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg">
                    <path fill-rule="evenodd"
                        d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
                        clip-rule="evenodd"></path>
                </svg>
            </button>
        </div>
    </div>
</template>
