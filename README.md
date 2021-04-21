# Labouffe

This repository contains my **personal** cooking website.
The repository is public because it's easier but I won't accept any pull-requests (at least on the `data/` folder).

Every notable changes are documented in the [CHANGELOG](./CHANGELOG.md) file.

## Local development

1. Get the repository on your local computer
   ```
   git clone git@github.com:kdisneur/labouffe
   ```
2. Go to the repository
   ```
   cd ./labouffe
   ```
3. Make sure the local version of your repository is up-to-date
   ```
   git pull
   ```
4. Start the local server
   ```
   make live-reload
   ```
5. Visit the [local version of labouffe](http://127.0.0.1:8080)
6. Add a new recipe in `data/recipes`:
   - the name of the file should match this pattern `[a-z]+[a-z0-9-]*[a-z0-9]+.yaml`.
     It will be used to generate the URL of the recipe, so it should be chosen wisely.
   - copy/paste an existing recipe to get an example of the template (more documentation should be added here)
   - edit the recipe to add your own steps and ingredients.
     read **carefully** the log messages from `make live-reload` as it will indicate precise errors when an invalid value is entered.

     If an ingredient is missing add it to `data/ingredients.yaml`
7. After you're sure the recipe is working as intended (meaning you've read it on [your local version](http://127.0.0.1:8080)) and there are no errors in the `make live-reload` output, add your changes.
    ```
    git add ./data
    git commit -m "Add recipe: <name of the recipe>"
    git push
    ```
