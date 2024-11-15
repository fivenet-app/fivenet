import Dexie, { type Table } from 'dexie';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
import type { Message } from '~~/gen/ts/resources/mailer/message';
import type { UserSettings } from '~~/gen/ts/resources/mailer/settings';
import type { Thread, ThreadStateUser } from '~~/gen/ts/resources/mailer/thread';
import type { UserShort } from '~~/gen/ts/resources/users/users';
import type {
    CreateThreadRequest,
    CreateThreadResponse,
    DeleteThreadRequest,
    DeleteThreadResponse,
    GetEmailSettingsRequest,
    GetEmailSettingsResponse,
    LeaveThreadResponse,
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
        users: UserShort[];
    };
    selectedThread: Thread | undefined;
    settings: UserSettings;
}

export const useMailerStore = defineStore('mailer', {
    state: () =>
        ({
            draft: {
                title: '',
                content: '',
                users: [],
            },
            selectedThread: undefined,
            settings: {
                userId: 0,
                blockedUsers: [],
            },
        }) as MailerState,
    persist: false,
    actions: {
        async handleEvent(event: MailerEvent): Promise<void> {
            if (event.data.oneofKind === 'threadUpdate') {
                await mailerDB.threads.put(event.data.threadUpdate);
            } else if (event.data.oneofKind === 'threadDelete') {
                await mailerDB.threads.delete(event.data.threadDelete);
            } else if (event.data.oneofKind === 'messageUpdate') {
                const msg = await mailerDB.messages.get(event.data.messageUpdate.id);
                await mailerDB.messages.put(event.data.messageUpdate);

                if (!msg) {
                    // Only set unread state when message isn't from user himself
                    const { activeChar } = useAuth();
                    if (event.data.messageUpdate.creatorId !== activeChar.value?.userId) {
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
            } else {
                logger.debug('Unknown MailerEvent received:', event.data.oneofKind);
            }
        },

        // Emails
        async listEmails(): Promise<Email[]> {
            try {
                const call = getGRPCMailerClient().listEmails({});
                const { response } = await call;

                return response.emails;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Thread
        async getThread(threadId: string): Promise<Thread | undefined> {
            const { activeChar } = useAuth();

            try {
                const call = getGRPCMailerClient().getThread({
                    threadId: threadId,
                });
                const { response } = await call;

                if (response.thread) {
                    if (!response.thread.userState) {
                        response.thread.userState = {
                            userId: activeChar.value?.userId ?? 0,
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

        async leaveThread(threadId: string): Promise<LeaveThreadResponse> {
            try {
                const call = getGRPCMailerClient().leaveThread({
                    threadId: threadId,
                });
                const { response } = await call;

                await mailerDB.threads.delete(threadId);

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Thread User State
        async getThreadUserState(threadId: string): Promise<ThreadStateUser | undefined> {
            return (await mailerDB.threads.get(threadId))?.userState;
        },

        async setThreadState(state: Partial<ThreadStateUser>, local?: boolean): Promise<ThreadStateUser | undefined> {
            const thread = await mailerDB.threads.get(state!.threadId);
            if (!thread) {
                return;
            }

            let update = false;
            if (!thread.userState) {
                update = true;

                thread.userState = {
                    threadId: state.threadId!,
                    userId: 0,
                    lastRead: toTimestamp(),
                    unread: false,
                    important: false,
                    favorite: false,
                    muted: false,
                    archived: false,
                };
            } else {
                if (state.lastRead !== undefined && thread.userState.lastRead?.timestamp !== state.lastRead.timestamp) {
                    update = true;
                    thread.userState.lastRead = state.lastRead;
                }
                if (state.unread !== undefined && thread.userState.unread !== state.unread) {
                    update = true;
                    thread.userState.unread = state.unread;
                }
                if (state.important !== undefined && thread.userState.important !== state.important) {
                    update = true;
                    thread.userState.important = state.important;
                }
                if (state.favorite !== undefined && thread.userState.favorite !== state.favorite) {
                    update = true;
                    thread.userState.favorite = state.favorite;
                }
                if (state.muted !== undefined && thread.userState.muted !== state.muted) {
                    update = true;
                    thread.userState.muted = state.muted;
                }
                if (state.archived !== undefined && thread.userState.archived !== state.archived) {
                    update = true;
                    thread.userState.archived = state.archived;
                }
            }

            if (update) {
                mailerDB.threads.update(thread.id, thread);

                if (!local) {
                    await getGRPCMailerClient().setThreadState({
                        state: {
                            oneofKind: 'user',
                            user: thread.userState,
                        },
                    });
                }
            }

            return thread.userState;
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
        async getUserSettings(req: GetEmailSettingsRequest): Promise<GetEmailSettingsResponse> {
            try {
                const call = getGRPCMailerClient().getUserSettings(req);
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
            threads: 'id',
            messages: 'id, threadId',
        });
    }
}

export const mailerDB = new MailerDexie();

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useMailerStore, import.meta.hot));
}
