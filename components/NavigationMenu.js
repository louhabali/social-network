
import Link from "next/link";
import '../app/styles/menu.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHouse,faEnvelope,faBell,faUser,faUserGroup,faRightFromBracket} from '@fortawesome/free-solid-svg-icons';

export default function NavigationMenu(){
   async function handlelogout() {
       try {
           const res = await fetch('http://localhost:8080/api/logout', {
                method: 'POST',
                credentials: "include",
                
            })
            if (res.ok) {
                window.location.href = "/login";
            } else {
                console.error("Logout failed");
            }
              
        } catch (error) {
            console.error("error:66", error);
        }

      
    }
    return (
        <nav className="navbar">
              <div style={{ paddingTop: "20px" ,paddingLeft: "15px" }}>
                    <h1 style={{ color: "rgba(155, 155, 155, 0.68)" }}>01<span style={{ color: "rgb(255, 255, 255)" }}>Connectify</span></h1>
            </div>
            <ul>
              
                <li>
                    <Link className="menulinks" href="/"><FontAwesomeIcon icon={faHouse} /><span>Home</span></Link>
                </li>
             
                <li>
                    <Link className="menulinks" href="/groupcreation"><FontAwesomeIcon icon={faUserGroup} /><span>groups</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/messages"><FontAwesomeIcon icon={faEnvelope} /><span>messages</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/notifications"><FontAwesomeIcon icon={faBell} /><span>notifications</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/profile"><FontAwesomeIcon icon={faUser} /><span>profile</span></Link>
                </li>
                <li>
                    <div className="menulinks" onClick={handlelogout} style={{ cursor:"pointer" }}><FontAwesomeIcon icon={faRightFromBracket} /><span>logout</span></div>
                </li>
            </ul>
           
        </nav>
    );  
}