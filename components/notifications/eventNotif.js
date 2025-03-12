import Link from "next/link"

export default function EventNotif({ notification, markRead }) {
  return (
    <div className="notification">
      <p>New event! visit {notification.related_name} for more info</p>
      <Link href={`/groups/${notification.related_id}`}>
        <button onClick={() => markRead(notification)}>Visit</button>
      </Link>
    </div>
  )
}