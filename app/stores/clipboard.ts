import { defineStore } from 'pinia';
import { stringToDate } from '~/utils/time';
import { ContentType } from '~~/gen/ts/resources/common/content/content';
import type { Category } from '~~/gen/ts/resources/documents/category';
import type { Document, DocumentShort } from '~~/gen/ts/resources/documents/documents';
import type { ObjectSpecs, TemplateData } from '~~/gen/ts/resources/documents/templates';
import type { User, UserShort } from '~~/gen/ts/resources/users/users';
import type { Vehicle } from '~~/gen/ts/resources/vehicles/vehicles';

/**
 * Represents a user in the clipboard.
 * @typedef {Object} ClipboardUser
 * @property {number | undefined} userId - The ID of the user.
 * @property {string | undefined} job - The job of the user.
 * @property {string | undefined} jobLabel - The label for the user's job.
 * @property {number | undefined} jobGrade - The grade of the user's job.
 * @property {string | undefined} jobGradeLabel - The label for the user's job grade.
 * @property {string} firstname - The first name of the user.
 * @property {string} lastname - The last name of the user.
 * @property {string | undefined} dateofbirth - The date of birth of the user.
 * @property {string | undefined} phoneNumber - The phone number of the user.
 * @property {string | undefined} profilePicture - The profile picture URL of the user.
 */
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
    public profilePicture: string | undefined;

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
        this.profilePicture = u.profilePicture;

        return this;
    }
}

/**
 * Represents a document in the clipboard.
 * @typedef {Object} ClipboardDocument
 * @property {number} id - The ID of the document.
 * @property {string | undefined} createdAt - The creation date of the document.
 * @property {string} title - The title of the document.
 * @property {ClipboardUser} creator - The creator of the document.
 * @property {Category | undefined} category - The category of the document.
 * @property {Object} meta - Metadata about the document.
 * @property {boolean} meta.closed - Whether the document is closed.
 * @property {boolean} meta.draft - Whether the document is a draft.
 * @property {boolean} meta.public - Whether the document is public.
 * @property {string} meta.state - The state of the document.
 * @property {boolean} meta.approved - Whether the document is approved.
 */
export class ClipboardDocument {
    public id: number;
    public createdAt?: string;
    public title: string;
    public creator: ClipboardUser;
    public category: Category | undefined;
    public meta: {
        closed: boolean;
        draft: boolean;
        public: boolean;
        state: string;
        approved: boolean;
    };

    constructor(d: Document) {
        this.id = d.id;
        this.createdAt = d.createdAt ? toDate(d.createdAt).toJSON() : undefined;
        this.category = d.category;
        this.title = d.title;
        this.meta = {
            closed: d.meta?.closed || false,
            draft: d.meta?.draft || false,
            public: d.meta?.public || false,
            state: d.meta?.state || '',
            approved: d.meta?.approved || false,
        };
        this.creator = new ClipboardUser(d.creator!);
    }
}

/**
 * Represents a vehicle in the clipboard.
 * @typedef {Object} ClipboardVehicle
 * @property {string} plate - The license plate of the vehicle.
 * @property {string | undefined} model - The model of the vehicle.
 * @property {string} type - The type of the vehicle.
 * @property {ClipboardUser} owner - The owner of the vehicle.
 */
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

/**
 * Converts a ClipboardVehicle object back to a Vehicle object.
 * @param {ClipboardVehicle} obj - The ClipboardVehicle object to convert.
 * @returns {Vehicle} The converted Vehicle object.
 */
export function getVehicle(obj: ClipboardVehicle): Vehicle {
    return {
        plate: obj.plate,
        model: obj.model,
        type: obj.type,
        owner: getUser(obj.owner),
    };
}

/**
 * Converts a ClipboardUser object back to a User object.
 * @param {ClipboardUser} obj - The ClipboardUser object to convert.
 * @returns {User} The converted User object.
 */
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
        profilePicture: obj.profilePicture,
    };

    return u;
}

/**
 * Converts a ClipboardDocument object back to a DocumentShort object.
 * @param {ClipboardDocument} obj - The ClipboardDocument object to convert.
 * @returns {DocumentShort} The converted DocumentShort object.
 */
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
        meta: {
            documentId: obj.id,
            closed: obj.meta.closed,
            draft: obj.meta.draft,
            public: obj.meta.public,
            state: obj.meta.state,
            approved: obj.meta.approved,
        },
    };
    if (obj.createdAt !== undefined) {
        doc.createdAt = toTimestamp(stringToDate(obj.createdAt));
    }
    return doc;
}

