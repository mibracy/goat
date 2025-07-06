package models

// UserStorage defines the interface for user data operations.
type UserStorage interface {
	List() []*User          // List retrieves all users.
	Get(int) *User          // Get retrieves a single user by ID.
	Update(int, User) *User // Update modifies an existing user.
	Create(User)            // Create adds a new user.
	Delete(int) *User       // Delete removes a user by ID.
}

// User represents a user entity with an ID and a Name.
type User struct {
	ID   int    `json:"id"`   // Unique identifier for the user.
	Name string `json:"name"` // Name of the user.
}

// UserStore implements the UserStorage interface using an in-memory slice.
type UserStore struct {
}

// users is an in-memory slice that acts as a mock database for User objects.
var users = []*User{
	{
		ID:   1,
		Name: "John Doe",
	},
	{
		ID:   2,
		Name: "Jane Smith",
	},
}

// Get retrieves a user by their ID from the in-memory store.
// It returns a pointer to the User if found, otherwise returns nil.
func (u UserStore) Get(id int) *User {
	for _, user := range users {
		if user.ID == id {
			return user
		}
	}
	return nil
}

// List returns all users currently in the in-memory store.
func (u UserStore) List() []*User {
	return users
}

// Create adds a new user to the in-memory store.
func (u UserStore) Create(user User) {
	users = append(users, &user)
}

// Delete removes a user by their ID from the in-memory store.
// It returns a pointer to an empty User struct if the user was found and deleted, otherwise returns nil.
func (u UserStore) Delete(id int) *User {
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], (users)[i+1:]...)
			return &User{}
		}
	}
	return nil
}

// Update modifies an existing user in the in-memory store.
// It takes the ID of the user to update and a User object with the new data.
// It returns a pointer to the updated User if found, otherwise returns nil.
func (u UserStore) Update(id int, userUpdate User) *User {
	for i, user := range users {
		if user.ID == id {
			users[i] = &userUpdate
			return users[i]
		}
	}
	return nil
}
