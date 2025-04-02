import { createContext, useRef, useEffect, useState } from "react";
import "./notification.css";

//+ create the context return an object to manage the values u wanna pass aroun
//+ the most importent properties of the obj returned by creatContext are:
//+ => Provider: a component used to store the values u want to pass around
//+ => Consumer: a component user to read the values in other components that
//      they are chlids of the component that provide context
//      (it not used a lot now the equivalent of it is useContext that also
//          read context values in a simple way)
//+ the creation of context sould be global as is meant to be accessed in 
// deffrent component
const WebSocketContext = createContext();

// after creating context we need a wraper component that will provide values
// to its childs 
function WebSocketProvider({ children }) {
    // here the useRef as its value stay presist and does not affect rendring 
    // also rather than use a regular var as variables when component re-render
    //  they re-initialized
    const ws = useRef(null);

    const [notification, setNotification] = useState(null);
    const timeoutRef = useRef(null);

    const showNotification = (message, type) => {
        setNotification({ message, type });
        if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
        }
        timeoutRef.current = setTimeout(() => {
            setNotification(null);
        }, 3000);
    };

    const sendMessage = (msg) => {
        ws.current.send(JSON.stringify(msg))
    }

    useEffect(() => {
        // the current property of ref obj is what actually hold the referenced thing
        ws.current = new WebSocket(`ws://localhost:8080/api/ws`);
        
        ws.current.onopen = () => {
            console.log("ws connected");   
        }

        ws.current.onclose = () => {
            console.log("ws disconnected")
        }
    }, [])

    // now we need to store the values we want them accessible in Provider
    // component witch is a property of context obj
    return (
        <WebSocketContext.Provider value={{ws, sendMessage, showNotification}}>
            {children}
            {notification && (
                <div className={`notification ${notification.type}`}>
                    {notification.message}
                </div>
            )}
        </WebSocketContext.Provider>
    )
}

export {WebSocketContext, WebSocketProvider}