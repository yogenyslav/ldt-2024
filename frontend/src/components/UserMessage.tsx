import { Avatar, AvatarFallback } from './ui/avatar';

type UserMessageProps = {
    message: string;
};

const UserMessage = ({ message }: UserMessageProps) => {
    return (
        <div className='flex items-start gap-4'>
            <Avatar className='border w-8 h-8'>
                <AvatarFallback>m.t</AvatarFallback>
            </Avatar>
            <div className='grid gap-1'>
                <div className='font-bold'>Вы</div>
                <div className='prose prose-stone'>
                    <p>{message}</p>
                </div>
            </div>
        </div>
    );
};

export default UserMessage;
