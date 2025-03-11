"use client"
import { useState, useEffect } from "react";
import { FollowNotif } from "../../components/notifications";

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

        //   setConversations(data);
      } catch (error) {
        console.error("Failed to fetch conversations", error);
      }
    }

    fetchNotifications();
  }, []);

  const handleAccept = async (requestId) => {
      const response = await fetch(`http://localhost:8080/api/follow/accept`, {
        credentials: "include",
        method: "post",
        body: JSON.stringify({ requestId }),
      });
  
      if(response.status == 200) {
        setNotifications(notifications.filter(notification => notification.request_id !== requestId));
      } else {
        //error
      }
  };

  const handleReject = async (requestId) => {
    const response = await fetch(`http://localhost:8080/api/follow/reject`, {
      credentials: "include",
      method: "DELETE",
      body: JSON.stringify({ requestId }),
    });

    if(response.status == 200) {
      setNotifications(notifications.filter(notification => notification.request_id !== requestId));
    } else {
      //error
    }
  };

  return (
    <div className="notifications-page">
      <h3>Notifications</h3>
      {notifications.map((request) => (
        request.type == "follow_request" ?
          <FollowNotif
            key={request.request_id}
            request={request}
            onAccept={handleAccept}
            onReject={handleReject}
          /> : ''
      ))}
    </div>
  );
}
