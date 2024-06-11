import ModelMessage from '@/components/ModelMessage';
import UserMessage from '@/components/UserMessage';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useToast } from '@/components/ui/use-toast';
import { useStores } from '@/hooks/useStores';
import { ArrowUpIcon } from 'lucide-react';
import { useEffect } from 'react';

const Chat = () => {
    const { rootStore } = useStores();
    const { toast } = useToast();

    useEffect(() => {
        rootStore.getSessions().catch(() => {
            toast({
                title: 'Ошибка',
                description: 'Не удалось загрузить сессии',
                variant: 'destructive',
            });
        });
    }, [rootStore, toast]);

    return (
        <>
            <div className='flex items-center'>
                <h1 className='text-lg font-semibold md:text-2xl'>Чат</h1>
            </div>

            <div className='flex flex-col h-screen'>
                <div className='max-w-5xl flex-1 mx-auto flex flex-col items-start gap-8 px-4 py-8'>
                    <UserMessage
                        message='Can you explain airplane turbulence to someone who has never
                                    flown before? Make it conversational and concise.'
                    />

                    <ModelMessage
                        message='Airplane turbulence is caused by changes in air pressure and
                                    temperature. It is a normal part of flying and is nothing to worry
                                    about. Pilots are trained to handle turbulence, and modern aircraft
                                    are designed to withstand it. Just sit back, relax, and enjoy the
                                    flight!'
                    />
                </div>
                <div className='max-w-5xl w-full sticky bottom-0 mx-auto py-4 flex flex-col gap-2 px-4 dark:bg-[#0f172a]'>
                    <div className='relative'>
                        <Textarea
                            placeholder='Message Acme AI...'
                            name='message'
                            id='message'
                            rows={1}
                            className='min-h-[48px] rounded-2xl resize-none p-4 border border-gray-200 border-neutral-400 shadow-sm pr-16 dark:border-gray-800'
                        />
                        <Button
                            type='submit'
                            size='icon'
                            className='absolute top-3 right-3 w-8 h-8'
                            disabled
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
};

export default Chat;
