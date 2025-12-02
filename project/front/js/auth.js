// auth.js
let token = localStorage.getItem('token');

// Очистка полей аутентификации
function clearAuthFields() {
    document.getElementById('username').value = '';
    document.getElementById('password').value = '';
}

// Регистрация
async function signUp() {
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;
    const name = username;

    if (!username || !password) {
        showNotification('Заполните все поля', 'error');
        return;
    }

    if (password.length < 6) {
        showNotification('Пароль должен быть не менее 6 символов', 'error');
        return;
    }

    const user = {
        username: username,
        password: password,
        name: name
    };

    try {
        const response = await fetch('/auth/sign-up', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(user)
        });

        if (response.ok) {
            showNotification('Регистрация успешна! Теперь войдите в систему', 'success');
            clearAuthFields();
        } else {
            const error = await response.json();
            showNotification(`Ошибка регистрации: ${error.message}`, 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Вход
async function signIn() {
    const username = document.getElementById('username').value.trim();
    const password = document.getElementById('password').value;

    if (!username || !password) {
        showNotification('Заполните все поля', 'error');
        return;
    }

    const credentials = {
        username: username,
        password: password
    };

    try {
        const response = await fetch('/auth/sign-in', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(credentials)
        });

        if (response.ok) {
            const data = await response.json();
            token = data.token;
            localStorage.setItem('token', token);
            showNotification('Вход выполнен успешно!', 'success');
            clearAuthFields();
            showMainSection();
        } else {
            showNotification('Неверный логин или пароль', 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Выход
function logout() {
    localStorage.removeItem('token');
    token = null;
    document.getElementById('auth-section').classList.remove('hidden');
    document.getElementById('main-section').classList.add('hidden');
    clearAuthFields();
    showNotification('Вы вышли из системы', 'info');
}

// Показать основную секцию
function showMainSection() {
    document.getElementById('auth-section').classList.add('hidden');
    document.getElementById('main-section').classList.remove('hidden');
    loadFiles();
}

export { token, signUp, signIn, logout, showMainSection };