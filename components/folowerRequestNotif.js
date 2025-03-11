
export default function Notification({ request, onAccept, onReject }) {
  return (
    <div className="notification">
      <p>{request.follower_name} wants to follow you!</p>
      <button onClick={() => onAccept(request.id)}>Accept</button>
      <button onClick={() => onReject(request.id)}>Reject</button>
    </div>
  );
}