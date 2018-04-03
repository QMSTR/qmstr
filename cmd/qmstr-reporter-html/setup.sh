#!/usr/bin/env bash
#
# This script sets up the data and template directories for the QMSTR
# HTML reporter. It expects three arguments:
# * the target shared data,
# * the path to the qmstr repository, and
# * the path to the html-reporter-theme repository.
# If the theme repository is not specified, it will be checked out from Github.
#
# directory. It is usually located in /usr/share/qmstr. Other options
# are /usr/loca/share/qmstr or /opt/share/qmstr.
#
# Example: ./setup.sh /usr/share/qmstr

# echo $@ $#

if [ $# -lt "2" ]; then
    echo "Please specify the setup target directory (/usr/share/qmstr) and the path to the qmstr repository!"
    exit 1
fi

TARGET_DIR="$1"
echo "Setting up QMSTR HTML reporter in $TARGET_DIR..."
REPO_DIR="$2"

if [ ! -d $REPODIR ]; then
    echo "No qmstr repository found at $REPO_DIR!"
    exit 1
fi
echo "Using qmstr repository at $REPO_DIR."

mkdir -p $TARGET_DIR || {
    echo "Error creating target directory $TARGET_DIR. Do you have permission?"
    exit
}

MODULE_DIR=$TARGET_DIR/reporter-html
mkdir -p $MODULE_DIR || {
    echo "Unable to create module data directory $MODULE_DIR."
    exit 1
}
cd $MODULE_DIR
echo "HTML reporter module directory is at $MODULE_DIR."

# Set up the theme directory:
if [ $# -eq "3" ]; then
    THEME_REPO_DIR="$3"
    ln -s $THEME_REPO_DIR theme
    echo "Linking theme located in $THEME_REPO_DIR to theme/."
else
    git clone --quiet --depth 1 https://github.com/QMSTR/html-reporter-theme.git theme/
    echo "Created shallow clone of the theme repo in theme/."
fi

for DIR in skeleton templates; do
    cp -Rfp $REPO_DIR/pkg/reporter/htmlreporter/share/$DIR . || {
	echo "Error copying the $DIR directory into the module shared data directory."
	exit 1
    }
done

echo "HTML reporter shared data directory set up at $MODULE_DIR."
