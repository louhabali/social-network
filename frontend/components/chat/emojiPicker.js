import {useRef, useEffect} from "react"
import styles from "./picker.module.css"

export default function Picker({setMessage, message}) {
    const emojis = useRef(["😀", "😄", "🤑", "🥵", "🤬", "🤯", "🥶", "🥰",
        "😭","😱","😖","😣","😞","😓","😩","😫","🥱","😤","😡", "🤚", "🤣"
        ,"👽","👾","🤖","😺","😸","😹","😻","😼","😽","🙀","😿"
        ,"🤚","🖐️","✋","🖖","👌","🤏","✌️","🤞"
        ,"💕","💟","❣️","💔","❤️","🧡","💛","💚","💙"
        ,"🐵","🐒","🦍","🦧"
    ])
    const emojisContainer = useRef(null)
    
    useEffect(() => {
        const handleClick = (e) => {
            if (e.target != emojisContainer.current) {
                console.log(e.target);
                setMessage((prevMessage) => prevMessage + e.target.innerText);
            }
        };

        if (emojisContainer.current) {
            emojisContainer.current.addEventListener("click", (e) => {
                e.stopPropagation();
                handleClick(e);
            });
        }

        return () => {
            if (emojisContainer.current) {
                emojisContainer.current.removeEventListener("click", handleClick);
            }
        };
    }, [])

    return(
        <div className={styles.picker}  ref={emojisContainer}>
            {
                emojis.current.map((emoji, i) => (
                    <span className={styles.emoji} key={i}>{emoji}</span>
                ))
            }
        </div>
    )
}