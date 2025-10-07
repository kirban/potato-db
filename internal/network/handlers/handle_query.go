package handlers

import (
	"fmt"
	"github.com/kirban/potato-db/internal/db"
)

type DatabaseHandler struct {
	Db db.Executable
}

func (h *DatabaseHandler) HandleRequest(req string) (string, error) {
	resp, err := h.Db.ExecuteQuery(req)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", resp), nil
}
