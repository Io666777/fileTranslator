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

// Обработчик кликов по навигационным ссылкам
document.addEventListener('click', function(e) {
    if (e.target.matches('.nav-button')) {
        e.preventDefault();
        const href = e.target.getAttribute('href');
        showSection(href);
        history.pushState(null, null, href);
    }
});

// При загрузке страницы
window.addEventListener('DOMContentLoaded', () => {
    const hash = location.hash || '#translator';
    showSection(hash);
});

// При изменении хеша в URL
window.addEventListener('hashchange', () => {
    showSection(location.hash);
});

// Для поддержки кнопок "Назад/Вперед" браузера
window.addEventListener('popstate', () => {
    showSection(location.hash);
});

const removeFileBtn = document.getElementById('removeFileBtn');

removeFileBtn.addEventListener('click', () => {
    // Сброс выбранного файла
    fileInput.value = '';

    // Очистка текста
    fileName.textContent = '';
    fileSize.textContent = '';

    // Скрываем блок fileInfo
    fileInfo.hidden = true;

    // Показываем блок загрузки
    uploadArea.hidden = false;
});



//////////////////////////////////////////////////////////////////

// Получаем элементы DOM
const uploadArea = document.getElementById('uploadArea');
const fileInput = document.getElementById('fileInput');
const fileInfo = document.getElementById('fileInfo');
const fileName = document.getElementById('fileName');
const fileSize = document.getElementById('fileSize');
const selectFileBtn = document.querySelector('.select-file-btn');

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
}

// Функция для обработки выбранного файла
function handleFileSelect(file) {
    if (file) {
        displayFileInfo(file);
        // Здесь можно добавить логику для загрузки файла на сервер
        console.log('Выбран файл:', file);
    }
}

// Обработчик клика по кнопке выбора файла
selectFileBtn.addEventListener('click', function() {
    fileInput.click();
});

// Обработчик изменения input файла
fileInput.addEventListener('change', function(e) {
    const file = e.target.files[0];
    handleFileSelect(file);
});

// Обработчики для drag and drop
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

// Обработчик клика по области загрузки (опционально)
uploadArea.addEventListener('click', function(e) {
    if (e.target === uploadArea) {
        fileInput.click();
    }
});