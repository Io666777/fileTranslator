// auth.js - управление аутентификацией и аккаунтом

// ========== КОНФИГУРАЦИЯ ==========
const API_BASE = 'http://localhost:5500';

// ========== ПЕРЕМЕННЫЕ СОСТОЯНИЯ ==========
let currentUser = null;

// ========== API ФУНКЦИИ ==========

// Работа с пользователями
class UserAPI {
    static async createUser(userData) {
        const response = await fetch(`${API_BASE}/user`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        return await response.json();
    }

    static async getUser(id) {
        const response = await fetch(`${API_BASE}/user/${id}`);
        return await response.json();
    }

    static async updateUser(id, userData) {
        const response = await fetch(`${API_BASE}/user/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });
        return await response.json();
    }
}

// ========== УПРАВЛЕНИЕ ЭКРАНАМИ ==========

// Показать определенный экран
function showScreen(screenId) {
    // Скрываем все экраны
    const screens = ['welcomeScreen', 'registerScreen', 'loginScreen', 'accountScreen'];
    screens.forEach(screen => {
        const element = document.getElementById(screen);
        if (element) {
            element.style.display = 'none';
            element.classList.remove('visible');
            element.classList.add('hidden');
        }
    });

    // Показываем нужный экран
    const targetScreen = document.getElementById(screenId);
    if (targetScreen) {
        setTimeout(() => {
            targetScreen.style.display = 'block';
            targetScreen.classList.remove('hidden');
            targetScreen.classList.add('visible');
        }, 50);
    }
}

// ========== ОСНОВНЫЕ ФУНКЦИИ АККАУНТА ==========

// Инициализация аккаунта
function initializeAccount() {
    console.log('🎯 Инициализация секции аккаунта...');
    
    // Проверяем, есть ли пользователь в localStorage
    const userId = localStorage.getItem('userId');
    const userName = localStorage.getItem('userName');
    
    if (userId && userName) {
        // Пользователь уже вошел - показываем личный кабинет
        showAccountScreen(userName);
        loadUserDetails(userId);
    } else {
        // Показываем экран приветствия
        showScreen('welcomeScreen');
    }
    
    setupAccountEventListeners();
}

// Настройка обработчиков событий
function setupAccountEventListeners() {
    // Кнопки на экране приветствия
    document.getElementById('loginBtn')?.addEventListener('click', () => showScreen('loginScreen'));
    document.getElementById('registerBtn')?.addEventListener('click', () => showScreen('registerScreen'));
    
    // Кнопки "Назад"
    document.getElementById('backFromRegister')?.addEventListener('click', () => showScreen('welcomeScreen'));
    document.getElementById('backFromLogin')?.addEventListener('click', () => showScreen('welcomeScreen'));
    
    // Форма регистрации
    document.getElementById('registerForm')?.addEventListener('submit', handleRegister);
    
    // Форма входа
    document.getElementById('loginForm')?.addEventListener('submit', handleLogin);
    
    // Кнопка выхода
    document.getElementById('logoutBtn')?.addEventListener('click', handleLogout);
}

// Обработка регистрации
async function handleRegister(e) {
    e.preventDefault();
    
    const userData = {
        name: document.getElementById('regName').value,
        password: document.getElementById('regPassword').value
    };
    
    try {
        console.log('🔄 Регистрирую нового пользователя...');
        const newUser = await UserAPI.createUser(userData);
        
        // Сохраняем данные пользователя
        localStorage.setItem('userId', newUser.id);
        localStorage.setItem('userName', newUser.name);
        
        // Показываем личный кабинет
        showAccountScreen(newUser.name);
        loadUserDetails(newUser.id);
        
        console.log('✅ Пользователь зарегистрирован:', newUser.name);
        
    } catch (error) {
        console.error('❌ Ошибка регистрации:', error);
        alert('Ошибка регистрации: ' + error.message);
    }
}

// Обработка входа
async function handleLogin(e) {
    e.preventDefault();
    
    const name = document.getElementById('loginName').value;
    const password = document.getElementById('loginPassword').value;
    
    try {
        console.log('🔄 Вход в аккаунт...');
        
        // TODO: Реализовать нормальную аутентификацию по имени
        // Пока используем заглушку - создаем нового пользователя
        const userData = {
            name: name,
            password: password
        };
        
        const user = await UserAPI.createUser(userData);
        
        // Сохраняем данные пользователя
        localStorage.setItem('userId', user.id);
        localStorage.setItem('userName', user.name);
        
        // Показываем личный кабинет
        showAccountScreen(user.name);
        loadUserDetails(user.id);
        
        console.log('✅ Вход выполнен');
        
    } catch (error) {
        console.error('❌ Ошибка входа:', error);
        alert('Ошибка входа: ' + error.message);
    }
}

// Показать личный кабинет
function showAccountScreen(userName) {
    const greetingText = document.getElementById('greetingText');
    if (greetingText) {
        greetingText.textContent = `Привет, ${userName}!`;
    }
    showScreen('accountScreen');
}

// Загрузка деталей пользователя
async function loadUserDetails(userId) {
    try {
        const user = await UserAPI.getUser(userId);
        displayAccountDetails(user);
    } catch (error) {
        console.error('Ошибка загрузки деталей:', error);
        displayAccountDetails(null);
    }
}

// Отображение деталей аккаунта
function displayAccountDetails(user) {
    const accountDetails = document.getElementById('accountDetails');
    if (!accountDetails) return;

    if (user) {
        accountDetails.innerHTML = `
            <p><strong>ID:</strong> ${user.id}</p>
            <p><strong>Имя пользователя:</strong> ${user.name}</p>
            <p><strong>Дата регистрации:</strong> ${new Date().toLocaleDateString()}</p>
            <p><strong>Статус:</strong> <span style="color: #28a745;">Активен</span></p>
        `;
    } else {
        accountDetails.innerHTML = '<p style="color: #dc3545;">Не удалось загрузить информацию</p>';
    }
}

// Выход из аккаунта
function handleLogout() {
    // Очищаем localStorage
    localStorage.removeItem('userId');
    localStorage.removeItem('userName');
    
    // Сбрасываем формы
    document.getElementById('registerForm')?.reset();
    document.getElementById('loginForm')?.reset();
    
    // Показываем экран приветствия
    showScreen('welcomeScreen');
    
    console.log('👋 Выход из аккаунта');
}

// Проверка авторизации (для других модулей)
function isUserLoggedIn() {
    return !!localStorage.getItem('userId');
}

function getCurrentUserId() {
    return localStorage.getItem('userId');
}

function getCurrentUserName() {
    return localStorage.getItem('userName');
}

// ========== ЭКСПОРТ ФУНКЦИЙ ДЛЯ ИСПОЛЬЗОВАНИЯ В ДРУГИХ ФАЙЛАХ ==========
window.AuthManager = {
    initializeAccount,
    isUserLoggedIn,
    getCurrentUserId,
    getCurrentUserName,
    handleLogout
};