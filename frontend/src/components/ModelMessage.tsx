import { ClipboardIcon } from 'lucide-react';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';
import { ChatCommand, DisplayedIncomingMessage, IncomingMessageStatus } from '@/api/models';
import { useStores } from '@/hooks/useStores';
import { useToast } from './ui/use-toast';
import { useState } from 'react';

type ModelMessageProps = {
    incomingMessage: DisplayedIncomingMessage;
    isLastMessage: boolean;
};

const ModelMessage = ({ incomingMessage, isLastMessage }: ModelMessageProps) => {
    const { rootStore } = useStores();
    const { toast } = useToast();
    const [showInvalidButton, setShowInvalidButton] = useState(true);

    const getModelResonse = () => {
        switch (incomingMessage.status) {
            case IncomingMessageStatus.Pending:
                return (
                    <>
                        <div className='prose prose-stone'>
                            <p>{`продукт: ${incomingMessage.product}, период: ${incomingMessage.period}, тип: ${incomingMessage.type}`}</p>
                        </div>
                        {isLastMessage && (
                            <div className='flex items-center gap-2 flex-wrap'>
                                <Button
                                    onClick={() => {
                                        rootStore.sendMessage({
                                            prompt: '',
                                            command: ChatCommand.Valid,
                                        });
                                    }}
                                    variant='outline'
                                    className='flex-1'
                                >
                                    Продолжить
                                </Button>
                                {showInvalidButton && (
                                    <Button
                                        onClick={() => {
                                            toast({
                                                title: 'Введите уточняющий запрос',
                                                description:
                                                    'Например, "Получить запасы товара картофель""',
                                            });

                                            setShowInvalidButton(false);
                                            rootStore.setChatDisabled(false);
                                        }}
                                        variant='outline'
                                        className='flex-1'
                                    >
                                        Уточнить запрос
                                    </Button>
                                )}
                            </div>
                        )}
                    </>
                );

            case IncomingMessageStatus.Valid:
                return (
                    <>
                        <div className='prose prose-stone'>
                            <p>{incomingMessage.body}</p>
                        </div>
                        <div className='flex items-center gap-2 py-2'>
                            <Button
                                variant='ghost'
                                size='icon'
                                className='w-4 h-4 hover:bg-transparent text-stone-400 hover:text-stone-900'
                            >
                                <ClipboardIcon className='w-4 h-4' />
                                <span className='sr-only'>Копировать</span>
                            </Button>
                        </div>
                    </>
                );

            default:
                return 'Error';
        }
    };

    return (
        <div className='flex items-start gap-4'>
            <Avatar className='border w-8 h-8'>
                <AvatarFallback>MT</AvatarFallback>
            </Avatar>
            <div className='grid gap-1 mt-2'>
                <div className='font-bold'>Ответ модели</div>

                {getModelResonse()}
            </div>
        </div>
    );
};

export default ModelMessage;
