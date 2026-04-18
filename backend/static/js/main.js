// ===== Главная страница =====
document.addEventListener('DOMContentLoaded', () => {
    // Если пользователь авторизован, можно показать дополнительную информацию
    if (isLoggedIn()) {
        console.log('Пользователь авторизован');
    }
});
