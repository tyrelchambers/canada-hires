# Gemini Project Configuration

This file helps Gemini understand the project's context, conventions, and commands.

## Project Overview

*   **Description:** A web application for Canadian business demographics, with a Go backend and a React frontend.
*   **Technology Stack:**
    *   **Frontend:** React, TypeScript, Vite, TanStack Router, Tailwind CSS (based on `components.json` and typical shadcn/ui usage)
    *   **Backend:** Go
    *   **Database:** SQL

## Development Environment

*   **Setup:** The project is divided into two main parts: `api` (backend) and `web` (frontend).
*   **Key Files:**
    *   `api/main.go`: Backend entry point.
    *   `web/src/main.tsx`: Frontend entry point.
    *   `web/vite.config.ts`: Frontend build configuration.
    *   `docker-compose.yml`: (If present) for running services together.

## Commands

A list of common commands for building, testing, and running the project.

*   **Install Dependencies:**
    *   `api`: `cd api && go mod tidy`
    *   `web`: `cd web && npm install`
*   **Run Development Servers:**
    *   `api`: `cd api && make run`
    *   `web`: `cd web && npm run dev`
*   **Run Tests:**
    *   `api`: `cd api && make test`
    *   `web`: `cd web && npm run test`
*   **Lint/Format:**
    *   `api`: `cd api && golangci-lint run ./...`
    *   `web`: `cd web && npm run lint`

## Coding Style & Conventions

*   **Formatting:** Follow existing project formatting. Go standard format for the backend and Prettier/ESLint for the frontend.
*   **Naming Conventions:** Use camelCase for variables and functions in TypeScript, and PascalCase for React components. Follow Go conventions for the backend.
*   **Component Structure:** Follow the existing structure for React components in `web/src/components`.

## User Preferences

*(This section can be used to store user-specific preferences for how Gemini should behave in this project.)*
