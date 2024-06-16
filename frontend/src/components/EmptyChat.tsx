import { useStores } from '@/hooks/useStores';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Card } from './ui/card';

const promptExamples = [
    'Какие шкафы есть на складе?',
    'Нужно ли закупать шкафы в ближайшие 3 месяца?',
    'Прогноз закупок бумаги на пол года вперед',
];

const EmptyChat = () => {
    const { rootStore } = useStores();

    return (
        <div className='flex w-full justify-center'>
            <div className='space-y-4'>
                {promptExamples.map((prompt, index) => (
                    <Card
                        key={index}
                        className='rounded-2xl border border-gray-200 dark:border-gray-800 p-4 hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors cursor-pointer'
                        onClick={() => {
                            rootStore.sendMessage({
                                prompt,
                            });
                        }}
                    >
                        <div className='flex items-start gap-4'>
                            <Avatar>
                                <AvatarFallback></AvatarFallback>
                            </Avatar>
                            <div className='grid gap-1 flex-1'>
                                <div>{prompt}</div>
                            </div>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    );
};

export default EmptyChat;
