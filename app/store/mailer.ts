import Dexie, { type Table } from 'dexie';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
import type { Message } from '~~/gen/ts/resources/mailer/message';
import type { Thread, ThreadState } from '~~/gen/ts/resources/mailer/thread';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type {
    CreateOrUpdateEmailRequest,
    CreateOrUpdateEmailResponse,
    CreateThreadRequest,
    CreateThreadResponse,
    DeleteEmailRequest,
    DeleteEmailResponse,
    DeleteMessageRequest,
    DeleteMessageResponse,
    DeleteThreadRequest,
    DeleteThreadResponse,
    GetEmailSettingsRequest,
    GetEmailSettingsResponse,
    ListEmailsResponse,
    ListThreadMessagesRequest,
    ListThreadMessagesResponse,
    ListThreadsRequest,
    PostMessageRequest,
    PostMessageResponse,
    SetEmailSettingsRequest,
    SetEmailSettingsResponse,
} from '~~/gen/ts/services/mailer/mailer';
import { useNotificatorStore } from './notificator';

const logger = useLogger('💬 Mailer');

export interface MailerState {
    loaded: boolean;
    error: Error | undefined;

    draft: {
        title: string;
        content: string;
        recipients: { label: string }[];
    };

    emails: Email[];
    selectedEmail: Email | undefined;
    selectedThread: Thread | undefined;

    addressBook: { label: string; name?: string }[];
}

