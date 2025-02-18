document.addEventListener("DOMContentLoaded", () => {
    if (window.innerWidth <= 768) {
        const textArea = document.getElementById("textArea");
        textArea.style.width = "100%";
        textArea.style.height = "calc(100vh - 20px)";
        textArea.style.fontSize = "16px";
    }
});
