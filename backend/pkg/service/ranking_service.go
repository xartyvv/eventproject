package service

import (
	"encoding/json"
	"math"

	"github.com/xartyvv/eventproject/backend/pkg/domain"
	"github.com/xartyvv/eventproject/backend/pkg/repository"
)

// RankingService определяет интерфейс для ранжирования мероприятий по МАИ
type RankingService interface {
	ComputeWeights(matrix [10][10]float64) ([]float64, error)
	RankEvents(events []domain.Event, weights []float64) ([]domain.RankedEvent, error)
	SaveRankingProfile(userID uint, matrix [10][10]float64) (*domain.RankingProfile, error)
	GetRankingProfile(userID uint) (*domain.RankingProfile, error)
	GetRankedEvents(userID uint, category string, start, end string) ([]domain.RankedEvent, error)
}

// rankingService — реализация RankingService
type rankingService struct {
	rankingRepo  repository.RankingRepository
	eventRepo    repository.EventRepository
	favoriteRepo repository.FavoriteRepository
}

// NewRankingService создает новый экземпляр rankingService
func NewRankingService(
	rankingRepo repository.RankingRepository,
	eventRepo repository.EventRepository,
	favoriteRepo repository.FavoriteRepository,
) RankingService {
	return &rankingService{
		rankingRepo:  rankingRepo,
		eventRepo:    eventRepo,
		favoriteRepo: favoriteRepo,
	}
}

// ComputeWeights вычисляет веса критериев по методу анализа иерархий (МАИ)
// Использует метод собственных векторов (приближённый через геометрическое среднее)
func (s *rankingService) ComputeWeights(matrix [10][10]float64) ([]float64, error) {
	n := 10 // Количество критериев

	// Шаг 1: Вычисляем геометрическое среднее каждой строки
	geoMeans := make([]float64, n)
	for i := 0; i < n; i++ {
		product := 1.0
		for j := 0; j < n; j++ {
			product *= matrix[i][j]
		}
		geoMeans[i] = math.Pow(product, 1.0/float64(n))
	}

	// Шаг 2: Нормализуем (сумма весов = 1)
	sum := 0.0
	for _, gm := range geoMeans {
		sum += gm
	}

	if sum == 0 {
		return nil, &RankingError{Message: "сумма геометрических средних равна нулю"}
	}

	weights := make([]float64, n)
	for i := 0; i < n; i++ {
		weights[i] = geoMeans[i] / sum
	}

	// Шаг 3: Проверка согласованности (опционально, для отладки)
	// Consistency Ratio должен быть < 0.1
	// _ = s.checkConsistency(matrix, weights)

	return weights, nil
}

// RankEvents ранжирует мероприятия по МАИ с учётом весов критериев
func (s *rankingService) RankEvents(events []domain.Event, weights []float64) ([]domain.RankedEvent, error) {
	if len(events) == 0 {
		return []domain.RankedEvent{}, nil
	}

	// Шаг 1: Извлекаем матрицу значений критериев для всех мероприятий
	n := len(events)
	m := 10 // количество критериев

	rawMatrix := make([][]float64, n)
	for i := 0; i < n; i++ {
		rawMatrix[i] = s.extractCriteria(events[i])
	}

	// Шаг 2: Нормализуем критерии (приводим к диапазону 0-1)
	normalizedMatrix := s.normalizeCriteria(rawMatrix, m, n)

	// Шаг 3: Вычисляем итоговый скор для каждого мероприятия
	ranked := make([]domain.RankedEvent, n)
	for i := 0; i < n; i++ {
		score := 0.0
		for j := 0; j < m; j++ {
			// Для "плохих" критериев (cost, ageRestriction) инвертируем
			inverted := s.isInvertedCriterion(j)
			value := normalizedMatrix[i][j]
			if inverted {
				value = 1.0 - value
			}
			score += weights[j] * value
		}

		ranked[i] = domain.RankedEvent{
			EventID:     events[i].ID,
			Title:       events[i].Title,
			Description: events[i].Description,
			Date:        events[i].Date.Format("2006-01-02 15:04"),
			Location:    events[i].Location,
			Category:    events[i].Category,
			Score:       math.Round(score*10000) / 10000, // 4 знака после запятой
			Criteria:    rawMatrix[i],
		}
	}

	// Шаг 4: Сортируем по убыванию score
	for i := 0; i < len(ranked); i++ {
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].Score > ranked[i].Score {
				ranked[i], ranked[j] = ranked[j], ranked[i]
			}
		}
	}

	// Шаг 5: Присваиваем ранги
	for i := range ranked {
		ranked[i].Rank = i + 1
	}

	return ranked, nil
}

