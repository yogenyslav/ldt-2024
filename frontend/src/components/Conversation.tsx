import { ChatConversation } from '@/api/models';
import UserMessage from './UserMessage';
import ModelMessage from './ModelMessage';
import { observer } from 'mobx-react-lite';

type Props = {
    conversation: ChatConversation;
    isLastConversation: boolean;
};

const Conversation = observer(({ conversation, isLastConversation }: Props) => {
    return (
        <div>
            {conversation.outcomingMessage && (
                <UserMessage message={conversation.outcomingMessage.prompt} />
            )}

            {conversation.incomingMessage && (
                <ModelMessage
                    incomingMessage={conversation.incomingMessage}
                    isLastMessage={isLastConversation}
                />
            )}
        </div>
    );
});

export default Conversation;
