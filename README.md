# Grole
manage user permissions and roles


# Install
```bash
go get github.com/mousav1/grole
```

Then get the database driver because we use gorm library and to use it, we need one of the database driver to connect to the database.

```bash
# mysql 
go get gorm.io/driver/mysql 
# or postgres
go get gorm.io/driver/postgres
# or sqlite
go get gorm.io/driver/sqlite
# or sqlserver
go get gorm.io/driver/sqlserver
# or clickhouse
go get gorm.io/driver/clickhouse
```

# Initialize
Initialize the new Grole.

```go

// initialize the database. 
var DB *gorm.DB
dsn := "host=localhost user=postgres password=postgres dbname=grole port=5432 sslmode=disable TimeZone=Asia/Shanghai"
DB, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

//To initiate the grole, you need to pass the DB variable 

grole.New(grole.Options{
    DB: DB,
})
```

# Usage
After installed you can do things like this:

```go
// Permissions


// find Permission or Create Permission If not found
err = grole.FindOrCreatePermission(models.Permission{
		Name:        "manage-articles",
		Description: "test",
	})
// output (models.Permission, error) => {1 manage-articles test []} <nil>


// update the given permission
err = grole.UpdatePermission(1, models.Permission{
		Name:        "manage-articles",
		Description: "update test",
	})
// output (bool, error) => true <nil>


// delete the given permission
err = grole.DeletePermission(1)
// output (bool, error) => true <nil>


// Show All Permission With Role
err = grole.FindAllPermission()
// output ([]models.Permission, error) => {1 manage-articles test []} <nil>


// find Permission By Name and Show each with Role
err = grole.FindPermissionByName("manage-articles")
// output (models.Permission, error) => {1 manage-articles test []} <nil>


// find Permission By Id and Show each with Role
err = grole.FindPermissionById(1)
// output (models.Permission, error) => {1 manage-articles test []} <nil>


// Find roles of each permission
err = grole.Roles("manage-articles")
// output ([]models.Role) => [{1 admin test []}] <nil>


// Revoke the given role by id for permission
err = grole.RemoveRoleByIdFromPermission(1, 1)
// output (bool, error) => true <nil>


// Revoke the given role by name for permission
err = grole.RemoveRoleByNameFromPermission(1, "admin")
// output (bool, error) => true <nil>


// Revoke the given role for permission
err = grole.RemoveRoleFromPermission(1, "admin")
// output (bool, error) => true <nil>


// Return the number of permissions Role.
err = grole.CountRoleFromPermission(1)
// output (int64, error) => 1 <nil>


// Remove all current Role for Permission.
err = grole.RemoveAllRoleFromPermission(1)
// output (bool, error) => true <nil>


// Remove all current Permission role and set the given ones.
err = grole.SyncRolesFromPermission(1, "admin")
// output ([]models.Role, error) => [{1 admin test []}] <nil>



// Roles


// Find Or Create Role
err = grole.FindOrCreateRole(models.Role{
		Name:        "admin",
		Description: "test",
	})
// output (models.Permission, error) => {1 admin test []} <nil>


// update the given role
err = grole.UpdateRole(1, models.Role{
		Name:        "admin",
		Description: "teste test",
	})
// output (bool, error) => [{1 writer test []}] <nil>


// delete the given role
err = grole.DeleteRole(1)
// output (bool, error) => true <nil>

// Find All Role
err = grole.FindAllRole()
// output ([]models.Role, error) => [{1 admin test []} {2 user test []}] <nil>

// Find Role By Name
err = grole.FindRoleByName("user")
// output (models.Role, error) => {2 user test []} <nil>


// Find Role By Id
err = grole.FindRoleById(2)
// output (models.Role, error) => {2 user test []} <nil>


// Return all Permissions the Role.
err = grole.Permissions("admin")
// output ([]models.Permission) => [{1 manage-articles test []}] <nil>


// Revoke the given Permission by id for Role
err = grole.RemovePermissionByIdFromRole(1, 1)
// output (bool, error) => true <nil>


// Revoke the given Permission by name for Role
err = grole.RemovePermissionByNameFromRole("user", "manage-articles")
// output (bool, error) => true <nil>


// Return the number of Role permissions.
err = grole.CountPermissionFromRole(1)
// output (int64, error) => 1 <nil>


// Remove all current Permission for Role.
err = grole.RemoveAllPermissionFromRole(1)
// output (bool, error) => 1 <nil>


// Remove all current role Permission and set the given ones.
err = grole.SyncPermissionsFromRole(1, "manage-articles", "manage-user")
// output ([]models.Permission, error) => [{1 manage-articles test []} {2 manage-user test []}] <nil>


// Assign the given Permissions to the Role.
err = grole.AssignPermissionsFromRole(1, "manage-articles", "manage-user")
// output ([]models.Permission, error) => [{1 manage-articles test []} {2 manage-user test []}] <nil>


// Determine if the Role may perform the given permission.
err = grole.HasPermissionTo(1, "manage-user")
// output (models.Permission, error) => {2 manage-user test []} <nil>



// User

// Return all the Roles the user.
err = grole.GetRole(1)
// output ([]models.Role, error) => [{1 writer test []}] <nil>


// Return all the Roles Name the user.
err = grole.GetRoleNames(2)
// output ([]string, error) => [admin] <nil>


// Return all the permissions the user.
err = grole.GetAllPermissions(1)
// output ([]models.Permission, error) => [{1 manage-users test []}] <nil>


// Assign the given roles to the User.
err = grole.AssignRoles(1, "admin")
// output (bool, error) => true <nil>


// Revoke the given role by id for user
err = grole.RemoveRoleByIdFromUser(1, 2)
// output (bool, error) => true <nil>


// Revoke the given role by name for user
err = grole.RemoveRoleByNameFromUser(2, "admin")
// output (bool, error) => true <nil>


// Remove all current roles for user.
err = grole.RemoveAllRoleFromUser(1)
// output (bool, error) => true <nil>


// Remove all current user roles and set the given ones.
err = grole.SyncRolesFromUser(1, "writer", "writer2", "writer3")
// output (bool, error) => true <nil>


// Determine if the user has  of the given role id.
err = grole.HasRole(1, 2)
// output (bool, error) => true <nil>


// Determine if the user has of the given roles name.
err = grole.HasAnyRole(1, "admin", "writer")
// output (bool, error) => true <nil>


// Determine if the user has all of the given roles name.
err = grole.HasAllRole(1, "writer")
// output (bool, error) => true <nil>


// Determine if the user has all of the given permissions name.
err = grole.HasAllPermission(1, "manage-articles")
// output (bool, error) => true <nil>


// Determine if the User has of the given permissions name.
err = grole.HasAnyPermissions(1, "manage-articles", "manage-users")
// output (bool, error) => true <nil>


```
