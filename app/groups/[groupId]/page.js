import { cookies } from 'next/headers';
import GroupsClient from './GroupsClient';


async function fetchGroupById(groupId) {
    const cookieStore = cookies();
    const sessionToken = cookieStore.get('Token')?.value;

    const res = await fetch(`http://localhost:8080/api/groups/${groupId}`, {
        method: "GET",
        headers: {
            Cookie: sessionToken ? `Token=${sessionToken}` : "",
        },
        cache: "no-store",
    });
    if (!res.ok) return null;
    return res.json();
}

async function fetchUserInfo() {
    const cookieStore = cookies();
    const sessionToken = cookieStore.get('Token')?.value;

    const res = await fetch("http://localhost:8080/api/userdata", {
        method: "GET",
        headers: {
            Cookie: sessionToken ? `Token=${sessionToken}` : "",
        },
        cache: "no-store",
    });
    if (!res.ok) return null;
    return res.json();
}

async function getPosts() {
    const cookieStore = cookies();
    const sessionToken = cookieStore.get('Token')?.value;
    
    const res = await fetch(`http://localhost:8080/api/posts`, {
      headers: {
        Cookie: sessionToken ? `Token=${sessionToken}` : "",
      },
    });
  
    if (!res.ok) return [];
    return res.json();
}

export default async function GroupDetailsPage({ params }) {
    const group = await fetchGroupById(params.groupId);
    const initialPosts = await getPosts();
    const userInfo = await fetchUserInfo();

    console.log(group,params.groupId);
    
    return (
        <div>
            <h1>{group.title}</h1>
            <p>{group.description}</p>

            <GroupsClient  groupId={params.groupId} status={group.status} initialPosts={initialPosts} userInfo={userInfo}/>
        </div>
    );
}