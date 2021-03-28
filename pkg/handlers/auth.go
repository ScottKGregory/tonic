package handlers

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *AuthHandler) Callback() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *AuthHandler) Logout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (h *AuthHandler) Token() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
