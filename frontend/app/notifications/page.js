"use client"
import { useState, useEffect, useContext } from "react";
import EventNotif from "@/components/notifications/eventNotif.js";
import AcceptRejectNotif from "@/components/notifications/acceptRejectNotif.js";
import { WebSocketContext } from "@/context/WebSocketContext";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);
  const {ws} = useContext(WebSocketContext)
  const handlNotif = (notif) => {
    console.log("notifs =>", notifications);
    notif = JSON.parse(notif.data)
    if (notif.type == "notification") {
      setNotifications((notifications) => [notif.payload, ...notifications])
    }
  }

  useEffect(() => {
    async function fetchNotifications() {
      try {
        const response = await fetch(`${apiUrl}/api/notifications`, {
          credentials: "include",
        });
        const data = await response.json();
        console.log("notif", data);
        setNotifications(data ? data : []);
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }
    
    fetchNotifications();
    setTimeout(() => {
      ws.current.addEventListener("message", handlNotif)
    }, 500)
  }, []);

  const handleAccept = async (notif) => {
    let url
    switch (notif.type) {
      case "follow request":
        url = `${apiUrl}/api/follow/accept`
        break
      default:
        url = `${apiUrl}/api/group-join/accept`
    }

    console.log(notif);
    console.log("url ->", url);
    
      const response = await fetch(url, {
        credentials: "include",
        method: "POST",
        body: JSON.stringify(notif),
      });
  
      if(response.status == 200) {
        setNotifications(notifications.filter(notification => notification.id !== notif.id));
      } else {
        //error
      }
  };

  const handleReject = async (notif) => {
    let url
    console.log(notif.type);
    
    switch (notif.type) {
      case "follow request":
        url = `${apiUrl}/api/follow/reject`
        break
      default:
        url = `${apiUrl}/api/group-join/reject`
    }

    const response = await fetch(url, {
      credentials: "include",
      method: "DELETE",
      body: JSON.stringify(notif),
    });

    if(response.status == 200) {
      setNotifications(notifications.filter(notification => notification.id !== notif.id));
    } else {
      //error
    }
  };

  async function visit(notif) {
    const response = await fetch(`${apiUrl}/api/notification/read`, {
      credentials: "include",
      method: "POST",
      body: JSON.stringify(notif),
    });

    if(response.status == 200) {
      setNotifications(notifications.filter(notification => notification.id !== notif.id));
    } else {
      //error
    }
  }

  return (
    <div className="notifications-page">
      <h3>Notifications</h3>
      {notifications.map((notif) => (
        notif.type == "group event" ?
          <EventNotif
            key={notif.id}
            notification={notif}
            markRead={visit}
          /> :
          <AcceptRejectNotif
            key={notif.id}
            notification={notif}
            onAccept={handleAccept}
            onReject={handleReject}
          />
      ))}
    </div>
  );
}
