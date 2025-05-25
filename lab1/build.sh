#!/bin/sh
set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 source_file"
    exit 1
fi

source_file="$1"

if [ ! -f "$source_file" ]; then
    echo "Source file does not exist"
    exit 1
fi

output_line=$(grep -E '&Output:\s*(\S+)' "$source_file" | head -n 1)
if [ -z "$output_line" ]; then
    echo "Output comment not found"
    exit 2
fi

output_file=$(echo "$output_line" | sed -E 's/.*&Output:\s*(\S+).*/\1/')
if [ -z "$output_file" ]; then
    echo "Invalid output filename"
    exit 2
fi

case "$source_file" in
    *.c)
        compiler="gcc"
        ;;
    *.cpp|*.cc)
        compiler="g++"
        ;;
    *.tex)
        compiler="pdflatex"
        case "$output_file" in
            *.pdf) ;;
            *) output_file="${output_file}.pdf" ;;
        esac
        ;;
    *)
        echo "Unsupported file type"
        exit 1
        ;;
esac

temp_dir=$(mktemp -d)

trap 'rm -rf "$temp_dir"' EXIT

original_dir=$(pwd)

cp "$source_file" "$temp_dir/"

cd "$temp_dir"

if [ "$compiler" = "pdflatex" ]; then
    jobname=$(basename "$output_file" .pdf)
    pdflatex -jobname="$jobname" "$(basename "$source_file")" > /dev/null 2>&1
    pdflatex -jobname="$jobname" "$(basename "$source_file")" > /dev/null 2>&1
    if [ ! -f "$output_file" ]; then
        echo "Compilation failed"
        exit 3
    fi
else
    $compiler -o "$output_file" "$(basename "$source_file")" > /dev/null 2>&1
fi

mv "$output_file" "$original_dir/"

exit 0