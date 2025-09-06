package main

import (
	"fmt"
	"testing"
)

func TestNewDirectoryAPI(t *testing.T) {
	api := NewDirectoryAPI()
	
	if api == nil {
		t.Fatal("NewDirectoryAPI() returned nil")
	}
	
	if api.users == nil {
		t.Error("users map not initialized")
	}
	
	if api.groups == nil {
		t.Error("groups map not initialized")
	}
	
	if len(api.users) != 0 {
		t.Error("users map should be empty initially")
	}
	
	if len(api.groups) != 0 {
		t.Error("groups map should be empty initially")
	}
}

func TestCreateUser(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test successful user creation
	user, err := api.CreateUser("user1", "testuser", "test@example.com")
	if err != nil {
		t.Fatalf("CreateUser() failed: %v", err)
	}
	
	if user.ID != "user1" {
		t.Errorf("Expected user ID 'user1', got '%s'", user.ID)
	}
	
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	
	if len(user.Groups) != 0 {
		t.Errorf("Expected empty groups slice, got %v", user.Groups)
	}
	
	// Test duplicate user creation
	_, err = api.CreateUser("user1", "duplicate", "duplicate@example.com")
	if err == nil {
		t.Error("Expected error when creating duplicate user, got nil")
	}
}

func TestGetUser(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test getting non-existent user
	_, err := api.GetUser("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user, got nil")
	}
	
	// Create and get user
	originalUser, _ := api.CreateUser("user1", "testuser", "test@example.com")
	
	retrievedUser, err := api.GetUser("user1")
	if err != nil {
		t.Fatalf("GetUser() failed: %v", err)
	}
	
	if retrievedUser.ID != originalUser.ID {
		t.Errorf("Retrieved user ID doesn't match: expected '%s', got '%s'", originalUser.ID, retrievedUser.ID)
	}
	
	if retrievedUser.Username != originalUser.Username {
		t.Errorf("Retrieved username doesn't match: expected '%s', got '%s'", originalUser.Username, retrievedUser.Username)
	}
}

func TestCreateGroup(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test successful group creation
	group, err := api.CreateGroup("group1", "Test Group", "A test group")
	if err != nil {
		t.Fatalf("CreateGroup() failed: %v", err)
	}
	
	if group.ID != "group1" {
		t.Errorf("Expected group ID 'group1', got '%s'", group.ID)
	}
	
	if group.Name != "Test Group" {
		t.Errorf("Expected group name 'Test Group', got '%s'", group.Name)
	}
	
	if group.Description != "A test group" {
		t.Errorf("Expected description 'A test group', got '%s'", group.Description)
	}
	
	if len(group.Members) != 0 {
		t.Errorf("Expected empty members slice, got %v", group.Members)
	}
	
	// Test duplicate group creation
	_, err = api.CreateGroup("group1", "Duplicate Group", "Duplicate description")
	if err == nil {
		t.Error("Expected error when creating duplicate group, got nil")
	}
}

func TestGetGroup(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test getting non-existent group
	_, err := api.GetGroup("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent group, got nil")
	}
	
	// Create and get group
	originalGroup, _ := api.CreateGroup("group1", "Test Group", "A test group")
	
	retrievedGroup, err := api.GetGroup("group1")
	if err != nil {
		t.Fatalf("GetGroup() failed: %v", err)
	}
	
	if retrievedGroup.ID != originalGroup.ID {
		t.Errorf("Retrieved group ID doesn't match: expected '%s', got '%s'", originalGroup.ID, retrievedGroup.ID)
	}
	
	if retrievedGroup.Name != originalGroup.Name {
		t.Errorf("Retrieved group name doesn't match: expected '%s', got '%s'", originalGroup.Name, retrievedGroup.Name)
	}
}

func TestAddUserToGroup(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Create user and group
	api.CreateUser("user1", "testuser", "test@example.com")
	api.CreateGroup("group1", "Test Group", "A test group")
	
	// Test successful addition
	err := api.AddUserToGroup("user1", "group1")
	if err != nil {
		t.Fatalf("AddUserToGroup() failed: %v", err)
	}
	
	// Verify user is in group
	user, _ := api.GetUser("user1")
	if len(user.Groups) != 1 || user.Groups[0] != "group1" {
		t.Errorf("User groups not updated correctly: got %v", user.Groups)
	}
	
	// Verify group has user
	group, _ := api.GetGroup("group1")
	if len(group.Members) != 1 || group.Members[0] != "user1" {
		t.Errorf("Group members not updated correctly: got %v", group.Members)
	}
	
	// Test adding same user again
	err = api.AddUserToGroup("user1", "group1")
	if err == nil {
		t.Error("Expected error when adding user to group they're already in, got nil")
	}
	
	// Test adding non-existent user
	err = api.AddUserToGroup("nonexistent", "group1")
	if err == nil {
		t.Error("Expected error when adding non-existent user to group, got nil")
	}
	
	// Test adding user to non-existent group
	err = api.AddUserToGroup("user1", "nonexistent")
	if err == nil {
		t.Error("Expected error when adding user to non-existent group, got nil")
	}
}

