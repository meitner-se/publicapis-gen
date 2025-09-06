package main

import (
	"fmt"
	"log"
	"strings"
)

// User represents a user in the directory
type User struct {
	ID       string
	Username string
	Email    string
	Groups   []string
}

// Group represents a group in the directory
type Group struct {
	ID          string
	Name        string
	Description string
	Members     []string
}

// DirectoryAPI represents a simple directory management system
type DirectoryAPI struct {
	users  map[string]*User
	groups map[string]*Group
}

// NewDirectoryAPI creates a new directory API instance
func NewDirectoryAPI() *DirectoryAPI {
	return &DirectoryAPI{
		users:  make(map[string]*User),
		groups: make(map[string]*Group),
	}
}

// CreateUser creates a new user in the directory
func (d *DirectoryAPI) CreateUser(id, username, email string) (*User, error) {
	if _, exists := d.users[id]; exists {
		return nil, fmt.Errorf("user with ID %s already exists", id)
	}

	user := &User{
		ID:       id,
		Username: username,
		Email:    email,
		Groups:   []string{},
	}
	
	d.users[id] = user
	return user, nil
}

// GetUser retrieves a user by ID
func (d *DirectoryAPI) GetUser(id string) (*User, error) {
	user, exists := d.users[id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s not found", id)
	}
	return user, nil
}

// CreateGroup creates a new group in the directory
func (d *DirectoryAPI) CreateGroup(id, name, description string) (*Group, error) {
	if _, exists := d.groups[id]; exists {
		return nil, fmt.Errorf("group with ID %s already exists", id)
	}

	group := &Group{
		ID:          id,
		Name:        name,
		Description: description,
		Members:     []string{},
	}
	
	d.groups[id] = group
	return group, nil
}

// GetGroup retrieves a group by ID
func (d *DirectoryAPI) GetGroup(id string) (*Group, error) {
	group, exists := d.groups[id]
	if !exists {
		return nil, fmt.Errorf("group with ID %s not found", id)
	}
	return group, nil
}

// AddUserToGroup adds a user to a group
func (d *DirectoryAPI) AddUserToGroup(userID, groupID string) error {
	user, err := d.GetUser(userID)
	if err != nil {
		return err
	}

	group, err := d.GetGroup(groupID)
	if err != nil {
		return err
	}

	// Check if user is already in group
	for _, member := range group.Members {
		if member == userID {
			return fmt.Errorf("user %s is already a member of group %s", userID, groupID)
		}
	}

	// Add user to group members
	group.Members = append(group.Members, userID)
	
	// Add group to user's groups
	user.Groups = append(user.Groups, groupID)

	return nil
}

// ListUsers returns all users in the directory
func (d *DirectoryAPI) ListUsers() []*User {
	users := make([]*User, 0, len(d.users))
	for _, user := range d.users {
		users = append(users, user)
	}
	return users
}

// ListGroups returns all groups in the directory
func (d *DirectoryAPI) ListGroups() []*Group {
	groups := make([]*Group, 0, len(d.groups))
	for _, group := range d.groups {
		groups = append(groups, group)
	}
	return groups
}

func main() {
	fmt.Println("=== Directory API v1 Demo ===")
	fmt.Println("Initializing Directory API...")

	// Create a new directory API instance
	api := NewDirectoryAPI()

	// Create some sample users
	fmt.Println("\nCreating users...")
	users := []struct {
		id, username, email string
	}{
		{"user1", "alice", "alice@example.com"},
		{"user2", "bob", "bob@example.com"},
		{"user3", "charlie", "charlie@example.com"},
	}

	for _, u := range users {
		user, err := api.CreateUser(u.id, u.username, u.email)
		if err != nil {
			log.Printf("Error creating user %s: %v", u.username, err)
			continue
		}
		fmt.Printf("Created user: %s (%s)\n", user.Username, user.Email)
	}

	// Create some sample groups
	fmt.Println("\nCreating groups...")
	groups := []struct {
		id, name, description string
	}{
		{"group1", "Administrators", "System administrators with full access"},
		{"group2", "Users", "Regular users with standard access"},
		{"group3", "Developers", "Software developers with code access"},
	}

	for _, g := range groups {
		group, err := api.CreateGroup(g.id, g.name, g.description)
		if err != nil {
			log.Printf("Error creating group %s: %v", g.name, err)
			continue
		}
		fmt.Printf("Created group: %s - %s\n", group.Name, group.Description)
	}

	// Add users to groups
	fmt.Println("\nAssigning users to groups...")
	assignments := []struct {
		userID, groupID string
	}{
		{"user1", "group1"}, // Alice -> Administrators
		{"user1", "group3"}, // Alice -> Developers
		{"user2", "group2"}, // Bob -> Users
		{"user3", "group2"}, // Charlie -> Users
		{"user3", "group3"}, // Charlie -> Developers
	}

	for _, a := range assignments {
		err := api.AddUserToGroup(a.userID, a.groupID)
		if err != nil {
			log.Printf("Error adding user to group: %v", err)
			continue
		}
		
		user, _ := api.GetUser(a.userID)
		group, _ := api.GetGroup(a.groupID)
		fmt.Printf("Added %s to %s\n", user.Username, group.Name)
	}

	// Display final state
	fmt.Println("\n=== Final Directory State ===")
	
	fmt.Println("\nUsers:")
	for _, user := range api.ListUsers() {
		groupNames := []string{}
		for _, groupID := range user.Groups {
			if group, err := api.GetGroup(groupID); err == nil {
				groupNames = append(groupNames, group.Name)
			}
		}
		fmt.Printf("- %s (%s) - Groups: [%s]\n", user.Username, user.Email, strings.Join(groupNames, ", "))
	}

	fmt.Println("\nGroups:")
	for _, group := range api.ListGroups() {
		userNames := []string{}
		for _, userID := range group.Members {
			if user, err := api.GetUser(userID); err == nil {
				userNames = append(userNames, user.Username)
			}
		}
		fmt.Printf("- %s: %s - Members: [%s]\n", group.Name, group.Description, strings.Join(userNames, ", "))
	}

	fmt.Println("\nDirectory API v1 demo completed successfully!")
}