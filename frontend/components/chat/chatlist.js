import styles from './chatlist.module.css'

export default function ChatList({ conversations, onChatSelect, setConversations }) {
    function updateLastReadMsg(chat_id) {
        // fetch(`http://localhost:8080/api/messages/read`, {
        //     method: "POST",
        //     credentials: "include",
        //     body: JSON.stringify({ chat_id }),
        // });
    }


    return (
        <div className={styles.chatList}>
            <div className={styles.chatHeader}></div>
            {
                conversations.map((conversation) => (
                    <div 
                        className={
                            conversation.unread_count > 0 ?
                            `${styles.chatItem} ${styles.unreaded}`:
                            `${styles.chatItem} ${styles.readed}`
                        }
                        key={conversation.chat_id}
                        onClick={() => {
                            onChatSelect(conversation);
                            updateLastReadMsg(conversation.chat_id);
                            setConversations(
                                conversations.map((conv) => {
                                    if (conv.chat_id === conversation.chat_id) {
                                        return {
                                            ...conv,
                                            unread_count: 0
                                        }
                                    }
                                    return conv;
                                })
                            );
                        }}
                    >
                        <div className={styles.uImage}></div>
                        <p className={styles.uName}>{conversation.chat_name}</p>
                        <p className={styles.lastMsg}>
                            {conversation.last_message ?
                             conversation.last_message.length > 20 ?
                             conversation.last_message.slice(0, 20) + '...' :
                             conversation.last_message :
                             '' 
                            }
                        </p>
                        {
                            conversation.unread_count > 0 &&
                            <span className={styles.count}>
                            {
                                conversation.unread_count > 9 &&
                                <span className={styles.plus}>+</span>
                            }
                            <span className={styles.num}>{conversation.unread_count >= 9 ? 9 : conversation.unread_count}</span>
                            </span>
                        }
                        
                    </div>
                ))
            }
        </div>
    );
}