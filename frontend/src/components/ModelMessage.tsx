import { ClipboardIcon } from 'lucide-react';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';
import {
    ChatCommand,
    DisplayedIncomingMessage,
    IncomingMessageStatus,
    IncomingMessageType,
} from '@/api/models';
import { useStores } from '@/hooks/useStores';
import { useToast } from './ui/use-toast';
import { useState } from 'react';
import Prediction from './Prediction';
import PrompterResult from './PrompterResult';
import MarkdownPreview from '@uiw/react-markdown-preview';
import { LoaderButton } from './ui/loader-button';
import FavoritesApiService from '@/api/FavoritesApiService';
import StocksGroup from './StocksGroup';

type ModelMessageProps = {
    incomingMessage: DisplayedIncomingMessage;
    isLastMessage: boolean;
};

const ModelMessage = ({ incomingMessage, isLastMessage }: ModelMessageProps) => {
    const { rootStore } = useStores();
    const { toast } = useToast();
    const [showInvalidButton, setShowInvalidButton] = useState(true);
    const [isSaving, setIsSaving] = useState(false);

    console.log(incomingMessage.type);

    const downloadFile = () => {
        if (incomingMessage.outputJson) {
            const data = JSON.stringify(incomingMessage.outputJson);

            const blob = new Blob([data], { type: 'application/json' });
            const url = URL.createObjectURL(blob);

            const a = document.createElement('a');
            a.href = url;
            a.download = 'result.json';
            document.body.appendChild(a);
            a.click();
            a.remove();
        }
    };

    const getModelResonse = () => {
        switch (incomingMessage.status) {
            case IncomingMessageStatus.Pending:
                return (
                    <>
                        <div className='prose prose-stone'>
                            <PrompterResult
                                product={incomingMessage.product}
                                period={incomingMessage.period}
                                type={incomingMessage.type}
                            />
                        </div>
                        {isLastMessage && (
                            <div className='flex items-center gap-2 flex-wrap mt-2'>
                                <Button
                                    onClick={() => {
                                        rootStore.sendMessage({
                                            command: ChatCommand.Valid,
                                        });
                                    }}
                                    variant='default'
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
                            {incomingMessage.stocks && (
                                <StocksGroup stocksGroup={incomingMessage.stocks} />
                            )}
                            {incomingMessage.type === IncomingMessageType.Prediction && (
                                <div className='flex gap-2 flex-wrap'>
                                    <LoaderButton
                                        variant='outline'
                                        onClick={() => {
                                            console.log(incomingMessage.outputJson);

                                            if (incomingMessage.outputJson) {
                                                setIsSaving(true);

                                                FavoritesApiService.createFavorite({
                                                    id: 1,
                                                    response: incomingMessage.outputJson,
                                                })
                                                    .then(() => {
                                                        setIsSaving(false);
                                                        toast({
                                                            title: 'Успех',
                                                            description:
                                                                'Прогноз сохранен в избранное',
                                                        });
                                                    })
                                                    .catch(() => {
                                                        setIsSaving(false);
                                                        toast({
                                                            title: 'Ошибка',
                                                            description:
                                                                'Не удалось сохранить ответ в избранное',
                                                            variant: 'destructive',
                                                        });
                                                    });
                                            }
                                        }}
                                        isLoading={isSaving}
                                    >
                                        Сохранить план закупки
                                    </LoaderButton>

                                    <Button variant='outline' onClick={downloadFile}>
                                        Загрузить план закупки (.json)
                                    </Button>
                                </div>
                            )}
                            <div className='prose prose-stone'>
                                <MarkdownPreview
                                    source={incomingMessage.body}
                                    style={{ padding: 16 }}
                                />
                            </div>
                        </div>
                        <div className='flex items-center gap-2 py-2'>
                            <Button
                                variant='ghost'
                                size='icon'
                                className='w-4 h-4 hover:bg-transparent text-stone-400 hover:text-stone-900'
                                onClick={() => {
                                    navigator.clipboard.writeText(incomingMessage.body);
                                    toast({
                                        title: 'Скопировано',
                                        description: 'Текст ответа скопирован в буфер обмена',
                                    });
                                }}
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
                            <PrompterResult
                                product={incomingMessage.product}
                                period={incomingMessage.period}
                                type={incomingMessage.type}
                            />
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
