package test

import (
	"testing"

	"github.com/mousav1/grole"
	"github.com/mousav1/grole/models"
	"github.com/stretchr/testify/require"
)

func TestFindOrCreatePermission(t *testing.T) {
	grole.New(grole.Options{
		DB: db,
	})

	permission, err := grole.FindOrCreatePermission(models.Permission{
		Name:        "manage-articles",
		Description: "test",
	})

	require.NoError(t, err)
	require.NotEmpty(t, permission)
	require.Equal(t, "manage-articles", permission.Name)
	require.Equal(t, "test", permission.Description)

	require.NotZero(t, permission.ID)
	db.Where("name = ?", "manage-articles").Delete(models.Permission{})
}
