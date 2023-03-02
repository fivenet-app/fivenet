<script lang="ts">
import { Character } from '@arpanet/gen/common/character_pb';
import { OrderBy } from '@arpanet/gen/common/database_pb';
import { defineComponent } from 'vue';

import authInterceptor from '../grpcauth';
import { RpcError } from 'grpc-web';
import { UsersServiceClient } from '@arpanet/gen/users/UsersServiceClientPb';
import { FindUsersRequest } from '@arpanet/gen/users/users_pb';
import TablePagination from './TablePagination.vue';

export default defineComponent({
    data() {
        return {
            client: new UsersServiceClient("https://localhost:8181", null, {
                unaryInterceptors: [authInterceptor],
                streamInterceptors: [authInterceptor],
            }),
            loading: false,
            searchFirstname: "",
            searchLastname: "",
            orderBys: [] as Array<OrderBy>,
            users: [] as Array<Character>,
            offset: 0,
            totalCount: 0,
            listEnd: 0,
        };
    },
    methods: {
        findUsers: function (offset: number) {
            if (offset < 0) {
                return;
            }
            if (this.loading)
                return;
            this.loading = true;
            const req = new FindUsersRequest();
            req.setCurrent(offset);
            req.setFirstname(this.searchFirstname);
            req.setLastname(this.searchLastname);
            req.setOrderbyList(this.orderBys);
            this.client.
                findUsers(req, null).
                then((resp) => {
                    this.users = resp.getUsersList();
                    this.totalCount = resp.getTotalcount();
                    this.offset = resp.getCurrent();
                    this.listEnd = resp.getEnd();
                    this.loading = false;
                }).catch((err: RpcError) => {
                    this.loading = false;
                    authInterceptor.handleError(err, this.$route);
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
                }
                else {
                    orderBy.setDesc(true);
                }
            }
            else {
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
            defaultOrderBy.setColumn("firstname");
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
    components: {
        TablePagination,
    },
});
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="findUsers(offset)">
                        <div class="grid grid-cols-2 gap-4">
                            <div class="form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">First Name</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="searchFirstname" v-on:keyup.enter="findUsers(offset)" type="text"
                                        name="search" id="search"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                            <div class="form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-white">Last Name</label>
                                <div class="relative mt-2 flex items-center">
                                    <input v-model="searchFirstname" v-on:keyup.enter="findUsers(offset)" type="text"
                                        name="search" id="search"
                                        class="block w-full rounded-md border-0 py-1.5 pr-14 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6" />
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-8 flow-root">
                <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <table class="min-w-full divide-y divide-gray-700">
                            <thead>
                                <tr>
                                    <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">Name
                                    </th>
                                    <th v-on:click="toggleOrderBy('job')" scope="col"
                                        class="py-3.5 px-2 text-left text-sm font-semibold text-white">Job
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Sex
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Date of
                                        Birth
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Height
                                    </th>
                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                        <span class="sr-only">Edit</span>
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="divide-y divide-gray-800">
                                <tr v-for="user in users" :key="user.getIdentifier()">
                                    <td class="whitespace-nowrap py-2 pl-4 pr-3 text-sm font-medium text-white sm:pl-0">
                                        {{ user.getFirstname() }}, {{ user.getLastname() }}
                                    </td>
                                    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
                                        {{ toTitleCase(user.getJob()) }}
                                    </td>
                                    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
                                        {{ user.getSex().toUpperCase() }}
                                    </td>
                                    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
                                        {{ user.getDateofbirth() }}
                                    </td>
                                    <td class="whitespace-nowrap px-2 py-2 text-sm text-gray-300">
                                        {{ user.getHeight() }}cm
                                    </td>
                                    <td
                                        class="relative whitespace-nowrap py-2 pl-3 pr-4 text-right text-sm font-medium sm:pr-0">
                                        <a href="#" class="text-indigo-400 hover:text-indigo-300">VIEW</a>
                                    </td>
                                </tr>
                            </tbody>
                            <thead>
                                <tr>
                                    <th v-on:click="toggleOrderBy('firstname')" scope="col"
                                        class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-white sm:pl-0">Name
                                    </th>
                                    <th v-on:click="toggleOrderBy('job')" scope="col"
                                        class="py-3.5 px-2 text-left text-sm font-semibold text-white">Job
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Sex
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Date of
                                        Birth
                                    </th>
                                    <th scope="col" class="py-3.5 px-2 text-left text-sm font-semibold text-white">Height
                                    </th>
                                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-0">
                                        <span class="sr-only">Edit</span>
                                    </th>
                                </tr>
                            </thead>
                        </table>

                        <TablePagination :start="offset" :entries="users.length" :end="listEnd" :total="totalCount"
                            :callback="findUsers" />
                    </div>
            </div>
        </div>
    </div>
</div></template>
