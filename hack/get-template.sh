#!/usr/bin/env sh

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

versions='0.2.22.Final'
tmp_dir='/tmp/apicurio/templates'

while getopts "v:d:" arg; do
    case $arg in
    v)
        versions=$OPTARG
        ;;
    d)
        tmp_dir=$OPTARG
        ;;
  esac
done

for version in $(echo $versions | tr "," "\n")
do
    echo "[INFO] Starting process to get required files for release $version"
    mkdir -p $tmp_dir

    echo "[INFO] Downloading release $version"
    url="https://codeload.github.com/Apicurio/apicurio-studio/zip/v$version"
    curl $url -s -o $tmp_dir/$version.zip

    echo "[INFO] Unzipping release $version.zip"
    unzip -qq  "$tmp_dir/$version.zip" -d $tmp_dir

    echo "[INFO] Preparing resources dir for version $version"
    tmpl_dir=$tmp_dir/apicurio-studio-$version/distro/openshift
    mkdir -p $DIR/../res/$version

    echo "[INFO] Copying resource files from release $version"
    cp $tmpl_dir/apicurio-standalone-template.yml $DIR/../res/$version/template.yml
    cp $tmpl_dir/apicurio-auth-template.yml $DIR/../res/$version/auth-template.yml
    cp $tmpl_dir/apicurio-postgres-template.yml $DIR/../res/$version/postgres-template.yml
    cp $tmpl_dir/auth/realm.json $DIR/../res/$version/realm.json

    echo "[INFO] Applying template api group fix for version $version"
    for fname in $DIR/../res/$version/*.yml; do
        sed -i -f $DIR/sed-template.txt $fname
    done

    echo "[INFO] Cleaning up temporary files for release $version"
    rm -rf "$tmp_dir/apicurio-studio-$version"

    echo "[INFO] Checking copied files for release $version"
    for fname in 'realm.json' 'template.yml'
    do
        if [[ ! -e "$DIR/../res/$version/$fname" ]];
        then
            echo -e "${RED}==> Failed to check release file: $version/$fname${NC}"
            exit 1
        fi
    done

    echo -e "${GREEN}[SUCCESS] Finished retrieving and patching apicurio for version $version${NC}"
    echo
done
