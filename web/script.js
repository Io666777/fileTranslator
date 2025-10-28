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