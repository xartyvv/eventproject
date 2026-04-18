// ===== Аутентификация =====

/**
 * Сохраняет JWT токен в localStorage
 * @param {string} token
 */
function saveToken(token) {
    localStorage.setItem('jwt_token', token);
}

/**
 * Получает JWT токен из localStorage
 * @returns {string|null}
 */
function getToken() {
    return localStorage.getItem('jwt_token');
}

/**
 * Удаляет JWT токен (выход из системы)
 */
function clearToken() {
    localStorage.removeItem('jwt_token');
}

/**
 * Проверяет, авторизован ли пользователь
 * @returns {boolean}
 */
function isLoggedIn() {
    return !!getToken();
}

/**
 * Обновляет навигацию в зависимости от статуса авторизации
 */
function updateNavigation() {
    const loggedIn = isLoggedIn();
    const loginNav = document.getElementById('nav-login');
    const logoutNav = document.getElementById('nav-logout');
    const createNav = document.getElementById('nav-create');
    const rankingNav = document.getElementById('nav-ranking');
    const favoritesNav = document.getElementById('nav-favorites');
    const profileNav = document.getElementById('nav-profile');

    if (loginNav) loginNav.style.display = loggedIn ? 'none' : 'block';
    if (logoutNav) logoutNav.style.display = loggedIn ? 'block' : 'none';
    if (createNav) createNav.style.display = loggedIn ? 'block' : 'none';
    if (rankingNav) rankingNav.style.display = loggedIn ? 'block' : 'none';
    if (favoritesNav) favoritesNav.style.display = loggedIn ? 'block' : 'none';
    if (profileNav) profileNav.style.display = loggedIn ? 'block' : 'none';
}

/**
 * Выход из системы
 */
function logout() {
    clearToken();
    window.location.href = '/';
}

/**
 * Обрабатывает форму входа
 */
document.addEventListener('DOMContentLoaded', () => {
    updateNavigation();

    // Форма входа
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const errorDiv = document.getElementById('error-message');

            try {
                const data = await apiPost('/auth/login', { email, password });
                saveToken(data.token);
                window.location.href = '/';
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
            }
        });
    }

    // Форма регистрации
    const registerForm = document.getElementById('register-form');
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            const errorDiv = document.getElementById('error-message');

            try {
                await apiPost('/auth/register', { email, username, password });
                // Автоматический вход после регистрации
                const data = await apiPost('/auth/login', { email, password });
                saveToken(data.token);
                window.location.href = '/';
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
            }
        });
    }
});
