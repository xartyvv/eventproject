const API_BASE = 'http://localhost:8081/api';
let favoriteIds = new Set();

async function fetchJSON(url, options = {}) {
    const response = await fetch(url, options);
    return response.json();
}

async function loadEvents() {
    const eventsElement = document.getElementById('events');
    if (!eventsElement) return;

    const events = await fetchJSON(`${API_BASE}/events`);
    renderEvents(events);
}

async function loadFavorites() {
    const favoritesElement = document.getElementById('favorites');
    if (!favoritesElement) return;

    const data = await fetchJSON(`${API_BASE}/favorites`);
    favoriteIds = new Set(data.favorites.map(event => event.id));
    renderFavorites(data.favorites);
}

function renderEvents(events) {
    const container = document.getElementById('events');
    if (!container) return;
    container.innerHTML = '';

    events.forEach(event => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
            <h3>${escapeHtml(event.title)}</h3>
            <div class="meta">${escapeHtml(event.date)} · ${escapeHtml(event.location)}</div>
            <p>${escapeHtml(event.description)}</p>
            <button class="button" data-id="${event.id}">${favoriteIds.has(event.id) ? 'Убрано в избранное' : 'Добавить в избранное'}</button>
        `;
        const button = card.querySelector('button');
        button.addEventListener('click', () => toggleFavorite(event.id));
        container.appendChild(card);
    });
}

function renderFavorites(events) {
    const container = document.getElementById('favorites');
    const empty = document.getElementById('favorites-empty');
    if (!container || !empty) return;
    container.innerHTML = '';

    if (events.length === 0) {
        empty.style.display = 'block';
        return;
    }

    empty.style.display = 'none';
    events.forEach(event => {
        const card = document.createElement('div');
        card.className = 'card';
        card.innerHTML = `
            <h3>${escapeHtml(event.title)}</h3>
            <div class="meta">${escapeHtml(event.date)} · ${escapeHtml(event.location)}</div>
            <p>${escapeHtml(event.description)}</p>
            <button class="button secondary" data-id="${event.id}">Убрать из избранного</button>
        `;
        const button = card.querySelector('button');
        button.addEventListener('click', () => removeFavorite(event.id));
        container.appendChild(card);
    });
}

async function handleAuth(action) {
    const emailInput = document.getElementById('auth-email');
    const passwordInput = document.getElementById('auth-password');
    if (!emailInput || !passwordInput) return;

    const email = emailInput.value.trim();
    const password = passwordInput.value.trim();
    if (!email || !password) {
        alert('Введите email и пароль.');
        return;
    }

    const endpoint = action === 'login' ? 'login' : 'register';
    const data = await fetchJSON(`${API_BASE}/${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
    });

    alert(data.message || 'Функция пока не реализована.');
}

async function toggleFavorite(eventId) {
    if (favoriteIds.has(eventId)) {
        await removeFavorite(eventId);
        return;
    }
    await addFavorite(eventId);
}

async function addFavorite(eventId) {
    await fetchJSON(`${API_BASE}/favorites`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ event_id: eventId }),
    });
    favoriteIds.add(eventId);
    await refreshViews();
}

async function removeFavorite(eventId) {
    await fetch(`${API_BASE}/favorites/${eventId}`, { method: 'DELETE' });
    favoriteIds.delete(eventId);
    await refreshViews();
}

async function refreshViews() {
    if (document.getElementById('favorites')) {
        await loadFavorites();
    }
    if (document.getElementById('events')) {
        await loadEvents();
    }
}

function escapeHtml(text) {
    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}

window.addEventListener('DOMContentLoaded', async () => {
    const loginButton = document.getElementById('auth-login');
    const registerButton = document.getElementById('auth-register');

    if (loginButton) {
        loginButton.addEventListener('click', () => handleAuth('login'));
    }
    if (registerButton) {
        registerButton.addEventListener('click', () => handleAuth('register'));
    }

    await refreshViews();
});