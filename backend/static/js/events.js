// ===== Мероприятия =====

/**
 * Загружает и отображает список мероприятий
 */
async function loadEvents() {
    const listDiv = document.getElementById('events-list');
    const loadingDiv = document.getElementById('events-loading');
    const emptyDiv = document.getElementById('events-empty');

    if (!listDiv) return;

    try {
        loadingDiv.style.display = 'block';
        listDiv.innerHTML = '';

        const params = {
            page: 1,
            page_size: 50
        };

        const data = await apiGet('/events', params);
        loadingDiv.style.display = 'none';

        if (data.events.length === 0) {
            emptyDiv.style.display = 'block';
            return;
        }

        emptyDiv.style.display = 'none';
        data.events.forEach(event => {
            listDiv.appendChild(createEventCard(event));
        });
    } catch (error) {
        loadingDiv.textContent = 'Ошибка загрузки мероприятий: ' + error.message;
    }
}

/**
 * Создаёт HTML-карточку мероприятия
 * @param {Object} event
 * @returns {HTMLElement}
 */
function createEventCard(event) {
    const card = document.createElement('div');
    card.className = 'event-card';

    const date = new Date(event.date).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    const criteria = [];
    if (event.cost > 0) criteria.push(`${event.cost}₽`);
    if (event.is_free) criteria.push('Бесплатно');
    if (event.distance > 0) criteria.push(`${event.distance} км`);
    if (event.duration > 0) criteria.push(`${event.duration}ч`);
    if (event.is_online) criteria.push('Онлайн');
    if (event.is_weekend) criteria.push('Выходной');

    const favoriteBtn = isLoggedIn() ? `<button class="btn btn-secondary btn-small" onclick="addToFavorite(${event.id})">Добавить в избранное</button>` : '';

    card.innerHTML = `
        <h3>${escapeHtml(event.title)}</h3>
        <div class="event-meta">${date} | ${escapeHtml(event.location)}</div>
        <div class="event-description">${escapeHtml(event.description)}</div>
        <div class="event-criteria">
            ${criteria.map(c => `<span class="criteria-tag">${c}</span>`).join('')}
        </div>
        <div class="event-actions">
            ${favoriteBtn}
        </div>
    `;

    return card;
}

/**
 * Применяет фильтры
 */
async function applyFilters() {
    const category = document.getElementById('filter-category')?.value || '';
    const searchQuery = document.getElementById('filter-search')?.value.trim().toLowerCase() || '';
    const startDate = document.getElementById('filter-start-date')?.value || '';
    const endDate = document.getElementById('filter-end-date')?.value || '';

    const listDiv = document.getElementById('events-list');
    const loadingDiv = document.getElementById('events-loading');

    try {
        loadingDiv.style.display = 'block';
        listDiv.innerHTML = '';

        const params = {
            page: 1,
            page_size: 50
        };
        if (category) params.category = category;
        if (startDate) params.start_date = startDate;
        if (endDate) params.end_date = endDate;

        const data = await apiGet('/events/filter', params);
        loadingDiv.style.display = 'none';

        let events = data.events || [];
        if (searchQuery) {
            events = events.filter(event => {
                const text = [event.title, event.description, event.location, event.category].filter(Boolean).join(' ').toLowerCase();
                return text.includes(searchQuery);
            });
        }

        if (events.length === 0) {
            document.getElementById('events-empty').style.display = 'block';
            return;
        }

        document.getElementById('events-empty').style.display = 'none';
        events.forEach(event => {
            listDiv.appendChild(createEventCard(event));
        });
    } catch (error) {
        loadingDiv.textContent = 'Ошибка: ' + error.message;
    }
}

/**
 * Сбрасывает фильтры
 */
function resetFilters() {
    const category = document.getElementById('filter-category');
    const startDate = document.getElementById('filter-start-date');
    const endDate = document.getElementById('filter-end-date');

    if (category) category.value = '';
    if (startDate) startDate.value = '';
    if (endDate) endDate.value = '';

    loadEvents();
}

/**
 * Добавляет мероприятие в избранное
 * @param {number} eventId
 */
async function addToFavorite(eventId) {
    try {
        await apiPost('/favorites/add', { event_id: eventId });
        alert('Добавлено в избранное!');
    } catch (error) {
        alert('Ошибка: ' + error.message);
    }
}

/**
 * Обрабатывает форму создания мероприятия
 */
document.addEventListener('DOMContentLoaded', () => {
    // Загрузка мероприятий на странице events.html
    if (document.getElementById('events-list')) {
        loadEvents();
    }

    // Форма создания мероприятия
    const createForm = document.getElementById('create-event-form');
    if (createForm) {
        createForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const errorDiv = document.getElementById('form-error');
            const successDiv = document.getElementById('form-success');

            try {
                const formData = {
                    title: document.getElementById('title').value,
                    description: document.getElementById('description').value,
                    date: document.getElementById('date').value,
                    location: document.getElementById('location').value,
                    category: document.getElementById('category').value,
                    cost: parseFloat(document.getElementById('cost').value) || 0,
                    distance: parseFloat(document.getElementById('distance').value) || 0,
                    duration: parseFloat(document.getElementById('duration').value) || 2,
                    rating: parseFloat(document.getElementById('rating').value) || 5,
                    capacity: parseFloat(document.getElementById('capacity').value) || 100,
                    is_weekend: document.getElementById('is_weekend').checked,
                    is_online: document.getElementById('is_online').checked,
                    age_restriction: parseInt(document.getElementById('age_restriction').value) || 0,
                    requires_registration: document.getElementById('requires_registration').checked,
                    organizer_rating: parseFloat(document.getElementById('organizer_rating').value) || 5,
                    is_free: document.getElementById('is_free').checked,
                    time_of_day: parseInt(document.getElementById('time_of_day').value) || 3,
                    accessibility: parseFloat(document.getElementById('accessibility').value) || 5,
                    popularity: parseFloat(document.getElementById('popularity').value) || 0,
                    interactivity: parseFloat(document.getElementById('interactivity').value) || 5
                };

                await apiPost('/events/create', formData);
                successDiv.textContent = 'Мероприятие успешно создано!';
                successDiv.style.display = 'block';
                errorDiv.style.display = 'none';
                createForm.reset();
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
                successDiv.style.display = 'none';
            }
        });
    }
});
