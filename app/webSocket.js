export const websoketProvider = {
    connect() {
        return new Promise((resolve, reject) => {
            if (!this.ws) {
                const ws = new WebSocket(`ws://localhost:8080/api/ws`);

                ws.onopen = () => {
                    console.log("ws connected");
                    this.ws = ws;
                    resolve(ws);
                }

                ws.onclose = () => {
                    console.log("ws closed");
                    reject(new Error("WebSocket connection closed unexpectedly"));
                }

                ws.onerror = (error) => {
                    console.log("ws error:", error);
                    reject(new Error("Failed to connect to WebSocket"));
                }

                // ws.onmessage = (event) => {
                //     this.onmessage(event)
                // }
            } else resolve(this.ws)
        })
    },

    disconnect() {
        if (this.ws) this.ws.close();
    },

    onmessage(event) {
        console.log("New message:", event.data);
    },

    sendMessage(message) {
        console.log(message);
        
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
            console.log("sent");
            
        } else {
            console.warn("WebSocket not connected");
        }
    }
}