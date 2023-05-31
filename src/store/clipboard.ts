import { StoreDefinition, defineStore } from 'pinia';
import { fromString } from '~/utils/time';
import * as google_protobuf_timestamp from '~~/gen/ts/google/protobuf/timestamp';
import { Document, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import { ObjectSpecs, TemplateData } from '~~/gen/ts/resources/documents/templates';
import { Timestamp } from '~~/gen/ts/resources/timestamp/timestamp';
import { User, UserShort } from '~~/gen/ts/resources/users/users';
import { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

export interface ClipboardData {
    documents: ClipboardDocument[];
    users: ClipboardUser[];
    vehicles: ClipboardVehicle[];
}

export interface ClipboardState extends ClipboardData {
    activeStack: ClipboardData;
}

export type ListType = 'users' | 'documents' | 'vehicles';

export const useClipboardStore = defineStore('clipboard', {
    state: () =>
        ({
            users: [],
            documents: [],
            vehicles: [],
            activeStack: {
                users: [],
                documents: [],
                vehicles: [],
            } as ClipboardData,
        } as ClipboardState),
    persist: true,
    actions: {
        getTemplateData(): TemplateData {
            const data: TemplateData = {
                documents: [],
                users: [],
                vehicles: [],
            };

            if (this.activeStack.documents) {
                this.activeStack.documents.forEach((v: ClipboardDocument) => {
                    data.documents.push(getDocument(v));
                });
            }
            if (this.activeStack.users) {
                this.activeStack.users.forEach((v: ClipboardUser) => {
                    data.users.push(getUser(v));
                });
            }
            if (this.activeStack.vehicles) {
                this.activeStack.vehicles.forEach((v: ClipboardVehicle) => {
                    data.vehicles.push(v.getVehicle());
                });
            }

            return data;
        },

        clearActiveStack(): void {
            this.activeStack.documents.length = 0;
            this.activeStack.users.length = 0;
            this.activeStack.vehicles.length = 0;
        },
        // Documents
        addDocument(document: Document): void {
            const docId = document.id.toString();
            const idx = this.documents.findIndex((o: ClipboardDocument) => {
                return o.id === docId;
            });
            if (idx === -1) {
                this.documents.unshift(new ClipboardDocument(document));
            }
        },
        removeDocument(id: bigint): void {
            const docId = id.toString();
            const idx = this.documents.findIndex((o: ClipboardDocument) => {
                return o.id === docId;
            });
            if (idx > -1) {
                this.documents.splice(idx, 1);
            }
        },
        clearDocuments(): void {
            this.documents.splice(0, this.documents.length);
        },
        // Users
        addUser(user: User): void {
            const idx = this.users.findIndex((o: ClipboardUser) => {
                return o.id === user.userId;
            });
            if (idx === -1) {
                this.users.unshift(new ClipboardUser(user!));
            }
        },
        removeUser(id: number): void {
            const idx = this.users.findIndex((o: ClipboardUser) => {
                return o.id === id;
            });
            if (idx > -1) {
                this.users.splice(idx, 1);
            }
        },
        clearUsers(): void {
            this.users.splice(0, this.users.length);
        },
        // Vehicles
        addVehicle(vehicle: Vehicle): void {
            const idx = this.vehicles.findIndex((o: ClipboardVehicle) => {
                return o.plate === vehicle.plate;
            });
            if (idx === -1) {
                this.vehicles.unshift(new ClipboardVehicle(vehicle));
            }
        },
        removeVehicle(plate: string): void {
            const idx = this.vehicles.findIndex((o: ClipboardVehicle) => {
                return o.plate === plate;
            });
            if (idx > -1) {
                this.vehicles.splice(idx, 1);
            }
        },
        clearVehicles(): void {
            this.vehicles.splice(0, this.vehicles.length);
        },

        clear(): void {
            this.clearDocuments();
            this.clearUsers();
            this.clearVehicles();
            this.clearActiveStack();
        },

        checkRequirements(reqs: ObjectSpecs, listType: ListType): boolean {
            const list = this.$state[listType];

            if (reqs.required && list.length <= 0) {
                return false;
            } else if (reqs.min && reqs.max && list.length > reqs.max && list.length < reqs.min) {
                return false;
            }

            return true;
        },

        promoteToActiveStack(listType: ListType): void {
            const list = this.$state[listType];

            switch (listType) {
                case 'documents':
                    this.$state.activeStack.documents = list as ClipboardDocument[];
                    break;
                case 'users':
                    this.$state.activeStack.users = list as ClipboardUser[];
                    break;
                case 'vehicles':
                    this.$state.activeStack.vehicles = list as ClipboardVehicle[];
                    break;
            }
        },
    },
});

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useClipboardStore as unknown as StoreDefinition, import.meta.hot));
}

