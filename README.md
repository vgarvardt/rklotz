# This branch holds the version of this project

Please dont delete

How to create a version branch

```bash
git checkout --orphan version
git rm --cached -r .
rm -rf *
rm .gitignore .gitmodules
touch README.md
echo "0.0.0" > version
git add .
git commit -m "version branch"
git push origin version
```

