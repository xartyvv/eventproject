// ===== API клиент =====
const API_BASE = '/api';

/**
 * Выполняет GET запрос к API
 * @param {string} endpoint - endpoint API
 * @param {Object} params - query параметры
 * @returns {Promise<Object>}
 */
async function apiGet(endpoint, params = {}) {
    const url = new URL(`${API_BASE}${endpoint}`, window.location.origin);
    Object.keys(params).forEach(key => {
        if (params[key]) {
            url.searchParams.append(key, params[key]);
        }
    });

    const response = await fetch(url.toString(), {
        method: 'GET',
        headers: getHeaders()
    });

    return handleResponse(response);
}

/**
 * Выполняет POST запрос к API
 * @param {string} endpoint - endpoint API
 * @param {Object} data - тело запроса
 * @returns {Promise<Object>}
 */
async function apiPost(endpoint, data) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        method: 'POST',
        headers: getHeaders(),
        body: JSON.stringify(data)
    });

    return handleResponse(response);
}

/**
 * Выполняет PUT запрос к API
 * @param {string} endpoint - endpoint API
 * @param {Object} data - тело запроса
 * @returns {Promise<Object>}
 */
async function apiPut(endpoint, data) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify(data)
    });

    return handleResponse(response);
}

/**
 * Выполняет DELETE запрос к API
 * @param {string} endpoint - endpoint API
 * @returns {Promise<Object>}
 */
async function apiDelete(endpoint) {
    const response = await fetch(`${API_BASE}${endpoint}`, {
        method: 'DELETE',
        headers: getHeaders()
    });

    return handleResponse(response);
}

/**
 * Возвращает заголовки с JWT токеном
 * @returns {Object}
 */
function getHeaders() {
    const headers = {
        'Content-Type': 'application/json'
    };
    const token = getToken();
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
}

/**
 * Обрабатывает ответ от API
 * @param {Response} response
 * @returns {Promise<Object>}
 */
async function handleResponse(response) {
    const data = await response.json();
    if (!response.ok) {
        throw new Error(data.error || 'Произошла ошибка');
    }
    return data;
}

/**
 * Экранирует HTML для предотвращения XSS
 * @param {string} text
 * @returns {string}
 */
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}
