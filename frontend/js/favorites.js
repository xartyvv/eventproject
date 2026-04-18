// ===== Избранное =====

/**
 * Загружает и отображает избранные мероприятия
 */
async function loadFavorites() {
    const listDiv = document.getElementById('favorites-list');
    const loadingDiv = document.getElementById('favorites-loading');
    const emptyDiv = document.getElementById('favorites-empty');

    if (!listDiv) return;

    try {
        loadingDiv.style.display = 'block';
        listDiv.innerHTML = '';

        const data = await apiGet('/favorites', { page: 1, page_size: 50 });
        loadingDiv.style.display = 'none';

        if (!data.favorites || data.favorites.length === 0) {
            emptyDiv.style.display = 'block';
            return;
        }

        emptyDiv.style.display = 'none';
        data.favorites.forEach(fav => {
            if (fav.Event) {
                listDiv.appendChild(createFavoriteCard(fav.Event, fav.id));
            }
        });
    } catch (error) {
        loadingDiv.textContent = 'Ошибка загрузки избранного: ' + error.message;
    }
}

/**
 * Создаёт карточку избранного мероприятия
 */
function createFavoriteCard(event, favId) {
    const card = document.createElement('div');
    card.className = 'event-card';

    const date = new Date(event.date).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric'
    });

    card.innerHTML = `
        <h3>${escapeHtml(event.title)}</h3>
        <div class="event-meta">${date} | ${escapeHtml(event.location)}</div>
        <div class="event-description">${escapeHtml(event.description)}</div>
        <div class="event-actions">
            <button class="btn btn-secondary btn-small" onclick="removeFromFavorites(${favId}, ${event.id})">
                Удалить
    `;

    return card;
}

/**
 * Удаляет мероприятие из избранного
 */
async function removeFromFavorites(favId, eventId) {
    try {
        await apiDelete(`/favorites/${eventId}`);
        // Перезагружаем список
        loadFavorites();
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (document.getElementById('favorites-list')) {
        loadFavorites();
    }
});
