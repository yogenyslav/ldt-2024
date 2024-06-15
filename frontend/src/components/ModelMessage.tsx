import { ClipboardIcon } from 'lucide-react';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';
import { ChatCommand, DisplayedIncomingMessage, IncomingMessageStatus } from '@/api/models';
import { useStores } from '@/hooks/useStores';
import { useToast } from './ui/use-toast';
import { useState } from 'react';
import Prediction from './Prediction';
import Stocks from './Stocks';

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
                                                    'Например, "Получить запасы товара картофель"',
                                            });

                                            rootStore.sendMessage({
                                                command: ChatCommand.Invalid,
                                            });

                                            setShowInvalidButton(false);
                                            rootStore.setChatDisabled(false);
                                            rootStore.setIsModelAnswering(false);
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
                        <div className='flex w-full flex-col gap-5'>
                            {' '}
                            {incomingMessage.prediction && (
                                <Prediction
                                    history={incomingMessage.prediction.history}
                                    forecast={incomingMessage.prediction.forecast}
                                />
                            )}
                            {incomingMessage.stocks && <Stocks stocks={incomingMessage.stocks} />}
                            <div className='prose prose-stone'>
                                <p>{incomingMessage.body}</p>
                            </div>
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

            case IncomingMessageStatus.Invalid:
                return (
                    <>
                        <div className='prose prose-stone'>
                            <p>{`продукт: ${incomingMessage.product}, период: ${incomingMessage.period}, тип: ${incomingMessage.type}`}</p>
                        </div>
                    </>
                );

            default:
                return 'Error';
        }
    };

    return (
        <div className='flex items-start gap-4 w-full'>
            <Avatar className='border w-8 h-8'>
                <AvatarFallback>MT</AvatarFallback>
            </Avatar>
            <div className='grid gap-1 mt-2 w-full'>
                <div className='font-bold'>Ответ модели</div>

                {getModelResonse()}
            </div>
        </div>
    );
};

export default ModelMessage;
