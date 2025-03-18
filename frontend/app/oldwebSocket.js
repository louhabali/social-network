export const websoketProvider = {
    connect() {
        return new Promise((resolve, reject) => {
            // if (!this.ws) {
            //     const ws = new WebSocket(`ws://localhost:8080/api/ws`);

            //     ws.onopen = () => {
            //         console.log("ws connected");
            //         this.ws = ws;
            //         window.websoketProvider = websoketProvider;
            //         window.websoketProvider = null;
            //         resolve(ws);
            //     }

            //     ws.onclose = () => {
            //         console.log("ws closed");
            //         reject(new Error("WebSocket connection closed unexpectedly"));
            //     }

            //     ws.onerror = (error) => {
            //         console.log("ws error:", error);
            //         reject(new Error("Failed to connect to WebSocket"));
            //     }
                
            // } else resolve(this.ws)
        })
    },

    disconnect() {
        if (this.ws) this.ws.close();
        window.websoketProvider = null;
    },

    sendMessage(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.log("WebSocket not connected");
        }
    }
}

