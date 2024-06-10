import { ClipboardIcon } from 'lucide-react';
import { Avatar, AvatarFallback } from './ui/avatar';
import { Button } from './ui/button';

type ModelMessageProps = {
    message: string;
};

const ModelMessage = ({ message }: ModelMessageProps) => {
    return (
        <div className='flex items-start gap-4'>
            <Avatar className='border w-8 h-8'>
                <AvatarFallback>AC</AvatarFallback>
            </Avatar>
            <div className='grid gap-1'>
                <div className='font-bold'>Acme AI</div>
                <div className='prose prose-stone'>
                    <p>{message}</p>
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
            </div>
        </div>
    );
};

export default ModelMessage;
