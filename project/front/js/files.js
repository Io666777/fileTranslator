// files.js
import { token } from './auth.js';
import { showNotification } from './ui.js';
import { formatFileSize, getFileType, getFileIcon, formatDate } from './utils.js';

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

export { uploadFile, loadFiles, translateFile, downloadFile, deleteFile };