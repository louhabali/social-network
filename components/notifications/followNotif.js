export default function AcceptRejectNotif({ notification }) {
    function onAccept(relatedId) {
        // to implement
    }

    function onReject(relatedId) {
        // to implement
    }

    let  msg
    // a switch case on notification type to create a convintion msg
            
    

    return (
      <div className="notification">
        <p>{msg}</p>
        <button onClick={() => onAccept(notification.related_id)}>Accept</button>
        <button onClick={() => onReject(notification.related_id)}>Reject</button>
      </div>
    )
}