"use client"
import { useEffect, useState } from "react";
import ChatList from "../../components/chat/chatlist";
import ChatSection from "../../components/chat/chatSection";

export default function messages() {
  const [conversations, setConversations] = useState([]);
  const [selectedChat, setSelectedChat] = useState(null);
  const [messages, setMessages] = useState([]);
  const [offset, setOffset] = useState(0);

  useEffect(() => {
    const storedConv = localStorage.getItem("conv")
    if (storedConv) {
    setSelectedChat(JSON.parse(storedConv))
  }
  },[])
  
  // Fetch conversations
  useEffect(() => {
    async function fetchConversations() {
      try {
        const response = await fetch(`http://localhost:8080/api/conversations`, {
          credentials: "include",
        });
        const data = await response.json();
        console.log("vonvs",data);
        setConversations(data ? data : [] );

        //   setConversations(data);
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchConversations();
  }, []);

  useEffect(() => {
    if (selectedChat) {
      // console.log("sdk");
      
      async function fetchMessages() {
        try {
          const response = await fetch(`http://localhost:8080/api/messages?chat_id=${selectedChat.chat_id}&offset=${offset}`, {
            credentials: "include",
          });
          const data = await response.json();
          console.log("msgs => ", data);
          
          setMessages(data ? data : []);
        } catch (error) {
          console.error("Failed to fetch messages", error);
        }
      }

      fetchMessages();
    }
  }, [selectedChat, offset]);

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
    localStorage.setItem("conv", JSON.stringify(conv));
    setSelectedChat(conv);
    setOffset(0);
  };


  return (
    <div className="messages-page">
        <ChatList
          conversations={conversations}
          onChatSelect={handleChatSelect}
        />
        {(selectedChat 
          && <ChatSection
          messages={messages}
          selectedChat={selectedChat}
        />)}
    </div>
  );
}