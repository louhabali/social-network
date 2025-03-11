
export function FollowNotif({ request, onAccept, onReject }) {
  console.log(request);
  
  return (
    <div className="notification">
      <p>{request.first_name + ' '+ request.last_name} wants to follow you!</p>
      <button onClick={() => onAccept(request.request_id)}>Accept</button>
      <button onClick={() => onReject(request.request_id)}>Reject</button>
    </div>
  );
}