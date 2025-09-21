# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **snap-share**, a wedding photo sharing web application that allows guests to upload and share photos in real-time using QR codes. The architecture follows a clear separation with:

- **Frontend**: Next.js (TypeScript) - deployed to Vercel
- **Backend**: Go API server using Supabase SDK - deployed to Cloud Run
- **Database & Storage**: Supabase (Postgres + Storage + Realtime)

## Architecture

### Directory Structure
- `frontend/` - Next.js application for guest and admin interfaces
- `backend/` - Go API server handling photo uploads and event management

### Core Data Models
- **events**: Stores wedding events with unique codes for QR access
- **photos**: Stores photo metadata with uploader names and Supabase storage URLs

### Key Workflows
1. Guests scan QR code → access `/e/{event_code}`
2. Guest enters name → session created → photo upload interface
3. Photos uploaded to Supabase Storage with metadata saved to Postgres
4. Real-time photo display using Supabase Realtime
5. Admin interface for bulk download and photo management

## Development Commands

### Frontend (Next.js)
Since the frontend directory is currently empty, initialize with:
```bash
cd frontend/
npx create-next-app@latest . --typescript --tailwind --eslint --app --src-dir --import-alias "@/*"
npm run dev        # Start development server
npm run build      # Build for production
npm run lint       # Run ESLint
npm run type-check # Run TypeScript checks (if configured)
```

### Backend (Go)
Since the backend directory is currently empty, initialize with:
```bash
cd backend/
go mod init snap-share-backend
go run main.go     # Start development server
go build          # Build binary
go test ./...     # Run tests
```

## Development Setup

The codebase structure shows empty frontend/ and backend/ directories that need initialization:

### Frontend Setup
- Use Next.js 14+ with App Router
- Configure Tailwind CSS for mobile-first responsive design
- Install Supabase client: `npm install @supabase/supabase-js`
- Set up TypeScript strict mode for better type safety

### Backend Setup
- Initialize Go module in backend/ directory
- Install Supabase Go SDK: `go get github.com/supabase/supabase-go`
- Use environment variables for Supabase credentials
- Implement Chi router or Gin for REST API endpoints

## Key Requirements
- Mobile-first responsive UI (primary users are mobile wedding guests)
- Real-time photo updates using Supabase Realtime
- QR code-based access with event codes
- Guest name input (no pre-registration required)
- Bulk photo download for event organizers
- Support for ~1000 photos per event

# AI Dev Tasks
Use these files when I request structured feature development using PRDs:
/.claude/rules/create-prd.md
/.claude/rules/generate-tasks.md
/.claude/rules/process-task-list.md