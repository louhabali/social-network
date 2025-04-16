"use client"
import { useEffect, useState, useContext } from "react";
import ChatList from "../../components/chat/chatlist";
import ChatSection from "../../components/chat/chatSection";
import styles from "./chatPage.module.css";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

function sortConvs(convs) {
  const newConvs = [...convs]
  newConvs.sort((a, b) => {
    console.log(new Date(b.last_message_time) , new Date(a.last_message_time));
    
    if (a.last_message_time && b.last_message_time) {
      return new Date(b.last_message_time) - new Date(a.last_message_time);
    } else if (b.last_message_time) {
      return 1;
    } else if (a.last_message_time) {
      return -1;
    }
  })
  return newConvs
}

export default function messages() {
  const [conversations, setConversations] = useState([]);
  const [selectedChat, setSelectedChat] = useState(null);

  useEffect(() => {
    async function fetchConversations() {
      try {
        const response = await fetch(`${apiUrl}/api/conversations`, {
          credentials: "include",
        });
        let data = await response.json();
        // setConvs(data ? data : []);
        setConversations(data ? sortConvs(data) : [])
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchConversations();
  }, [selectedChat]);

  // function setConvs(convs) {
  //   // const nConvs = sortConvs(convs)
    
  // }

  function handleNewMsgChatList(msg) {
    const data = JSON.parse(msg.data);
    if (data.type === "new_message") {
        const updatedConvs = conversations.map(conv => {
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
      setConvs(updatedConvs)
    }
  }


  const handleChatSelect = (conv) => {
    setSelectedChat(conv);
  };

  function closeChat() {
    setSelectedChat(null);
  }

  return (
    <div className={styles.chatPage}>
      {selectedChat ?
        <ChatSection
          selectedChat={selectedChat}
          closeChat={closeChat}
        /> :
        <ChatList
          sortConvs={sortConvs}
          setConversations={setConversations}
          conversations={conversations}
          onChatSelect={handleChatSelect}
        />
      }
    </div>
  );
}