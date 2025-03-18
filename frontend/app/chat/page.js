"use client"
import { useEffect, useState, useContext } from "react";
import ChatList from "../../components/chat/chatlist";
import ChatSection from "../../components/chat/chatSection";
import styles from "./chatPage.module.css";

export default function messages() {
  const [conversations, setConversations] = useState([]);
  const [selectedChat, setSelectedChat] = useState(null);

  useEffect(() => {
    const storedConv = localStorage.getItem("conv")
    if (storedConv) {
      setSelectedChat(JSON.parse(storedConv))
    }
  },[])

  function updateLastReadMsg(chat_id) {
    fetch(`http://localhost:8080/api/messages/read`, {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({ chat_id }),
    });
}

  
  useEffect(() => {
    async function fetchConversations() {
      try {
        const response = await fetch(`http://localhost:8080/api/conversations`, {
          credentials: "include",
        });
        const data = await response.json();
        console.log("vonvs",data);
        setConversations(() => data ? data : [] );
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchConversations();
  }, [selectedChat]);



  const handleChatSelect = (conv) => {
    localStorage.setItem("conv", JSON.stringify(conv));
    setSelectedChat(conv);
  };

  function closeChat() {
    localStorage.removeItem("conv");
    setSelectedChat(null);
  }

  return (
    <div className={styles.chatPage}>
      { selectedChat ?
          <ChatSection
            selectedChat={selectedChat}
            closeChat={closeChat}
            updateLastReadMsg={updateLastReadMsg}
          /> :
          <ChatList
          setConversations={setConversations}
          conversations={conversations}
          onChatSelect={handleChatSelect}
          updateLastReadMsg={updateLastReadMsg}
        />
      }
    </div>
  );
}