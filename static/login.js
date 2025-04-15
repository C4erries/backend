document.addEventListener("DOMContentLoaded", function() {
const forma = document.getElementById("loginForm")
forma.addEventListener("submit", (e) => {

    e.preventDefault() // отмено перезагрузко

    const formData = new FormData(e.target)
        // Преобразуем FormData в обычный объект
    const formObject = Object.fromEntries(formData.entries());
    fetch(e.target.action, {
        method: e.target.method,
        headers: {
            "Content-Type": "application/json" 
        },
        body: JSON.stringify(formObject),
    }).then(res => {
        if(res.redirected){
            window.location.href = res.url;
        }
    });
})
})




