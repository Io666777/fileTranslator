// utils.js
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

export { formatFileSize, getFileType, getFileIcon, formatDate };