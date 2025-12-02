// ui.js
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

// Закрытие модального окна
function closeModal() {
    const modal = document.getElementById('modal');
    modal.classList.add('hidden');
}

export { showNotification, updateFileInfo, removeSelectedFile, resetFileInfo, closeModal };