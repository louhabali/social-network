import Link from "next/link"
import styles from "./notifs.module.css";
import { faCheck, faXmark } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

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
    <div className={styles.notif}>
      <p>
        <Link className={styles.link} href={`/profile/${notification.sender_id}`}>
          <span>{notification.sender_name}</span>
        </Link>
        {msg}
        {notification.type != "follow reqeust"
          &&
          <Link className={styles.link} href={`/groups/${notification.group_id}`}>
            <span>{notification.group_name}</span>
          </Link>
        }
      </p>
      <div className={styles.btns}>
        <FontAwesomeIcon onClick={() => onAccept(notification)} className={`${styles.btn} ${styles.mg}`} icon={faCheck}/>
        <FontAwesomeIcon onClick={() => onReject(notification)} className={styles.btn} icon={faXmark}/>
      </div>
    </div>
  )
}