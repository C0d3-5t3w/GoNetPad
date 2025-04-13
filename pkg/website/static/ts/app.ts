declare const hljs: {
    highlightAll: () => void;
    highlightElement: (element: HTMLElement) => void;
};

document.addEventListener("DOMContentLoaded", () => {
    hljs.highlightAll();
    
    const textArea = document.getElementById("textArea") as HTMLTextAreaElement;
    
    if (textArea && window.innerWidth <= 768) {
        textArea.style.width = "100%";
        textArea.style.height = "calc(100vh - 20px)";
        textArea.style.fontSize = "16px";
    }

    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    const ws = new WebSocket(`${wsProtocol}//${wsHost}/ws`);
    const outputDiv = document.getElementById('output');

    ws.onmessage = (event: MessageEvent) => {
        try {
            if (typeof event.data === 'string') {
                const data: string = event.data;
                
                if (data.startsWith('iVBOR')) {
                    return;
                }
                
                if (data.startsWith('[{') && data.endsWith('}]')) {
                    console.error('Received error message:', data);
                    return;
                }
                
                if (textArea) {
                    textArea.value = data;
                }
                
                if (outputDiv) {
                    const codeBlock = document.createElement('pre');
                    const code = document.createElement('code');
                    code.className = 'language-go';
                    code.textContent = data;
                    codeBlock.appendChild(code);
                    
                    outputDiv.appendChild(codeBlock);
                    hljs.highlightElement(code);
                    
                    outputDiv.scrollTop = outputDiv.scrollHeight;
                }
            }
        } catch (error) {
            console.error('Error processing message:', error instanceof Error ? error.message : String(error));
        }
    };

    ws.onclose = () => {
        console.log("WebSocket connection closed");
        if (textArea) {
            textArea.value = "Connection lost. Attempting to reconnect...";
        }
        setTimeout(() => {
            location.reload();
        }, 3000);
    };

    ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        if (textArea) {
            textArea.value = "Connection error. Please check your connection.";
        }
    };
});
