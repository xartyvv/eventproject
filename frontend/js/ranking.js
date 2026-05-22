// ===== Ранжирование МАИ =====

// Названия 10 критериев
const CRITERIA_NAMES = [
    'Стоимость билета',
    'Длительность',
    'Вместимость',
    'Выходной день',
    'Онлайн формат',
    'Возрастное ограничение',
    'Требуется регистрация',
    'Рейтинг организатора',
    'Время проведения',
    'Интерактивность'
];

const CRITERIA_COUNT = 10;

/**
 * Инициализирует матрицу парных сравнений
 */
function initMatrix() {
    const container = document.getElementById('matrix-container');
    if (!container) return;

    let html = '<table class="matrix-table">';

    // Заголовок
    html += '<thead><tr><th>Критерий</th>';
    for (let j = 0; j < CRITERIA_COUNT; j++) {
        html += `<th title="${CRITERIA_NAMES[j]}">${abbreviate(j + 1)}</th>`;
    }
    html += '</tr></thead>';

    // Тело матрицы
    html += '<tbody>';
    for (let i = 0; i < CRITERIA_COUNT; i++) {
        html += `<tr><td title="${CRITERIA_NAMES[i]}">${abbreviate(i + 1)}</td>`;
        for (let j = 0; j < CRITERIA_COUNT; j++) {
            if (i === j) {
                html += '<td><input type="number" value="1" disabled></td>';
            } else if (i < j) {
                // Верхняя треугольная матрица - редактируемая
                html += `<td><input type="number" id="m_${i}_${j}" min="1" max="9" step="1" value="1" onchange="updateSymmetric(${i}, ${j})"></td>`;
            } else {
                // Нижняя треугольная матрица - автоматически заполняется
                html += `<td><input type="number" id="m_${i}_${j}" disabled></td>`;
            }
        }
        html += '</tr>';
    }
    html += '</tbody></table>';

    container.innerHTML = html;
}

/**
 * Обновляет симметричное значение матрицы (matrix[j][i] = 1/matrix[i][j])
 */
function updateSymmetric(i, j) {
    const input = document.getElementById(`m_${i}_${j}`);
    const value = parseFloat(input.value) || 1;
    const symmetricInput = document.getElementById(`m_${j}_${i}`);

    if (symmetricInput && value > 0) {
        symmetricInput.value = (1 / value).toFixed(2);
    }
}

/**
 * Сокращает название критерия для отображения в матрице
 */
function abbreviate(index) {
    const abbreviations = [
        'Стоим.', 'Длит.', 'Вмест.', 'Выход.', 'Онлайн',
        'Возр.', 'Регист.', 'Орган.', 'Время', 'Интер.'
    ];
    return abbreviations[index - 1] || index;
}

/**
 * Собирает матрицу из полей ввода
 * @returns {Array<Array<number>>}
 */
function collectMatrix() {
    const matrix = [];
    for (let i = 0; i < CRITERIA_COUNT; i++) {
        matrix[i] = [];
        for (let j = 0; j < CRITERIA_COUNT; j++) {
            if (i === j) {
                matrix[i][j] = 1;
            } else {
                const input = document.getElementById(`m_${i}_${j}`);
                matrix[i][j] = parseFloat(input.value) || 1;
            }
        }
    }
    return matrix;
}

/**
 * Сохраняет матрицу и вычисляет веса
 */
async function saveMatrix() {
    const errorDiv = document.getElementById('ranking-error');
    const weightsSection = document.getElementById('weights-section');
    const rankedSection = document.getElementById('ranked-section');

    try {
        errorDiv.textContent = '';
        errorDiv.style.display = 'none';

        const matrix = collectMatrix();

        const data = await apiPost('/ranking/matrix', { matrix });

        // Отображаем веса
        displayWeights(data.weights, data.criteria);
        weightsSection.style.display = 'block';

        // Загружаем ранжированные мероприятия
        await loadRankedEvents();
        rankedSection.style.display = 'block';

    } catch (error) {
        errorDiv.textContent = error.message;
        errorDiv.style.display = 'block';
    }
}

/**
 * Отображает вычисленные веса критериев
 */
function displayWeights(weights, criteria) {
    const container = document.getElementById('weights-list');
    container.innerHTML = '';

    for (let i = 0; i < weights.length; i++) {
        const item = document.createElement('div');
        item.className = 'weight-item';
        const percentage = (weights[i] * 100).toFixed(1);
        item.innerHTML = `
            <span>${criteria[i]}</span>
            <span><strong>${percentage}%</strong></span>
        `;
        container.appendChild(item);
    }
}

function setMatrixValues(matrix) {
    for (let i = 0; i < CRITERIA_COUNT; i++) {
        for (let j = 0; j < CRITERIA_COUNT; j++) {
            const input = document.getElementById(`m_${i}_${j}`);
            if (!input) continue;

            if (i === j) {
                input.value = 1;
            } else {
                input.value = matrix[i][j] || 1;
            }
        }
    }
}

async function loadSavedMatrix() {
    const errorDiv = document.getElementById('ranking-error');
    const weightsSection = document.getElementById('weights-section');
    const rankedSection = document.getElementById('ranked-section');

    try {
        const data = await apiGet('/ranking/weights');

        if (data.matrix) {
            setMatrixValues(data.matrix);
        }

        if (data.weights) {
            displayWeights(data.weights, data.criteria);
            weightsSection.style.display = 'block';
            await loadRankedEvents();
            rankedSection.style.display = 'block';
        }
    } catch (error) {
        if (!error.message.includes('ranking profile not found')) {
            errorDiv.textContent = error.message;
            errorDiv.style.display = 'block';
        }
    }
}

/**
 * Загружает и отображает ранжированные мероприятия
 */
async function loadRankedEvents() {
    const container = document.getElementById('ranked-events');

    try {
        const data = await apiGet('/ranking/events');

        if (!data.ranked_events || data.ranked_events.length === 0) {
            container.innerHTML = '<p class="empty-state">Нет мероприятий для ранжирования</p>';
            return;
        }

        container.innerHTML = '';
        data.ranked_events.slice(0, 20).forEach(event => {
            const item = document.createElement('div');
            item.className = 'ranked-item';
            item.innerHTML = `
                <div class="rank-badge">${event.rank}</div>
                <div class="rank-content">
                    <h3>${escapeHtml(event.title)}</h3>
                    <p class="event-meta">${event.date} | ${escapeHtml(event.location)}</p>
                    <p class="event-description">${escapeHtml(event.description)}</p>
                    <span class="rank-score">Балл: ${event.score.toFixed(4)}</span>
                </div>
            `;
            container.appendChild(item);
        });
    } catch (error) {
        container.innerHTML = `<p class="error-message">${escapeHtml(error.message)}</p>`;
    }
}

/**
 * Авто-заполнение матрицы для тестирования
 */
function autoFillMatrix() {
    // Пример: Стоимость самая важная (9), потом остальные критерии по важности
    const importance = [9, 4, 3, 2, 2, 3, 2, 6, 3, 5];

    for (let i = 0; i < CRITERIA_COUNT; i++) {
        for (let j = i + 1; j < CRITERIA_COUNT; j++) {
            const input = document.getElementById(`m_${i}_${j}`);
            if (input) {
                const ratio = importance[i] / importance[j];
                const value = Math.max(1, Math.min(9, Math.round(ratio * 2) / 2)); // Округляем до 0.5
                input.value = value;
                updateSymmetric(i, j);
            }
        }
    }
}

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', async () => {
    if (document.getElementById('matrix-container')) {
        initMatrix();
        await loadSavedMatrix();
    }
});
