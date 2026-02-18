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
                rows, err := db.Query(context.Background(), "SELECT id, level, message, details, created_at FROM bot_logs ORDER BY created_at DESC LIMIT 100")
                if err != nil {
                        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
                }
                defer rows.Close()

                logs := []map[string]interface{}{}
                for rows.Next() {
                        var id int
                        var level, message string
                        var details *string
                        var createdAt time.Time
                        if err := rows.Scan(&id, &level, &message, &details, &createdAt); err != nil {
                                continue
                        }
                        logs = append(logs, map[string]interface{}{
                                "id":        id,
                                "level":     level,
                                "message":   message,
                                "details":   details,
                                "createdAt": createdAt,
                        })
                }
                return c.JSON(http.StatusOK, logs)
        })

        // Contacts
        e.GET("/contacts", func(c echo.Context) error {
                rows, err := db.Query(context.Background(), "SELECT id, name, description, phone, email FROM snt_contacts ORDER BY name")
                if err != nil {
                        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
                }
                defer rows.Close()

                contacts := []map[string]interface{}{}
                for rows.Next() {
                        var id int
                        var name string
                        var description, phone, email *string
                        if err := rows.Scan(&id, &name, &description, &phone, &email); err != nil {
                                continue
                        }
                        contacts = append(contacts, map[string]interface{}{
                                "id":          id,
                                "name":        name,
                                "description": description,
                                "phone":       phone,
                                "email":       email,
                        })
                }
                return c.JSON(http.StatusOK, contacts)
        })

        e.POST("/contacts", func(c echo.Context) error {
                var input struct {
                        Name        string  `json:"name"`
                        Description *string `json:"description"`
                        Phone       *string `json:"phone"`
                        Email       *string `json:"email"`
                }
                if err := c.Bind(&input); err != nil {
                        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
                }

                var id int
                err := db.QueryRow(context.Background(), "INSERT INTO snt_contacts (name, description, phone, email) VALUES ($1, $2, $3, $4) RETURNING id", input.Name, input.Description, input.Phone, input.Email).Scan(&id)
                if err != nil {
                        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
                }

                return c.JSON(http.StatusCreated, map[string]interface{}{
                        "id":          id,
                        "name":        input.Name,
                        "description": input.Description,
                        "phone":       input.Phone,
                        "email":       input.Email,
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
