import { useState } from "react";
import styles from "./chatSection.module.css"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
// import { websoketProvider } from "../webSocket";

export default function ChatSection({ messages, selectedChat, sendMessage }) {
    const [message, setMessage] = useState("");

    const handleInputChange = (e) => {
        setMessage(e.target.value);
    };

    const handleSendMessage = () => {
        // if (message.trim()) {
        //     if (window.ws) {
        //         console.log(selectedChat);

        //         websoketProvider.sendMessage({
        //             type: "new_message",
        //             payload: {
        //                 chat_id: selectedChat.chat_id,
        //                 reciver_id: selectedChat.chat_id ? null : selectedChat.user_id,
        //                 content: message,
        //             }
        //         })
        //     }

        //     setMessage("");
        // }
    };

    return (
        <div className={styles.chatSection}>
            <div className={styles.header}>
                <span>{selectedChat.chat_name}</span>
                <span>X</span>
            </div>
            <div className={styles.messages}>
                {messages.map((msg) => (
                    <div key={msg.id} className={styles.message} >
                        <p>{msg.content}</p>
                        <span>{formatTime(msg.created_at)}</span>
                    </div>
                ))}
            </div>
            <div className={styles.messageInput}>
                <input
                    type="text"
                    value={message}
                    onChange={handleInputChange}
                    placeholder="Type a message..."
                />
                <button>
                    Emojis
                    {/* <FontAwesomeIcon icon={faPaperPlane} /> */}
                </button>
                <button onClick={sendMessage}>
                    Send
                    {/* <FontAwesomeIcon icon={faPaperPlane} /> */}
                </button>
            </div>
        </div>
    );
}

function formatTime(dateStr) {
    const date = new Date(dateStr);
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
  
    return `${hours}:${minutes}`;
  }