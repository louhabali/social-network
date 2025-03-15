import Link from "next/link"

export default function AcceptRejectNotif({ notification, onAccept, onReject }) {
  let msg
  switch (notification.type) {
    case "follow reqeust":
      msg = ` send you a follow request`
      break
    case "group join request":
      msg = ` want to join `
      break
    default:
      msg = ` send you an invitation to join `
  }

  return (
    <div>
      <p>
        <Link href={`/profile/${notification.sender_id}`}>
          <span>{notification.sender_name}</span>
        </Link>
        {msg}
        {notification.type != "follow reqeust"
          && 
          <Link href={`/groups/${notification.group_id}`}>
          <span>{notification.group_name}</span>
          </Link>
        }
      </p>
      <button onClick={() => onAccept(notification)}>Accept</button>
      <button onClick={() => onReject(notification)}>Reject</button>
    </div>
  )
}