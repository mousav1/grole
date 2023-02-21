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

	grole.DeletePermission(permission.ID)
}

func TestUpdatePermission(t *testing.T) {
	grole.New(grole.Options{
		DB: db,
	})

	oldPermission, errCreateOldPermission := grole.FindOrCreatePermission(models.Permission{
		Name:        "manage-articles",
		Description: "test",
	})

	_, errUpdatePermission := grole.UpdatePermission(oldPermission.ID, models.Permission{
		Name:        "manage-articles-update",
		Description: "update test update",
	})

	newPermission, errNewPermission := grole.FindPermissionById(oldPermission.ID)

	require.NoError(t, errCreateOldPermission)
	require.NoError(t, errUpdatePermission)
	require.NoError(t, errNewPermission)
	require.NotEmpty(t, newPermission)
	require.NotZero(t, newPermission.ID)
	require.Equal(t, "manage-articles-update", newPermission.Name)
	require.Equal(t, "update test update", newPermission.Description)
	require.NotEqual(t, oldPermission.Name, newPermission.Name)
	require.NotEqual(t, oldPermission.Description, newPermission.Description)
	require.NotZero(t, oldPermission.ID)

	grole.DeletePermission(oldPermission.ID)
}

func TestDeletePermission(t *testing.T) {
	grole.New(grole.Options{
		DB: db,
	})

	permission, error := grole.FindOrCreatePermission(models.Permission{
		Name:        "manage-articles",
		Description: "test",
	})

	findPermission, errFindPermission := grole.FindPermissionById(permission.ID)

	require.NoError(t, error)
	require.NoError(t, errFindPermission)
	require.NotEmpty(t, findPermission)
	require.NotZero(t, findPermission.ID)
	require.Equal(t, permission.Name, findPermission.Name)
	require.Equal(t, permission.Description, findPermission.Description)
	require.NotZero(t, permission.ID)

	grole.DeletePermission(permission.ID)

	_, errNewPermissionDelete := grole.FindPermissionById(permission.ID)
	require.EqualError(t, errNewPermissionDelete, "RECORD NOT FOUND")
}
