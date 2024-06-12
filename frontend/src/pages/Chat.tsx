import { Button } from '@/components/ui/button';
import { Skeleton } from '@/components/ui/skeleton';
import { Textarea } from '@/components/ui/textarea';
import { useToast } from '@/components/ui/use-toast';
import { useStores } from '@/hooks/useStores';
import { ArrowUpIcon, FilePenIcon } from 'lucide-react';
import { observer } from 'mobx-react-lite';
import { ChangeEvent, KeyboardEvent, useEffect, useRef, useState } from 'react';
import debounce from 'lodash/debounce';
import { useNavigate, useParams } from 'react-router-dom';
import { Pages } from '@/router/constants';
import Conversation from '@/components/Conversation';

const Chat = observer(() => {
    const { rootStore } = useStores();
    const { toast } = useToast();
    const { sessionId } = useParams();
    const navigate = useNavigate();
    const [message, setMessage] = useState('');

    const titleInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        console.log('sessionId', sessionId);

        if (!sessionId) {
            rootStore
                .createSession()
                .then(() => {
                    rootStore.getSessions();
                })
                .catch(() => {
                    toast({
                        title: 'Ошибка',
                        description: 'Не удалось создать сессию',
                        variant: 'destructive',
                    });
                });
        } else {
            rootStore.setActiveSessionId(sessionId);

            rootStore.getSession({ id: sessionId }).catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось загрузить сессию',
                    variant: 'destructive',
                });

                navigate(`/${Pages.Chat}`, { replace: true });
            });
        }
    }, [rootStore, toast, sessionId, navigate]);

    useEffect(() => {
        rootStore.getSessions().catch(() => {
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить сессии',
                variant: 'destructive',
            });
        });
    }, [rootStore, toast]);

    const handleKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            sendMessage();
        }
    };

    const sendMessage = () => {
        if (message.trim()) {
            rootStore.sendMessage({
                prompt: message.trim(),
            });
            setMessage('');
        }
    };

    const debouncedRenameSession = debounce((sessionId: string, title: string) => {
        rootStore
            .renameSession({ id: sessionId, title })
            .then(() => {
                toast({
                    title: 'Успех',
                    description: 'Сессия успешно переименована',
                    variant: 'default',
                });
            })
            .catch(() => {
                toast({
                    title: 'Ошибка',
                    description: 'Не удалось переименовать сессию',
                    variant: 'destructive',
                });
            });
    }, 1000);

    const onTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
        const title = event.target?.value;
        if (rootStore.activeSessionId) {
            debouncedRenameSession(rootStore.activeSessionId, title);
        }
    };

    return (
        <>
            <div className='flex items-center'>
                <div className='flex items-center gap-2 group'>
                    <input
                        ref={titleInputRef}
                        type='text'
                        className='bg-transparent text-lg font-medium focus:outline-none'
                        defaultValue={rootStore.activeSession?.title || 'Новый чат'}
                        onChange={(event) => onTitleChange(event)}
                    />
                    <Button
                        variant='ghost'
                        size='icon'
                        className='rounded-full hover:bg-gray-100 dark:hover:bg-[#1e293b] transition-colors'
                        onClick={() => {
                            titleInputRef.current?.focus();
                        }}
                    >
                        <FilePenIcon className='w-5 h-5 text-gray-500 dark:text-gray-400 group-hover:text-gray-700 dark:group-hover:text-gray-300' />
                    </Button>
                </div>
            </div>

            <div className='flex flex-col h-screen'>
                <div className='max-w-5xl w-full flex-1 mx-auto flex flex-col items-start gap-8 px-4 py-8'>
                    {rootStore.activeSessionLoading
                        ? Array.from({ length: 3 }).map((_, i) => (
                              <div key={i} className='flex items-start gap-4 animate-pulse w-full'>
                                  <Skeleton className='bg-slate-200 w-8 h-8 rounded-full' />
                                  <div className='grid gap-1 flex-1'>
                                      <Skeleton className='bg-slate-200 h-8 w-24' />
                                      <Skeleton
                                          className={`bg-slate-200 w-full ${
                                              Math.random() > 0.5 ? 'h-8' : 'h-12'
                                          }`}
                                      />
                                  </div>
                              </div>
                          ))
                        : rootStore.activeDisplayedSession?.messages.map((conversation, i) => (
                              <Conversation
                                  key={i}
                                  conversation={conversation}
                                  isLastConversation={
                                      i ===
                                      (rootStore.activeDisplayedSession?.messages.length || 0) - 1
                                  }
                              />
                          ))}
                </div>

                <div className='max-w-5xl w-full sticky bottom-0 mx-auto py-4 flex flex-col gap-2 px-4 dark:bg-[#0f172a]'>
                    <div className='relative'>
                        <Textarea
                            onChange={(e) => setMessage(e.target.value)}
                            onKeyDown={(event) => handleKeyDown(event)}
                            value={message}
                            placeholder='Напишите в чат...'
                            name='message'
                            id='message'
                            rows={1}
                            className='min-h-[48px] rounded-2xl resize-none p-4 border border-gray-300 shadow-sm pr-16 dark:border-gray-800'
                        />
                        <Button
                            type='submit'
                            size='icon'
                            className='absolute top-3 right-3 w-8 h-8'
                            onClick={sendMessage}
                        >
                            <ArrowUpIcon className='w-4 h-4' />
                            <span className='sr-only'>Отправить</span>
                        </Button>
                    </div>
                    <p className='text-xs text-center text-neutral-700 font-medium dark:text-gray-400'>
                        Предсказания модели являются лишь рекомендациями. Пожалуйста, проверяйте
                        важную информацию.
                    </p>
                </div>
            </div>
        </>
    );
});

export default Chat;
