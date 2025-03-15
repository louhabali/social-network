import styles from './chatlist.module.css'

export default function ChatList({ conversations, onChatSelect }) {
    
    return (
        <div className={styles.chatList}>
            {
                conversations.map((conversation) => (
                    <div 
                        className={
                            conversation.unread_count > 0 ?
                            `${styles.chatItem} ${styles.unreaded}`:
                            `${styles.chatItem} ${styles.readed}`
                        }
                        key={conversation.chat_id}
                        onClick={() => onChatSelect(conversation)}
                    >
                        <div className={styles.uImage}></div>
                        <p className={styles.uName}>{conversation.chat_name}</p>
                        <p className={styles.lastMsg}>
                            {conversation.last_message ?
                             conversation.last_message.length > 20 ?
                             conversation.last_message.slice(0, 19) + '...' :
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
                            <span className={styles.num}>{conversation.unread_count}</span>
                            </span>
                        }
                        
                    </div>
                ))
            }
        </div>
    );
}
