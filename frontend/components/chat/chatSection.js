import { useState, useRef, useEffect, useContext } from "react";
import styles from "./chatSection.module.css"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPaperPlane, faFaceSmile, faRotate } from "@fortawesome/free-solid-svg-icons";
// import EmojiPicker from './emojiPicker';
import { WebSocketContext } from "@/context/WebSocketContext";
import Picker from "./emojiPicker";
import Link from "next/link";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function ChatSection({ selectedChat, closeChat }) {
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);
    const [isHasMore, setIsHasMore] = useState(true);
    const [shouldScroll, setShouldScroll] = useState(false);
    const [showEmojis, setShowEmojis] = useState(false);
    const [newMsg, setNewMsg] = useState(false)
    const { userId, ws, sendMessage,showNotification } = useContext(WebSocketContext)
    const msgsContainer = useRef(null);
    const textarea = useRef(null);
    const offset = useRef(0);
    
    async function fetchMoremsgs() {
        const data = await fetchMsgs();
        setMessages([...data, ...messages]);
        if (data.length < 10) {
            setIsHasMore(false);
        }
        offset.current += data.length;
    }

    async function fetchMsgs() {
        const response = await fetch(`${apiUrl}/api/messages?chat_id=${selectedChat.chat_id}&offset=${offset.current}`, {
            credentials: "include",
        });
        let data = await response.json();
        data = data ? data.reverse() : []
        console.log(data);
        
        return data;
    }

    function sendMsg() {
        const msg = textarea.current.value

        if (msg.trim() && msg.length <= 500) {
            sendMessage({
                type: 'new_message',
                payload: {
                    chat_id: selectedChat.chat_id,
                    content: msg
                }
            })
            setMessage("")
        }
    }

    async function updateLastMsgReaded(id) {
        const resp = await fetch(`${apiUrl}/api/messages/read`, {
            credentials: "include",
            method: 'POST',
            body: JSON.stringify({id})
        })

        if (!resp.ok) {

        }
    }

    useEffect(() => {
        //initial msgs
        async function initialFetch() {
            const data = await fetchMsgs();
            setMessages(data);
            setShouldScroll(true);
            if (data.length < 10) {
                setIsHasMore(false);
            }
            offset.current = data.length;
        }

        function handleMessage(msg) {
            // console.log("new msg from chatSection!!!!!!!!");
            
            const data = JSON.parse(msg.data);
            if (data.type === "new_message" && data.payload.chat_id == selectedChat.chat_id) {
                setMessages((messages) => [...messages, data.payload])
                setNewMsg((newMsg) => !newMsg)
                offset.current += 1;
                console.log("update ->", data.payload);
                
                updateLastMsgReaded(data.payload.chat_id)
            }else if (data.type === "error_message") {
                showNotification("At least one follow the other constraint", "error");
            }
        }
        initialFetch();

        function handleKeyPress(e) {
            if (e.key === 'Enter' && !e.shiftKey) {
                // console.log("comme");
                e.preventDefault();
                sendMsg();
            }
        }

        window.addEventListener("keypress", handleKeyPress)
        if (ws.current) ws.current.addEventListener("message", handleMessage)
        return () => {
            if (ws.current) ws.current.removeEventListener("message", handleMessage)
            removeEventListener("keypress", handleKeyPress)
        }
    }, []);

    function scrollDown() {
        msgsContainer.current.scrollTop = msgsContainer.current.scrollHeight
    }

    useEffect(scrollDown, [shouldScroll, newMsg]);
    
    function emojiClick(emoji) {
        setMessage(message => message + emoji)
    }

    return (
        <div className={styles.chatSection}>
            <div className={styles.header}>
                <span>{selectedChat.chat_name}</span>
                <span onClick={closeChat} className={styles.close}>X</span>
            </div>
            <div ref={msgsContainer} className={`${styles.messages}`}>
                {isHasMore &&
                    <span className={styles.loadMore}>
                        <FontAwesomeIcon onClick={fetchMoremsgs} className={styles.icon} icon={faRotate} />
                    </span>
                }
                {messages.map((msg) => (
                    <span key={msg.id} className={`${styles.message} ${userId.current == msg.sender_id ? styles.sender : styles.receiver}`} >
                        {msg.is_group && !(userId.current == msg.sender_id) ? <Link href={`/profile/${msg.sender_id}`}><span className={styles.senderName}>{msg.sender_name}</span></Link> : ''}
                        <span>{msg.content}</span>
                        <span className={styles.msgTime}>{formatTime(msg.created_at)}</span>
                    </span>
                ))}
            </div>
            <div className={styles.messageInput}>
                <textarea
                    ref={textarea}
                    className={styles.msgInput}
                    value={message}
                    onChange={(e) => { setMessage(e.target.value) }}
                    placeholder="Type a message..."
                >
                </textarea>

                {showEmojis &&
                    <Picker
                        onEmojiClick={emojiClick}
                    />
                }

                <button onClick={() => setShowEmojis(!showEmojis)} className={styles.emogisBtn}>
                    <FontAwesomeIcon className={styles.icon} icon={faFaceSmile} />
                </button>
                <button onClick={sendMsg} className={styles.emogisBtn}>
                    <FontAwesomeIcon className={styles.icon} icon={faPaperPlane} />
                </button>
            </div>
        </div>
    );
}

function formatTime(dateStr) {
    const date = new Date(dateStr);
    // console.log(dateStr,"//", date);
    
    const hours = (date.getHours()).toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');

    return `${hours}:${minutes}`;
}