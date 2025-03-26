import { cookies } from 'next/headers';
import GroupsClient from './GroupsClient';
import ErrorPage from "@/components/ErrorPage";
import './GroupClient.css';


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
    if (!res.ok) {
        const error = {
            error: true,
            message: res.status === 404 ? "Group not found" : res.status === 500 ? "Internal server error" : `Unexpected error: ${res.status}`,
            status: res.status
        };
        return error;
    }
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

async function getPosts(groupId) {
    const cookieStore = cookies();
    const sessionToken = cookieStore.get('Token')?.value;

    const res = await fetch(`http://localhost:8080/api/posts?groupId=${groupId}`, {
        headers: {
            Cookie: sessionToken ? `Token=${sessionToken}` : "",
        },
    });

    if (!res.ok) return [];
    return res.json();
}

export default async function GroupDetailsPage({ params }) {
    const group = await fetchGroupById(params.groupId);
    const initialPosts = await getPosts(params.groupId);
    const userInfo = await fetchUserInfo();

    if (group.error) {        
        return <ErrorPage error={group} />;
    }

    return (
        <div className="groupContainer">
            <div className="groupHeader">
                <h1 className="groupTitle">{group.title}</h1>
                <p className="groupDescription">{group.description}</p>
            </div>

            <GroupsClient
                groupId={params.groupId}
                status={group.status}
                initialPosts={initialPosts}
                userInfo={userInfo}
            />
        </div>);
}