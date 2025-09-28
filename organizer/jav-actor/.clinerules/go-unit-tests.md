# Go Unit Testing Rules

- Split "Success" and "Error" case
- Use Table-Driven Tests for Error Cases
- Use `stretchr/testify` lib instead of raw comparsion
- Use `testing.T.Cleanup` instead of defer cleanup
