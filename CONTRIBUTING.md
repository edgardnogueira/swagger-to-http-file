# Contributing to Swagger to HTTP File Converter

Thank you for your interest in contributing to this project! Here's how you can help make this project better.

## Development Setup

1. Fork and clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Build the project:
   ```bash
   make build
   ```

## Development Workflow

We follow [Gitflow](https://nvie.com/posts/a-successful-git-branching-model/) for branch management.

1. **Main Branch**: Contains the latest stable release.
2. **Feature Branches**: For new features, create a branch from `main` named `feature/issue-X-description`.
3. **Bugfix Branches**: For bugfixes, create a branch from `main` named `fix/issue-X-description`.
4. **Release Branches**: Created from `main` when preparing a new release, named `release/X.Y.Z`.

### Making Changes

1. Create a new branch for your feature or fix:
   ```bash
   git checkout -b feature/my-new-feature
   ```
2. Make your changes and commit using conventional commit messages:
   ```bash
   git commit -m "feat: add new feature"
   ```
   
   Commit types:
   - `feat`: A new feature
   - `fix`: A bug fix
   - `docs`: Documentation only changes
   - `style`: Changes that do not affect the meaning of the code
   - `refactor`: A code change that neither fixes a bug nor adds a feature
   - `test`: Adding missing tests or correcting existing tests
   - `chore`: Changes to the build process or auxiliary tools

3. Push your branch and create a pull request:
   ```bash
   git push origin feature/my-new-feature
   ```

## Testing

Before submitting a pull request, ensure all tests pass:

```bash
make test
```

For test coverage:

```bash
make coverage
```

## Code Style

Follow Go's standard code style. Run linters and formatters:

```bash
make fmt
make lint
make vet
```

## Pull Request Process

1. Ensure your code adheres to the project's style and passes all tests
2. Update the documentation if necessary
3. Link the pull request to any related issues
4. The PR title should be descriptive and follow the same conventions as commit messages
5. Wait for a review from a maintainer

## Releasing

Releases are handled by maintainers using the following process:

1. Update version information in appropriate files
2. Tag the release in git:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. GitHub Actions will automatically build and publish the release

## License

By contributing to this project, you agree that your contributions will be licensed under the same [MIT License](LICENSE) that covers the project.
