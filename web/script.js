// Конфигурация
const COMPONENTS_PATH = './components/';
const API_BASE = 'http://localhost:5500';

// ========== СИСТЕМА КОМПОНЕНТОВ И НАВИГАЦИИ ==========

// Загрузка компонентов
async function loadComponent(componentName, targetElementId) {
    try {
        const response = await fetch(`${COMPONENTS_PATH}${componentName}.html`);
        const html = await response.text();
        document.getElementById(targetElementId).innerHTML = html;
    } catch (error) {
        console.error(`Ошибка загрузки компонента ${componentName}:`, error);
    }
}

// Загрузка страницы
async function loadPage(pageName) {
    try {
        const response = await fetch(`${COMPONENTS_PATH}${pageName}.html`);
        const html = await response.text();
        document.getElementById('page-content').innerHTML = html;
        
        // Инициализация страницы
        initializePage(pageName);
    } catch (error) {
        console.error(`Ошибка загрузки страницы ${pageName}:`, error);
    }
}

// Показать секцию (твоя существующая функция)
function showSection(id) {
    // Скрываем все секции
    document.querySelectorAll('.page').forEach(section => {
        section.classList.remove('active');
    });

    // Убираем активный класс со всех кнопок
    document.querySelectorAll('.nav-button').forEach(button => {
        button.classList.remove('active');
    });

    // Показываем целевую секцию
    const target = document.querySelector(id);
    if (target) {
        target.classList.add('active');
    }

    // Активируем соответствующую кнопку
    const activeButton = document.querySelector(`[href="${id}"]`);
    if (activeButton) {
        activeButton.classList.add('active');
    }
}

// Навигация
function setupNavigation() {
    // Обработчики для кнопок навигации
    document.addEventListener('click', function(e) {
        if (e.target.matches('.nav-button')) {
            e.preventDefault();
            const href = e.target.getAttribute('href');
            showSection(href);
            history.pushState(null, null, href);
        }
    });

    // При изменении хеша в URL
    window.addEventListener('hashchange', () => {
        showSection(location.hash);
    });

    // Для поддержки кнопок "Назад/Вперед" браузера
    window.addEventListener('popstate', () => {
        showSection(location.hash);
    });
}

// Инициализация страницы
function initializePage(pageName) {
    switch(pageName) {
        case 'translator':
            initializeTranslator();
            break;
        case 'account':
            initializeAccount();
            break;
        case 'library':
            initializeLibrary();
            break;
    }
}

// Инициализация приложения
async function initializeApp() {
    // Загружаем статические компоненты
    await loadComponent('header', 'header');
    await loadComponent('nav', 'navigation');
    
    // Настраиваем навигацию
    setupNavigation();
    
    // Загружаем начальную страницу
    const initialPage = window.location.hash.substring(1) || 'translator';
    await loadPage(initialPage);
    
    // Показываем начальную секцию
    showSection('#' + initialPage);
}

// ========== СЕКЦИЯ ПЕРЕВОДЧИКА ==========

