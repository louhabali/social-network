import Link from "next/link"
import styles from "./notifs.module.css";

const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function EventNotif({notification}) {
  function deleteNotif() {
    fetch(`${apiUrl}/api/notification/delete`, {
      credentials: "include",
      method:"DELETE",
      body: JSON.stringify({id: notification.id})
    })
  }
  return (
    <div className={styles.notif}>
      <p>
        <Link className={styles.link} href={`/profile/${notification.sender_id}`}>
          {notification.sender_name}
        </Link>
        {' '}
        have create a new event visit
        {' '}
        <Link
          className={styles.link}
          href={`/groups/${notification.related_id}`}
          onClick={() => deleteNotif()}
        >
          {notification.group_name}
        </Link>
        {' '}
        for more info!
      </p>

    </div>
  )
}