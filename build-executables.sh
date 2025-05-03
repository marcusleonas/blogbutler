#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_name="blogbutler"
dist_dir="dist"

rm -rf "$dist_dir"
mkdir "$dist_dir"

platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

echo "Building executables..."

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    output_path="$dist_dir/$output_name"

    env GOOS=$GOOS GOARCH=$GOARCH go build -o "$output_path" $package
    if [ $? -ne 0 ]; then
           echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi

    echo "Built $output_path"
done
