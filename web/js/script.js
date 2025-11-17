// Конфигурация
const COMPONENTS_PATH = './components/';

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
    console.log('🎯 Загружаю страницу:', pageName);
    
    try {
        const response = await fetch(`${COMPONENTS_PATH}${pageName}.html`);
        const html = await response.text();
        document.getElementById('page-content').innerHTML = html;
        console.log('✅ Страница загружена:', pageName);
        
        // Ждем обновления DOM перед показом секции
        setTimeout(() => {
            showSection('#' + pageName);
            // Инициализация страницы
            initializePage(pageName);
        }, 10);
        
    } catch (error) {
        console.error(`❌ Ошибка загрузки страницы ${pageName}:`, error);
    }
}

// Показать секцию
function showSection(id) {
    console.log('🔄 Пытаюсь показать секцию:', id);
    
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
        console.log('✅ Секция найдена и активирована:', id);
    } else {
        console.error('❌ Секция не найдена:', id);
    }

    // Активируем соответствующую кнопку
    const activeButton = document.querySelector(`[href="${id}"]`);
    if (activeButton) {
        activeButton.classList.add('active');
        console.log('✅ Кнопка активирована:', id);
    } else {
        console.error('❌ Кнопка не найдена:', id);
    }
}

// Навигация
function setupNavigation() {
    // Обработчики для кнопок навигации
    document.addEventListener('click', function(e) {
        if (e.target.matches('.nav-button')) {
            e.preventDefault();
            const href = e.target.getAttribute('href');
            const pageName = href.substring(1);
            
            loadPage(pageName);
            history.pushState(null, null, href);
        }
    });

    // При изменении хеша в URL
    window.addEventListener('hashchange', () => {
        const pageName = location.hash.substring(1) || 'translator';
        loadPage(pageName);
    });

    // Для поддержки кнопок "Назад/Вперед" браузера
    window.addEventListener('popstate', () => {
        const pageName = location.hash.substring(1) || 'translator';
        loadPage(pageName);
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
            if (fileInput) fileInput.value = '';
            if (fileName) fileName.textContent = '';
            if (fileSize) fileSize.textContent = '';
            if (fileInfo) fileInfo.hidden = true;
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

// Инициализация аккаунта (теперь используем AuthManager из auth.js)
function initializeAccount() {
    if (window.AuthManager && window.AuthManager.initializeAccount) {
        window.AuthManager.initializeAccount();
    } else {
        console.error('❌ AuthManager не загружен');
        // Fallback: показываем простой текст
        document.getElementById('page-content').innerHTML = '<p>Модуль аутентификации не загружен</p>';
    }
}

// ========== СЕКЦИЯ БИБЛИОТЕКИ ==========

// Инициализация библиотеки
function initializeLibrary() {
    loadUserFiles();
}

// ========== ФУНКЦИИ ДЛЯ ФАЙЛОВ ==========

// Обработка перевода файла
async function handleTranslation() {
    // Используем AuthManager для проверки авторизации
    if (!window.AuthManager || !window.AuthManager.isUserLoggedIn()) {
        alert('❌ Войдите в аккаунт для перевода файлов');
        return;
    }

    const userId = window.AuthManager.getCurrentUserId();
    const fileName = document.getElementById('fileName')?.textContent;
    
    if (!fileName || !fileName.trim()) {
        alert('Пожалуйста, выберите файл для перевода');
        return;
    }

    // Получаем выбранные настройки
    const langFrom = document.querySelector('.right_select_language select')?.value || 'auto';
    const langTo = document.querySelector('.right_select_translate select')?.value || 'en';

    // Создаем запись о файле
    const fileData = {
        namefile: fileName,
        author: userId,
        cloudkey: `cloud_${Date.now()}`,
        langfrom: langFrom,
        langto: langTo
    };

    try {
        // TODO: Перенести FileAPI в отдельный модуль
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
    if (!window.AuthManager || !window.AuthManager.isUserLoggedIn()) {
        return;
    }

    const userId = window.AuthManager.getCurrentUserId();
    
    try {
        // TODO: Перенести FileAPI в отдельный модуль
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

// ========== ЗАПУСК ПРИЛОЖЕНИЯ ==========

// Запуск приложения
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    // Инициализация пользователя теперь в auth.js
});