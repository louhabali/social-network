import { useState } from "react";
import { websoketProvider } from "../webSocket";

export default function MessageDisplay({ messages, onSendMessage, selectedChat, chatName }) {
    const [message, setMessage] = useState("");

    const handleInputChange = (e) => {
        setMessage(e.target.value);
    };
    // console.log(window.ws);

    const handleSendMessage = () => {
        if (message.trim()) {
            if (window.ws) {
                console.log(selectedChat);

                websoketProvider.sendMessage({
                    type: "new_message",
                    payload: {
                        chat_id: selectedChat.chat_id,
                        reciver_id: selectedChat.chat_id ? null : selectedChat.user_id,
                        content: message,
                    }
                })
            }

            setMessage("");
        }
    };

    return (
        <div className="message-display">
            <div className="chat-header">
                <h3>{selectedChat.caht_name}</h3>
            </div>
            <div className="messages">
                {messages.map((msg) => (
                    <div key={msg.id} className="message">
                        <p>{msg.content}</p>
                        <br />
                        <span>{msg.created_at}</span>
                    </div>
                ))}
            </div>
            <div className="message-input">
                <input
                    type="text"
                    value={message}
                    onChange={handleInputChange}
                    placeholder="Type a message..."
                />
                <button onClick={handleSendMessage}>Send</button>
            </div>
        </div>
    );
}
