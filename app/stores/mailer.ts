import type { NotificationActionI18n } from '~/types/notifications';
import { getMailerMailerClient } from '~~/gen/ts/clients';
import type { Email } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
import type { Message, MessageAttachment } from '~~/gen/ts/resources/mailer/message';
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

const logger = useLogger('💬 Mailer');

export const useMailerStore = defineStore(
    'mailer',
    () => {
        const notifications = useNotificationsStore();

        // State
        /**
         * Indicates whether the mailer store has finished loading its data.
         */
        const loaded = ref<boolean>(false);

        /**
         * Stores any error encountered during operations in the mailer store.
         */
        const error = ref<Error | undefined>(undefined);

        /**
         * Holds the current email draft being composed.
         * @property {string} title - The title of the email draft.
         * @property {string} content - The content of the email draft.
         * @property {Array<{label: string}>} recipients - The list of recipients for the email draft.
         * @property {Array<MessageAttachment>} attachments - The list of attachments for the email draft.
         */
        const draft = ref({
            title: '',
            content: '',
            recipients: [] as { label: string }[],
            attachments: [] as MessageAttachment[],
        });

        /**
         * Contains the list of email accounts available to the user.
         */
        const emails = ref<Email[]>([]);

        /**
         * The ID of the currently selected email account.
         */
        const selectedEmailId = ref<number | undefined>(undefined);

        /**
         * The currently selected email account object.
         */
        const selectedEmail = ref<Email | undefined>(undefined);

        /**
         * The currently selected email thread.
         */
        const selectedThread = ref<Thread | undefined>(undefined);

        /**
         * Keeps track of thread IDs that have unread messages.
         */
        const unreadThreadIds = ref<number[]>([]);

        /**
         * Contains the list of threads for the selected email account.
         */
        const threads = ref<ListThreadsResponse | undefined>(undefined);

        /**
         * Contains the list of messages in the currently selected thread.
         */
        const messages = ref<ListThreadMessagesResponse | undefined>(undefined);

        /**
         * Stores a list of email addresses and their optional names.
         * @property {string} label - The email address.
         * @property {string} [name] - The optional name associated with the email address.
         */
        const addressBook = ref<{ label: string; name?: string }[]>([]);

        const notificationSound = useSounds('/sounds/notification.mp3');

        /**
         * Handle updates to a thread.
         *
         * @param {Thread} data - The thread data to process.
         */
        const handleThreadUpdate = async (data: Thread): Promise<void> => {
            logger.debug('threadUpdate', data);

            if (data.creatorEmail?.email && checkIfEmailBlocked(data.creatorEmail?.email)) {
                await setThreadState({ threadId: data.id, archived: true, muted: true });
                return;
            }

            if (data.creatorEmailId === selectedEmail.value?.id || data.creatorEmail?.email === selectedEmail.value?.email) {
                if (data.state) {
                    data.state.unread = false;
                } else {
                    data.state = {
                        emailId: selectedEmail.value?.id ?? 0,
                        threadId: data.id,
                        unread: false,
                    };
                }
                return;
            }

            await setThreadState({ threadId: data.id, unread: true });

            const threadIdx = threads.value?.threads.findIndex((t) => t.id === data.id);
            if (threadIdx !== undefined && threadIdx > -1) {
                const thread = threads.value!.threads[threadIdx]!;
                threads.value!.threads.splice(threadIdx, 1);
                threads.value!.threads.unshift(thread);
            }

            notifications.add({
                title: { key: 'notifications.mailer.new_email.title', parameters: {} },
                description: {
                    key: 'notifications.mailer.new_email.content',
                    parameters: {
                        title: data.title,
                        from: data.creatorEmail?.email ?? 'N/A',
                    },
                },
                type: NotificationType.INFO,
                actions: getNotificationActions(data.id),
            });
            notificationSound.play();
        };

        /**
         * Handle updates to a message.
         *
         * @param {Message} data - The message data to process.
         */
        const handleMessageUpdate = async (data: Message): Promise<void> => {
            const threadIdx = threads.value?.threads.findIndex((t) => t.id === data.threadId);
            if (threadIdx !== undefined && threadIdx > -1) {
                const thread = threads.value!.threads[threadIdx]!;
                thread.updatedAt = toTimestamp(new Date());
                threads.value!.threads.splice(threadIdx, 1);
                threads.value!.threads.unshift(thread);
            }

            if (selectedThread.value?.id === data.threadId) {
                selectedThread.value!.updatedAt = toTimestamp(new Date());
                messages.value?.messages?.unshift({
                    id: 0, // Placeholder ID
                    threadId: data.threadId,
                    senderId: data.senderId ?? 0,
                    title: data.title,
                });
            }

            logger.debug('messageUpdate', data);

            // Handle email sent by blocked email
            if (data.sender?.email && checkIfEmailBlocked(data.sender?.email)) {
                // Make sure to set thread state accordingly (locally)
                await setThreadState({ threadId: data.threadId, archived: true, muted: true });
                return;
            }

            if (data.senderId === selectedEmail.value?.id) return;

            // Only set unread state when message isn't from same email and the user isn't active on that thread
            const state = await setThreadState({ threadId: data.threadId, unread: data.threadId !== selectedThread.value?.id });
            if (state?.muted) return;

            notifications.add({
                title: { key: 'notifications.mailer.new_email.title', parameters: {} },
                description: {
                    key: 'notifications.mailer.new_email.content',
                    parameters: {
                        title: data.title,
                        from: data.sender?.email ?? 'N/A',
                    },
                },
                type: NotificationType.INFO,
                actions: getNotificationActions(data.threadId),
            });
            notificationSound.play();
        };

        // Actions
        // `handleEvent` processes incoming mailer events and updates the store accordingly.
        /**
         * Process incoming mailer events and update the store accordingly.
         *
         * @param {MailerEvent} event - The mailer event to handle.
         */
        const handleEvent = async (event: MailerEvent): Promise<void> => {
            logger.debug('Received change - oneofKind:', event.data.oneofKind, event.data);

            if (event.data.oneofKind === 'emailUpdate') {
                const searchId = event.data.emailUpdate.id;
                const idx = emails.value.findIndex((e) => e.id === searchId);
                if (idx > -1) {
                    emails.value[idx] = event.data.emailUpdate;
                }
            } else if (event.data.oneofKind === 'emailDelete') {
                const searchId = event.data.emailDelete;
                const idx = emails.value.findIndex((e) => e.id === searchId);
                if (idx > -1) {
                    emails.value.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'emailSettingsUpdated') {
                const searchId = event.data.emailSettingsUpdated.emailId;
                const idx = emails.value.findIndex((e) => e.id === searchId);
                if (idx > -1 && emails.value[idx]) {
                    emails.value[idx].settings = event.data.emailSettingsUpdated;
                }
            } else if (event.data.oneofKind === 'threadUpdate') {
                await handleThreadUpdate(event.data.threadUpdate);
            } else if (event.data.oneofKind === 'threadDelete') {
                const id = event.data.threadDelete;
                if (selectedThread.value?.id === id) {
                    selectedThread.value = undefined;
                }

                // Remove thread if it is currently in our threads list
                const idx = threads.value?.threads.findIndex((t) => t.id === id);
                if (idx !== undefined && idx > -1) {
                    threads.value?.threads.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'messageUpdate') {
                await handleMessageUpdate(event.data.messageUpdate);
            } else if (event.data.oneofKind === 'messageDelete') {
                // Remove message if it is currently in our messages list
                const id = event.data.messageDelete;
                const idx = messages.value?.messages.findIndex((t) => t.id === id);
                if (idx !== undefined && idx > -1) {
                    messages.value?.messages.splice(idx, 1);
                }
            } else if (event.data.oneofKind === 'threadStateUpdate') {
                const newState = event.data.threadStateUpdate;

                // Add/Remove thread from unreadThreadIds
                const unreadThreadIdx = unreadThreadIds.value.findIndex((t) => t === newState.threadId);
                if (!newState.unread) {
                    if (unreadThreadIdx > -1) {
                        unreadThreadIds.value.splice(unreadThreadIdx, 1);
                    }
                } else {
                    if (unreadThreadIdx === -1) {
                        unreadThreadIds.value.push(newState.threadId);
                    }
                }

                if (selectedEmail.value?.id !== newState.emailId) return;

                if (selectedThread.value?.id === newState.threadId) {
                    selectedThread.value.state = newState;
                }

                const thread = threads.value?.threads.find((t) => t.id === newState.threadId);
                if (thread) {
                    thread.state = newState;
                }
            } else {
                logger.debug('Unknown MailerEvent type received:', event.data.oneofKind);
            }
        };

        /**
         * Check and update the list of emails.
         */
        const checkEmails = async (): Promise<void> => {
            try {
                if (emails.value.length === 0) {
                    // Reset unread thread ids list
                    unreadThreadIds.value.length = 0;
                    await listEmails(true, 0, false);
                }

                // Still no email addresses? Return here.
                if (emails.value.length === 0) return;

                const { isSuperuser } = useAuth();
                // If is superuser and doesn't have a private email to check
                if (isSuperuser.value && getPrivateEmail.value === undefined) return;

                const emailIds = isSuperuser.value ? [getPrivateEmail.value!.id] : emails.value.map((e) => e.id);
                // Truncate email ids to 10 if needed
                emailIds.length = Math.min(emailIds.length, 10);

                // Load unread threads for all emails
                const threadsResponse = await listThreads(
                    {
                        pagination: {
                            offset: 0,
                        },
                        emailIds: emailIds,
                        unread: true,
                    },
                    false,
                );

                unreadThreadIds.value = threadsResponse?.threads.map((t) => t.id) ?? [];
            } catch (_) {
                // empty
            }
        };

        // Emails
        /**
         * Fetch the list of email accounts and update the store.
         *
         * @param {boolean} [all=false] - Whether to fetch all email accounts.
         * @param {number} [offset=0] - The pagination offset.
         * @param {boolean} [redirect=true] - Whether to redirect if no emails are found.
         * @returns {Promise<ListEmailsResponse>} - The response containing the list of emails.
         */
        const listEmails = async (
            all: boolean = false,
            offset: number = 0,
            redirect: boolean = true,
        ): Promise<ListEmailsResponse> => {
            error.value = undefined;

            if (addressBook.value.length > 30) {
                addressBook.value.length = 30;
            }

            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.listEmails({
                    pagination: {
                        offset,
                    },
                    all,
                });
                const { response } = await call;

                emails.value = response.emails;
                if (emails.value.length === 0 || !hasPrivateEmail.value) {
                    if (redirect) {
                        await navigateTo({
                            name: 'mail-manage',
                            query: {
                                tab: 'new',
                            },
                            hash: '#',
                        });
                    }
                } else if (emails.value.length > 0) {
                    // Check if previously selected email is available
                    const previousEmail = emails.value.find((e) => e.id === selectedEmail.value?.id);
                    if (previousEmail) {
                        selectedEmail.value = previousEmail;
                    } else if (emails.value[0] && emails.value[0].settings === undefined) {
                        selectedEmail.value = await getEmail(emails.value[0].id);
                    } else {
                        selectedEmail.value = emails.value[0];
                    }

                    selectedEmailId.value = selectedEmail.value?.id;
                }

                loaded.value = true;
                return response;
            } catch (e) {
                error.value = e as RpcError;
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Fetch details of a specific email by its ID.
         *
         * @param {number} id - The ID of the email to fetch.
         * @returns {Promise<Email | undefined>} - The email details, if found.
         */
        const getEmail = async (id: number): Promise<Email | undefined> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.getEmail({
                    id,
                });
                const { response } = await call;

                const emailObj = emails.value.find((e) => e.id === id);
                if (emailObj) {
                    emailObj.settings = response.email?.settings;
                    emailObj.access = response.email?.access;
                }
                if (selectedEmail.value && selectedEmail.value.id === response.email?.id) {
                    selectedEmail.value.settings = response.email.settings;
                    selectedEmail.value.access = response.email.access;
                }

                return response.email;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Create or update an email account.
         *
         * @param {CreateOrUpdateEmailRequest} req - The request data for creating or updating the email.
         * @returns {Promise<CreateOrUpdateEmailResponse>} - The response containing the created or updated email.
         */
        const createOrUpdateEmail = async (req: CreateOrUpdateEmailRequest): Promise<CreateOrUpdateEmailResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.createOrUpdateEmail(req);
                const { response } = await call;

                if (response.email) {
                    const idx = emails.value.findIndex((e) => e.id === response.email!.id);
                    if (idx === -1) {
                        emails.value.unshift(response.email);
                    } else {
                        emails.value[idx] = response.email;
                    }

                    if (!selectedEmail.value) {
                        selectedEmail.value = response.email;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Delete an email account.
         *
         * @param {DeleteEmailRequest} req - The request data for deleting the email.
         * @returns {Promise<DeleteEmailResponse>} - The response confirming the deletion.
         */
        const deleteEmail = async (req: DeleteEmailRequest): Promise<DeleteEmailResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.deleteEmail(req);
                const { response } = await call;

                if (selectedEmail.value?.id === req.id) {
                    selectedEmail.value = undefined;
                }

                const idx = emails.value.findIndex((e) => e.id === req.id);
                if (idx > -1) {
                    emails.value.splice(idx, 1);
                }

                notifications.restartStream();

                notifications.add({
                    title: { key: 'notifications.action_successful.title', parameters: {} },
                    description: { key: 'notifications.action_successful.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Threads
        // `listThreads` fetches the list of threads for the selected email account.
        /**
         * Fetch the list of threads for the selected email account.
         *
         * @param {ListThreadsRequest} req - The request data for fetching threads.
         * @param {boolean} [store=true] - Whether to store the fetched threads in the store.
         * @returns {Promise<ListThreadsResponse | undefined>} - The response containing the list of threads.
         */
        const listThreads = async (
            req: ListThreadsRequest,
            store: boolean = true,
        ): Promise<ListThreadsResponse | undefined> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.listThreads(req);
                const { response } = await call;

                // If response is at offset 0 and request is not for archived threads, update unread threads list
                if (response.pagination?.offset === 0 && !!req.unread) {
                    for (let i = 0; i < response.threads.length; i++) {
                        const thread = response.threads[i]!;

                        // Keep unreadThreadIds up to date
                        const idx = unreadThreadIds.value.findIndex((t) => t === thread.id);
                        if (thread.state?.unread !== true) {
                            if (idx > -1) {
                                unreadThreadIds.value.splice(idx, 1);
                            }
                            continue;
                        }

                        if (idx === -1) {
                            unreadThreadIds.value.push(thread.id);
                        }
                    }
                }

                if (store) {
                    threads.value = response;

                    // Add selected thread to list to ensure there is no flickering between tab switches
                    if (selectedThread.value) {
                        const existing = response?.threads.filter((t) => t.id === selectedThread.value?.id);
                        if (!existing || existing.length === 0) {
                            response?.threads.unshift(selectedThread.value);
                        }
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Fetch details of a specific thread by its ID.
         *
         * @param {number} threadId - The ID of the thread to fetch.
         * @returns {Promise<Thread | undefined>} - The thread details, if found.
         */
        const getThread = async (threadId: number): Promise<Thread | undefined> => {
            if (!selectedEmail.value) return;

            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.getThread({
                    emailId: selectedEmail.value.id,
                    threadId,
                });
                const { response } = await call;

                if (response.thread && !response.thread.state) {
                    response.thread.state = {
                        emailId: selectedEmail.value.id,
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
                    if (selectedThread.value?.id === threadId) {
                        selectedThread.value = undefined;
                    }
                }
            }
        };

        /**
         * Create a new thread.
         *
         * @param {CreateThreadRequest} req - The request data for creating the thread.
         * @returns {Promise<CreateThreadResponse>} - The response containing the created thread.
         */
        const createThread = async (req: CreateThreadRequest): Promise<CreateThreadResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.createThread(req);
                const { response } = await call;

                if (response.thread) {
                    req.recipients.forEach((r) => addToAddressBook(r));

                    threads.value?.threads.unshift(response.thread);
                    if (threads.value?.pagination) {
                        threads.value.pagination.totalCount++;
                    }
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Delete a thread.
         *
         * @param {DeleteThreadRequest} req - The request data for deleting the thread.
         * @returns {Promise<DeleteThreadResponse>} - The response confirming the deletion.
         */
        const deleteThread = async (req: DeleteThreadRequest): Promise<DeleteThreadResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.deleteThread(req);
                const { response } = await call;

                if (selectedThread.value?.id === req.threadId) {
                    selectedThread.value = undefined;
                }

                notifications.add({
                    title: { key: 'notifications.action_successful.title', parameters: {} },
                    description: { key: 'notifications.action_successful.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Fetch the state of a specific thread by its ID.
         *
         * @param {number} threadId - The ID of the thread to fetch the state for.
         * @returns {Promise<ThreadState | undefined>} - The thread state, if found.
         */
        const getThreadState = async (threadId: number): Promise<ThreadState | undefined> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.getThreadState({
                    emailId: selectedEmail.value!.id,
                    threadId,
                });
                const { response } = await call;

                return response.state;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Update the state of a thread locally.
         *
         * @param {number} threadId - The ID of the thread to update.
         * @param {ThreadState} newState - The new state to apply to the thread.
         */
        const updateThreadState = (threadId: number, newState: ThreadState): void => {
            const thread = threads.value?.threads.find((t) => t.id === threadId);
            if (thread) {
                thread.state = { ...thread.state, ...newState } as ThreadState;
            }

            if (selectedThread.value?.id === threadId) {
                selectedThread.value.state = { ...selectedThread.value.state, ...newState };

                // Reset selected thread when archived
                if (newState.archived === true) {
                    selectedThread.value = undefined;
                }
            }
        };

        /**
         * Update the list of unread thread IDs.
         *
         * @param {number} threadId - The ID of the thread to update.
         * @param {boolean} unread - Whether the thread is unread.
         */
        const updateUnreadThreadIds = (threadId: number, unread: boolean): void => {
            const idx = unreadThreadIds.value.findIndex((id) => id === threadId);
            if (unread && idx === -1) {
                unreadThreadIds.value.push(threadId);
            } else if (!unread && idx > -1) {
                unreadThreadIds.value.splice(idx, 1);
            }
        };

        /**
         * Set the state of a thread.
         *
         * @param {Partial<ThreadState>} state - The partial state to apply to the thread.
         * @param {boolean} [notify=false] - Whether to show a notification after updating the state.
         * @returns {Promise<ThreadState | undefined>} - The updated thread state, if successful.
         */
        const setThreadState = async (
            state: Partial<ThreadState>,
            notify: boolean = false,
        ): Promise<ThreadState | undefined> => {
            if (!selectedEmail.value) return;

            const mailerMailerClient = await getMailerMailerClient();

            const { response } = await mailerMailerClient.setThreadState({
                state: {
                    ...state,
                    threadId: state.threadId!,
                    emailId: selectedEmail.value.id,
                },
            });

            if (response.state) {
                updateThreadState(state.threadId!, response.state);
                updateUnreadThreadIds(state.threadId!, response.state.unread ?? false);
            }

            if (notify) {
                notifications.add({
                    title: { key: 'notifications.action_successful.title', parameters: {} },
                    description: { key: 'notifications.action_successful.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });
            }

            return response.state;
        };

        // Messages
        /**
         * Fetch the list of messages for a specific thread.
         *
         * @param {ListThreadMessagesRequest} req - The request data for fetching thread messages.
         * @returns {Promise<ListThreadMessagesResponse | undefined>} - The response containing the list of messages.
         */
        const listThreadMessages = async (req: ListThreadMessagesRequest): Promise<ListThreadMessagesResponse | undefined> => {
            if (!selectedEmail.value) return;

            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.listThreadMessages(req);
                const { response } = await call;

                messages.value = response;

                return response;
            } catch (e) {
                const error = e as RpcError;
                await handleGRPCError(error);

                // Switch away from thread if inaccessible
                if (error?.message?.includes('.ErrThreadAccessDenied')) {
                    if (selectedThread.value?.id === req.threadId) {
                        selectedThread.value = undefined;
                    }
                }
            }
        };

        /**
         * Post a new message to a thread.
         *
         * @param {PostMessageRequest} req - The request data for posting the message.
         * @returns {Promise<PostMessageResponse>} - The response containing the posted message.
         */
        const postMessage = async (req: PostMessageRequest): Promise<PostMessageResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.postMessage(req);
                const { response } = await call;

                if (response.message) {
                    req.recipients.forEach((r) => addToAddressBook(r));

                    messages.value?.messages?.unshift(response.message);
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Delete a message from a thread.
         *
         * @param {DeleteMessageRequest} req - The request data for deleting the message.
         * @returns {Promise<DeleteMessageResponse>} - The response confirming the deletion.
         */
        const deleteMessage = async (req: DeleteMessageRequest): Promise<DeleteMessageResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.deleteMessage(req);
                const { response } = await call;

                notifications.add({
                    title: { key: 'notifications.action_successful.title', parameters: {} },
                    description: { key: 'notifications.action_successful.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // User Settings
        /**
         * Fetch the email settings for a specific email account.
         *
         * @param {GetEmailSettingsRequest} req - The request data for fetching email settings.
         * @returns {Promise<GetEmailSettingsResponse>} - The response containing the email settings.
         */
        const getEmailSettings = async (req: GetEmailSettingsRequest): Promise<GetEmailSettingsResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.getEmailSettings(req);
                const { response } = await call;

                if (response.settings && selectedEmail.value) {
                    selectedEmail.value.settings = response.settings;
                }

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        /**
         * Update the email settings for a specific email account.
         *
         * @param {SetEmailSettingsRequest} req - The request data for updating email settings.
         * @returns {Promise<SetEmailSettingsResponse>} - The response confirming the update.
         */
        const setEmailSettings = async (req: SetEmailSettingsRequest): Promise<SetEmailSettingsResponse> => {
            const mailerMailerClient = await getMailerMailerClient();

            try {
                const call = mailerMailerClient.setEmailSettings(req);
                const { response } = await call;

                if (response.settings && selectedEmail.value) {
                    selectedEmail.value.settings = response.settings;
                }

                notifications.add({
                    title: { key: 'notifications.action_successful.title', parameters: {} },
                    description: { key: 'notifications.action_successful.content', parameters: {} },
                    type: NotificationType.SUCCESS,
                });

                return response;
            } catch (e) {
                handleGRPCError(e as RpcError);
                throw e;
            }
        };

        // Utility
        /**
         * Check if an email address is blocked.
         *
         * @param {string} emailAddress - The email address to check.
         * @returns {boolean} - Whether the email address is blocked.
         */
        const checkIfEmailBlocked = (emailAddress: string): boolean => {
            if (!selectedEmail.value?.settings?.blockedEmails) {
                return false;
            }
            return selectedEmail.value.settings.blockedEmails.includes(emailAddress.toLowerCase());
        };

        /**
         * Get notification actions for a specific thread.
         *
         * @param {number} [threadId] - The ID of the thread to generate actions for.
         * @returns {NotificationActionI18n[]} - The list of notification actions.
         */
        const getNotificationActions = (threadId?: number): NotificationActionI18n[] => {
            return [
                {
                    label: { key: 'common.click_here' },
                    to: threadId ? { name: 'mail-thread', query: { thread: threadId }, hash: '#' } : { name: 'mail-thread' },
                },
            ];
        };

        // Address book
        /**
         * Add an email address to the address book.
         *
         * @param {string} emailAddress - The email address to add.
         * @param {string} [label] - An optional label for the email address.
         */
        const addToAddressBook = (emailAddress: string, label?: string): void => {
            const email = emailAddress.trim();
            const name = label?.trim();

            const idx = addressBook.value.findIndex((a) => a.label === email);
            if (idx > -1 && addressBook.value[idx]) {
                addressBook.value[idx].label = email;
                addressBook.value[idx].name = name;
                return;
            }

            addressBook.value.unshift({ label: email, name });
        };

        // Getters
        const hasPrivateEmail = computed<boolean>(() => {
            const { activeChar } = useAuth();
            return !!emails.value.find((e) => e.userId === activeChar.value?.userId);
        });

        const getPrivateEmail = computed<Email | undefined>(() => {
            const { activeChar } = useAuth();
            return emails.value.find((e) => e.userId === activeChar.value!.userId);
        });

        const unreadCount = computed<number>(() => unreadThreadIds.value.length);

        return {
            // State
            loaded,
            error,
            draft,
            emails,
            selectedEmailId,
            selectedEmail,
            selectedThread,
            unreadThreadIds,
            threads,
            messages,
            addressBook,

            // Actions
            handleEvent,
            checkEmails,
            listEmails,
            getEmail,
            createOrUpdateEmail,
            deleteEmail,
            listThreads,
            getThread,
            createThread,
            deleteThread,
            getThreadState,
            setThreadState,
            listThreadMessages,
            postMessage,
            deleteMessage,
            getEmailSettings,
            setEmailSettings,
            checkIfEmailBlocked,
            getNotificationActions,
            addToAddressBook,

            // Getters
            hasPrivateEmail,
            getPrivateEmail,
            unreadCount,
        };
    },
    {
        persist: {
            pick: ['draft', 'addressBook'],
        },
    },
);

if (import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(useMailerStore, import.meta.hot));
}
