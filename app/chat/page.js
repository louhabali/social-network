"use client"
import { useEffect, useState } from "react";
import ChatList from "./chatlist";
import MessageDisplay from "./diplayMsgs";

export default function messages() {
  const [conversations, setConversations] = useState([]);
  const [selectedChat, setSelectedChat] = useState(null);
  const [messages, setMessages] = useState([]);

  // Fetch conversations
  useEffect(() => {
    async function fetchConversations() {
      try {
        const response = await fetch(`http://localhost:8080/api/conversations`, {
          credentials: "include",
        }); // API call to get chats
        const data = await response.json();
        console.log(data);
        setConversations(data);

        //   setConversations(data);
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchConversations();
  }, []);

  useEffect(() => {
    if (selectedChat) {
      async function fetchMessages() {
        try {
          const response = await fetch(`http://localhost:8080/api/messages?chat_id=${selectedChat}`, {
            credentials: "include",
          });
          const data = await response.json();
          setMessages(data);
        } catch (error) {
          console.error("Failed to fetch messages", error);
        }
      }

      // fetch msgs
      // fetchMessages();
    }
  }, [selectedChat]);

  window.ws.onmessage = (event) => {
    console.log(event.data);
  }

  const handleChatSelect = (conv) => {
    setSelectedChat(conv);
  };


  return (
    <div className="messages-page">
      <div className="left-section">
        <ChatList
          conversations={conversations}
          onChatSelect={handleChatSelect}
        />
      </div>
      <div className="right-section">
        {(selectedChat 
          && <MessageDisplay
          messages={messages}
          selectedChat={selectedChat}
        />)}
      </div>
    </div>
  );
}