/**
 * Represents the clipboard data structure.
 * @typedef {Object} ClipboardData
 * @property {ClipboardDocument[]} documents - The list of documents in the clipboard.
 * @property {ClipboardUser[]} users - The list of users in the clipboard.
 * @property {ClipboardVehicle[]} vehicles - The list of vehicles in the clipboard.
 */
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

        /**
         * Retrieves template data from the active stack.
         * @returns {TemplateData} The template data containing documents, users, and vehicles.
         */
        const getTemplateData = (): TemplateData => ({
            documents: activeStack.value.documents.map(getDocument),
            users: activeStack.value.users.map(getUser),
            vehicles: activeStack.value.vehicles.map(getVehicle),
        });

        /**
         * Promotes a specific list type to the active stack.
         * @param {ListType} listType - The type of list to promote (e.g., 'documents', 'citizens', 'vehicles').
         */
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

        /**
         * Clears the active stack.
         */
        const clearActiveStack = (): void => {
            activeStack.value.documents.length = 0;
            activeStack.value.users.length = 0;
            activeStack.value.vehicles.length = 0;
        };

        /**
         * Adds a document to the clipboard.
         * @param {Document} document - The document to add.
         */
        const addDocument = (document: Document): void => {
            if (!documents.value.find((o) => o.id === document.id)) {
                documents.value.unshift(new ClipboardDocument(unref(document)));
            }
        };

        /**
         * Removes a document from the clipboard by ID.
         * @param {number} id - The ID of the document to remove.
         */
        const removeDocument = (id: number): void => {
            documents.value = documents.value.filter((o) => o.id !== id);
        };

        /**
         * Clears all documents from the clipboard.
         */
        const clearDocuments = (): void => {
            documents.value = [];
        };

        /**
         * Adds a user to the clipboard.
         * @param {User} user - The user to add.
         * @param {boolean} [active] - Whether to promote the user to the active stack.
         */
        const addUser = (user: User, active?: boolean): void => {
            if (!users.value.find((o) => o.userId === user.userId)) {
                users.value.unshift(new ClipboardUser(unref(user)));
            }
            if (active) promoteToActiveStack('citizens');
        };

        /**
         * Removes a user from the clipboard by ID.
         * @param {number} id - The ID of the user to remove.
         */
        const removeUser = (id: number): void => {
            users.value = users.value.filter((o) => o.userId !== id);
        };

        /**
         * Clears all users from the clipboard.
         */
        const clearUsers = (): void => {
            users.value = [];
        };

        /**
         * Adds a vehicle to the clipboard.
         * @param {Vehicle} vehicle - The vehicle to add.
         */
        const addVehicle = (vehicle: Vehicle): void => {
            if (!vehicles.value.find((o) => o.plate === vehicle.plate)) {
                vehicles.value.unshift(new ClipboardVehicle(unref(vehicle)));
            }
        };

        /**
         * Removes a vehicle from the clipboard by plate number.
         * @param {string} plate - The plate number of the vehicle to remove.
         */
        const removeVehicle = (plate: string): void => {
            vehicles.value = vehicles.value.filter((o) => o.plate !== plate);
        };

        /**
         * Clears all vehicles from the clipboard.
         */
        const clearVehicles = (): void => {
            vehicles.value = [];
        };

        /**
         * Clears all data from the clipboard, including documents, users, vehicles, and the active stack.
         */
        const clear = (): void => {
            clearDocuments();
            clearUsers();
            clearVehicles();
            clearActiveStack();
        };

        /**
         * Checks if the clipboard meets the specified requirements.
         * @param {ObjectSpecs} reqs - The requirements to check against.
         * @param {ListType} listType - The type of list to check (e.g., 'documents', 'citizens', 'vehicles').
         * @returns {boolean} True if the requirements are met, false otherwise.
         */
        const checkRequirements = (reqs: ObjectSpecs, listType: ListType): boolean => {
            const listLength =
                listType === 'documents' ? documents.value : listType === 'citizens' ? users.value : vehicles.value;
            const length = listLength.length;

            // Check if the list is required and empty
            if (reqs.required && length === 0) {
                return false;
            }

            // Check minimum length requirement
            if (typeof reqs.min === 'number' && length < reqs.min) {
                return false;
            }

            // Check maximum length requirement
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
