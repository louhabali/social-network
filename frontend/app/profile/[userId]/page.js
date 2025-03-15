import ProfileClient from "./ProfileClient";
import ErrorPage from "../error";
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
  if (!res.ok) {
    // if (res.status === 404) throw new Error("Profile not found");
    // if (res.status === 405) throw new Error("Method not allowed");
    // if (res.status === 403) throw new Error("This profile is private");
    // if (res.status === 500) throw new Error("Internal server error");
    // throw new Error(`Unexpected error: ${res.status}`);
  }
  return res.json();
}

export default async function ProfilePage({ params }) {
  let profile = null;

  // try {
  profile = await getProfile(params.userId);
  // } catch (error) {
  //   return <ErrorPage error={error} />; 
  // }

  return <ProfileClient profile={profile} userId={params.userId} />;
}