func TestListUsers(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test empty list
	users := api.ListUsers()
	if len(users) != 0 {
		t.Errorf("Expected empty user list, got %d users", len(users))
	}
	
	// Create some users
	api.CreateUser("user1", "alice", "alice@example.com")
	api.CreateUser("user2", "bob", "bob@example.com")
	
	users = api.ListUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	
	// Verify all created users are in the list
	userIDs := make(map[string]bool)
	for _, user := range users {
		userIDs[user.ID] = true
	}
	
	if !userIDs["user1"] || !userIDs["user2"] {
		t.Error("Not all created users found in list")
	}
}

func TestListGroups(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Test empty list
	groups := api.ListGroups()
	if len(groups) != 0 {
		t.Errorf("Expected empty group list, got %d groups", len(groups))
	}
	
	// Create some groups
	api.CreateGroup("group1", "Admins", "Administrator group")
	api.CreateGroup("group2", "Users", "Regular users")
	
	groups = api.ListGroups()
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
	
	// Verify all created groups are in the list
	groupIDs := make(map[string]bool)
	for _, group := range groups {
		groupIDs[group.ID] = true
	}
	
	if !groupIDs["group1"] || !groupIDs["group2"] {
		t.Error("Not all created groups found in list")
	}
}

func TestCompleteWorkflow(t *testing.T) {
	api := NewDirectoryAPI()
	
	// Create users
	alice, err := api.CreateUser("alice", "alice", "alice@example.com")
	if err != nil {
		t.Fatalf("Failed to create Alice: %v", err)
	}
	
	bob, err := api.CreateUser("bob", "bob", "bob@example.com")
	if err != nil {
		t.Fatalf("Failed to create Bob: %v", err)
	}
	
	// Create groups
	admins, err := api.CreateGroup("admins", "Administrators", "System administrators")
	if err != nil {
		t.Fatalf("Failed to create admins group: %v", err)
	}
	
	users, err := api.CreateGroup("users", "Users", "Regular users")
	if err != nil {
		t.Fatalf("Failed to create users group: %v", err)
	}
	
	// Add users to groups
	err = api.AddUserToGroup("alice", "admins")
	if err != nil {
		t.Fatalf("Failed to add Alice to admins: %v", err)
	}
	
	err = api.AddUserToGroup("alice", "users")
	if err != nil {
		t.Fatalf("Failed to add Alice to users: %v", err)
	}
	
	err = api.AddUserToGroup("bob", "users")
	if err != nil {
		t.Fatalf("Failed to add Bob to users: %v", err)
	}
	
	// Verify final state
	alice, _ = api.GetUser("alice")
	if len(alice.Groups) != 2 {
		t.Errorf("Alice should be in 2 groups, got %d", len(alice.Groups))
	}
	
	bob, _ = api.GetUser("bob")
	if len(bob.Groups) != 1 {
		t.Errorf("Bob should be in 1 group, got %d", len(bob.Groups))
	}
	
	admins, _ = api.GetGroup("admins")
	if len(admins.Members) != 1 || admins.Members[0] != "alice" {
		t.Errorf("Admins group should have 1 member (alice), got %v", admins.Members)
	}
	
	users, _ = api.GetGroup("users")
	if len(users.Members) != 2 {
		t.Errorf("Users group should have 2 members, got %d", len(users.Members))
	}
	
	// Verify both alice and bob are in users group
	memberMap := make(map[string]bool)
	for _, member := range users.Members {
		memberMap[member] = true
	}
	
	if !memberMap["alice"] || !memberMap["bob"] {
		t.Error("Both Alice and Bob should be members of users group")
	}
}

// Benchmark test for creating users
func BenchmarkCreateUser(b *testing.B) {
	api := NewDirectoryAPI()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i)
		username := fmt.Sprintf("username%d", i)
		email := fmt.Sprintf("user%d@example.com", i)
		
		_, err := api.CreateUser(userID, username, email)
		if err != nil {
			b.Fatalf("CreateUser failed: %v", err)
		}
	}
}

// Benchmark test for adding users to groups
func BenchmarkAddUserToGroup(b *testing.B) {
	api := NewDirectoryAPI()
	
	// Setup
	api.CreateGroup("testgroup", "Test Group", "A test group")
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i)
		api.CreateUser(userID, fmt.Sprintf("username%d", i), fmt.Sprintf("user%d@example.com", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i)
		err := api.AddUserToGroup(userID, "testgroup")
		if err != nil {
			b.Fatalf("AddUserToGroup failed: %v", err)
		}
	}
}