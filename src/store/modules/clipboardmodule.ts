import { GetterTree, Module } from 'vuex';
import { RootState } from '../store';
import { TemplateData } from '@arpanet/gen/resources/documents/templates/templates_pb';
import { User, UserShort } from '@arpanet/gen/resources/users/users_pb';
import { Document } from '@arpanet/gen/resources/documents/documents_pb';
import { Vehicle } from '@arpanet/gen/resources/vehicles/vehicles_pb';

export interface ClipboardModuleState {
    users: ClipboardUser[];
    documents: ClipboardDocument[];
    vehicles: ClipboardVehicle[];
}

export type Getters = {
    getTemplateData(state: ClipboardModuleState): TemplateData;
};

export const getters: GetterTree<ClipboardModuleState, RootState> & Getters = {
    getTemplateData(state: ClipboardModuleState): TemplateData {
        const data = new TemplateData();

        if (state.users) {
            const users = new Array<User>();
            state.users.forEach((v: ClipboardUser) => {
                users.push(getUser(v));
            });
            data.setUsersList(users);
        }

        return data;
    },
};

const clipboardModule: Module<ClipboardModuleState, RootState> = {
    namespaced: true,
    state: {
        users: [],
        documents: [],
        vehicles: [],
    },
    actions: {
        clear({ commit }) {
            commit('clearUsers');
            commit('clearVehicles');
        },
        // Users
        addUser({ commit }, user: User) {
            commit('addUser', user);
        },
        removeUser({ commit }, id: number) {
            commit('removeUser', id);
        },
        clearUsers({ commit }) {
            commit('clearUsers');
        },
        // Documents
        addDocument({ commit }, user: Document) {
            commit('addDocument', user);
        },
        removeDocument({ commit }, id: number) {
            commit('removeDocument', id);
        },
        clearDocuments({ commit }) {
            commit('clearDocuments');
        },
        // Vehicles
        addVehicle({ commit }, user: User) {
            commit('addVehicle', user);
        },
        removeVehicle({ commit }, id: number) {
            commit('removeVehicle', id);
        },
        clearVehicles({ commit }) {
            commit('clearVehicles');
        },
    },
    mutations: {
        // Users
        addUser(state: ClipboardModuleState, user: User): void {
            const idx = state.users.findIndex((o: ClipboardUser) => {
                return o.id === user.getUserId();
            });
            if (idx === -1) {
                state.users.push((new ClipboardUser()).setUser(user));
            }
        },
        removeUser(state: ClipboardModuleState, id: number): void {
            const idx = state.users.findIndex((o: ClipboardUser) => {
                return o.id === id;
            });
            if (idx > -1) {
                state.users.splice(idx, 1);
            }
        },
        clearUsers(state: ClipboardModuleState): void {
            state.users.splice(0, state.users.length);
        },
        // Documents
        addDocument(state: ClipboardModuleState, document: Document): void {
            const idx = state.documents.findIndex((o: ClipboardDocument) => {
                return o.id === document.getId();
            });
            if (idx === -1) {
                state.documents.push(new ClipboardDocument(document));
            }
        },
        removeDocument(state: ClipboardModuleState, id: number): void {
            const idx = state.documents.findIndex((o: ClipboardDocument) => {
                return o.id === id;
            });
            if (idx > -1) {
                state.documents.splice(idx, 1);
            }
        },
        clearDocuments(state: ClipboardModuleState): void {
            state.documents.splice(0, state.documents.length);
        },
        // Vehicles
        addVehicle(state: ClipboardModuleState, vehicle: Vehicle): void {
            const idx = state.vehicles.findIndex((o: ClipboardVehicle) => {
                return o.plate === vehicle.getPlate();
            });
            if (idx === -1) {
                state.vehicles.push(new ClipboardVehicle(vehicle));
            }
        },
        removeVehicle(state: ClipboardModuleState, plate: string): void {
            const idx = state.vehicles.findIndex((o: ClipboardVehicle) => {
                return o.plate === plate;
            });
            if (idx > -1) {
                state.vehicles.splice(idx, 1);
            }
        },
        clearVehicles(state: ClipboardModuleState): void {
            state.vehicles.splice(0, state.vehicles.length);
        },
    },
    getters: getters,
};

export default clipboardModule;

export class ClipboardUser {
    public id: number | undefined;
    public job: string | undefined;
    public jobLabel: string | undefined;
    public jobGrade: number | undefined;
    public jobGradeLabel: string | undefined;
    public firstname: string | undefined;
    public lastname: string | undefined;
    public sex: string | undefined;
    public dateofbirth: string | undefined;

    setUser(u: User): ClipboardUser {
        this.id = u.getUserId();
        this.job = u.getJob();
        this.jobLabel = u.getJobLabel();
        this.jobGrade = u.getJobGrade();
        this.jobGradeLabel = u.getJobGradeLabel();
        this.firstname = u.getFirstname();
        this.lastname = u.getLastname();
        this.sex = u.getSex();
        this.dateofbirth = u.getDateofbirth();

        return this;
    }

    setUserShort(u: UserShort): ClipboardUser {
        this.id = u.getUserId();
        this.job = u.getJob();
        this.jobLabel = u.getJobLabel();
        this.jobGrade = u.getJobGrade();
        this.jobGradeLabel = u.getJobGradeLabel();
        this.firstname = u.getFirstname();
        this.lastname = u.getLastname();

        return this;
    }
}

function getUser(obj: ClipboardUser): User {
    const u = new User();
    u.setUserId(obj['id']!);
    u.setJob(obj['job']!);
    u.setJobLabel(obj['jobLabel']!);
    u.setJobGrade(obj['jobGrade']!);
    u.setJobGradeLabel(obj['jobGradeLabel']!);
    u.setFirstname(obj['firstname']!);
    u.setLastname(obj['lastname']!);
    if (obj['sex']) {
        u.setSex(obj['sex']!);
    }
    if (obj['dateofbirth']) {
        u.setDateofbirth(obj['dateofbirth']!);
    }

    return u;
}

export class ClipboardDocument {
    public id: number;
    public title: string;
    public creator: ClipboardUser;

    constructor(d: Document) {
        this.id = d.getId();
        this.title = d.getTitle();
        const creator = new ClipboardUser();
        if (d.getCreator()) {
            creator.setUserShort(d.getCreator()!);
        }
        this.creator = creator;
    }
}

export class ClipboardVehicle {
    public plate: string;
    public model: string;
    public type: string;
    public owner: ClipboardUser;

    constructor(v: Vehicle) {
        this.plate = v.getPlate();
        this.model = v.getModel();
        this.type = v.getType();
        const owner = new ClipboardUser();
        if (v.getOwner()) {
            owner.setUserShort(v.getOwner()!);
        }
        this.owner = owner;
    }
}
