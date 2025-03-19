import Link from "next/link"
import styles from "./notifs.module.css";

export default function EventNotif({ notification, markRead }) {
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
          onClick={() => markRead(notification)}
        >
          {notification.group_name}
        </Link>
        {' '}
        for more info!
      </p>

    </div>
  )
}