import {useRef, useEffect} from "react"
import styles from "./picker.module.css"

export default function Picker({setMessage}) {
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
                setMessage((prevMessage) => prevMessage + e.target.innerText);
            }
        };

        emojisContainer.current.addEventListener("click", handleClick);

        return () => {
            if (emojisContainer.current) {
                emojisContainer.current.removeEventListener("click", handleClick);
            }
        };
    })

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