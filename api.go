package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var startTime time.Time

func StartAPIServer(cfg ServerConfig) {
	startTime = time.Now()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Status
	e.GET("/status", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "running",
			"uptime":    time.Since(startTime).Seconds(),
			"lastCheck": time.Now().Format(time.RFC3339),
		})
	})

	// Logs
	e.GET("/logs", func(c echo.Context) error {
		rows, err := db.Query(context.Background(), "SELECT id, level, message, details, created FROM snt_logs ORDER BY created DESC LIMIT 100")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		defer rows.Close()

		logs := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var level, message string
			var details *string
			var created time.Time
			if err := rows.Scan(&id, &level, &message, &details, &created); err != nil {
				continue
			}
			logs = append(logs, map[string]interface{}{
				"id":      id,
				"level":   level,
				"message": message,
				"details": details,
				"created": created,
			})
		}
		return c.JSON(http.StatusOK, logs)
	})

	// Contacts
	e.GET("/contacts", func(c echo.Context) error {
		rows, err := db.Query(context.Background(), "SELECT id, prior, type, value, adds, comment FROM snt_contacts ORDER BY prior")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		defer rows.Close()

		contacts := []map[string]interface{}{}
		for rows.Next() {
			var id, prior int
			var cType, value string
			var adds, comment *string
			if err := rows.Scan(&id, &prior, &cType, &value, &adds, &comment); err != nil {
				continue
			}
			contacts = append(contacts, map[string]interface{}{
				"id":      id,
				"prior":   prior,
				"type":    cType,
				"value":   value,
				"adds":    adds,
				"comment": comment,
			})
		}
		return c.JSON(http.StatusOK, contacts)
	})

	e.POST("/contacts", func(c echo.Context) error {
		var input struct {
			Prior   int     `json:"prior"`
			Type    string  `json:"type"`
			Value   string  `json:"value"`
			Adds    *string `json:"adds"`
			Comment *string `json:"comment"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		var prior int
		err := db.QueryRow(context.Background(), "INSERT INTO snt_contacts (prior, type, value, adds, comment) VALUES ($1, $2, $3, $4, $5) RETURNING id", input.Prior, input.Type, input.Value, input.Adds, input.Comment).Scan(&prior)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"prior":   prior,
			"type":    input.Type,
			"value":   input.Value,
			"adds":    input.Adds,
			"comment": input.Comment,
		})
	})

	e.PUT("/contacts/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
		}

		var input struct {
			Prior   int     `json:"prior"`
			Type    string  `json:"type"`
			Value   string  `json:"value"`
			Adds    *string `json:"adds"`
			Comment *string `json:"comment"`
		}
		if err := c.Bind(&input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		_, err = db.Exec(context.Background(), "UPDATE snt_contacts SET prior = $1, type = $2, value = $3, adds = $4, comment = $5, modified = CURRENT_TIMESTAMP WHERE id = $6", input.Prior, input.Type, input.Value, input.Adds, input.Comment, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"id":      id,
			"prior":   input.Prior,
			"type":    input.Type,
			"value":   input.Value,
			"adds":    input.Adds,
			"comment": input.Comment,
		})
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
		}

		tag, err := db.Exec(context.Background(), "DELETE FROM snt_contacts WHERE id = $1", id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}
		if tag.RowsAffected() == 0 {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "Contact not found"})
		}

		return c.NoContent(http.StatusNoContent)
	})

	// Serve static files (if build exists)
	e.Static("/", "client/dist")

	// Fallback for SPA (index.html)
	e.File("/*", "client/dist/index.html")

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(cfg.Port)))
}
