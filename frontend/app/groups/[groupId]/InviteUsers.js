"use client"
import { useEffect, useState, useContext } from "react";
import { WebSocketContext } from "@/context/WebSocketContext";
import styles from "./inviteUsers.module.css"


const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function InviteUsers({ groupId }) {
  const [showPopup, setShowPopup] = useState(false);
  const [users, setUsers] = useState([]);
  const [search, setSearch] = useState("");
  const [offset, setOffset] = useState(0);
  const [invitedUsers, setInvitedUsers] = useState({});
  const { showNotification } = useContext(WebSocketContext);


  useEffect(() => {
    if (showPopup) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "auto";
    }

    return () => {
      document.body.style.overflow = "auto";
    };
  }, [showPopup]);


  const fetchUsers = async (reset = false) => {
    const res = await fetch(
      `${apiUrl}/api/groups/${groupId}/invite?search=${search}&offset=${reset ? 0 : offset}`,
      { credentials: "include" }
    );

    if (!res.ok) {
      showNotification("Error updating event response", "error");
      return;
    }

    const data = await res.json();
    if (reset) {
      setUsers(data);
      setOffset(6);
    } else {
      setUsers((prev) => data ? [...prev, ...data] : [...prev]);
      setOffset((prev) => prev + 6);
    }
  };

  useEffect(() => {
    if (showPopup) {
      fetchUsers(true);
    }
  }, [showPopup, groupId, search]);

  const inviteUser = async (userId) => {
    const res = await fetch(`${apiUrl}/api/groups/${groupId}/join`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ userId }),
    });
    if (!res.ok) {
      showNotification("Error inviting user", "error");
      return;
    }
    setInvitedUsers((prev) => ({ ...prev, [userId]: true }));
  };

  return (
    <div>
      <button className={styles.button} onClick={() => setShowPopup(true)}>Invite Users</button>

      {showPopup && (
        <div className={styles.popup}>
          <div className={styles["popup-content"]}>
            <h3 className={styles.header}>Invite Users</h3>
            <input
              className={styles.searchInput}
              type="text"
              placeholder="Search users..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <ul className={styles.userList} onScroll={(e) => {
              if (e.target.scrollHeight - e.target.scrollTop === e.target.clientHeight) {
                fetchUsers();
              }
            }}>
              {users?.map((user) => (
                <li key={user.ID} className={styles.userItem}>
                  {user.Firstname} {user.Lastname}
                  {invitedUsers[user.ID] ? (
                    <button disabled className={styles.inviteButton}>Done</button>
                  ) : (
                    <button className={styles.inviteButton} onClick={() => inviteUser(user.ID)}>Invite</button>
                  )}
                </li>
              ))}
            </ul>
            <button className={styles.closeButton} onClick={() => setShowPopup(false)}>Close</button>
          </div>
        </div>
      )}
    </div>
  );
}