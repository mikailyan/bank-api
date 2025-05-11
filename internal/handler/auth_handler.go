package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/therealadik/bank-api/internal/dto"
	"github.com/therealadik/bank-api/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService service.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Register обрабатывает запрос на регистрацию пользователя
// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя с уникальными email и username
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Данные для регистрации"
// @Success 201 {string} string "Пользователь успешно зарегистрирован"
// @Failure 400 {string} string "Ошибка валидации данных"
// @Failure 409 {string} string "Пользователь с таким email или username уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Ошибка декодирования запроса регистрации")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID, err := h.authService.Register(r.Context(), req)
	if err != nil {
		h.logger.WithError(err).Warn("Ошибка при регистрации пользователя")

		if errors.Is(err, service.ErrUserExists) {
			http.Error(w, "Пользователь с таким email или username уже существует", http.StatusConflict)
			return
		}

		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message": "Пользователь успешно зарегистрирован",
		"user_id": userID,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("Ошибка при формировании ответа")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}

// Login обрабатывает запрос на вход в систему
// @Summary Вход в систему
// @Description Аутентифицирует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.AuthResponse "JWT токен"
// @Failure 400 {string} string "Ошибка валидации данных"
// @Failure 401 {string} string "Неверные учетные данные"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Warn("Ошибка декодирования запроса авторизации")
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email и пароль обязательны", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		h.logger.WithError(err).Warn("Ошибка при авторизации пользователя")

		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Ошибка авторизации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := dto.AuthResponse{Token: token}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.WithError(err).Error("Ошибка при формировании ответа авторизации")
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
}