// Инициализация переводчика
function initializeTranslator() {
    // Получаем элементы DOM
    const uploadArea = document.getElementById('uploadArea');
    const fileInput = document.getElementById('fileInput');
    const fileInfo = document.getElementById('fileInfo');
    const fileName = document.getElementById('fileName');
    const fileSize = document.getElementById('fileSize');
    const selectFileBtn = document.querySelector('.select-file-btn');
    const removeFileBtn = document.getElementById('removeFileBtn');

    // Функция для форматирования размера файла
    function formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    // Функция для отображения информации о файле
    function displayFileInfo(file) {
        fileName.textContent = file.name;
        fileSize.textContent = formatFileSize(file.size);
        fileInfo.hidden = false;
        uploadArea.hidden = true;
    }

    // Функция для обработки выбранного файла
    function handleFileSelect(file) {
        if (file) {
            displayFileInfo(file);
            console.log('Выбран файл:', file);
        }
    }

    // Обработчик клика по кнопке выбора файла
    if (selectFileBtn) {
        selectFileBtn.addEventListener('click', function() {
            fileInput.click();
        });
    }

    // Обработчик изменения input файла
    if (fileInput) {
        fileInput.addEventListener('change', function(e) {
            const file = e.target.files[0];
            handleFileSelect(file);
        });
    }

    // Обработчики для drag and drop
    if (uploadArea) {
        uploadArea.addEventListener('dragover', function(e) {
            e.preventDefault();
            uploadArea.classList.add('drag-over');
        });

        uploadArea.addEventListener('dragleave', function(e) {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
        });

        uploadArea.addEventListener('drop', function(e) {
            e.preventDefault();
            uploadArea.classList.remove('drag-over');
            
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                handleFileSelect(files[0]);
            }
        });

        uploadArea.addEventListener('click', function(e) {
            if (e.target === uploadArea) {
                fileInput.click();
            }
        });
    }

    // Обработчик кнопки очистки
    if (removeFileBtn) {
        removeFileBtn.addEventListener('click', () => {
            // Сброс выбранного файла
            if (fileInput) fileInput.value = '';

            // Очистка текста
            if (fileName) fileName.textContent = '';
            if (fileSize) fileSize.textContent = '';

            // Скрываем блок fileInfo
            if (fileInfo) fileInfo.hidden = true;

            // Показываем блок загрузки
            if (uploadArea) uploadArea.hidden = false;
        });
    }

    // Кнопка перевода
    const translateBtn = document.querySelector('.buttonTranslate');
    if (translateBtn) {
        translateBtn.addEventListener('click', handleTranslation);
    }
}

// ========== СЕКЦИЯ АККАУНТА ==========

// Инициализация аккаунта
function initializeAccount() {
    // Загружаем текущего пользователя
    loadCurrentUser();
    
    // Форма обновления
    const updateForm = document.getElementById('updateUserForm');
    if (updateForm) {
        updateForm.addEventListener('submit', function(e) {
            e.preventDefault();
            const updateData = {
                name: document.getElementById('updateName').value,
                email: document.getElementById('updateEmail').value,
                password: document.getElementById('updatePassword').value
            };
            updateUserInfo(updateData);
        });
    }
}

// ========== СЕКЦИЯ БИБЛИОТЕКИ ==========

// Инициализация библиотеки
function initializeLibrary() {
    loadUserFiles();
}

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

// Работа с файлами
class FileAPI {
    static async createFile(fileData) {
        const response = await fetch(`${API_BASE}/file`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(fileData)
        });
        return await response.json();
    }

    static async getUserFiles(userId) {
        const response = await fetch(`${API_BASE}/files/${userId}`);
        return await response.json();
    }

    static async getFile(id) {
        const response = await fetch(`${API_BASE}/file/${id}`);
        return await response.json();
    }
}

// ========== БИЗНЕС-ЛОГИКА ==========

// Обработка перевода файла
async function handleTranslation() {
    const userId = localStorage.getItem('userId');
    if (!userId) {
        alert('Пользователь не найден. Обновите страницу.');
        return;
    }

    const fileName = document.getElementById('fileName')?.textContent;
    
    if (!fileName || !fileName.trim()) {
        alert('Пожалуйста, выберите файл для перевода');
        return;
    }

    // Получаем выбранные настройки
    const langFrom = document.querySelector('.right_select_language select')?.value || 'auto';
    const langTo = document.querySelector('.right_select_translate select')?.value || 'en';
    const format = document.querySelectorAll('.right_select_translate select')[1]?.value || 'same';

    // Создаем запись о файле
    const fileData = {
        namefile: fileName,
        author: userId,
        cloudkey: `cloud_${Date.now()}`,
        langfrom: langFrom,
        langto: langTo
    };

    try {
        const newFile = await FileAPI.createFile(fileData);
        alert(`Файл "${fileName}" отправлен на перевод! ID: ${newFile.id}`);
        
        // Обновляем библиотеку
        loadUserFiles();
    } catch (error) {
        console.error('Ошибка при создании файла:', error);
        alert('Ошибка при отправке файла на перевод');
    }
}

