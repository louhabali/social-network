import styles from './chatlist.module.css'
import { useContext, useEffect } from 'react';
import { WebSocketContext } from "@/context/WebSocketContext";



export default function ChatList({ sortConvs, conversations, onChatSelect, setConversations }) {
    const {ws} = useContext(WebSocketContext);
    function newMsg(msg) {
        // console.log("new msg from chatSection!!!!!!!!");
        
        const data = JSON.parse(msg.data);
        // console.log("sortConvs 00000000000000", data);
        if (data.type === "new_message") {
            setConversations((preConvs) => {
                const updatedConvs = preConvs.map(conv => {
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
                
                return sortConvs(updatedConvs) 
            })

            // setConversations(() => [convx, ...newConvs])
            // console.log(conversations);
        }
    }


    useEffect(() => {
        setTimeout(()=>{
            console.log(ws.current);
            ws.current.addEventListener("message", newMsg)
        },500)
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