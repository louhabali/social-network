import styles from './chatlist.module.css'
import { useContext, useEffect } from 'react';
import { WebSocketContext } from "@/components/websocket";


export default function ChatList({ conversations, onChatSelect, setConversations, updateLastReadMsg }) {
    const {ws} = useContext(WebSocketContext);

    function newMsg(msg) {
        const data = JSON.parse(msg.data);
        if (data.type === "new_message") {
            setConversations(preConvs => {
                return preConvs.map(conv => {
                    if (conv.chat_id === data.payload.chat_id) {
                        return {
                            ...conv,
                            last_message: data.payload.content,
                            last_message_time: data.payload.created_at,
                            unread_count: data.payload.is_sender ? conv.unread_count : conv.unread_count + 1
                        }
                    }
                    return conv
                })
            })
        }
    }


    useEffect(() => {
        if (ws.current) ws.current.addEventListener("message", newMsg)
        return (
            () => {
                if (ws.current) ws.current.removeEventListener("message", newMsg)
            }
        )
    }, [])


    return (
        <div className={styles.chatList}>
            <div className={styles.chatHeader}></div>
            {
                conversations.map((conversation) => (
                    <div
                        className={
                            conversation.unread_count > 0 ?
                                `${styles.chatItem} ${styles.unreaded}` :
                                `${styles.chatItem} ${styles.readed}`
                        }
                        key={conversation.chat_id}
                        onClick={() => {
                            onChatSelect(conversation);
                            console.log(conversation.unread_count);
                            
                            if(conversation.unread_count > 0) updateLastReadMsg(conversation.chat_id)
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
                                <span className={styles.num}>{conversation.unread_count >= 9 ? 9 : conversation.unread_count}</span>
                            </span>
                        }
                    </div>
                ))
            }
        </div>
    );
}