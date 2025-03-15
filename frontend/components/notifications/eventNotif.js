import Link from "next/link"

export default function EventNotif({ notification, markRead }) {
  return (
    <div>
      <p>
        <Link href={`/profile/${notification.sender_id}`}>
          {notification.sender_name}
        </Link>
        {' '}
        have create a new event visit
        {' '}
        <Link 
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