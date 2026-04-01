#!/usr/bin/env bash
# Exit on error
set -o errexit

echo "Building frontend..."
cd frontend
npm install
npm run build
cd ..

echo "Building backend..."
go build -o edumentor-server .

echo "Build complete."
