import type { Email } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
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
    ListThreadsResponse,
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
    selectedEmailId: string | undefined;
    selectedEmail: Email | undefined;
    selectedThread: Thread | undefined;

    unreadCount: number;

    threads: ListThreadsResponse | undefined;
    messages: ListThreadMessagesResponse | undefined;

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
            selectedEmailId: undefined,
            selectedEmail: undefined,
            selectedThread: undefined,

            unreadCount: 0,

            threads: undefined,
            messages: undefined,

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
                console.log('threadUpdate', event.data.threadUpdate);

                // Handle email sent by blocked email
                if (
                    event.data.threadUpdate.creatorEmail?.email &&
                    this.checkIfEmailBlocked(event.data.threadUpdate.creatorEmail?.email)
                ) {
                    // Make sure to set thread state accordingly (locally)
                    await this.setThreadState({
                        threadId: event.data.threadUpdate.id,
                        archived: true,
                        muted: true,
                    });
                    return;
                }

                if (event.data.threadUpdate.creatorEmailId === this.selectedEmail?.id) {
                    return;
                }

                await this.setThreadState({
                    threadId: event.data.threadUpdate.id,
                    unread: true,
                });

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
                useSound().play({ name: 'notification' });
            } else if (event.data.oneofKind === 'threadDelete') {
                const id = event.data.threadDelete;
                if (this.selectedThread?.id === id) {
                    this.selectedThread = undefined;
                }

                // Remove thread if it is currently in our messagess list
                const idx = this.threads?.threads.findIndex((t) => t.id === id);
                if (idx !== undefined && idx > -1) {
                    this.threads?.threads.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'messageUpdate') {
                const data = event.data.messageUpdate;
                // Update thread updated at time
                const thread = this.threads?.threads.find((t) => t.id === data.threadId);
                if (thread) {
                    thread.updatedAt = toTimestamp(new Date());
                }

                // Handle email sent by blocked email
                if (data.sender?.email && this.checkIfEmailBlocked(data.sender?.email)) {
                    // Make sure to set thread state accordingly (locally)
                    await this.setThreadState({
                        threadId: data.threadId,
                        archived: true,
                        muted: true,
                    });
                    return;
                }

                if (data.senderId === this.selectedEmail?.id) {
                    return;
                }

                console.log('messageUpdate', data);

                // Only set unread state when message isn't from same email and the user isn't active on that thread
                const state = await this.setThreadState({
                    threadId: data.threadId,
                    unread: data.threadId !== this.selectedThread?.id,
                });
                if (state?.muted) {
                    return;
                }

                useNotificatorStore().add({
                    title: { key: 'notifications.mailer.new_email.title', parameters: {} },
                    description: {
                        key: 'notifications.mailer.new_email.content',
                        parameters: {
                            title: data.title,
                            from: data.sender?.email ?? 'N/A',
                        },
                    },
                    type: NotificationType.INFO,
                    actions: this.getNotificationActions(data.threadId),
                });
                useSound().play({ name: 'notification' });
            } else if (event.data.oneofKind === 'messageDelete') {
                // Remove message if it is currently in our messagess list
                const id = event.data.messageDelete;
                const idx = this.messages?.messages.findIndex((t) => t.id === id);
                if (idx !== undefined && idx > -1) {
                    this.messages?.messages.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'threadStateUpdate') {
                const state = event.data.threadStateUpdate;
                if (this.selectedEmail?.id !== state.emailId) {
                    return;
                }

                if (this.selectedThread?.id === state.threadId) {
                    this.selectedThread.state = state;
                }

                const thread = this.threads?.threads.find((t) => t.id === state.threadId);
                if (thread) {
                    thread.state = state;
                }
            } else {
                logger.debug('Unknown MailerEvent type received:', event.data.oneofKind);
            }
        },

        async checkEmails(): Promise<void> {
            try {
                if (this.emails.length === 0) {
                    await this.listEmails(true, 0, false);
                }

                if (this.emails.length === 0) {
                    this.unreadCount = 0;
                    return;
                }

                const threads = await this.listThreads(
                    {
                        pagination: {
                            offset: 0,
                        },
                        emailIds: this.emails.map((e) => e.id),
                        unread: true,
                    },
                    false,
                );
                this.unreadCount = threads?.pagination?.totalCount ?? 0;
            } catch (_) {
                /* empty */
            }
        },

        // Emails
        async listEmails(all?: boolean, offset?: number, redirect: boolean = true): Promise<ListEmailsResponse> {
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
                    if (redirect) {
                        await navigateTo({
                            name: 'mail-manage',
                            query: {
                                tab: 'new',
                            },
                            hash: '#',
                        });
                    }
                } else if (this.emails.length > 0) {
                    // Check if previously selected email is available
                    const previousEmail = this.emails.find((e) => e.id === this.selectedEmail?.id);
                    if (previousEmail) {
                        this.selectedEmail = previousEmail;
                    } else if (this.emails[0] && this.emails[0].settings === undefined) {
                        this.selectedEmail = await this.getEmail(this.emails[0].id);
                    } else {
                        this.selectedEmail = this.emails[0];
                    }

                    this.selectedEmailId = this.selectedEmail?.id;
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
        async listThreads(req: ListThreadsRequest, store: boolean = true): Promise<ListThreadsResponse | undefined> {
            try {
                const call = getGRPCMailerClient().listThreads(req);
                const { response } = await call;

                if (response.pagination?.offset === 0) {
                    let count = 0;
                    response.threads.forEach((t) => {
                        if (t.state?.unread !== true) {
                            return;
                        }

                        count++;
                    });
                    this.unreadCount = count;
                }

                if (store) {
                    this.threads = response;

                    // Add selected thread to list to ensure there is no flickering between tab switches
                    if (this.selectedThread) {
                        const thread = response?.threads.filter((t) => t.id === this.selectedThread?.id);
                        if (!thread) {
                            response?.threads.unshift(this.selectedThread);
                        }
                    }
                }

                return response;
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

                if (response.thread && !response.thread.state) {
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

                return response.thread;
            } catch (e) {
                const error = e as RpcError;
                handleGRPCError(error);

                // Switch away from thread if inaccessible
                if (error?.message?.includes('.ErrThreadAccessDenied')) {
                    if (this.selectedThread?.id === threadId) {
                        this.selectedThread = undefined;
                    }
                }
            }
        },

        async createThread(req: CreateThreadRequest): Promise<CreateThreadResponse> {
            try {
                const call = getGRPCMailerClient().createThread(req);
                const { response } = await call;

                if (response.thread) {
                    req.recipients.forEach((r) => this.addToAddressBook(r));

                    this.threads?.threads.unshift(response.thread);
                    if (this.threads?.pagination) {
                        this.threads.pagination.totalCount++;
                    }
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

                if (this.selectedThread?.id === req.threadId) {
                    this.selectedThread = undefined;
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

        // Thread User State
        async getThreadState(threadId: string): Promise<ThreadState | undefined> {
            try {
                const call = getGRPCMailerClient().getThreadState({
                    emailId: this.selectedEmail!.id,
                    threadId: threadId,
                });
                const { response } = await call;

                return response.state;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        },

        async setThreadState(state: Partial<ThreadState>, notify: boolean = false): Promise<ThreadState | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            const { response } = await getGRPCMailerClient().setThreadState({
                state: {
                    threadId: state.threadId!,
                    emailId: this.selectedEmail?.id,
                    ...state,
                },
            });

            if (this.selectedThread && this.selectedThread?.id === state.threadId) {
                this.selectedThread.state = response.state;
            }
            const thread = this.threads?.threads.find((t) => t.id === state.threadId);
            if (thread) {
                thread.state = response.state;
            }

            if (notify) {
                useNotificatorStore().add({
                    title: { key: 'notifications.action_successfull.title', parameters: {} },
                    description: { key: 'notifications.action_successfull.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });
            }

            return response.state;
        },

        // Messages
        async listThreadMessages(req: ListThreadMessagesRequest): Promise<ListThreadMessagesResponse | undefined> {
            if (!this.selectedEmail) {
                return;
            }

            try {
                const call = getGRPCMailerClient().listThreadMessages(req);
                const { response } = await call;

                this.messages = response;

                return response;
            } catch (e) {
                const error = e as RpcError;
                await handleGRPCError(error);

                // Switch away from thread if inaccessible
                if (error?.message?.includes('.ErrThreadAccessDenied')) {
                    if (this.selectedThread?.id === req.threadId) {
                        this.selectedThread = undefined;
                    }
                }
            }
        },

        async postMessage(req: PostMessageRequest): Promise<PostMessageResponse> {
            try {
                const call = getGRPCMailerClient().postMessage(req);
                const { response } = await call;

                if (response.message) {
                    req.recipients.forEach((r) => this.addToAddressBook(r));

                    this.messages?.messages?.unshift(response.message);
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
            if (!this.selectedEmail?.settings?.blockedEmails) {
                return false;
            }

            return this.selectedEmail.settings.blockedEmails.includes(email.toLowerCase());
        },

        getNotificationActions(threadId?: string): NotificationActionI18n[] {
            return [
                {
                    label: { key: 'common.click_here' },
                    to: threadId ? { name: 'mail', query: { thread: threadId }, hash: '#' } : { name: 'mail' },
                },
            ];
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

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useMailerStore, import.meta.hot));
}
