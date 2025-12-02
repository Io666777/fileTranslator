// main.js
import { token, showMainSection } from './auth.js';
import { initDragAndDrop } from './dragdrop.js';
import { resetFileInfo } from './ui.js';

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

// Экспортируем функции для использования в HTML
window.signUp = signUp;
window.signIn = signIn;
window.logout = logout;
window.uploadFile = uploadFile;
window.translateFile = translateFile;
window.downloadFile = downloadFile;
window.deleteFile = deleteFile;
window.removeSelectedFile = removeSelectedFile;
window.closeModal = closeModal;