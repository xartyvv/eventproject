// ===== Профиль пользователя =====

/**
 * Загружает и отображает информацию о профиле
 */
async function loadProfile() {
    const usernameEl = document.getElementById('profile-username');
    const emailEl = document.getElementById('profile-email');

    if (!usernameEl) return;

    let user = null;

    try {
        const data = await apiGet('/auth/me');
        user = data.user;

        usernameEl.textContent = user.username;
        emailEl.textContent = user.email;

        // Загружаем статистику
        await loadProfileStats(user.id);
    } catch (error) {
        usernameEl.textContent = 'Ошибка загрузки';
    }

    if (user) {
        await loadMyEvents(user.id);
    } else {
        const listDiv = document.getElementById('my-events');
        if (listDiv) {
            listDiv.innerHTML = '<div class="empty-state">Не удалось загрузить профиль. Пожалуйста, войдите снова.</div>';
        }
    }
}

/**
 * Загружает список мероприятий пользователя
 */
async function loadMyEvents(userId) {
    const listDiv = document.getElementById('my-events');
    if (!listDiv) return;

    try {
        listDiv.innerHTML = '';
        const params = {
            page: 1,
            page_size: 100
        };

        const data = await apiGet('/events/mine', params);
        const userEvents = data.events || [];

        if (userEvents.length === 0) {
            listDiv.innerHTML = '<div class="empty-state">Вы ещё не создали ни одного мероприятия.</div>';
            return;
        }

        userEvents.forEach(event => {
            listDiv.appendChild(createUserEventCard(event));
        });
    } catch (error) {
        listDiv.innerHTML = '<div class="empty-state">Ошибка загрузки мероприятий.</div>';
        console.error('Ошибка загрузки моих мероприятий:', error);
    }
}

function createUserEventCard(event) {
    const card = createEventCard(event);
    const actionsDiv = card.querySelector('.event-actions');
    if (actionsDiv) {
        actionsDiv.innerHTML += `
            <button class="btn btn-secondary btn-small" onclick="goToEditEvent(${event.id})">Редактировать</button>
            <button class="btn btn-danger btn-small" onclick="deleteEvent(${event.id})">Удалить</button>
        `;
    }
    return card;
}

function goToEditEvent(eventId) {
    window.location.href = `/create-event.html?event_id=${eventId}`;
}

async function deleteEvent(eventId) {
    if (!confirm('Вы уверены, что хотите удалить это мероприятие?')) {
        return;
    }

    try {
        await apiDelete(`/events/${eventId}`);
        await loadProfile();
    } catch (error) {
        console.error('Ошибка удаления мероприятия:', error);
        alert('Не удалось удалить мероприятие: ' + error.message);
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
