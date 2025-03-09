export default function ChatList({ conversations, onChatSelect }) {
    return (
        <div className="chat-list">
            { conversations ? 
                conversations.map((conversation) => (
                    <div
                        key={conversation.chat_id}
                        className="chat-item"
                        onClick={() => onChatSelect(conversation)}
                    >
                        <p>{conversation.chat_name}</p>
                        <p>{conversation.last_message}</p>
                        {conversation.unread_count > 0 && (
                            <span className="unread-count">{conversation.unread_count}</span>
                        )}
                    </div>
                ))
            :''}
        </div>
    );
}
