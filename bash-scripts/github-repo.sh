#!/bin/bash

# list all repositories for the authenticated user or organization
# gh repo list --json name,createdAt --limit 1000

# GitHub Bulk Repository Deletion Script using 'gh' CLI
# Usage: ./bulk_delete_gh_repos.sh

# Set your GitHub username or organization name
OWNER="MKS-01"

# List of repositories to delete (space-separated, no quotes)
REPOS=(
  # Add more repositories as needed
)

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI ('gh') is not installed. Please install it from https://cli.github.com/"
    exit 1
fi

# Check if user is authenticated
if ! gh auth status &> /dev/null; then
    echo "Error: Not authenticated with GitHub CLI. Run 'gh auth login' first."
    exit 1
fi

# Main script
for repo in "${REPOS[@]}"; do
  read -p "Are you sure you want to delete $repo? (y/n) " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Deleting repository: $OWNER/$repo"
    if gh repo delete "$OWNER/$repo" --confirm; then
      echo "Successfully deleted $repo"
    else
      echo "Failed to delete $repo"
    fi
  else
    echo "Skipping $repo"
  fi
done

echo "Done."
