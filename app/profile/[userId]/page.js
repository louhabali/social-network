import ProfileClient from "./ProfileClient";
import { notFound } from "next/navigation";
import { cookies } from 'next/headers';


async function getProfile(userId) {
  const cookieStore = cookies();
  const sessionToken = cookieStore.get('Token')?.value;
  
  const res = await fetch(`http://localhost:8080/api/profile/${userId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      Cookie: sessionToken ? `Token=${sessionToken}` : "",
    },
    cache: "no-store",
  });
  if (!res.ok) return null;
  return res.json();
}

export default async function ProfilePage({ params }) {
  const profile = await getProfile(params.userId);

  return <ProfileClient profile={profile} userId={params.userId} />;
}
