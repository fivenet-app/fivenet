import { defineStore } from 'pinia';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Document, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import type { ObjectSpecs, TemplateData } from '~~/gen/ts/resources/documents/templates';
import type { User, UserShort } from '~~/gen/ts/resources/users/users';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

export class ClipboardUser {
    public userId: number | undefined;
    public job: string | undefined;
    public jobLabel: string | undefined;
    public jobGrade: number | undefined;
    public jobGradeLabel: string | undefined;
    public firstname: string;
    public lastname: string;
    public dateofbirth: string | undefined;
    public phoneNumber: string | undefined;
    public avatar: string | undefined;

    constructor(u: UserShort | User) {
        this.userId = u.userId;
        this.job = u.job;
        this.jobLabel = u.jobLabel;
        this.jobGrade = u.jobGrade;
        this.jobGradeLabel = u.jobGradeLabel;
        this.firstname = u.firstname;
        this.lastname = u.lastname;
        this.dateofbirth = u.dateofbirth;
        this.phoneNumber = u.phoneNumber;
        this.avatar = u.avatar;

        return this;
    }
}

export class ClipboardDocument {
    public id: number;
    public createdAt?: string;
    public title: string;
    public creator: ClipboardUser;
    public category: Category | undefined;
    public state: string;
    public closed: boolean;
    public draft: boolean;
    public public: boolean;

    constructor(d: Document) {
        this.id = d.id;
        this.createdAt = d.createdAt ? toDate(d.createdAt).toJSON() : undefined;
        this.category = d.category;
        this.title = d.title;
        this.state = d.state;
        this.creator = new ClipboardUser(d.creator!);
        this.closed = d.closed;
        this.draft = d.draft;
        this.public = d.public;
    }
}

export class ClipboardVehicle {
    public plate: string;
    public model: string | undefined;
    public type: string;
    public owner: ClipboardUser;

    constructor(v: Vehicle) {
        this.plate = v.plate;
        this.model = v.model;
        this.type = v.type;
        this.owner = new ClipboardUser(v.owner!);
    }
}

export function getVehicle(obj: ClipboardVehicle): Vehicle {
    return {
        plate: obj.plate,
        model: obj.model,
        type: obj.type,
        owner: getUser(obj.owner),
    };
}

export function getUser(obj: ClipboardUser): User {
    const u: User = {
        userId: obj.userId!,
        job: obj.job!,
        jobLabel: obj.jobLabel ?? '',
        jobGrade: obj.jobGrade!,
        jobGradeLabel: obj.jobGradeLabel ?? '',
        firstname: obj.firstname!,
        lastname: obj.lastname!,
        dateofbirth: obj.dateofbirth ?? '',
        phoneNumber: obj.phoneNumber ?? '',
        licenses: [],
        avatar: obj.avatar,
    };

    return u;
}

export function getDocument(obj: ClipboardDocument): DocumentShort {
    const user = getUser(obj.creator);

    const doc: DocumentShort = {
        id: obj.id,
        categoryId: obj.category && obj.category.id ? obj.category.id : 0,
        category: obj.category,
        title: obj.title,
        contentType: ContentType.HTML,
        content: {
            rawContent: '',
        },
        creatorId: user.userId,
        creator: user,
        creatorJob: user.job,
        state: obj.state,
        closed: obj.closed,
        draft: obj.draft,
        public: obj.public,
    };
    if (obj.createdAt !== undefined) {
        doc.createdAt = toTimestamp(fromString(obj.createdAt));
    }
    return doc;
}

export interface ClipboardData {
    documents: ClipboardDocument[];
    users: ClipboardUser[];
    vehicles: ClipboardVehicle[];
}

export type ListType = 'citizens' | 'documents' | 'vehicles';

export const useClipboardStore = defineStore(
    'clipboard',
    () => {
        const users = ref<ClipboardUser[]>([]);
        const documents = ref<ClipboardDocument[]>([]);
        const vehicles = ref<ClipboardVehicle[]>([]);
        const activeStack = ref<ClipboardData>({
            users: [],
            documents: [],
            vehicles: [],
        });

        const getTemplateData = (): TemplateData => ({
            documents: activeStack.value.documents.map(getDocument),
            users: activeStack.value.users.map(getUser),
            vehicles: activeStack.value.vehicles.map(getVehicle),
        });

        const promoteToActiveStack = (listType: ListType): void => {
            switch (listType) {
                case 'documents':
                    activeStack.value.documents = JSON.parse(JSON.stringify(documents.value)) as ClipboardDocument[];
                    break;
                case 'citizens':
                    activeStack.value.users = JSON.parse(JSON.stringify(users.value)) as ClipboardUser[];
                    break;
                case 'vehicles':
                    activeStack.value.vehicles = JSON.parse(JSON.stringify(vehicles.value)) as ClipboardVehicle[];
                    break;
            }
        };

        const clearActiveStack = (): void => {
            activeStack.value.documents.length = 0;
            activeStack.value.users.length = 0;
            activeStack.value.vehicles.length = 0;
        };

        const addDocument = (document: Document): void => {
            if (!documents.value.find((o) => o.id === document.id)) {
                documents.value.unshift(new ClipboardDocument(unref(document)));
            }
        };

        const removeDocument = (id: number): void => {
            documents.value = documents.value.filter((o) => o.id !== id);
        };

        const clearDocuments = (): void => {
            documents.value = [];
        };

        const addUser = (user: User, active?: boolean): void => {
            if (!users.value.find((o) => o.userId === user.userId)) {
                users.value.unshift(new ClipboardUser(unref(user)));
            }
            if (active) promoteToActiveStack('citizens');
        };

        const removeUser = (id: number): void => {
            users.value = users.value.filter((o) => o.userId !== id);
        };

        const clearUsers = (): void => {
            users.value = [];
        };

        const addVehicle = (vehicle: Vehicle): void => {
            if (!vehicles.value.find((o) => o.plate === vehicle.plate)) {
                vehicles.value.unshift(new ClipboardVehicle(unref(vehicle)));
            }
        };

        const removeVehicle = (plate: string): void => {
            vehicles.value = vehicles.value.filter((o) => o.plate !== plate);
        };

        const clearVehicles = (): void => {
            vehicles.value = [];
        };

        const clear = (): void => {
            clearDocuments();
            clearUsers();
            clearVehicles();
            clearActiveStack();
        };

        const checkRequirements = (reqs: ObjectSpecs, listType: ListType): boolean => {
            const length = (listType === 'documents' ? documents.value : listType === 'citizens' ? users.value : vehicles.value)
                .length;
            if (reqs.max !== undefined) {
                reqs.max = reqs.min;
            }

            if (reqs.required && length === 0) {
                return false;
            }
            if (typeof reqs.min === 'number' && length < reqs.min) {
                return false;
            }
            if (typeof reqs.max === 'number' && length > reqs.max) {
                return false;
            }
            return true;
        };

        return {
            users,
            documents,
            vehicles,
            activeStack,

            getTemplateData,
            promoteToActiveStack,
            clearActiveStack,
            addDocument,
            removeDocument,
            clearDocuments,
            addUser,
            removeUser,
            clearUsers,
            addVehicle,
            removeVehicle,
            clearVehicles,
            clear,
            checkRequirements,
        };
    },
    {
        persist: true,
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useClipboardStore, import.meta.hot));
}
