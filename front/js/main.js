// Глобальные переменные
let token = localStorage.getItem('token');
let currentUser = null;

// Показать уведомление
function showNotification(message, type = 'info') {
    const colors = {
        success: '#38a169',
        error: '#f56565',
        info: '#4299e1',
        warning: '#ed8936'
    };
    
    const notification = document.createElement('div');
    notification.className = 'notification';
    notification.innerHTML = `
        <div style="
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            background: ${colors[type] || colors.info};
            color: white;
            border-radius: 8px;
            box-shadow: 0 5px 15px rgba(0,0,0,0.2);
            z-index: 1000;
            animation: slideIn 0.3s ease;
            max-width: 300px;
        ">
            <i class="fas fa-${type === 'success' ? 'check-circle' : 
                            type === 'error' ? 'exclamation-circle' : 
                            type === 'warning' ? 'exclamation-triangle' : 'info-circle'}"></i>
            ${message}
        </div>
    `;
    
    document.body.appendChild(notification);
    
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
    
    // Добавляем стили для анимации
    if (!document.querySelector('#notification-styles')) {
        const style = document.createElement('style');
        style.id = 'notification-styles';
        style.textContent = `
            @keyframes slideIn {
                from { transform: translateX(100%); opacity: 0; }
                to { transform: translateX(0); opacity: 1; }
            }
            @keyframes slideOut {
                from { transform: translateX(0); opacity: 1; }
                to { transform: translateX(100%); opacity: 0; }
            }
        `;
        document.head.appendChild(style);
    }
}

// Очистка полей аутентификации
function clearAuthFields() {
    document.getElementById('username').value = '';
    document.getElementById('password').value = '';
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    console.log('Страница загружена');
    
    if (token) {
        showMainSection();
    }
    
    // Инициализируем drag and drop
    initDragAndDrop();
    
    // Добавляем обработчик изменения файла
    const fileInput = document.getElementById('file-input');
    if (fileInput) {
        fileInput.addEventListener('change', function() {
            if (this.files.length > 0) {
                updateFileInfo(this.files[0]);
            } else {
                resetFileInfo();
            }
        });
    }
});

// Инициализация drag and drop
function initDragAndDrop() {
    const uploadArea = document.getElementById('upload-area');
    if (!uploadArea) return;
    
    // Предотвращаем стандартное поведение браузера
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        uploadArea.addEventListener(eventName, preventDefaults, false);
        document.body.addEventListener(eventName, preventDefaults, false);
    });
    
    function preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }
    
    // Подсветка при перетаскивании
    ['dragenter', 'dragover'].forEach(eventName => {
        uploadArea.addEventListener(eventName, highlight, false);
    });
    
    ['dragleave', 'drop'].forEach(eventName => {
        uploadArea.addEventListener(eventName, unhighlight, false);
    });
    
    function highlight() {
        uploadArea.classList.add('highlight');
    }
    
    function unhighlight() {
        uploadArea.classList.remove('highlight');
    }
    
    // Обработка drop
    uploadArea.addEventListener('drop', handleDrop, false);
    
    function handleDrop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        
        if (files.length > 0) {
            const fileInput = document.getElementById('file-input');
            fileInput.files = files;
            updateFileInfo(files[0]);
        }
    }
}

// Показать основную секцию
function showMainSection() {
    document.getElementById('auth-section').classList.add('hidden');
    document.getElementById('main-section').classList.remove('hidden');
    loadFiles();
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
    currentUser = null;
    document.getElementById('auth-section').classList.remove('hidden');
    document.getElementById('main-section').classList.add('hidden');
    clearAuthFields();
    showNotification('Вы вышли из системы', 'info');
}

