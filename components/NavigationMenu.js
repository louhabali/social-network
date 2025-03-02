import Link from "next/link";
import '../app/styles/menu.css'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHouse,faUserPlus,faEnvelope,faBell,faUser,faUserGroup,faRightFromBracket} from '@fortawesome/free-solid-svg-icons';

export default function NavigationMenu(){
    
    return (
        <nav className="navbar">
            <ul>
                <li>
                    <h1 style={{ color: "rgba(155, 155, 155, 0.68)" }}>01<span style={{ color: "rgb(255, 255, 255)" }}>Connectify</span></h1>
                </li>
                <li>
                    <Link className="menulinks" href="/"><FontAwesomeIcon icon={faHouse} /><span>Home</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/groups"><FontAwesomeIcon icon={faUserGroup} /><span>groups</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/messages"><FontAwesomeIcon icon={faEnvelope} /><span>messages</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/notifications"><FontAwesomeIcon icon={faBell} /><span>notifications</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/suggestions"><FontAwesomeIcon icon={faUserPlus} /><span>suggestions</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/profile"><FontAwesomeIcon icon={faUser} /><span>profile</span></Link>
                </li>
                <li>
                    <Link className="menulinks" href="/logout"><FontAwesomeIcon icon={faRightFromBracket} /><span>logout</span></Link>
                </li>
            </ul>
        </nav>
    );  
}