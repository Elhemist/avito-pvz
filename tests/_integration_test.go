package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"pvz-test/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api"

func TestIntegration_ProcessPVZReception(t *testing.T) {

	t.Skip("Skipping test in short mode.")
	var pvzID uuid.UUID
	var token string

	t.Run("Moderator dummy login", func(t *testing.T) {
		reqBody := models.DummyLoginRequest{Role: models.RoleModerator}
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&token)

		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("Create PVZ", func(t *testing.T) {
		reqBody := models.PVZRequest{City: "Москва"}
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, baseURL+"/pvz", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		var createdPVZ models.PVZ
		err = json.NewDecoder(resp.Body).Decode(&createdPVZ)
		require.NoError(t, err)
		assert.Equal(t, "Москва", createdPVZ.City)

		pvzID = createdPVZ.ID
	})

	t.Run("Employee dummy login", func(t *testing.T) {
		reqBody := models.DummyLoginRequest{Role: models.RoleEmployee}
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewBuffer(body))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&token)

		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("Create Reception", func(t *testing.T) {
		reqBody := models.CreateReceptionRequest{PvzID: pvzID}
		body, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, baseURL+"/receptions", bytes.NewBuffer(body))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
		defer resp.Body.Close()

		var createdReception models.Reception
		err = json.NewDecoder(resp.Body).Decode(&createdReception)
		assert.NoError(t, err)
		require.NotEmpty(t, token)
		assert.Equal(t, pvzID, createdReception.PVZID)
	})

	t.Run("Add 50 Items", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			reqBody := models.AddProductRequest{PvzID: pvzID, Type: string(models.ItemTypeElectronics)}
			body, err := json.Marshal(reqBody)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, baseURL+"/products", bytes.NewBuffer(body))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			defer resp.Body.Close()

			require.Equal(t, http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("Close Reception", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, baseURL+"/pvz/"+pvzID.String()+"/close_last_reception", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var closedReception models.Reception
		err = json.NewDecoder(resp.Body).Decode(&closedReception)
		assert.NoError(t, err)
		assert.Equal(t, "close", closedReception.Status)
	})
}
