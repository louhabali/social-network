import Link from "next/link"

export default function AcceptRejectNotif({ notification, onAccept, onReject }) {
  let msg
  switch (notification.type) {
    case "follow reqeust":
      msg = ` want to follow you`
      break
    case "group join requst":
      msg = ` want to join the group`
      break
    default:
      msg = ` send you a group invitation`
  }

  return (
    <div className="notification">
      <p>
        <Link href={`/profile/${notification.user_id}`}>
          <span>{notification.related_name}</span>
        </Link>
        {msg}
      </p>
      <button onClick={() => onAccept(notification)}>Accept</button>
      <button onClick={() => onReject(notification)}>Reject</button>
    </div>
  )
}