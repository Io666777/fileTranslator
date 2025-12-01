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
});

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
            const error = await response.text();
            showNotification(`Ошибка регистрации: ${error}`, 'error');
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
            fileInput.value = '';
            const fileInfo = document.getElementById('file-info');
            if (fileInfo) fileInfo.innerHTML = '';
            loadFiles();
        } else {
            const error = await response.text();
            showNotification(`Ошибка загрузки: ${error}`, 'error');
        }
    } catch (error) {
        showNotification('Ошибка сети', 'error');
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
            a.download = `file_${fileId}`;
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