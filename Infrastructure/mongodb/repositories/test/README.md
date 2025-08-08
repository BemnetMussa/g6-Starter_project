# Repository Tests

This directory contains comprehensive integration tests for all MongoDB repositories in the Blog API project.

## Test Structure

The tests are organized using a clean architecture approach with real MongoDB integration:

- **Real Database**: Tests use actual MongoDB connections to ensure realistic testing
- **Isolated Tests**: Each test runs in isolation with clean database state
- **Comprehensive Coverage**: Tests cover all CRUD operations and edge cases

## Test Files

### 1. `test_config.go`

- Centralized test configuration and utilities
- Database connection setup and teardown
- Helper functions for creating test data
- Assertion utilities

### 2. `user_repository_test.go`

Tests for `UserRepository` covering:

- User creation with various field combinations
- User retrieval by ID, email, username
- User updates and verification status changes
- Password reset token management
- User deletion and count operations
- Name-based search functionality

### 3. `blog_repository_test.go`

Tests for `BlogRepository`, `BlogInteractionRepository`, and `CommentRepository` covering:

- Blog creation and management
- Blog search and filtering (by author, tags, title)
- Blog interaction tracking (likes, dislikes, views)
- Comment creation and counting
- Pagination and advanced queries

### 4. `token_repository_test.go`

Tests for `TokenRepository` covering:

- Token creation and storage
- Token retrieval by user ID
- Token updates and modifications
- Token deletion and cleanup
- Complete token lifecycle management

### 5. `chat_repository_test.go`

Tests for `ChatRepository` covering:

- AI chat creation and storage
- Chat retrieval by ID and user ID
- Chat deletion and cleanup
- Multi-user chat management

## Running Tests

### Prerequisites

1. **MongoDB**: Ensure MongoDB is running or accessible via your connection string
2. **Environment Variables**: The tests will automatically load from your `.env` file:
   ```bash
   MONGODB_URI=mongodb+srv://leulgedion:kT5JsmzjYL8hVBrc@cluster0.1y2cmpf.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
   MONGODB_DATABASE=blog_api
   MONGODB_COLLECTION=posts
   ```

### Running All Repository Tests

```bash
# From the project root
go test ./Infrastructure/mongodb/repositories/test/...

# Or from the test directory
cd Infrastructure/mongodb/repositories/test
go test ./...
```

### Running Specific Test Files

```bash
# Run only user repository tests
go test -v -run TestUserRepository

# Run only blog repository tests
go test -v -run TestBlogRepository

# Run only token repository tests
go test -v -run TestTokenRepository

# Run only chat repository tests
go test -v -run TestChatRepository
```

### Running Specific Test Cases

```bash
# Run a specific test function
go test -v -run TestCreateUser

# Run tests with a specific pattern
go test -v -run "Test.*Create.*"

# Run tests with verbose output
go test -v ./...
```

## Test Configuration

### Default Configuration

- **MongoDB URI**: Uses your `.env` file's `MONGODB_URI`
- **Database Name**: Uses your `.env` file's `MONGODB_DATABASE` or defaults to `test_blog_api`
- **Collection Name**: Uses your `.env` file's `MONGODB_COLLECTION` or defaults to `users`

### Environment Variables

The tests automatically load from your `.env` file:

- `MONGODB_URI`: Your MongoDB connection string
- `MONGODB_DATABASE`: Database name (will create test database if different)
- `MONGODB_COLLECTION`: Default collection name

## Test Features

### 1. Real Database Integration

- Uses your actual MongoDB Atlas cluster
- Tests real database operations
- Validates actual data persistence

### 2. Isolated Test Environment

- Each test runs with a clean database state
- Collections are cleared before each test
- No test interference with your production data

### 3. Comprehensive Coverage

- **CRUD Operations**: Create, Read, Update, Delete
- **Edge Cases**: Invalid IDs, non-existent records
- **Error Handling**: Database errors, validation errors
- **Business Logic**: Complex queries, filtering, pagination

### 4. Test Data Management

- Helper functions for creating test entities
- Realistic test data with various field combinations
- Consistent test data across all tests

### 5. Assertion Utilities

- Custom assertion functions for common patterns
- User-friendly error messages
- Comprehensive field validation

## Test Patterns

### 1. Setup and Teardown

```go
func TestExample(t *testing.T) {
    ts := setupTestSuite(t)
    defer ts.teardown()

    // Test implementation
}
```

### 2. Test Data Creation

```go
// Create basic test user
user := CreateTestUser()

// Create user with custom fields
user := CreateTestUserWithCustomFields("John Doe", "johndoe", "john@example.com")

// Create verified user
user := CreateVerifiedUser()
```

### 3. Assertion Patterns

```go
// Basic assertions
assert.NoError(t, err)
assert.NotNil(t, result)
assert.Equal(t, expected, actual)

// Custom assertions
AssertUserFields(t, expectedUser, actualUser)
AssertUserExists(t, repo, userID)
AssertUserNotExists(t, repo, userID)
```

## Best Practices

### 1. Test Organization

- Group related tests in the same function
- Use descriptive test names
- Follow the Arrange-Act-Assert pattern

### 2. Test Data

- Use helper functions for consistent test data
- Create realistic test scenarios
- Test edge cases and error conditions

### 3. Cleanup

- Always use `defer ts.teardown()` to ensure cleanup
- Don't rely on test order
- Each test should be independent

### 4. Assertions

- Use specific assertions for better error messages
- Test both success and failure cases
- Validate all relevant fields

## Troubleshooting

### Common Issues

1. **MongoDB Connection Failed**

   - Ensure your MongoDB Atlas cluster is accessible
   - Check your connection string in the `.env` file
   - Verify network connectivity and firewall settings

2. **Test Database Not Found**

   - MongoDB will create the test database automatically
   - Check if your MongoDB user has create database permissions

3. **Collection Access Denied**
   - Ensure your MongoDB user has read/write permissions
   - Check if collections exist and are accessible

### Debug Mode

Run tests with verbose output for debugging:

```bash
go test -v -run TestCreateUser ./...
```

### Test Coverage

To check test coverage:

```bash
go test -cover ./...
```

## Next Steps

After completing repository tests, the next testing layers should be:

1. **Service Layer Tests** (Mock external dependencies)
2. **Use Case Layer Tests** (Mock repositories)
3. **Handler Layer Tests** (Mock use cases)
4. **Integration Tests** (End-to-end with real services)

This layered approach ensures comprehensive testing coverage while maintaining test isolation and performance.
