// ===== Профиль пользователя =====

/**
 * Загружает и отображает информацию о профиле
 */
async function loadProfile() {
    const usernameEl = document.getElementById('profile-username');
    const emailEl = document.getElementById('profile-email');

    if (!usernameEl) return;

    try {
        const data = await apiGet('/auth/me');
        const user = data.user;

        usernameEl.textContent = user.username;
        emailEl.textContent = user.email;

        // Загружаем статистику
        await loadProfileStats(user.id);
    } catch (error) {
        usernameEl.textContent = 'Ошибка загрузки';
    }
}

/**
 * Загружает статистику профиля
 */
async function loadProfileStats(userId) {
    try {
        // Загружаем мероприятия пользователя
        const eventsData = await apiGet('/events', { page: 1, page_size: 100 });
        const userEvents = eventsData.events.filter(e => e.creator_id === userId);
        document.getElementById('stat-events').textContent = userEvents.length;

        // Загружаем избранное
        const favData = await apiGet('/favorites', { page: 1, page_size: 100 });
        document.getElementById('stat-favorites').textContent = favData.total || 0;

        // Проверяем профиль МАИ
        try {
            await apiGet('/ranking/weights');
            document.getElementById('stat-ranking').textContent = 'Создан';
        } catch {
            document.getElementById('stat-ranking').textContent = 'Нет';
        }
    } catch (error) {
        console.error('Ошибка загрузки статистики:', error);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('profile-info')) {
        loadProfile();
    }
});
