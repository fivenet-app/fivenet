import { Module, Store } from 'vuex';
import { RootState } from '../store';
import { TemplateData } from '@arpanet/gen/resources/documents/templates/templates_pb';
import { User } from '@arpanet/gen/resources/users/users_pb';

export interface ClipboardModuleState {
    usersList: Array<ClipboardUser>;
}

export class ClipboardUser {
    private id: number;
    private data: {
        sex: string;
        dateofbirth: string;
        job: string;
        jobLabel: string;
        jobGrade: number;
        jobGradeLabel: string;
        firstname: string;
        lastname: string;
    };

    constructor(user: User) {
        this.id = user.getUserId();
        this.data = {
            sex: user.getSex(),
            dateofbirth: user.getDateofbirth(),
            job: user.getJob(),
            jobLabel: user.getJobLabel(),
            jobGrade: user.getJobGrade(),
            jobGradeLabel: user.getJobGradeLabel(),
            firstname: user.getFirstname(),
            lastname: user.getLastname(),
        };
    }

    getId(): number {
        return this.id;
    }

    getUser(): User {
        const user = new User();
        user.setUserId(this.id);
        user.setSex(this.data.sex);
        user.setDateofbirth(this.data.dateofbirth);
        user.setJob(this.data.job);
        user.setJobLabel(this.data.jobLabel);
        user.setJobGrade(this.data.jobGrade);
        user.setJobGradeLabel(this.data.jobGradeLabel);
        user.setFirstname(this.data.firstname);
        user.setLastname(this.data.lastname);

        return user;
    }
}

const clipboardModule: Module<ClipboardModuleState, RootState> = {
    namespaced: true,
    state: {
        usersList: new Array<ClipboardUser>(),
    },
    actions: {
        addUser({ commit }, user: User) {
            commit('addUser', user);
        },
        removeUser({ commit }, id: number) {
            commit('removeUser', id);
        },
        clearUsers({ commit }) {
            commit('clearUsers');
        },
    },
    mutations: {
        addUser(state: ClipboardModuleState, user: User): void {
            const charIdx = state.usersList.findIndex((o: ClipboardUser) => {
                return o.getId() === user.getUserId();
            });
            if (charIdx === -1) {
                state.usersList.push(new ClipboardUser(user));
            }
        },
        removeUser(state: ClipboardModuleState, id: number): void {
            const charIdx = state.usersList.findIndex((o: ClipboardUser) => {
                return o.getId() === id;
            });
            if (charIdx > 0) {
                state.usersList.splice(charIdx, 1);
            }
        },
        clearUsers(state: ClipboardModuleState) {
            state.usersList.splice(0, state.usersList.length);
        },
    },
    getters: {
        getTemplateData(state: ClipboardModuleState, getters): TemplateData {
            const data = new TemplateData();
            if (state.usersList) {
                const usersList = new Array<User>();
                state.usersList.forEach((v: ClipboardUser) => {
                    usersList.push(v.getUser());
                });
                data.setUsersList(usersList);
            }
            console.log('GETTING TEMPLATE DATA 3');
            return data;
        },
    },
};

export default clipboardModule;
