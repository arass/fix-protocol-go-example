echo "Updating..."

git fetch
git reset --hard origin
git pull -X theirs

chmod 700 update.sh run.sh run-fix-cert.sh

echo "Done."