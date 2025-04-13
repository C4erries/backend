const fielderrors = {
    Fio: "Поле должно содержать полное ФИО (фамилия, имя, отчество) на кириллице, состоящее минимум из трех слов, разделенных пробелами; допустимы только буквы (например, Иванов Иван Иванович).",
    Tel: "Поле должно содержать номер телефона в формате +XXXXXXXXXXX; допустимы только цифры.",
    Email: "Поле должно содержать корректный адрес электронной почты в формате example@domain.com; допустимы буквы, цифры, точки, дефисы и символ @.",
    Date: "Поле должно содержать дату в формате ГГГГ-ММ-ДД;",
    Gender: "Поле обязательно для выбора; выберите один из вариантов: \"М\" (мужской) или \"Ж\" (женский).",
    Favlangs: "Поле должно содержать хотя бы один выбранный язык программирования; можно выбрать несколько вариантов из списка.",
    Bio: "Допустимы латинские и кириллические буквы, цифры, пробельные символы и символы \";,.:-!?",
    Familiar: "Поле обязательно для подтверждения; поставьте галочку, чтобы подтвердить ознакомление с контрактом."
};

document.addEventListener("DOMContentLoaded", function() {
    addErrors();
const forma = document.getElementById("forma")

forma.addEventListener("submit", (e) => {

    e.preventDefault() // отмено перезагрузко
    forma.querySelectorAll('.error').forEach(input => {
        input.classList.remove('error');
    });
    forma.querySelectorAll('.error-message').forEach(el => {
        el.style.display = 'none';
    });
    
    const formData = new FormData(e.target)
    
    const validationResult = validate(formData);
    if (validationResult) {
        // Показ ошибок
        
        showErrors(validationResult);
    } else {
        // Отправка формы

        // Преобразуем FormData в обычный объект
        const formObject = Object.fromEntries(formData.entries());
        console.log(formObject)
        formObject.Favlangs= []
        formData.getAll("Favlangs").forEach(val => formObject.Favlangs.push(parseInt(val)))
        
        formObject.Familiar = formData.get("Familiar") || null;
        console.log(formObject)

        //здесь могла быть ваша -р-е-к-л-а-м-а- обработка ответа
        fetch(e.target.action, {
            method: e.target.method,
            headers: {
                "Content-Type": "application/json" 
            },
            body: JSON.stringify(formObject),
        }).then(res => {
            if (res.redirected) {
              window.location.href = res.url;
            }
          });
    }

});
})

// Показать ошибки
const addErrors = () =>{
    for (const [field, message] of Object.entries(fielderrors)) {
        const input = forma.querySelector(`[name="${field}"]`);
        const errorEl = document.createElement("div");
        if (!input) {
            console.warn(`Поле с name="${field}" не найдено`);
            continue;
        }
        // Добавляем классы и сообщение

        errorEl.className = "error-message";
        errorEl.textContent = fielderrors[field];
        errorEl.dataset.for = field;
        errorEl.style.color = "red";
        errorEl.style.fontSize = "0.8rem";
        errorEl.style.display = "none";
        // Вставляем сообщение перед полем
        input.parentNode.insertBefore(errorEl, input);
    }
}
const showErrors = (errors) => {
    for (const [field, message] of Object.entries(errors)) {
        console.log(field)
        const input = forma.querySelector(`[name="${field}"]`);
        const errorEl = forma.querySelector(`.error-message[data-for="${field}"]`);
        
        if (input && errorEl) {
          input.classList.add('error');
          errorEl.style.display = 'block';
        } else {
          console.warn(`Не найдены элементы для поля: ${field}`);
        }
      }
}


const validate = (data) => {
    const errors = {};
    console.log(data)
    // Валидация Fio
    if (data.get("Fio")) {
        const fioRegex = /^[A-Za-zА-Яа-яЁё\s]{1,150}$/;
        if (!fioRegex.test(data.get("Fio"))) {
            errors.Fio = "Invalid fio";
        }
    } else {
        errors.Fio = "Fio is required";
    }

    // Валидация Tel
    if (data.get("Tel")) {
        const telRegex = /^\+[0-9]{1,29}$/;
        if (!telRegex.test(data.get("Tel"))) {
            errors.Tel = "Invalid telephone";
        }
    } else {
        errors.Tel = "Tel is required";
    }

    // Валидация Email
    if (data.get("Email")) {
        const emailRegex = /^[A-Za-z0-9._%+-]{1,30}@[A-Za-z0-9.-]{1,20}\.[A-Za-z]{1,10}$/;
        if (!emailRegex.test(data.get("Email"))) {
            errors.Email = "Invalid email";
        }
    } else {
        errors.Email = "Email is required";
    }

    // Валидация Date
    if (data.get("Date")) {
        const dateRegex = /^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$/;
        if (!dateRegex.test(data.get("Date"))) {
            errors.Date = "Invalid date";
        }
    } else {
        errors.Date = "Date is required";
    }

    // Валидация Gender
    if (!["Male", "Female"].includes(data.get("Gender"))) {
        errors.Gender = "Invalid gender";
    }

    // Валидация Familiar
    if (data.get("Familiar") !== "on") {
        errors.Familiar = "Invalid familiar checkbox";
    }

    // Валидация Favlangs
    if (data.getAll("Favlangs")?.length > 0) {
        const validLangs = data.getAll("Favlangs").every(lang => {
            const num = parseInt(lang, 10);
            return !isNaN(num) && num >= 1 && num <= 11;
        });
        
        if (!validLangs) {
            errors.Favlangs = "Invalid favourite langs";
        }
    } else {
        errors.Favlangs = "Favourite langs are required";
    }

    if (data.get("Bio")) {
        const bioRegex = /^[A-Za-zА-Яа-яЁё;,.:0-9\-!?""\s]{0,}$/;
        if (!bioRegex.test(data.get("Bio"))) {
            errors.Bio = "Invalid Bio";
        }
    } 

    return Object.keys(errors).length === 0 ? null : errors;
}


