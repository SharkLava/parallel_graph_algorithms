#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <directory> <output_markdown_file>"
    exit 1
fi

# Input directory and output markdown file
input_directory="$1"
output_file="$2"

# Check if the provided directory exists
if [ ! -d "$input_directory" ]; then
    echo "Directory $input_directory does not exist."
    exit 1
fi

# Create or clear the output markdown file
> "$output_file"

# Find all .go files in the directory recursively
find "$input_directory" -type f -name "*.go" | while read -r file; do
    echo "# $file" >> "$output_file"
    echo '```go' >> "$output_file"
    cat "$file" >> "$output_file"
    echo '```' >> "$output_file"
    echo >> "$output_file"  # Add a newline for readability
done

echo "Markdown file created: $output_file"

