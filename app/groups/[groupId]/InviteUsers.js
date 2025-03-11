import { useEffect, useState } from "react";
import "./invite-users.css"

export default function InviteUsers({ groupId }) {
  const [showPopup, setShowPopup] = useState(false);
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState("");
  const [offset, setOffset] = useState(0);


  console.log(users);
  
  const fetchUsers = async (reset = false) => {
    const res = await fetch(
      `http://localhost:8080/api/groups/${groupId}/invite?search=${search}&offset=${reset ? 0 : offset}`,
      { credentials: "include" }
    );
    const data = await res.json();
    if (reset) {
      setUsers(data);
      setOffset(4);
    } else {
      setUsers((prev) => data ? [...prev, ...data]:[...prev]);
      setOffset((prev) => prev + 4);
    }
  };
  
  useEffect(() => {
    if (showPopup) {
      fetchUsers(true);
    }
  }, [showPopup, groupId, search]);

  const inviteUser = async (userId) => {
    await fetch(`http://localhost:8080/api/groups/${groupId}/join`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ userId }),
    });
  };

  return (
    <div>
      <button onClick={() => setShowPopup(true)}>Invite Users</button>

      {showPopup && (
        <div className="popup">
          <div className="popup-content">
            <h3>Invite Users</h3>
            <input
              type="text"
              placeholder="Search users..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <ul onScroll={(e) => {
            if (e.target.scrollHeight - e.target.scrollTop === e.target.clientHeight) {
              fetchUsers();
            }
          }}>
              {users?.map((user) => (
                <li key={user.id}>
                  {user.name}
                  <button onClick={() => inviteUser(user.id)}>Invite</button>
                </li>
              ))}
            </ul>
            <button onClick={() => setShowPopup(false)}>Close</button>
          </div>
        </div>
      )}

      <style>{`
        .popup {
          position: fixed;
          top: 0; left: 0; width: 100%; height: 100%;
          background: rgba(0, 0, 0, 0.5);
          color: black;
          display: flex; justify-content: center; align-items: center;
        }
        .popup-content {
          background: white; padding: 20px; border-radius: 10px;
          width: 300px;
        }
        .popup-content ul { list-style: none; padding: 0; max-height: 130px; overflow-y: auto; }
        li { display: flex; justify-content: space-between; padding: 5px 0; }
      `}</style>
    </div>
  );
}