// extractCriteria извлекает 10 критериев из мероприятия
func (s *rankingService) extractCriteria(event domain.Event) []float64 {
	return []float64{
		float64(event.Cost),                       // 0: Стоимость (меньше = лучше, инвертируется)
		float64(event.Duration),                   // 1: Длительность
		float64(event.Capacity),                   // 2: Вместимость
		boolToFloat64(event.IsWeekend),            // 3: Выходной день
		boolToFloat64(event.IsOnline),             // 4: Онлайн формат
		float64(event.AgeRestriction),             // 5: Возрастное ограничение (меньше = лучше, инвертируется)
		boolToFloat64(event.RequiresRegistration), // 6: Требуется регистрация
		float64(event.OrganizerRating),            // 7: Рейтинг организатора
		float64(event.TimeOfDay),                  // 8: Время проведения
		float64(event.Interactivity),              // 9: Интерактивность
	}
}

// normalizeCriteria нормализует критерии к диапазону 0-1 (min-max нормализация)
func (s *rankingService) normalizeCriteria(matrix [][]float64, m, n int) [][]float64 {
	normalized := make([][]float64, n)

	for j := 0; j < m; j++ {
		// Находим min и max для критерия j
		minVal := matrix[0][j]
		maxVal := matrix[0][j]
		for i := 1; i < n; i++ {
			if matrix[i][j] < minVal {
				minVal = matrix[i][j]
			}
			if matrix[i][j] > maxVal {
				maxVal = matrix[i][j]
			}
		}

		// Нормализуем
		rangeVal := maxVal - minVal
		for i := 0; i < n; i++ {
			if normalized[i] == nil {
				normalized[i] = make([]float64, m)
			}
			if rangeVal == 0 {
				normalized[i][j] = 0.5 // Все значения одинаковые
			} else {
				normalized[i][j] = (matrix[i][j] - minVal) / rangeVal
			}
		}
	}

	return normalized
}

// isInvertedCriterion возвращает true для критериев, где "меньше = лучше"
func (s *rankingService) isInvertedCriterion(index int) bool {
	switch index {
	case 0, 5: // Cost, AgeRestriction
		return true
	default:
		return false
	}
}

// SaveRankingProfile сохраняет матрицу парных сравнений и вычисленные веса
func (s *rankingService) SaveRankingProfile(userID uint, matrix [10][10]float64) (*domain.RankingProfile, error) {
	// Вычисляем веса
	weights, err := s.ComputeWeights(matrix)
	if err != nil {
		return nil, err
	}

	// Сериализуем в JSON
	matrixJSON, _ := json.Marshal(matrix)
	weightsJSON, _ := json.Marshal(weights)

	profile := &domain.RankingProfile{
		UserID:         userID,
		PairwiseMatrix: string(matrixJSON),
		Weights:        string(weightsJSON),
	}

	if err := s.rankingRepo.CreateOrUpdate(profile); err != nil {
		return nil, err
	}

	return profile, nil
}

// GetRankingProfile получает профиль ранжирования пользователя
func (s *rankingService) GetRankingProfile(userID uint) (*domain.RankingProfile, error) {
	return s.rankingRepo.GetByUserID(userID)
}

// GetRankedEvents получает ранжированный список мероприятий для пользователя
func (s *rankingService) GetRankedEvents(userID uint, category string, startDate, endDate string) ([]domain.RankedEvent, error) {
	// Получаем профиль ранжирования пользователя
	profile, err := s.rankingRepo.GetByUserID(userID)
	if err != nil {
		return nil, &RankingError{Message: "у вас нет сохранённого профиля ранжирования. Создайте матрицу сравнений."}
	}

	// Десериализуем веса
	var weights []float64
	if err := json.Unmarshal([]byte(profile.Weights), &weights); err != nil {
		return nil, err
	}

	// Получаем мероприятия (с фильтрами)
	var events []domain.Event
	events, err = s.eventRepo.GetAllForRanking()
	if err != nil {
		return nil, err
	}

	// Применяем фильтры через eventRepo (упрощённо — все мероприятия)
	// В реальном варианте можно использовать s.eventRepo.GetByFilters

	// Ранжируем
	ranked, err := s.RankEvents(events, weights)
	if err != nil {
		return nil, err
	}

	return ranked, nil
}

// boolToFloat64 конвертирует bool в float64 (true=1, false=0)
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// RankingError — ошибка ранжирования
type RankingError struct {
	Message string
}

func (e *RankingError) Error() string {
	return e.Message
}
