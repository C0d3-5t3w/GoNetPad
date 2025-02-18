document.addEventListener('DOMContentLoaded', () => {
    hljs.highlightAll();
    
    const conn = new WebSocket(`ws://${window.location.host}/ws`);
    const outputDiv = document.getElementById('output');

    conn.onmessage = function(evt) {
        const message = evt.data;
        const codeBlock = document.createElement('pre');
        const code = document.createElement('code');
        code.className = 'language-go';
        code.textContent = message;
        codeBlock.appendChild(code);
        
        outputDiv.appendChild(codeBlock);
        hljs.highlightElement(code);
        
        outputDiv.scrollTop = outputDiv.scrollHeight;
    };

    conn.onclose = function() {
        console.log('WebSocket connection closed');
    };
});
