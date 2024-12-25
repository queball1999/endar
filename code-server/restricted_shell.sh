#!/bin/bash

# Allowed commands
allowed_commands=("cd" "ls" "git" "mkdir" "rm" "rmdir")

# Base directory
base_dir="/home/coder/scripts"

# Function to check if the current directory is within the base directory
check_directory() {
  if [[ $(pwd) != $base_dir* ]]; then
    cd $base_dir
    echo "Navigation outside $base_dir is not allowed."
  fi
}

# Function to handle mkdir command
handle_mkdir() {
  local args="$1"
  local target_dir
  if [[ $args == /* ]]; then
    target_dir=$(realpath -m "$args")
  else
    target_dir=$(realpath -m "$PWD/$args")
  fi
  echo "Target directory: $target_dir"
  if [[ $target_dir == $base_dir* ]]; then
    mkdir "$args"
  else
    echo "mkdir is only allowed within $base_dir"
  fi
}

# Function to handle rm and rmdir commands
handle_rm_rmdir() {
  local cmd="$1"
  local args="$2"
  local target_dir
  if [[ $args == /* ]]; then
    target_dir=$(realpath -m "$args")
  else
    target_dir=$(realpath -m "$PWD/$args")
  fi
  echo "Target directory: $target_dir"
  if [[ $target_dir == $base_dir* ]]; then
    read -p "Are you sure you want to delete $args? (y/n): " confirm
    if [[ $confirm == "y" ]]; then
      $cmd "$args"
    else
      echo "Deletion cancelled."
    fi
  else
    echo "$cmd is only allowed within $base_dir"
  fi
}

# Start in the base directory
cd $base_dir

# Enable command history and readline features
HISTFILE=~/.restricted_shell_history
HISTSIZE=1000
SAVEHIST=1000
set -o emacs

while true; do
  # Use 'readline' to read the command line with history and editing support
  read -e -p "$USER@restricted-shell: $(pwd) $ " input

  # Split the input into command and arguments
  cmd=$(echo $input | awk '{print $1}')
  args=$(echo $input | awk '{for (i=2; i<=NF; i++) printf $i " "; print ""}')

  # Allow empty commands
  if [[ -z $cmd ]]; then
    continue
  fi

  # Check if the command is allowed
  if [[ ! " ${allowed_commands[@]} " =~ " ${cmd} " ]]; then
    echo "Command not allowed: $cmd"
    continue
  fi

  # Handle specific commands
  case $cmd in
    mkdir)
      handle_mkdir "$args"
      ;;
    rm|rmdir)
      handle_rm_rmdir "$cmd" "$args"
      ;;
    *)
      $cmd $args
      ;;
  esac

  check_directory
done