export class ClipboardUser {
    public id: number | undefined;
    public identifier: string | undefined;
    public job: string | undefined;
    public jobLabel: string | undefined;
    public jobGrade: number | undefined;
    public jobGradeLabel: string | undefined;
    public dateofbirth: string | undefined;
    public firstname: string | undefined;
    public lastname: string | undefined;

    constructor(u: UserShort | User) {
        this.id = u.userId;
        this.identifier = u.identifier;
        this.job = u.job;
        this.jobLabel = u.jobLabel;
        this.jobGrade = u.jobGrade;
        this.jobGradeLabel = u.jobGradeLabel;
        this.firstname = u.firstname;
        this.lastname = u.lastname;
        if ('dateofbirth' in u) {
            this.dateofbirth = u.dateofbirth;
        }

        return this;
    }
}

export function getUser(obj: ClipboardUser): User {
    const u: User = {
        userId: obj.id!,
        identifier: obj.identifier!,
        job: obj.job!,
        jobLabel: obj.jobLabel!,
        jobGrade: obj.jobGrade!,
        jobGradeLabel: obj.jobGradeLabel!,
        firstname: obj.firstname!,
        lastname: obj.lastname!,
        dateofbirth: obj.dateofbirth!,
        licenses: [],
    };
    if ('dateofbirth' in u) {
        u.dateofbirth = obj.dateofbirth!;
    }

    return u;
}

export class ClipboardDocument {
    public id: string;
    public createdAt: string;
    public title: string;
    public state: string;
    public creator: ClipboardUser;
    public closed: boolean;

    constructor(d: Document) {
        this.id = d.id.toString();
        this.createdAt = google_protobuf_timestamp.Timestamp.toDate(d.createdAt?.timestamp!).toLocaleDateString();
        this.title = d.title;
        this.state = d.state;
        this.creator = new ClipboardUser(d.creator!);
        this.closed = d.closed;
    }
}

export function getDocument(obj: ClipboardDocument): DocumentShort {
    const ts: Timestamp = {
        timestamp: google_protobuf_timestamp.Timestamp.fromDate(fromString(obj.createdAt)!),
    };

    const user = getUser(obj.creator);

    return {
        id: BigInt(obj.id),
        createdAt: {
            timestamp: {
                nanos: ts.timestamp?.nanos!,
                seconds: ts.timestamp?.seconds!,
            },
        },
        title: obj.title,
        state: obj.state,
        creator: user,
        creatorId: user.userId,
        closed: obj.closed,
        categoryId: BigInt(0),
    };
}

export class ClipboardVehicle {
    public plate: string;
    public model: string;
    public type: string;
    public owner: ClipboardUser;

    constructor(v: Vehicle) {
        this.plate = v.plate;
        this.model = v.model;
        this.type = v.type;
        this.owner = new ClipboardUser(v.owner!);
    }

    getVehicle(): Vehicle {
        return {
            plate: this.plate,
            model: this.model,
            type: this.type,
            owner: getUser(this.owner),
        };
    }
}
