# Gemini Project Configuration

This file helps Gemini understand the project's context, conventions, and commands.

## Project Overview

**JobWatch Canada** is a community-driven business demographics reporting platform that provides transparency about Temporary Foreign Worker (TFW) program usage in Canadian businesses. The platform allows verified users to submit reports about hiring practices, creating a public directory with confidence-scored ratings to help Canadians make informed decisions about which businesses to support.

**Data Disclaimer**: All government and business data is sourced from official databases and public records. While we make every effort to ensure accuracy, the complexity of these data systems may result in occasional errors. Our goal is to make this information more accessible and transparent. If you notice any inaccuracies or have feedback, please contact us through our feedback page - we value your input and will work promptly to address any issues.

See `docs/FUNCTIONAL_OUTLINE.md` and `docs/BUSINESS_DEMOGRAPHICS_FEATURE_SPEC.md` for detailed project requirements and goals.

## Project Architecture

This is a full-stack application with separate Go backend API and React frontend:

- **`api/`** - Go backend with PostgreSQL database
- **`web/`** - React frontend using Vite build system

## Development Commands

### Frontend (web/)
- `npm install` - Install dependencies
- `npm run dev` - Start development server (runs on port 5173)
- `npm run build` - Build for production
- `npm run lint` - Run ESLint
- `npm run preview` - Preview production build

### Backend (api/)
- `go mod tidy` - Install/update dependencies
- `go run main.go` - Start API server on port 8000
- `go build` - Build binary
- `make build` - Build using makefile (if available)

### Database Setup
- Ensure PostgreSQL is running locally
- Database migrations run automatically on server startup
- Check `.env` file for database connection settings

## Technology Stack

- **Backend**: Go with Chi router, CORS middleware, PostgreSQL database
- **Authentication**: JWT tokens, passwordless email login, session management
- **Database**: PostgreSQL with sqlx, automated migrations, repository pattern
- **Email**: SMTP email service for authentication links
- **Dependency Injection**: Uber Dig container for service management
- **Frontend**: React 18, TypeScript, Vite, TailwindCSS
- **Routing**: TanStack Router with file-based routing and authentication guards
- **State Management**: TanStack Query (React Query) with axios API client
- **Components**: Shadcn/ui component library, Radix UI primitives
- **Styling**: TailwindCSS v4 with Vite plugin, responsive design
- **Development**: Vite dev server, ESLint, TypeScript, FontAwesome icons

## Code Conventions

- Frontend uses TypeScript with strict mode
- Path alias `@/*` for src imports
- ESLint configuration with React hooks and TypeScript rules
- TanStack Router with auto-generated route tree
- Components follow Shadcn/ui patterns (based on components.json)
- Go backend follows repository pattern with dependency injection
- Database operations use sqlx with prepared statements
- Error handling with structured logging using charmbracelet/log
- **API Calls**: Always use TanStack Query (useQuery) hooks for API calls instead of manual fetch/useEffect patterns