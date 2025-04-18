# Application Tracker Project Guidelines

This document outlines the coding standards, project structure, and workflow guidelines for the Application Tracker project.

## Table of Contents

1. [Project Structure](#project-structure)
2. [Code Style Guidelines](#code-style-guidelines)
3. [Naming Conventions](#naming-conventions)
4. [Git Workflow](#git-workflow)
5. [Testing Guidelines](#testing-guidelines)
6. [Documentation Standards](#documentation-standards)

## Project Structure

The Application Tracker project follows this directory structure:

```
ApplicationTracker/
├── api/                # API handlers and routes
├── data/               # Data storage files
├── models/             # Data models and business logic
├── schemas/            # JSON schemas for validation
├── static/             # Static assets (CSS, JS, images)
│   ├── css/            # CSS files (Tailwind)
│   ├── js/             # JavaScript files (HTMX)
│   └── images/         # Image assets
├── storage/            # Storage implementation
├── templates/          # HTML templates
│   ├── layouts/        # Base layout templates
│   ├── partials/       # Reusable template components
│   ├── pages/          # Full page templates
│   └── htmx/           # HTMX partial templates
├── ui/                 # UI-related Go code
└── main.go             # Application entry point
```

### Component Responsibilities

- **api/**: Contains all API endpoint handlers and route definitions
- **models/**: Defines data structures and business logic
- **storage/**: Implements data persistence (currently using JSON files)
- **ui/**: Implements UI-specific handlers and routes
- **templates/**: Contains all HTML templates organized by purpose

## Code Style Guidelines

### Go Code

1. Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) for general style guidance
2. Use `gofmt` or `goimports` to format code
3. Follow standard Go naming conventions (camelCase for private, PascalCase for exported)
4. Add comments for all exported functions, types, and constants
5. Keep functions focused and small (under 50 lines when possible)
6. Use error wrapping with `fmt.Errorf("context: %w", err)` for better error context
7. Use consistent error handling patterns

### HTML/Templates

1. Use consistent indentation (2 spaces)
2. Use semantic HTML elements
3. Define templates with clear names: `{{ define "templateName" }}`
4. Keep templates modular and reusable
5. Use consistent naming for template variables

### CSS (Tailwind)

1. Use Tailwind utility classes directly in HTML
2. Follow a consistent order for utility classes (layout, typography, colors, etc.)
3. Extract common patterns to components when they're reused frequently

### JavaScript/HTMX

1. Use HTMX attributes for most interactive features
2. Keep custom JavaScript minimal and focused
3. Follow standard JavaScript naming conventions (camelCase)

## Naming Conventions

### Go Code

- **Packages**: Short, lowercase, single-word names (e.g., `api`, `storage`)
- **Variables**: Descriptive camelCase names (e.g., `applicationID`, `userInput`)
- **Constants**: PascalCase for exported, camelCase for package-level
- **Interfaces**: PascalCase, often ending with 'er' (e.g., `Storage`, `Renderer`)
- **Structs**: PascalCase, descriptive nouns (e.g., `Application`, `Response`)
- **Methods**: PascalCase for exported, camelCase for private

### Templates

- **Template Files**: Descriptive, lowercase with hyphens (e.g., `application-list.html`)
- **Template Definitions**: Descriptive, lowercase with hyphens (e.g., `{{ define "application-card" }}`)

### Routes

- **API Routes**: RESTful conventions (e.g., `/api/applications`, `/api/applications/:id`)
- **UI Routes**: Simple, descriptive paths (e.g., `/applications`, `/applications/new`)

## Git Workflow

### Branching Strategy

1. `main` - Production-ready code
2. `develop` - Integration branch for features
3. Feature branches - Named as `feature/short-description`
4. Bug fix branches - Named as `fix/issue-description`

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `test`: Adding or modifying tests
- `chore`: Changes to the build process or auxiliary tools

Examples:
- `feat(ui): add application search functionality`
- `fix(api): correct status update validation`
- `docs: update README with setup instructions`

### Pull Request Process

1. Create a feature/fix branch from `develop`
2. Make changes and commit with descriptive messages
3. Push branch and create a pull request to `develop`
4. Ensure tests pass and code meets guidelines
5. Request review from at least one team member
6. Merge after approval and delete the branch

## Testing Guidelines

### Unit Tests

1. Write tests for all business logic and utility functions
2. Use table-driven tests for functions with multiple test cases
3. Mock external dependencies for isolation
4. Aim for high test coverage of critical paths

### Integration Tests

1. Test API endpoints with realistic data
2. Verify correct HTTP status codes and response formats
3. Test both success and error cases

### UI Tests

1. Test critical user flows
2. Verify that UI components render correctly
3. Test HTMX interactions

## Documentation Standards

### Code Documentation

1. Add comments for all exported functions, types, and constants
2. Use godoc-compatible comments for Go code
3. Document complex algorithms or business logic
4. Keep comments up-to-date with code changes

### Project Documentation

1. Maintain a comprehensive README with:
   - Project overview
   - Setup instructions
   - Usage examples
   - Architecture overview
2. Document API endpoints with examples
3. Keep UI implementation plans up-to-date
4. Document significant design decisions

### Changelog

Maintain a CHANGELOG.md file following the [Keep a Changelog](https://keepachangelog.com/) format to track notable changes between versions.

---

These guidelines are meant to evolve with the project. Suggestions for improvements are welcome through pull requests or issues.