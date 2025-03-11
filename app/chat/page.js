"use client"
import { useEffect, useState } from "react";
import ChatList from "./chatlist";
import MessageDisplay from "./diplayMsgs";

export default function messages() {
  const [conversations, setConversations] = useState([]);
  const [selectedChat, setSelectedChat] = useState(null);
  const [messages, setMessages] = useState([]);
  const [offset, setOffset] = useState(0);

  // Fetch conversations
  useEffect(() => {
    async function fetchConversations() {
      try {
        const response = await fetch(`http://localhost:8080/api/conversations`, {
          credentials: "include",
        });
        const data = await response.json();
        console.log("vonvs",data);
        setConversations(data);

        //   setConversations(data);
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchConversations();
  }, [selectedChat, offset]);

  useEffect(() => {
    if (selectedChat) {
      console.log("sdk");
      
      async function fetchMessages() {
        try {
          const response = await fetch(`http://localhost:8080/api/messages?chat_id=${selectedChat.chat_id}&offset=${offset}`, {
            credentials: "include",
          });
          const data = await response.json();
          console.log("msgs => ", data);
          
          setMessages(data);
        } catch (error) {
          console.error("Failed to fetch messages", error);
        }
      }

      fetchMessages();
    }
  }, []);

  // useEffect(() => {
  //   if (window.ws) {
  //     window.ws.onmessage = (event) => {
  //       const newMessage = JSON.parse(event.data);
  //       console.log("New Message => ", newMessage);

  //       // Add the new message if it's for the selected chat
  //       if (selectedChat && newMessage.chat_id === selectedChat.chat_id) {
  //         setMessages((prevMessages) => [newMessage, ...prevMessages]);
  //       }
  //     };
  //   }
  // }, [selectedChat]);

  const handleChatSelect = (conv) => {
    setSelectedChat(conv);
    setOffset(0);
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