export const useMailerStore = defineStore('mailer', {
    state: () =>
        ({
            loaded: false,
            error: undefined,

            draft: {
                title: '',
                content: '',
                recipients: [],
            },

            emails: [],
            selectedEmail: undefined,
            selectedThread: undefined,

            addressBook: [],
        }) as MailerState,
    persist: {
        pick: ['draft', 'addressBook'],
    },
    actions: {
        async handleEvent(event: MailerEvent): Promise<void> {
            logger.debug('Received change - Kind:', event.data.oneofKind, event.data);

            if (event.data.oneofKind === 'emailUpdate') {
                const searchId = event.data.emailUpdate.id;
                const idx = this.emails.findIndex((e) => e.id === searchId);
                if (idx > -1) {
                    this.emails[idx] = event.data.emailUpdate;
                }
            } else if (event.data.oneofKind === 'emailDelete') {
                const searchId = event.data.emailDelete;
                const idx = this.emails.findIndex((e) => e.id === searchId);
                if (idx > -1) {
                    this.emails.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'emailSettingsUpdated') {
                const searchId = event.data.emailSettingsUpdated.emailId;
                const idx = this.emails.findIndex((e) => e.id === searchId);
                if (idx > -1 && this.emails[idx]) {
                    this.emails[idx].settings = event.data.emailSettingsUpdated;
                }
            } else if (event.data.oneofKind === 'threadUpdate') {
                await mailerDB.threads.put(event.data.threadUpdate);

                // Handle email sent by blocked email
                if (
                    event.data.threadUpdate.creatorEmail?.email &&
                    this.checkIfEmailBlocked(event.data.threadUpdate.creatorEmail?.email)
                ) {
                    // Make sure to set thread state accordingly (locally)
                    await this.setThreadState(
                        {
                            archived: true,
                            muted: true,
                        },
                        true,
                    );
                    return;
                }

                useNotificatorStore().add({
                    title: { key: 'notifications.mailer.new_email.title', parameters: {} },
                    description: {
                        key: 'notifications.mailer.new_email.content',
                        parameters: {
                            title: event.data.threadUpdate.title,
                            from: event.data.threadUpdate.creatorEmail?.email ?? 'N/A',
                        },
                    },
                    type: NotificationType.INFO,
                    actions: this.getNotificationActions(event.data.threadUpdate.id),
                });
            } else if (event.data.oneofKind === 'threadDelete') {
                await mailerDB.threads.delete(event.data.threadDelete);
            } else if (event.data.oneofKind === 'messageUpdate') {
                await mailerDB.messages.put(event.data.messageUpdate);

                // Update thread updated at time
                await mailerDB.threads.update(event.data.messageUpdate.threadId, {
                    updatedAt: event.data.messageUpdate.updatedAt,
                });

                // Handle email sent by blocked email
                if (
                    event.data.messageUpdate.sender?.email &&
                    this.checkIfEmailBlocked(event.data.messageUpdate.sender?.email)
                ) {
                    // Make sure to set thread state accordingly (locally)
                    await this.setThreadState(
                        {
                            archived: true,
                            muted: true,
                        },
                        true,
                    );
                    return;
                }

                if (event.data.messageUpdate.senderId === this.selectedEmail?.id) {
                    return;
                }

                // Only set unread state when message isn't from same email and the user isn't active on that thread
                if (event.data.messageUpdate.threadId !== this.selectedThread?.id) {
                    await this.setThreadState(
                        {
                            threadId: event.data.messageUpdate.threadId,
                            unread: true,
                        },
                        true,
                    );
                }

                const threadState = await this.getThreadState(event.data.messageUpdate.threadId);
                if (threadState?.muted) {
                    return;
                }

                useNotificatorStore().add({
                    title: { key: 'notifications.mailer.new_email.title', parameters: {} },
                    description: {
                        key: 'notifications.mailer.new_email.content',
                        parameters: {
                            title: event.data.messageUpdate.title,
                            from: event.data.messageUpdate.sender?.email ?? 'N/A',
                        },
                    },
                    type: NotificationType.INFO,
                    actions: this.getNotificationActions(event.data.messageUpdate.threadId),
                });
                useSound().play({ name: 'notification' });
            } else if (event.data.oneofKind === 'messageDelete') {
                await mailerDB.messages.delete(event.data.messageDelete);
            } else if (event.data.oneofKind === 'threadStateUpdate') {
                this.setThreadState(event.data.threadStateUpdate, true, false);
            } else {
                logger.debug('Unknown MailerEvent type received:', event.data.oneofKind);
            }
        },

        // Emails
        async listEmails(all?: boolean, offset?: number): Promise<ListEmailsResponse> {
            this.error = undefined;

            if (this.addressBook.length > 30) {
                this.addressBook.length = 30;
            }

            try {
                const call = getGRPCMailerClient().listEmails({
                    pagination: {
                        offset: offset ?? 0,
                    },
                    all: all ?? false,
                });
                const { response } = await call;

                this.emails = response.emails;
                if (this.emails.length === 0 || !this.hasPrivateEmail) {
                    await navigateTo({
                        name: 'mail-manage',
                    });
                } else if (this.emails[0]) {
                    if (this.emails[0].settings === undefined) {
                        this.selectedEmail = await this.getEmail(this.emails[0].id);
                    } else {
                        this.selectedEmail = this.emails[0];
                    }
                }

                this.loaded = true;
                return response;
            } catch (e) {
                this.error = e as RpcError;
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async getEmail(id: string): Promise<Email | undefined> {
            try {
                const call = getGRPCMailerClient().getEmail({
                    id: id,
                });
                const { response } = await call;

                if (this.selectedEmail && this.selectedEmail.id === response.email?.id) {
                    this.selectedEmail.settings = response.email.settings;
                    this.selectedEmail.access = response.email.access;
                } else {
                    const email = this.emails.find((e) => e.id === id);
                    if (email) {
                        email.settings = response.email?.settings;
                        email.access = response.email?.access;
                    }
                }

                return response.email;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async createOrUpdateEmail(req: CreateOrUpdateEmailRequest): Promise<CreateOrUpdateEmailResponse> {
            try {
                const call = getGRPCMailerClient().createOrUpdateEmail(req);
                const { response } = await call;

                if (response.email) {
                    const idx = this.emails.findIndex((e) => e.id === response.email!.id);
                    if (idx === -1) {
                        this.emails.unshift(response.email);
                    } else {
                        this.emails[idx] = response.email;
                    }

                    if (this.selectedEmail === undefined) {
                        this.selectedEmail = response.email;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async deleteEmail(req: DeleteEmailRequest): Promise<DeleteEmailResponse> {
            try {
                const call = getGRPCMailerClient().deleteEmail(req);
                const { response } = await call;

                if (this.selectedEmail?.id === req.id) {
                    this.selectedEmail = undefined;
                }

                const idx = this.emails.findIndex((e) => e.id === req.id);
                if (idx > -1) {
                    this.emails.slice(idx, 1);
                }

                useNotificatorStore().restartStream();

                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        // Thread
        async listThreads(req: ListThreadsRequest): Promise<Thread[] | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            try {
                const call = getGRPCMailerClient().listThreads(req);
                const { response } = await call;

                await mailerDB.threads.bulkPut(response.threads);

                return response.threads;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

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
                const error = e as RpcError;
                handleGRPCError(error);

                if (error?.message?.includes('.ErrThreadAccessDenied')) {
                    await Promise.all([
                        mailerDB.threads.delete(threadId),
                        mailerDB.messages.where('threadId').equals(threadId).delete(),
                    ]);
                }
            }
        },

        async createThread(req: CreateThreadRequest): Promise<CreateThreadResponse> {
            try {
                const call = getGRPCMailerClient().createThread(req);
                const { response } = await call;

                if (response.thread) {
                    await mailerDB.threads.put(response.thread);

                    req.recipients.forEach((r) => this.addToAddressBook(r));
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

                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

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

        async setThreadState(state: Partial<ThreadState>, local?: boolean, notify?: boolean): Promise<ThreadState | undefined> {
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
                    thread.state.lastRead = toTimestamp();
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
                await mailerDB.threads.put(thread, thread.id);

                if (!local) {
                    await getGRPCMailerClient().setThreadState({
                        state: thread.state,
                    });
                }
            }

            if (notify) {
                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });
            }

            return thread.state;
        },

        // Messages
        async listThreadMessages(req: ListThreadMessagesRequest): Promise<ListThreadMessagesResponse | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            try {
                const call = getGRPCMailerClient().listThreadMessages(req);
                const { response } = await call;

                await mailerDB.messages.bulkPut(response.messages);

                return response;
            } catch (e) {
                const error = e as RpcError;
                await handleGRPCError(error);

                if (error?.message?.includes('.ErrThreadAccessDenied')) {
                    await Promise.all([
                        mailerDB.threads.delete(req.threadId),
                        mailerDB.messages.where('threadId').equals(req.threadId).delete(),
                    ]);
                }
            }
        },

        async postMessage(req: PostMessageRequest): Promise<PostMessageResponse> {
            try {
                const call = getGRPCMailerClient().postMessage(req);
                const { response } = await call;

                if (response.message) {
                    await mailerDB.messages.put(response.message);

                    req.recipients.forEach((r) => this.addToAddressBook(r));
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async deleteMessage(req: DeleteMessageRequest): Promise<DeleteMessageResponse> {
            try {
                const call = getGRPCMailerClient().deleteMessage(req);
                const { response } = await call;

                await mailerDB.messages.delete(req.messageId);

                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

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

                if (response.settings && this.selectedEmail) {
                    this.selectedEmail.settings = response.settings;
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

                if (response.settings && this.selectedEmail) {
                    this.selectedEmail.settings = response.settings;
                }

                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        checkIfEmailBlocked(email: string): boolean {
            this.emails.find((e) => email);
            return false;
        },

        getNotificationActions(threadId?: string): NotificationActionI18n[] {
            return useRoute().name !== 'mail'
                ? [
                      {
                          label: { key: 'common.click_here' },
                          to: threadId ? { name: 'mail', query: { thread: threadId }, hash: '#' } : { name: 'mail' },
                      },
                  ]
                : [];
        },

        // Address book
        addToAddressBook(email: string, label?: string): void {
            email = email.trim();
            label = label?.trim();

            const idx = this.addressBook.findIndex((a) => a.label === email);
            if (idx > -1 && this.addressBook[idx]) {
                this.addressBook[idx].label = email;
                this.addressBook[idx].name = label;
                return;
            }

            this.addressBook.unshift({ label: email, name: label });
        },
    },
    getters: {
        hasPrivateEmail: (state) => {
            const { activeChar } = useAuth();
            return !!state.emails.find((e) => e.userId === activeChar.value?.userId);
        },

        getPrivateEmail: (state) => {
            const { activeChar } = useAuth();
            return state.emails.find((e) => e.userId === activeChar.value!.userId);
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
            messages: 'id, threadId',
        });
    }
}

export const mailerDB = new MailerDexie();

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useMailerStore, import.meta.hot));
}
