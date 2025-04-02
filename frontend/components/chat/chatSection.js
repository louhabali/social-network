import { useState, useRef, useEffect, useContext } from "react";
import styles from "./chatSection.module.css"
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faPaperPlane, faFaceSmile, faRotate } from "@fortawesome/free-solid-svg-icons";
import EmojiPicker from './emojiPicker';
import { WebSocketContext } from "@/context/WebSocketContext";

export default function ChatSection({ selectedChat, closeChat }) {
    const [message, setMessage] = useState("");
    const [messages, setMessages] = useState([]);
    const [isHasMore, setIsHasMore] = useState(true);
    const [shouldScroll, setShouldScroll] = useState(false);
    const [showEmojis, setShowEmojis] = useState(false);
    const [newMsg, setNewMsg] = useState(false)
    const {ws, sendMessage} = useContext(WebSocketContext)
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
        const response = await fetch(`http://localhost:8080/api/messages?chat_id=${selectedChat.chat_id}&offset=${offset.current}`, {
            credentials: "include",
        });
        let data = await response.json();
        data = data ? data.reverse() : []
        
        return data;
    }

    function sendMsg() {
        const msg = textarea.current.value
        
        if(msg.trim() && msg.length <= 500) {
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

    useEffect(() => {
        //initial msgs
        async function initialFetch() {
            const data = await fetchMsgs();
            console.log(data);
            
            setMessages(data);
            setShouldScroll(true);
            if (data.length < 10) {
                setIsHasMore(false);
            }
            offset.current = data.length;
        }
        
        function handleMessage(msg) {
            const data = JSON.parse(msg.data);
            if(data.type === "new_message" && data.payload.chat_id == selectedChat.chat_id) {
                setMessages((messages) => [...messages, data.payload])
                setNewMsg((newMsg) => !newMsg)
                offset.current += 1;
            }
        }
        initialFetch();

        function handleKeyPress(e){
            if (e.key === 'Enter' && !e.shiftKey) {
                console.log("comme");
                e.preventDefault();
                sendMsg();
            }
        }
        
        addEventListener("keypress", handleKeyPress)
        ws.current.addEventListener("message",handleMessage)
        return () => {
            ws.current.removeEventListener("message", handleMessage)
            removeEventListener("keypress", handleKeyPress)
        }
    }, []);

    function scrollDown() {
        msgsContainer.current.scrollTop = msgsContainer.current.scrollHeight
    }

    useEffect(scrollDown, [shouldScroll, newMsg]);

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
                    <span key={msg.id} className={`${styles.message} ${msg.is_sender ? styles.sender : styles.receiver}`} >
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
                    onChange={(e) => {setMessage(e.target.value)}}
                    placeholder="Type a message..."
                >
                </textarea>
                <button onClick={() => setShowEmojis(!showEmojis)} className={styles.emogisBtn}>
                    {
                        showEmojis && <EmojiPicker setMessage={setMessage} setShowEmojis={setShowEmojis}/>
                    }
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
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');

    return `${hours}:${minutes}`;
}