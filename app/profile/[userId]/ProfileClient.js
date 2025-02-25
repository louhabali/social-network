"use client";

import { useState } from "react";
import Post from "@/components/post";
import Image from "next/image";

export default function ProfileClient({ profile, userId }) {
  const [isPublic, setIsPublic] = useState(profile.user.IsPublic);
  const [isFollowing, setIsFollowing] = useState(profile.isFollowing);
  const [activeTab, setActiveTab] = useState("posts");

  const isOwnProfile = profile.user.ID === profile.userid;


  async function handleProfilePrivacyToggle() {
    const newStatus = !isPublic;
    const res = await fetch(`http://localhost:8080/api/profile/${userId}/privacy`, {
      method: "POST",
      credentials: "include",
      body: JSON.stringify({ isPublic: newStatus }),
    });

    if (res.ok) {
      setIsPublic(newStatus);
    }
  }

  async function onFollowToggle() {
    const method = isFollowing ? "DELETE" : "POST";
    const res = await fetch(`http://localhost:8080/api/profile/${userId}/follow`, {
      method,
      credentials: "include",
    });
    
    if (res.ok) {
      setIsFollowing(!isFollowing)
    }
  }
  


  return (
    <div>
      <div className="profile-header">
        <div className="profile-info">
          {profile.user.Avatar ? (
            <Image src={`data:image/png;base64,${profile.user.Avatar}`} alt="Avatar" width={64} height={64} />
          ) : (
            <Image src="/image/no-avatar.png" alt="No Avatar" width={64} height={64} />
          )}
          <h1>{profile.user.Firstname} {profile.user.Lastname}</h1>
          <p><strong>Email:</strong> {profile.user.Email}</p>
          <p>{profile.user.AboutMe || "No bio available"}</p>
          <p><strong>Date of Birth:</strong> {profile.user.DateOfBirth}</p>
          <p><strong>Followers:</strong> {profile.followers?.length || 0}</p>
          <p><strong>Following:</strong> {profile.following?.length || 0}</p>
          <p><strong>Posts:</strong> {profile.posts?.length || 0}</p>

          {isOwnProfile ? (
            <button onClick={handleProfilePrivacyToggle}>
              {isPublic ? "Switch to Private Profile" : "Switch to Public Profile"}
            </button>
          ) : (
            <button onClick={onFollowToggle} className="follow-button">
              {isFollowing ? 'Unfollow' : 'Follow'}
            </button>
          )}
        </div>

        <div className="profile-tabs">
          <button className={activeTab === "posts" ? "active" : ""} onClick={() => setActiveTab("posts")}>
            Posts
          </button>
          <button className={activeTab === "followers" ? "active" : ""} onClick={() => setActiveTab("followers")}>
            Followers
          </button>
          <button className={activeTab === "following" ? "active" : ""} onClick={() => setActiveTab("following")}>
            Following
          </button>
        </div>
      </div>

      <div className="profile-content">
        {activeTab === "posts" && (
          <div className="profile-posts">
            {(profile.posts?.length || 0) === 0 ? <p>No posts yet!</p> : profile.posts.map((post) => <Post key={post.Pid} post={post} />)}
          </div>
        )}

        {activeTab === "followers" && (
          <div className="profile-followers">
            <h3>Followers</h3>
            {profile.followers?.map((follower) => (
              <div key={follower.ID}>{follower.Firstname} {follower.Lastname}</div>
            ))}
          </div>
        )}

        {activeTab === "following" && (
          <div className="profile-following">
            <h3>Following</h3>
            {profile.following?.map((followingUser) => (
              <div key={followingUser.ID}>{followingUser.Firstname} {followingUser.Lastname}</div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
