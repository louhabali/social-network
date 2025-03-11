"use client"
import { useState, useEffect } from "react";
import Notification from "../../components/folowerRequestNotif";

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    async function fetchNotifications() {
        try {
          const response = await fetch(`http://localhost:8080/api/notifications`, {
            credentials: "include",
          });
          const data = await response.json();
          console.log("notif",data);
          setNotifications(data);
  
          //   setConversations(data);
        } catch (error) {
          console.error("Failed to fetch conversations", error);
        }
      }

    fetchNotifications();
  }, []);

  const handleAccept = async (requestId) => {
    await fetch(`/api/accept-follow-request`, {
      method: 'POST',
      body: JSON.stringify({ requestId }),
    });
    // Update the local state to reflect the accepted request
    setNotifications(notifications.filter(notification => notification.id !== requestId));
  };

  const handleReject = async (requestId) => {
    await fetch(`/api/reject-follow-request`, {
      method: 'DELETE',
      body: JSON.stringify({ requestId }),
    });
    // Update the local state to reflect the rejected request
    setNotifications(notifications.filter(notification => notification.id !== requestId));
  };

  return (
    <div className="notifications-page">
      <h3>Notifications</h3>
      {notifications.map((request) => (
        
        <Notification
          key={request.id}
          request={request}
          onAccept={handleAccept}
          onReject={handleReject}
        />
      ))}
    </div>
  );
}
