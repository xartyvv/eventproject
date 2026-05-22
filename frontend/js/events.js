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

function formatDateForInput(dateString) {
    const date = new Date(dateString);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
}

async function loadEventForEdit(eventId) {
    const data = await apiGet(`/events/${eventId}`);
    const event = data.event;
    if (!event) return;

    document.getElementById('event_id').value = event.id;
    document.getElementById('title').value = event.title || '';
    document.getElementById('description').value = event.description || '';
    document.getElementById('date').value = formatDateForInput(event.date);
    document.getElementById('location').value = event.location || '';
    document.getElementById('category').value = event.category || 'concert';
    document.getElementById('cost').value = event.cost || 0;
    document.getElementById('duration').value = event.duration || 0;
    document.getElementById('capacity').value = event.capacity || 0;
    document.getElementById('is_weekend').checked = !!event.is_weekend;
    document.getElementById('is_online').checked = !!event.is_online;
    document.getElementById('age_restriction').value = event.age_restriction || 0;
    document.getElementById('requires_registration').checked = !!event.requires_registration;
    document.getElementById('organizer_rating').value = event.organizer_rating || 0;
    document.getElementById('time_of_day').value = event.time_of_day || 3;
    document.getElementById('interactivity').value = event.interactivity || 0;
    document.querySelector('h1').textContent = 'Редактировать мероприятие';
    document.querySelector('button[type="submit"]').textContent = 'Сохранить изменения';
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
document.addEventListener('DOMContentLoaded', async () => {
    // Загрузка мероприятий на странице events.html
    if (document.getElementById('events-list')) {
        loadEvents();
    }

    // Форма создания/редактирования мероприятия
    const createForm = document.getElementById('create-event-form');
    if (createForm) {
        const urlParams = new URLSearchParams(window.location.search);
        const eventId = urlParams.get('event_id');

        if (eventId) {
            try {
                await loadEventForEdit(eventId);
            } catch (error) {
                console.error('Ошибка загрузки мероприятия для редактирования:', error);
            }
        }

        createForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const errorDiv = document.getElementById('form-error');
            const successDiv = document.getElementById('form-success');
            const currentEventId = document.getElementById('event_id').value;

            try {
                const formData = {
                    title: document.getElementById('title').value,
                    description: document.getElementById('description').value,
                    date: document.getElementById('date').value,
                    location: document.getElementById('location').value,
                    category: document.getElementById('category').value,
                    cost: parseFloat(document.getElementById('cost').value) || 0,
                    duration: parseFloat(document.getElementById('duration').value) || 2,
                    capacity: parseFloat(document.getElementById('capacity').value) || 100,
                    is_weekend: document.getElementById('is_weekend').checked,
                    is_online: document.getElementById('is_online').checked,
                    age_restriction: parseInt(document.getElementById('age_restriction').value) || 0,
                    requires_registration: document.getElementById('requires_registration').checked,
                    organizer_rating: parseFloat(document.getElementById('organizer_rating').value) || 5,
                    time_of_day: parseInt(document.getElementById('time_of_day').value) || 3,
                    interactivity: parseFloat(document.getElementById('interactivity').value) || 5
                };

                if (currentEventId) {
                    await apiPut(`/events/${currentEventId}`, formData);
                    successDiv.textContent = 'Мероприятие успешно обновлено!';
                } else {
                    await apiPost('/events/create', formData);
                    successDiv.textContent = 'Мероприятие успешно создано!';
                    createForm.reset();
                }

                successDiv.style.display = 'block';
                errorDiv.style.display = 'none';
            } catch (error) {
                errorDiv.textContent = error.message;
                errorDiv.style.display = 'block';
                successDiv.style.display = 'none';
            }
        });
    }
});
