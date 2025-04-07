"use client"
import { useState, useEffect } from "react";
import EventNotif from "@/components/notifications/eventNotif.js";
import AcceptRejectNotif from "@/components/notifications/acceptRejectNotif.js";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    async function fetchNotifications() {
      try {
        const response = await fetch(`http://localhost:8080/api/notifications`, {
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
  }, []);

  const handleAccept = async (notif) => {
    let url
    switch (notif.type) {
      case "follow request":
        url = 'http://localhost:8080/api/follow/accept'
        break
      default:
        url = 'http://localhost:8080/api/group-join/accept'
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
        url = 'http://localhost:8080/api/follow/reject'
        break
      default:
        url = 'http://localhost:8080/api/group-join/reject'
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
    const response = await fetch('http://localhost:8080/api/notification/read', {
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
