#!/bin/bash

#
# This script migrates Opera and Carmen to its public version.
#
# It removes all experimental, alternative and unfinished features
# and keep only to be disclosed ones.
#
# In particular, this version exports GoLang implementation of file-based Index/Store StateDB
# with LevelDB Archive database.
#
# This script checkouts fresh Opera, Tosca and Carmen repositories, filters out unnecessary parts
# combines them in one deliverable, and pushes it into a new repository.
#
# Input and output repositories and directories as well as other configurations are customisable
# via constants at the beginning of this file.
#
# The script uses 'Git', and 'Git filter-repo' commands, which must
# be installed before running this script.
#
#
# Carmen related directories and files to be included in the output version
# are configured in the file 'scripts/export/filter.txt', from the Carmen repository.
#
# Carmen and Tosca are switched from sub-modules to integral
# parts of the resulting repository.
#

#
# Configurable variables follow:
#

#
# Temporary folder to checkout into, and use it as workspace for
# modifications.
#
REPO_DIR=${TMPDIR-/tmp}/_opera_carmen_temp

#
# Source repository with Opera.
#
SOURCE_REPO=git@github.com:Fantom-foundation/go-opera-norma.git

#
# Source repository with Carmen.
#
CARMEN_SOURCE_REPO=git@github.com:Fantom-foundation/Carmen.git

#
# Source repository with Tosca.
#
TOSCA_SOURCE_REPO=git@github.com:Fantom-foundation/Tosca.git

#
# Destination repository, where modified Carmen will be stored into.
#
DEST_REPO=git@github.com:Fantom-foundation/go-opera-norma-public-test.git

#
# Git branch name to checkout specific version of Carmen.
#
CARMEN_SOURCE_BRANCH="kjezek/migration-scripts"

#
# Git branch name to checkout specific version of Tosca.
#
TOSCA_SOURCE_BRANCH="main"

#
# Git branch name to checkout specific version of Opera.
#
OPERA_SOURCE_BRANCH="kjezek/migration-script"


#
# Git branch where the modified Opera will be pushed into.
#
DEST_BRANCH="develop"

#
# Execution starts from here
#

ORIGINAL_DIR=$(pwd)

#
# Apply updates to Carmen.
# Checkout to a new local repository and filter out unnecessary parts
#
carmen_temp="$REPO_DIR-CARMEN"
mkdir $carmen_temp || exit
git clone $CARMEN_SOURCE_REPO $carmen_temp
cd $carmen_temp || exit
git checkout $CARMEN_SOURCE_BRANCH
git remote rm origin

#
# Filter out unnecessary parts of Carmen in a detached branch
#
cp scripts/export/filter.txt filter.txt
git filter-repo --force --path go/ --path-rename go/:
git filter-repo --force --paths-from-file filter.txt
git filter-repo --force --path-rename :carmen-embedded/

#
# Carmen ready in the local repo
#

#
# Apply updates to Tosca.
# Checkout to a new local repository and filter out unnecessary parts
#
tosca_temp="$REPO_DIR-TOSCA"
mkdir $tosca_temp || exit
git clone $TOSCA_SOURCE_REPO $tosca_temp
cd $tosca_temp || exit
git checkout $TOSCA_SOURCE_BRANCH
git remote rm origin

#
# Filter out unnecessary parts of Tosca in a detached branch
#
git filter-repo --force --path go/ --path go.mod --path go.sum
git filter-repo --force --path-rename :tosca-embedded/

#
# Tosca ready in the local repo
#


#
# Go to treat Opera/Norma
#
#
# Clone the Opera/Norma repo to a new directory
#
mkdir "$REPO_DIR" || exit
git clone $SOURCE_REPO "$REPO_DIR"

cd "$REPO_DIR" || exit
git checkout $OPERA_SOURCE_BRANCH
git remote rm origin

#
# Merge Carmen into this repository
#
git remote add carmen $carmen_temp
git fetch carmen

git merge -X ours --no-commit --allow-unrelated-histories "carmen/$CARMEN_SOURCE_BRANCH"
git rm --cached carmen
rm -rf carmen/.git
git commit -a -m "inlines carmen"
git remote rm carmen

#
# Update Opera dependencies
#
cd "$REPO_DIR"  || exit
sed -i -e 's/\.\/carmen\/go.*/.\/carmen-embedded/g' go.mod

##
## Remove C++ build from Makefile
##
sed -i -e 's/opera: carmen\/go\/lib\/libcarmen.so tosca/opera: tosca/g' Makefile

#
# Merge Tosca into this repository
#
git remote add tosca $tosca_temp
git fetch tosca

git merge -X ours --no-commit --allow-unrelated-histories "tosca/$TOSCA_SOURCE_BRANCH"
git rm --cached tosca
rm -rf tosca/.git
git commit -a -m "inlines tosca"
git remote rm tosca

#
# Update Opera dependencies
#
cd "$REPO_DIR"  || exit
sed -i -e 's/\.\/tosca\s*/.\/tosca-embedded/g' go.mod
sed -i -e 's/.*\/evmc\/.*//g' go.mod  # remove this line completely

go mod tidy -compat=1.20
git commit -a -m "publishes Opera"

#
# Push to the new repository, while adding required extra files.
#
git remote add origin $DEST_REPO
git commit -a -m "migrates to public repository at $d"
git push -f origin $OPERA_SOURCE_BRANCH:$DEST_BRANCH

#
# Clean-up
#

# test all can build
make

cd "$ORIGINAL_DIR" || exit

echo "Result stored in '$REPO_DIR'"