// Загрузка файлов пользователя для библиотеки
async function loadUserFiles() {
    const userId = localStorage.getItem('userId');
    if (!userId) return;

    try {
        const files = await FileAPI.getUserFiles(userId);
        displayFilesInLibrary(files);
    } catch (error) {
        console.error('Ошибка загрузки файлов:', error);
    }
}

// Отображение файлов в библиотеке
function displayFilesInLibrary(files) {
    const libraryContent = document.querySelector('.library-content');
    if (!libraryContent) return;

    if (files && files.length > 0) {
        libraryContent.innerHTML = files.map(file => `
            <div class="file-item">
                <h3>${file.namefile}</h3>
                <p>Перевод: ${file.langfrom} → ${file.langto}</p>
                <p>Статус: В процессе</p>
                <small>ID: ${file.id}</small>
            </div>
        `).join('');
    } else {
        libraryContent.innerHTML = '<p>У вас пока нет переведенных файлов</p>';
    }
}

// Функции для работы с пользователями
async function loadCurrentUser() {
    const userId = localStorage.getItem('userId');
    if (!userId) {
        showTestResult('Пользователь не найден. Создайте нового.');
        return;
    }

    try {
        const user = await UserAPI.getUser(userId);
        displayUserInfo(user);
        showTestResult(`Пользователь загружен: ${user.name} (${user.email})`);
    } catch (error) {
        showTestResult('Ошибка загрузки пользователя: ' + error.message);
    }
}

// Создание тестового пользователя
async function createTestUser() {
    const testUserData = {
        name: `Тестовый пользователь ${Date.now()}`,
        email: `test${Date.now()}@example.com`,
        password: 'test123'
    };

    try {
        const newUser = await UserAPI.createUser(testUserData);
        showTestResult(`Создан пользователь: ${newUser.name} (ID: ${newUser.id})`);
        
        // Сохраняем как текущего
        localStorage.setItem('userId', newUser.id);
        localStorage.setItem('userName', newUser.name);
        displayUserInfo(newUser);
    } catch (error) {
        showTestResult('Ошибка создания пользователя: ' + error.message);
    }
}

// Обновление информации о пользователе
async function updateUserInfo(userData) {
    const userId = localStorage.getItem('userId');
    if (!userId) {
        alert('Пользователь не найден');
        return;
    }

    try {
        const updatedUser = await UserAPI.updateUser(userId, userData);
        displayUserInfo(updatedUser);
        showTestResult('Данные пользователя обновлены!');
    } catch (error) {
        showTestResult('Ошибка обновления: ' + error.message);
    }
}

// Отображение информации о пользователе
function displayUserInfo(user) {
    const userInfoDiv = document.getElementById('userInfo');
    if (userInfoDiv && user) {
        userInfoDiv.innerHTML = `
            <p><strong>ID:</strong> ${user.id}</p>
            <p><strong>Имя:</strong> ${user.name}</p>
            <p><strong>Email:</strong> ${user.email}</p>
            <p><strong>Пароль:</strong> ${user.password}</p>
            <p><strong>Статус:</strong> Активен</p>
        `;
    }
}

// Показать результаты тестов
function showTestResult(message) {
    const resultsDiv = document.getElementById('testResults');
    if (resultsDiv) {
        resultsDiv.innerHTML = `<p>${message}</p>`;
    }
}

// Инициализация пользователя при загрузке
async function initializeUser() {
    let userId = localStorage.getItem('userId');
    
    if (!userId) {
        // Создаем нового пользователя при первом посещении
        const userData = {
            name: 'Новый пользователь',
            email: `user_${Date.now()}@example.com`,
            password: 'temp_password'
        };
        
        try {
            const newUser = await UserAPI.createUser(userData);
            localStorage.setItem('userId', newUser.id);
            localStorage.setItem('userName', newUser.name);
        } catch (error) {
            console.error('Ошибка создания пользователя:', error);
        }
    }
}

// ========== ЗАПУСК ПРИЛОЖЕНИЯ ==========

// Запуск приложения
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    initializeUser(); // Инициализация пользователя
});