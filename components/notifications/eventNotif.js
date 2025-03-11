export default function eventNotif({ notification }) {
    function visit(relatedId) {
        // to implement
    }

    return (
      <div className="notification">
        <p>New event! visit {notification.related_name} for more info</p>
        <button onClick={() => visit(notification.related_id)}>Accept</button>
     </div>
    )
}