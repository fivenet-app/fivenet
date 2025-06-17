import type { Email } from '~~/gen/ts/resources/mailer/email';
import type { MailerEvent } from '~~/gen/ts/resources/mailer/events';
import type { MessageAttachment } from '~~/gen/ts/resources/mailer/message';
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
        const { $grpc } = useNuxtApp();
        const notifications = useNotificationsStore();

        // State
        const loaded = ref<boolean>(false);
        const error = ref<Error | undefined>(undefined);

        const draft = ref({
            title: '',
            content: '',
            recipients: [] as { label: string }[],
            attachments: [] as MessageAttachment[],
        });

        const emails = ref<Email[]>([]);
        const selectedEmailId = ref<number | undefined>(undefined);
        const selectedEmail = ref<Email | undefined>(undefined);
        const selectedThread = ref<Thread | undefined>(undefined);

        const unreadThreadIds = ref<number[]>([]);

        const threads = ref<ListThreadsResponse | undefined>(undefined);
        const messages = ref<ListThreadMessagesResponse | undefined>(undefined);

        const addressBook = ref<{ label: string; name?: string }[]>([]);

        const notificationSound = useSounds('/sounds/notification.mp3');

        // Actions
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
                const data = event.data.threadUpdate;
                console.debug('threadUpdate', data);

                // Handle email sent by blocked email
                if (data.creatorEmail?.email && checkIfEmailBlocked(data.creatorEmail?.email)) {
                    // Make sure to set thread state accordingly (locally)
                    await setThreadState({
                        threadId: data.id,
                        archived: true,
                        muted: true,
                    });
                    return;
                }

                // Either creator id or email address matches
                if (
                    data.creatorEmailId === selectedEmail.value?.id ||
                    data.creatorEmail?.email === selectedEmail.value?.email
                ) {
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

                await setThreadState({
                    threadId: data.id,
                    unread: true,
                });

                // Update thread order in list
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
                const data = event.data.messageUpdate;
                // Update thread updatedAt time and move to beginning of list
                const threadIdx = threads.value?.threads.findIndex((t) => t.id === data.threadId);
                if (threadIdx !== undefined && threadIdx > -1) {
                    const thread = threads.value!.threads[threadIdx]!;
                    thread.updatedAt = toTimestamp(new Date());

                    threads.value!.threads.splice(threadIdx, 1);
                    threads.value!.threads.unshift(thread);
                }

                if (selectedThread.value?.id === data.threadId) {
                    selectedThread.value.updatedAt = toTimestamp(new Date());

                    messages.value?.messages?.unshift(data);
                }

                console.debug('messageUpdate', data);

                // Handle email sent by blocked email
                if (data.sender?.email && checkIfEmailBlocked(data.sender?.email)) {
                    // Make sure to set thread state accordingly (locally)
                    await setThreadState({
                        threadId: data.threadId,
                        archived: true,
                        muted: true,
                    });
                    return;
                }

                if (data.senderId === selectedEmail.value?.id) {
                    return;
                }

                // Only set unread state when message isn't from same email and the user isn't active on that thread
                const state = await setThreadState({
                    threadId: data.threadId,
                    unread: data.threadId !== selectedThread.value?.id,
                });
                if (state?.muted) {
                    return;
                }

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

                if (selectedEmail.value?.id !== newState.emailId) {
                    return;
                }

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

        const checkEmails = async (): Promise<void> => {
            try {
                if (emails.value.length === 0) {
                    // Reset unread thread ids list
                    unreadThreadIds.value.length = 0;
                    await listEmails(true, 0, false);
                }

                // Still no email addresses? Return here.
                if (emails.value.length === 0) {
                    return;
                }

                const { isSuperuser } = useAuth();
                // If is superuser and doesn't have a private email to check
                if (isSuperuser.value && getPrivateEmail.value === undefined) {
                    return;
                }

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
        const listEmails = async (
            all: boolean = false,
            offset: number = 0,
            redirect: boolean = true,
        ): Promise<ListEmailsResponse> => {
            error.value = undefined;

            if (addressBook.value.length > 30) {
                addressBook.value.length = 30;
            }

            try {
                const call = $grpc.mailer.mailer.listEmails({
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

        const getEmail = async (id: number): Promise<Email | undefined> => {
            try {
                const call = $grpc.mailer.mailer.getEmail({
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

        const createOrUpdateEmail = async (req: CreateOrUpdateEmailRequest): Promise<CreateOrUpdateEmailResponse> => {
            try {
                const call = $grpc.mailer.mailer.createOrUpdateEmail(req);
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

        const deleteEmail = async (req: DeleteEmailRequest): Promise<DeleteEmailResponse> => {
            try {
                const call = $grpc.mailer.mailer.deleteEmail(req);
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
        const listThreads = async (
            req: ListThreadsRequest,
            store: boolean = true,
        ): Promise<ListThreadsResponse | undefined> => {
            try {
                const call = $grpc.mailer.mailer.listThreads(req);
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

        const getThread = async (threadId: number): Promise<Thread | undefined> => {
            if (!selectedEmail.value) {
                return;
            }

            try {
                const call = $grpc.mailer.mailer.getThread({
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

        const createThread = async (req: CreateThreadRequest): Promise<CreateThreadResponse> => {
            try {
                const call = $grpc.mailer.mailer.createThread(req);
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

        const deleteThread = async (req: DeleteThreadRequest): Promise<DeleteThreadResponse> => {
            try {
                const call = $grpc.mailer.mailer.deleteThread(req);
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

        // Thread User State
        const getThreadState = async (threadId: number): Promise<ThreadState | undefined> => {
            try {
                const call = $grpc.mailer.mailer.getThreadState({
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

        const setThreadState = async (
            state: Partial<ThreadState>,
            notify: boolean = false,
        ): Promise<ThreadState | undefined> => {
            if (!selectedEmail.value) {
                return;
            }

            const { response } = await $grpc.mailer.mailer.setThreadState({
                state: {
                    ...state,
                    threadId: state.threadId!,
                    emailId: selectedEmail.value.id,
                },
            });

            if (selectedThread.value && selectedThread.value?.id === state.threadId) {
                selectedThread.value.state = response.state;

                // Reset selected thread when archived
                if (state.archived === true) {
                    selectedThread.value = undefined;
                }
            }

            const thread = threads.value?.threads.find((t) => t.id === state.threadId);
            if (thread) {
                thread.state = response.state;
            }

            // Add/Remove thread from unreadThreadIds
            const unreadThreadIdx = unreadThreadIds.value.findIndex((t) => t === state.threadId);
            if (!response.state?.unread) {
                if (unreadThreadIdx > -1) {
                    unreadThreadIds.value.splice(unreadThreadIdx, 1);
                }
            } else {
                if (unreadThreadIdx === -1) {
                    unreadThreadIds.value.push(response.state.threadId);
                }
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
        const listThreadMessages = async (req: ListThreadMessagesRequest): Promise<ListThreadMessagesResponse | undefined> => {
            if (!selectedEmail.value) {
                return;
            }

            try {
                const call = $grpc.mailer.mailer.listThreadMessages(req);
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

        const postMessage = async (req: PostMessageRequest): Promise<PostMessageResponse> => {
            try {
                const call = $grpc.mailer.mailer.postMessage(req);
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

        const deleteMessage = async (req: DeleteMessageRequest): Promise<DeleteMessageResponse> => {
            try {
                const call = $grpc.mailer.mailer.deleteMessage(req);
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
        const getEmailSettings = async (req: GetEmailSettingsRequest): Promise<GetEmailSettingsResponse> => {
            try {
                const call = $grpc.mailer.mailer.getEmailSettings(req);
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

        const setEmailSettings = async (req: SetEmailSettingsRequest): Promise<SetEmailSettingsResponse> => {
            try {
                const call = $grpc.mailer.mailer.setEmailSettings(req);
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
        const checkIfEmailBlocked = (emailAddress: string): boolean => {
            if (!selectedEmail.value?.settings?.blockedEmails) {
                return false;
            }
            return selectedEmail.value.settings.blockedEmails.includes(emailAddress.toLowerCase());
        };

        const getNotificationActions = (threadId?: number): NotificationActionI18n[] => {
            return [
                {
                    label: { key: 'common.click_here' },
                    to: threadId ? { name: 'mail', query: { thread: threadId }, hash: '#' } : { name: 'mail' },
                },
            ];
        };

        // Address book
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
