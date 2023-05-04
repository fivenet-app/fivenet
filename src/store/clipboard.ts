import { StoreDefinition, defineStore } from 'pinia';
import { TemplateData } from '@fivenet/gen/resources/documents/templates_pb';
import { User, UserShort } from '@fivenet/gen/resources/users/users_pb';
import { Document } from '@fivenet/gen/resources/documents/documents_pb';
import { Vehicle } from '@fivenet/gen/resources/vehicles/vehicles_pb';
import { fromString, toDateLocaleString } from '~/utils/time';
import { Timestamp } from '@fivenet/gen/resources/timestamp/timestamp_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';

export interface ClipboardData {
    documents: ClipboardDocument[];
    users: ClipboardUser[];
    vehicles: ClipboardVehicle[];
}

export interface ClipboardState extends ClipboardData {
    activeStack: ClipboardData;
}

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
    actions: {
        getTemplateData(): TemplateData {
            const data = new TemplateData();

            if (this.activeStack.documents) {
                const documents = new Array<Document>();
                this.activeStack.documents.forEach((v: ClipboardDocument) => {
                    documents.push(getDocument(v));
                });
                data.setDocumentsList(documents);
            }
            if (this.activeStack.users) {
                const users = new Array<User>();
                this.activeStack.users.forEach((v: ClipboardUser) => {
                    users.push(getUser(v));
                });
                data.setUsersList(users);
            }
            if (this.activeStack.vehicles) {
                const vehicles = new Array<Vehicle>();
                this.activeStack.vehicles.forEach((v: ClipboardVehicle) => {
                    vehicles.push(v.getVehicle());
                });
                data.setVehiclesList(vehicles);
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
            const idx = this.documents.findIndex((o: ClipboardDocument) => {
                return o.id === document.getId();
            });
            if (idx === -1) {
                this.documents.unshift(new ClipboardDocument(document));
            }
        },
        removeDocument(id: number): void {
            const idx = this.documents.findIndex((o: ClipboardDocument) => {
                return o.id === id;
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
                return o.id === user.getUserId();
            });
            if (idx === -1) {
                this.users.unshift(new ClipboardUser().setUser(user));
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
                return o.plate === vehicle.getPlate();
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

        clear() {
            this.clearActiveStack();
            this.clearDocuments();
            this.clearUsers();
            this.clearVehicles();
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
    public firstname: string | undefined;
    public lastname: string | undefined;

    setUser(u: UserShort | User): ClipboardUser {
        this.id = u.getUserId();
        this.identifier = u.getIdentifier();
        this.job = u.getJob();
        this.jobLabel = u.getJobLabel();
        this.jobGrade = u.getJobGrade();
        this.jobGradeLabel = u.getJobGradeLabel();
        this.firstname = u.getFirstname();
        this.lastname = u.getLastname();

        return this;
    }
}

export function getUser(obj: ClipboardUser): User {
    const u = new User();
    u.setUserId(obj['id']!);
    u.setIdentifier(obj['identifier']!);
    u.setJob(obj['job']!);
    u.setJobLabel(obj['jobLabel']!);
    u.setJobGrade(obj['jobGrade']!);
    u.setJobGradeLabel(obj['jobGradeLabel']!);
    u.setFirstname(obj['firstname']!);
    u.setLastname(obj['lastname']!);

    return u;
}

export class ClipboardDocument {
    public id: number;
    public createdAt: string;
    public title: string;
    public state: string;
    public creator: ClipboardUser;

    constructor(d: Document) {
        this.id = d.getId();
        this.createdAt = d.getCreatedAt()?.getTimestamp()?.toDate().toLocaleTimeString()!;
        this.title = d.getTitle();
        this.state = d.getState();
        const creator = new ClipboardUser();
        if (d.getCreator()) {
            creator.setUser(d.getCreator()!);
        }
        this.creator = creator;
    }
}

export function getDocument(obj: ClipboardDocument): Document {
    const tts = new google_protobuf_timestamp_pb.Timestamp();
    tts.fromDate(fromString(obj.createdAt)!);
    const ts = new Timestamp();
    ts.setTimestamp(tts);

    const d = new Document();
    d.setId(obj.id);
    d.setCreatedAt(ts);
    d.setTitle(obj.title);
    d.setState(obj.state);
    d.setCreator(getUser(obj.creator));

    return d;
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
            owner.setUser(v.getOwner()!);
        }
        this.owner = owner;
    }

    getVehicle(): Vehicle {
        const v = new Vehicle();
        v.setPlate(this.plate);
        v.setModel(this.model);
        v.setType(this.type);
        v.setOwner(getUser(this.owner));

        return v;
    }
}
