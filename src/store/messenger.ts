import Dexie, { type Table } from 'dexie';
import type { MessengerEvent } from '~~/gen/ts/resources/messenger/events';
import type { Message } from '~~/gen/ts/resources/messenger/message';
import type { Thread, ThreadUserState } from '~~/gen/ts/resources/messenger/thread';
import type {
    CreateOrUpdateThreadRequest,
    CreateOrUpdateThreadResponse,
    DeleteThreadRequest,
    DeleteThreadResponse,
    LeaveThreadResponse,
    PostMessageRequest,
    PostMessageResponse,
} from '~~/gen/ts/services/messenger/messenger';
import { useAuthStore } from './auth';

export interface MessengerState {
    selectedThread: Thread | undefined;
}

export const useMessengerStore = defineStore('messenger', {
    state: () =>
        ({
            selectedThread: undefined,
        }) as MessengerState,
    persist: false,
    actions: {
        async handleEvent(event: MessengerEvent): Promise<void> {
            if (event.data.oneofKind === 'threadUpdate') {
                await messengerDB.threads.put(event.data.threadUpdate);
            } else if (event.data.oneofKind === 'threadDelete') {
                await messengerDB.threads.delete(event.data.threadDelete);
            } else if (event.data.oneofKind === 'messageUpdate') {
                const msg = await messengerDB.messages.get(event.data.messageUpdate.id);
                await messengerDB.messages.put(event.data.messageUpdate);

                if (!msg) {
                    // Only set unread state when message isn't from user himself
                    const authStore = useAuthStore();
                    const { activeChar } = authStore;
                    if (event.data.messageUpdate.creatorId !== activeChar?.userId) {
                        useSound().play({ name: 'notification' });
                    }

                    if (this.selectedThread?.id !== event.data.messageUpdate.threadId) {
                        await this.setThreadUserState(
                            {
                                threadId: event.data.messageUpdate.threadId,
                                unread: true,
                            },
                            true,
                        );
                    }
                }
            } else if (event.data.oneofKind === 'messageDelete') {
                await messengerDB.messages.delete(event.data.messageDelete);
            } else {
                console.debug('Messenger: Unknown MessengerEvent received:', event.data.oneofKind);
            }
        },

        // Thread
        async getThread(threadId: string, local?: boolean): Promise<Thread | undefined> {
            if (local) {
                const thread = await messengerDB.threads.get(threadId);
                if (thread) {
                    return thread;
                }
            }

            try {
                const call = getGRPCMessengerClient().getThread({
                    threadId: threadId,
                });
                const { response } = await call;

                if (response.thread) {
                    if (!response.thread.userState) {
                        response.thread.userState = {
                            userId: useAuthStore().activeChar?.userId ?? 0,
                            threadId: response.thread.id,
                            unread: false,
                            favorite: false,
                            important: false,
                            muted: false,
                        };
                    }

                    await messengerDB.threads.put(response.thread);
                } else {
                    await messengerDB.threads.delete(threadId);
                }

                return response.thread;
            } catch (e) {
                handleGRPCError(e);
                throw e;
            }
        },

        async createOrUpdateThread(req: CreateOrUpdateThreadRequest): Promise<CreateOrUpdateThreadResponse> {
            try {
                const call = getGRPCMessengerClient().createOrUpdateThread(req);
                const { response } = await call;

                if (response.thread) {
                    await messengerDB.threads.put(response.thread);
                }

                return response;
            } catch (e) {
                handleGRPCError(e);
                throw e;
            }
        },

        async deleteThread(req: DeleteThreadRequest): Promise<DeleteThreadResponse> {
            try {
                const call = getGRPCMessengerClient().deleteThread(req);
                const { response } = await call;

                await messengerDB.threads.delete(req.threadId);

                return response;
            } catch (e) {
                handleGRPCError(e);
                throw e;
            }
        },

        async leaveThread(threadId: string): Promise<LeaveThreadResponse> {
            try {
                const call = getGRPCMessengerClient().leaveThread({
                    threadId: threadId,
                });
                const { response } = await call;

                await messengerDB.threads.delete(threadId);

                return response;
            } catch (e) {
                handleGRPCError(e);
                throw e;
            }
        },

        // Thread User State
        async getThreadUserState(threadId: string): Promise<ThreadUserState | undefined> {
            return (await messengerDB.threads.get(threadId))?.userState;
        },

        async setThreadUserState(state: Partial<ThreadUserState>, local?: boolean): Promise<ThreadUserState | undefined> {
            const thread = await messengerDB.threads.get(state!.threadId);
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
            }

            if (update) {
                messengerDB.threads.update(thread.id, thread);

                if (!local) {
                    await getGRPCMessengerClient().setThreadUserState({ state: thread.userState });
                }
            }

            return thread.userState;
        },

        // Messages
        async postMessage(req: PostMessageRequest): Promise<PostMessageResponse> {
            try {
                const call = getGRPCMessengerClient().postMessage(req);
                const { response } = await call;

                if (response.message) {
                    messengerDB.messages.add(response.message);
                }

                return response;
            } catch (e) {
                handleGRPCError(e);
                throw e;
            }
        },
    },
});

class MessengerDexie extends Dexie {
    threads!: Table<Thread>;
    messages!: Table<Message>;

    constructor() {
        super('messenger');
        this.version(1).stores({
            threads: 'id',
            messages: 'id, threadId',
        });
    }
}

export const messengerDB = new MessengerDexie();
