import Dexie, { type Table } from 'dexie';
import type { Email, EmailShort } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
import type { Message } from '~~/gen/ts/resources/mailer/message';
import type { EmailSettings } from '~~/gen/ts/resources/mailer/settings';
import type { Thread, ThreadState } from '~~/gen/ts/resources/mailer/thread';
import type {
    CreateThreadRequest,
    CreateThreadResponse,
    DeleteThreadRequest,
    DeleteThreadResponse,
    GetEmailSettingsRequest,
    GetEmailSettingsResponse,
    PostMessageRequest,
    PostMessageResponse,
    SetEmailSettingsRequest,
    SetEmailSettingsResponse,
} from '~~/gen/ts/services/mailer/mailer';

const logger = useLogger('💬 Mailer');

export interface MailerState {
    draft: {
        title: string;
        content: string;
        emails: string[];
    };
    emails: (Email | EmailShort)[];
    selectedEmail: EmailShort | undefined;
    selectedThread: Thread | undefined;
    settings: EmailSettings;
}

export const useMailerStore = defineStore('mailer', {
    state: () =>
        ({
            draft: {
                title: '',
                content: '',
                emails: [],
            },
            emails: [],
            selectedEmail: undefined,
            selectedThread: undefined,
            settings: {
                emailId: '0',
                blockedEmails: [],
            },
        }) as MailerState,
    persist: false,
    actions: {
        async handleEvent(event: MailerEvent): Promise<void> {
            console.log('MAILEREVENT', event);
            if (event.data.oneofKind === 'emailUpdate') {
                const idx = this.emails.findIndex((e) => {
                    if (event.data.oneofKind !== 'emailDelete') {
                        return false;
                    }

                    return e.id === event.data.emailDelete;
                });
                if (idx > -1) {
                    this.emails[idx] = event.data.emailUpdate;
                }
            } else if (event.data.oneofKind === 'emailDelete') {
                const idx = this.emails.findIndex((e) => {
                    if (event.data.oneofKind !== 'emailDelete') {
                        return false;
                    }

                    return e.id === event.data.emailDelete;
                });
                if (idx > -1) {
                    this.emails.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'threadUpdate') {
                await mailerDB.threads.put(event.data.threadUpdate);
            } else if (event.data.oneofKind === 'threadDelete') {
                await mailerDB.threads.delete(event.data.threadDelete);
            } else if (event.data.oneofKind === 'messageUpdate') {
                const msg = await mailerDB.messages.get(event.data.messageUpdate.id);
                await mailerDB.messages.put(event.data.messageUpdate);

                if (!msg) {
                    // Only set unread state when message isn't from same email
                    if (event.data.messageUpdate.senderId !== this.selectedEmail?.id) {
                        useSound().play({ name: 'notification' });
                    }

                    if (this.selectedThread?.id !== event.data.messageUpdate.threadId) {
                        await this.setThreadState(
                            {
                                threadId: event.data.messageUpdate.threadId,
                                unread: true,
                            },
                            true,
                        );
                    }
                }
            } else if (event.data.oneofKind === 'messageDelete') {
                await mailerDB.messages.delete(event.data.messageDelete);
            } else if (event.data.oneofKind === 'threadStateUpdate') {
                this.setThreadState(event.data.threadStateUpdate, true);
            } else {
                logger.debug('Unknown MailerEvent received:', event.data.oneofKind);
            }
        },

        // Emails
        async listEmails(): Promise<Email[]> {
            try {
                const call = getGRPCMailerClient().listEmails({});
                const { response } = await call;

                this.emails = response.emails;
                if (this.emails.length > 0) {
                    this.selectedEmail = this.emails[0];
                }

                return this.emails;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Thread
        async getThread(threadId: string): Promise<Thread | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            try {
                const call = getGRPCMailerClient().getThread({
                    emailId: this.selectedEmail?.id,
                    threadId: threadId,
                });
                const { response } = await call;

                if (response.thread) {
                    if (!response.thread.state) {
                        response.thread.state = {
                            emailId: this.selectedEmail.id,
                            threadId: response.thread.id,
                            unread: false,
                            favorite: false,
                            important: false,
                            muted: false,
                            archived: false,
                        };
                    }

                    await mailerDB.threads.put(response.thread);
                } else {
                    await mailerDB.threads.delete(threadId);
                }

                return response.thread;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async createThread(req: CreateThreadRequest): Promise<CreateThreadResponse> {
            try {
                const call = getGRPCMailerClient().createThread(req);
                const { response } = await call;

                if (response.thread) {
                    await mailerDB.threads.put(response.thread);
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async deleteThread(req: DeleteThreadRequest): Promise<DeleteThreadResponse> {
            try {
                const call = getGRPCMailerClient().deleteThread(req);
                const { response } = await call;

                await mailerDB.threads.delete(req.threadId);

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Thread User State
        async getThreadState(threadId: string): Promise<ThreadState | undefined> {
            return (await mailerDB.threads.get(threadId))?.state;
        },

        async setThreadState(state: Partial<ThreadState>, local?: boolean): Promise<ThreadState | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            const thread = await mailerDB.threads.get(state!.threadId);
            if (!thread) {
                return;
            }

            let update = false;
            if (!thread.state) {
                update = true;

                thread.state = {
                    threadId: state.threadId!,
                    emailId: this.selectedEmail?.id,
                    lastRead: toTimestamp(),
                    unread: false,
                    important: false,
                    favorite: false,
                    muted: false,
                    archived: false,
                };
            } else {
                if (state.lastRead !== undefined && thread.state.lastRead?.timestamp !== state.lastRead.timestamp) {
                    update = true;
                    thread.state.lastRead = state.lastRead;
                }
                if (state.unread !== undefined && thread.state.unread !== state.unread) {
                    update = true;
                    thread.state.unread = state.unread;
                }
                if (state.important !== undefined && thread.state.important !== state.important) {
                    update = true;
                    thread.state.important = state.important;
                }
                if (state.favorite !== undefined && thread.state.favorite !== state.favorite) {
                    update = true;
                    thread.state.favorite = state.favorite;
                }
                if (state.muted !== undefined && thread.state.muted !== state.muted) {
                    update = true;
                    thread.state.muted = state.muted;
                }
                if (state.archived !== undefined && thread.state.archived !== state.archived) {
                    update = true;
                    thread.state.archived = state.archived;
                }
            }

            if (update) {
                mailerDB.threads.update(thread.id, thread);

                if (!local) {
                    await getGRPCMailerClient().setThreadState({
                        state: thread.state,
                    });
                }
            }

            return thread.state;
        },

        // Messages
        async postMessage(req: PostMessageRequest): Promise<PostMessageResponse> {
            try {
                const call = getGRPCMailerClient().postMessage(req);
                const { response } = await call;

                if (response.message) {
                    mailerDB.messages.add(response.message);
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // User Settings
        async getEmailSettings(req: GetEmailSettingsRequest): Promise<GetEmailSettingsResponse> {
            try {
                const call = getGRPCMailerClient().getEmailSettings(req);
                const { response } = await call;

                if (response.settings) {
                    this.settings = response.settings;
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
        async setEmailSettings(req: SetEmailSettingsRequest): Promise<SetEmailSettingsResponse> {
            try {
                const call = getGRPCMailerClient().setEmailSettings(req);
                const { response } = await call;

                if (response.settings) {
                    this.settings = response.settings;
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },
    },
});

class MailerDexie extends Dexie {
    threads!: Table<Thread>;
    messages!: Table<Message>;

    constructor() {
        super('mailer');
        this.version(1).stores({
            threads: 'id, creatorEmailId',
            messages: 'id, threadId, emailId',
        });
    }
}

export const mailerDB = new MailerDexie();

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useMailerStore, import.meta.hot));
}
