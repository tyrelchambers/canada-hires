# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**JobWatch Canada** is a community-driven business demographics reporting platform that provides transparency about Temporary Foreign Worker (TFW) program usage in Canadian businesses. The platform allows verified users to submit reports about hiring practices, creating a public directory with confidence-scored ratings to help Canadians make informed decisions about which businesses to support.

**Data Disclaimer**: All government and business data is sourced from official databases and public records. While we make every effort to ensure accuracy, the complexity of these data systems may result in occasional errors. Our goal is to make this information more accessible and transparent. If you notice any inaccuracies or have feedback, please contact us through our feedback page - we value your input and will work promptly to address any issues.

Key features include:
- Passwordless email authentication system
- Business directory with location-based search
- Community reporting system for TFW usage
- AI-powered content generation using Gemini 2.5 Flash for Reddit posts
- Confidence scoring algorithm based on user verification tiers
- Public business ratings (Green/Yellow/Red) based on TFW percentage
- Feedback system for data accuracy and platform improvements

See `docs/FUNCTIONAL_OUTLINE.md` and `docs/BUSINESS_DEMOGRAPHICS_FEATURE_SPEC.md` for detailed project requirements and goals.

## Project Architecture

This is a full-stack application with separate Go backend API and React frontend:

- **`api/`** - Go backend with PostgreSQL database
  - `main.go` - HTTP server on port 8000 with Chi router, CORS middleware, and database initialization
  - `controllers/` - API controllers for auth, business, report, and user management
  - `services/` - Business logic services including email service and authentication
  - `router/` - Modular router configuration with separate route files
  - `models/` - Database models for users, sessions, login tokens
  - `repos/` - Repository pattern for database operations
  - `db/` - Database connection and migration management
  - `middleware/` - Session middleware and authentication
  - `migrations/` - SQL migration files for database schema
  - `container/` - Dependency injection container setup

- **`web/`** - React frontend using Vite build system
  - Built with React 18, TypeScript, and TanStack Router
  - Uses TailwindCSS for styling and Shadcn/ui components
  - React Query for state management and API calls
  - Authentication hooks and protected routes
  - Path alias `@/*` maps to `src/*`
  - Responsive design with mobile support

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

## Required Setup

### Environment Variables

#### API Environment (.env in api/ directory)
```
DATABASE_URL=postgres://username:password@localhost/dbname?sslmode=disable
EMAIL_SMTP_HOST=your_smtp_host
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USER=your_email@domain.com
EMAIL_SMTP_PASSWORD=your_app_password
JWT_SECRET=your_jwt_secret_key
GOOGLE_API_KEY=your_gemini_api_key_here
```

#### Frontend Environment (.env in web/ directory)
```
VITE_MAPBOX_ACCESS_TOKEN=your_mapbox_access_token_here
```

**Mapbox Setup**: Create a free account at [mapbox.com](https://mapbox.com) and get an access token for address search functionality. Without this token, the address search will fall back to a regular text input.

**Gemini AI Setup**: Create a free Google AI Studio account at [aistudio.google.com](https://aistudio.google.com) and generate an API key for Gemini content generation. This enables automatic generation of sarcastic Reddit post content for job approvals. Without this key, job approvals will work but without AI-generated content.

### Database Setup
1. Install and start PostgreSQL
2. Create a database for the project
3. Update `DATABASE_URL` in `.env` file
4. Migrations will run automatically when starting the API server

### Development Startup Sequence
1. **Start Database**: Ensure PostgreSQL is running
2. **Start API**: `cd api && go run main.go` (runs on port 8000)
3. **Start Frontend**: `cd web && npm run dev` (runs on port 5173)
4. **Access Application**: Frontend at http://localhost:5173, API at http://localhost:8000

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
- **Icons**: Use FontAwesome icons for consistency across the application

## Current Implementation Status

### Completed Features (Phase 1 MVP)
✅ **Authentication System**
- Email link authentication (passwordless login)
- JWT token management with refresh
- Session middleware and user context
- Login token generation and verification

✅ **Database Architecture**
- PostgreSQL database with automated migrations
- User, session, and login token models
- Repository pattern for data access
- Dependency injection container setup

✅ **API Foundation**
- Chi router with CORS configuration
- Modular route organization
- Authentication middleware
- Error handling and logging

✅ **Frontend Foundation**
- React 18 with TypeScript and TanStack Router
- Authentication hooks and protected routes
- Responsive UI with TailwindCSS and Shadcn/ui
- API client with authentication headers

✅ **Reporting System**
- Report model and repository with full CRUD operations
- Complete report submission endpoints
- Frontend report submission forms and management
- Report listing and filtering functionality

### In Development
🔄 **Business Directory**
- Business model and repository (structure exists)
- Basic business controller and service
- Frontend business directory page

### Planned (Phase 2)
📋 **Advanced Features**
- Confidence scoring algorithm
- Business rating system (Green/Yellow/Red)
- Advanced search and filtering
- Map integration
- Enhanced user verification tiers

📊 **Data Visualization & Analytics**
- LMIA job application trends charts (monthly/yearly)
- Regional employment pattern visualizations
- Industry-specific TFW usage charts
- Processing time trend analysis
- Wage subsidy program effectiveness metrics
- Interactive dashboards for policy impact analysis

**Technical Requirements for Charts:**
- Use Recharts library (already installed) for consistent styling
- Implement responsive charts that work on mobile
- Add data export functionality (CSV/PDF)
- Include filtering and time range selection
- Ensure accessibility with proper labels and alt text
- Consider server-side caching for large datasets

Refer to `docs/BUSINESS_DEMOGRAPHICS_FEATURE_SPEC.md` for complete feature specifications and implementation roadmap.