// Загрузка файла
async function uploadFile() {
    const fileInput = document.getElementById('file-input');
    const file = fileInput.files[0];
    
    if (!file) {
        showNotification('Выберите файл для загрузки', 'warning');
        return;
    }

    // Проверка размера файла (макс 10MB)
    if (file.size > 10 * 1024 * 1024) {
        showNotification('Файл слишком большой (макс 10MB)', 'error');
        resetFileInfo();
        return;
    }

    // Проверка типа файла
    const allowedTypes = ['.txt', '.doc', '.docx', '.pdf', '.rtf', '.odt', '.html', '.htm', '.json', '.xml', '.csv', '.md'];
    const fileExtension = '.' + file.name.split('.').pop().toLowerCase();
    
    if (!allowedTypes.includes(fileExtension)) {
        showNotification(`Неподдерживаемый тип файла. Разрешены: ${allowedTypes.join(', ')}`, 'error');
        resetFileInfo();
        return;
    }

    const formData = new FormData();
    formData.append('file', file);

    try {
        const response = await fetch('/api/files/upload', {
            method: 'POST',
            headers: { 'Authorization': 'Bearer ' + token },
            body: formData
        });

        if (response.ok) {
            showNotification('Файл успешно загружен!', 'success');
            resetFileInfo();
            fileInput.value = '';
            loadFiles();
        } else {
            const error = await response.text();
            showNotification(`Ошибка загрузки: ${error}`, 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Отображение информации о выбранном файле
function updateFileInfo(file) {
    const fileInfo = document.getElementById('file-info');
    const uploadArea = document.getElementById('upload-area');
    
    if (!file) {
        resetFileInfo();
        return;
    }

    const fileSize = formatFileSize(file.size);
    const fileType = getFileType(file.name);
    const fileIcon = getFileIcon(fileType);
    
    fileInfo.innerHTML = `
        <div class="selected-file-info">
            <i class="fas ${fileIcon}" style="color: #4299e1; margin-right: 12px; font-size: 20px;"></i>
            <div style="flex: 1;">
                <div style="font-weight: 600; color: #2d3748; margin-bottom: 4px; word-break: break-all;">
                    ${file.name}
                </div>
                <div style="font-size: 13px; color: #718096; display: flex; gap: 15px; flex-wrap: wrap;">
                    <span><i class="fas fa-weight-hanging"></i> ${fileSize}</span>
                    <span><i class="fas fa-file-alt"></i> ${fileType}</span>
                    <span><i class="far fa-calendar"></i> ${new Date(file.lastModified).toLocaleDateString('ru-RU')}</span>
                </div>
            </div>
            <button onclick="removeSelectedFile()" class="btn-remove-file" title="Удалить файл">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;
    
    // Обновляем стили upload area
    if (uploadArea) {
        uploadArea.classList.add('has-file');
        const uploadIcon = uploadArea.querySelector('.upload-icon');
        if (uploadIcon) {
            uploadIcon.className = `fas ${fileIcon} upload-icon`;
        }
        const uploadText = uploadArea.querySelector('p');
        if (uploadText) {
            uploadText.textContent = 'Файл выбран. Нажмите "Загрузить" или перетащите другой файл';
        }
    }
}

// Форматирование размера файла
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Определение типа файла по расширению
function getFileType(filename) {
    const extension = filename.split('.').pop().toLowerCase();
    
    const fileTypes = {
        'txt': 'Текстовый файл',
        'doc': 'Документ Word',
        'docx': 'Документ Word',
        'pdf': 'PDF документ',
        'rtf': 'RTF документ',
        'odt': 'OpenDocument',
        'html': 'HTML файл',
        'htm': 'HTML файл',
        'json': 'JSON файл',
        'xml': 'XML файл',
        'csv': 'CSV файл',
        'md': 'Markdown'
    };
    
    return fileTypes[extension] || extension.toUpperCase();
}

// Получение иконки для типа файла
function getFileIcon(fileType) {
    const icons = {
        'Текстовый файл': 'fa-file-alt',
        'Документ Word': 'fa-file-word',
        'PDF документ': 'fa-file-pdf',
        'RTF документ': 'fa-file-alt',
        'OpenDocument': 'fa-file-alt',
        'HTML файл': 'fa-file-code',
        'JSON файл': 'fa-file-code',
        'XML файл': 'fa-file-code',
        'CSV файл': 'fa-file-csv',
        'Markdown': 'fa-file-alt'
    };
    
    return icons[fileType] || 'fa-file';
}

// Удаление выбранного файла
function removeSelectedFile() {
    const fileInput = document.getElementById('file-input');
    fileInput.value = '';
    resetFileInfo();
}

// Сброс информации о файле
function resetFileInfo() {
    const fileInfo = document.getElementById('file-info');
    const uploadArea = document.getElementById('upload-area');
    
    fileInfo.innerHTML = '';
    
    if (uploadArea) {
        uploadArea.classList.remove('has-file');
        uploadArea.classList.remove('highlight');
        const uploadIcon = uploadArea.querySelector('.upload-icon');
        if (uploadIcon) {
            uploadIcon.className = 'fas fa-cloud-upload-alt upload-icon';
        }
        const uploadText = uploadArea.querySelector('p');
        if (uploadText) {
            uploadText.textContent = 'Перетащите файл сюда или нажмите для выбора';
        }
    }
}

// Загрузка списка файлов
async function loadFiles() {
    const filesList = document.getElementById('files-list');
    
    filesList.innerHTML = '<div class="loading"><i class="fas fa-spinner fa-spin"></i> Загрузка...</div>';

    try {
        const response = await fetch('/api/files/', {
            headers: { 'Authorization': 'Bearer ' + token }
        });

        if (response.ok) {
            const files = await response.json();
            displayFiles(files);
        } else if (response.status === 401) {
            showNotification('Сессия истекла. Войдите снова.', 'error');
            logout();
        } else {
            showNotification('Ошибка загрузки файлов', 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Форматирование даты
function formatDate(dateString) {
    try {
        const date = new Date(dateString);
        return date.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
    } catch (e) {
        return dateString;
    }
}

// Отображение файлов
function displayFiles(files) {
    const filesList = document.getElementById('files-list');
    
    if (!files || files.length === 0) {
        filesList.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-folder-open" style="font-size: 48px; color: #a0aec0; margin-bottom: 15px;"></i>
                <p style="color: #718096;">Нет загруженных файлов</p>
            </div>
        `;
        return;
    }

    filesList.innerHTML = files.map(file => `
        <div class="file-item" id="file-${file.id}">
            <div class="file-info-main">
                <div class="file-title">
                    <i class="fas fa-file" style="color: #4299e1; margin-right: 10px;"></i>
                    ${file.title || 'Без названия'}
                </div>
                <div class="file-meta">
                    <span><i class="far fa-calendar"></i> ${formatDate(file.created_at)}</span>
                    <span class="file-status status-${file.status || 'uploaded'}">
                        ${getStatusText(file.status)}
                    </span>
                </div>
            </div>
            <div class="file-actions-buttons">
                ${file.status === 'translated' || file.status === 'completed' ? 
                    `<button onclick="downloadFile(${file.id})" class="btn-action btn-download" title="Скачать">
                        <i class="fas fa-download"></i>
                    </button>` : ''
                }
                ${file.status === 'uploaded' ? 
                    `<button onclick="translateFile(${file.id})" class="btn-action btn-translate" title="Перевести">
                        <i class="fas fa-language"></i>
                    </button>` : ''
                }
                <button onclick="deleteFile(${file.id})" class="btn-action btn-delete" title="Удалить">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        </div>
    `).join('');
}

// Получить текст статуса
function getStatusText(status) {
    const statuses = {
        uploaded: 'Загружен',
        processing: 'В обработке',
        translated: 'Переведен',
        completed: 'Готов',
        error: 'Ошибка'
    };
    return statuses[status] || status || 'Неизвестно';
}

// Перевод файла
async function translateFile(fileId) {
    if (!confirm('Начать перевод файла?')) return;

    try {
        const response = await fetch(`/api/files/${fileId}/translate`, {
            method: 'POST',
            headers: { 'Authorization': 'Bearer ' + token }
        });

        if (response.ok) {
            showNotification('Перевод запущен! Статус обновится автоматически.', 'success');
            // Обновляем список файлов через 3 секунды
            setTimeout(loadFiles, 3000);
        } else {
            showNotification('Ошибка запуска перевода', 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Скачивание файла
async function downloadFile(fileId) {
    try {
        const response = await fetch(`/api/files/${fileId}/download`, {
            headers: { 'Authorization': 'Bearer ' + token }
        });

        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            
            // Получаем имя файла из заголовков
            const contentDisposition = response.headers.get('Content-Disposition');
            let filename = `file_${fileId}`;
            
            if (contentDisposition) {
                const matches = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition);
                if (matches != null && matches[1]) {
                    filename = matches[1].replace(/['"]/g, '');
                }
            }
            
            a.download = filename;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            a.remove();
            showNotification('Файл скачивается', 'success');
        } else {
            showNotification('Ошибка скачивания файла', 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Удаление файла
async function deleteFile(fileId) {
    if (!confirm('Вы уверены, что хотите удалить этот файл?')) return;

    try {
        const response = await fetch(`/api/files/${fileId}`, {
            method: 'DELETE',
            headers: { 'Authorization': 'Bearer ' + token }
        });

        if (response.ok) {
            showNotification('Файл успешно удален', 'success');
            loadFiles();
        } else {
            showNotification('Ошибка удаления файла', 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
    }
}

// Закрытие модального окна
function closeModal() {
    const modal = document.getElementById('modal');
    modal.classList.add('hidden');
}