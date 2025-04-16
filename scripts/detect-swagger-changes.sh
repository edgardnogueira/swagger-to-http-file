#!/bin/bash

# Script to detect changes in Swagger/OpenAPI files and update HTTP files
#
# This script:
# 1. Finds all Swagger/OpenAPI files in a directory (or specified files)
# 2. Checks if they've been modified since their corresponding HTTP files
# 3. If so, runs swagger-to-http-file to update the HTTP files

set -e

# Default settings
SEARCH_DIR="."
OUTPUT_DIR=""
RECURSIVE=true
VERBOSE=false
PATTERNS=("*.json" "*.yaml" "*.yml")
FILES=()

# Help message
show_help() {
  echo "Usage: $(basename "$0") [options] [files...]"
  echo ""
  echo "Detect changes in Swagger/OpenAPI files and update corresponding HTTP files."
  echo ""
  echo "Options:"
  echo "  -h, --help              Show this help message"
  echo "  -d, --directory DIR     Base directory to search for Swagger files (default: .)"
  echo "  -o, --output DIR        Output directory for HTTP files (default: same as input)"
  echo "  -r, --no-recursive      Don't search directories recursively"
  echo "  -v, --verbose           Show verbose output"
  echo ""
  echo "If specific files are provided, only those files will be processed."
  echo "Otherwise, all Swagger/OpenAPI files (*.json, *.yaml, *.yml) in the directory will be checked."
  echo ""
  echo "Example:"
  echo "  $(basename "$0") -d ./api -o ./http"
  echo "  $(basename "$0") ./api/swagger.json ./api/other-api.yaml"
  echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      show_help
      exit 0
      ;;
    -d|--directory)
      SEARCH_DIR="$2"
      shift 2
      ;;
    -o|--output)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    -r|--no-recursive)
      RECURSIVE=false
      shift
      ;;
    -v|--verbose)
      VERBOSE=true
      shift
      ;;
    -*)
      echo "Unknown option: $1"
      show_help
      exit 1
      ;;
    *)
      FILES+=("$1")
      shift
      ;;
  esac
done

# Check if swagger-to-http-file is installed
if ! command -v swagger-to-http-file &> /dev/null; then
  echo "Error: swagger-to-http-file not found in PATH."
  echo "To install: go install github.com/edgardnogueira/swagger-to-http-file/cmd/swagger-to-http-file@latest"
  exit 1
fi

# Function to find modified Swagger files
find_swagger_files() {
  local dir="$1"
  local recursive="$2"
  
  if [ "$recursive" = true ]; then
    find_opt="-name"
  else
    find_opt="-maxdepth 1 -name"
  fi
  
  for pattern in "${PATTERNS[@]}"; do
    if [ "$VERBOSE" = true ]; then
      echo "Searching for $pattern files in $dir..."
    fi
    
    # shellcheck disable=SC2086
    find "$dir" $find_opt "$pattern" -type f
  done
}

# Function to check if HTTP files need updating
check_and_update() {
  local swagger_file="$1"
  local output_dir="$2"
  
  # If no output directory specified, use the same directory as the Swagger file
  if [ -z "$output_dir" ]; then
    output_dir=$(dirname "$swagger_file")
  fi
  
  # Calculate HTTP file paths based on Swagger file
  local base_name=$(basename "$swagger_file" | sed 's/\.[^.]*$//')
  local tag_dir="$output_dir/$base_name"
  local http_file="$output_dir/$base_name.http"
  
  # Check if HTTP files exist and are newer than Swagger file
  local needs_update=false
  
  # Check single HTTP file
  if [ -f "$http_file" ]; then
    if [ "$swagger_file" -nt "$http_file" ]; then
      needs_update=true
      if [ "$VERBOSE" = true ]; then
        echo "$swagger_file is newer than $http_file"
      fi
    fi
  else
    needs_update=true
    if [ "$VERBOSE" = true ]; then
      echo "$http_file does not exist"
    fi
  fi
  
  # Check tag directory
  if [ -d "$tag_dir" ]; then
    local any_updated=false
    while IFS= read -r -d '' tag_file; do
      if [ "$swagger_file" -nt "$tag_file" ]; then
        any_updated=true
        if [ "$VERBOSE" = true ]; then
          echo "$swagger_file is newer than $tag_file"
        fi
      fi
    done < <(find "$tag_dir" -name "*.http" -type f -print0)
    
    if [ "$any_updated" = true ]; then
      needs_update=true
    fi
  elif [ -d "$output_dir" ]; then
    local tag_files=$(find "$output_dir" -name "*.http" -type f)
    if [ -z "$tag_files" ]; then
      needs_update=true
      if [ "$VERBOSE" = true ]; then
        echo "No HTTP files found in $output_dir"
      fi
    fi
  else
    needs_update=true
    if [ "$VERBOSE" = true ]; then
      echo "Output directory $output_dir does not exist"
    fi
  fi
  
  # Update HTTP files if needed
  if [ "$needs_update" = true ]; then
    echo "Updating HTTP files for $swagger_file..."
    swagger-to-http-file -i "$swagger_file" -o "$output_dir" -w
    if [ "$?" -eq 0 ]; then
      echo "✓ Successfully updated HTTP files for $swagger_file"
    else
      echo "✗ Failed to update HTTP files for $swagger_file"
      return 1
    fi
  else
    if [ "$VERBOSE" = true ]; then
      echo "HTTP files for $swagger_file are up to date"
    fi
  fi
  
  return 0
}

# Main execution
echo "Checking for Swagger/OpenAPI file changes..."

# Process specified files or find all Swagger files
if [ ${#FILES[@]} -gt 0 ]; then
  for file in "${FILES[@]}"; do
    if [ -f "$file" ]; then
      check_and_update "$file" "$OUTPUT_DIR"
    else
      echo "Warning: File not found: $file"
    fi
  done
else
  # Find and process all Swagger files
  while IFS= read -r file; do
    check_and_update "$file" "$OUTPUT_DIR"
  done < <(find_swagger_files "$SEARCH_DIR" "$RECURSIVE")
fi

echo "Done."
exit 0
