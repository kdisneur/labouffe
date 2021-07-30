# CHANGELOG

All notable changes to this project will be documented in this file

## Upcoming changes

n/a

## 2021-07

- [BUGFIX] Update recipe list page to sort recipes based on their name instead of their slugs.

## 2021-05

- [FEATURE] Add support for `reheat` on the recipe. It has also been added in the list of available filters.
  ```yaml
  # file:recipes/carbonade.yaml
  title: Carbonade
  # ...
  reheat: true
  instructions:
    # - ...
  ```

## 2021-04

- [ENHANCEMENT] Rename the `category` cake to dessert to be more flexible
- [ENHANCEMENT] Improve the recipe ingredient to make the sub-fields optional
  ```yaml
  # file:recipes/shepherd-pie.yaml
  title: Shepherd's pie
  # ...
  ingredients:
    # BEFORE
    salt:
      quantity: 0 # to make it pass the compilation but no
                  # value because it depends of the taste
    # NOW
    salt:
  ```
- [ENHANCEMENT] Add support for a new `category`: sauce
- [FEATURE] Add support for `warning` on the recipe, it supports multi-line. It can be used to have reminders the quantity wasn't right and try something else next time.
  ```yaml
  # file:recipes/shepherd-pie.yaml
  title: Shepherd's pie
  # ...
  warning: Something you should be aware of!
  instructions:
    # - ...
  ```
- [FEATURE] Add the missing `resting` time to the recipe. It was often added to the `preparation` or burried inside the instructions which wasn't ideal
  ```yaml
  # file:recipes/french-fries.yaml
  title: French fries
    cooking: 9m
    resting: 10m
    preparation: 20m
  ```

## 2020-11

- [ENHANCEMENT] Ensure that only accepted fields are present in the YAML to avoid typos
- [ENHANCEMENT] Add the ability to filter per duration range "< 30m", "> 30m et < 1h", "> 1h et < 1h30", "> 1h30" based on the total recipe duration
- [FEATURE] Add the ability to link a recipe instruction to another recipe. The parser validates the linked recipe exists.
  ```yaml
  # file:recipes/shepherd-pie.yaml
  title: Shepherd's pie
  # ...
  instructions:
    - instruction: Prepare the mashed potatoes
      recipe: mashed-potatoes
    - Slice the onions
  ```
- [FEATURE] Add the ability to suggest alternative ingredients in a recipe. The parser validates the ingredient exists.
  ```yaml
  # file:recipes/helenettes
  title: Helenettes
  # ...
  ingredients:
    almonds:
      quantity: 100g
      details: powder
      alternatives:
        - hazelnut
  ```
- [FEATURE] Add the ability to link an ingredient to a specific recipe. The parser validates the recipe exists.
  ```yaml
  # file:ingredients.yaml
  mashed-potatoes:
    title: Mashed Potatoes
    recipe: mashed-potatoes
  ```
- [ENHANCEMENT] Upgrade the UI to be more responsive friendly (and also hide the filtering by default)
- [ENHANCEMENT] Change recipe ingredient quantity to become float (from integer):
  ```yaml
  # file:recipes/french-fries.yaml
  title: French fries
   # ...
  ingredients:
    beef-fat:
      quantity: 1.2kg
  ```
- [Feature] Add an optional `details` field on the recipe's ingredient
  ```yaml
  # file:recipes/chocolate-cookies.yaml
  title: Cookies with chocolate chunks
  # ...
  ingredients:
    chocolate:
      quantity: 300g
      details: chunks
  ```
- [BUGFIX] Fix an issue with the duration where leading zero where added when not necessary (i.e. `5h0` instead of `5h`)
- [ENHANCEMENT] Replace the `guests` field on the recipes to `servings`. It can be one of:
  ```yaml
  # file:recipes/spinach-pasta.yaml
  title: Spinach Pasta
  # ...
  # when we count per number of guests
  servings:
    quantity: 3
    type: personnes
  # when we count per number of items like "5 biscuits"
  servings:
    quantity: 5
    type: unités
  ```
- [FEATURE] Add the ability to filter list on the recipes list page based on category, price, difficulty,
- [ENHANCEMENT] Rename the recipe format to use french instead of english
   - `difficulty` changed from easy, average, hard to facile, moyen, difficile
   - `price` changed from cheap, abordable, expensive to économique, abordable and cher
- [FEATURE] Generate all recipe page using the pattern `/recipes/<recipe-file-name>`
- [FEATURE] Generate recipes list page at `/recipes`
- [FEATURE] Add the `category` to the recipe format. It can be one of: biscuit, gâteau, plat froid, plat chaud, entrée
- [FEATURE] Generate ingredient list page at `/ingredients`
- [FEATURE] Automtatically load the `data` folder when running the binary:
   - The recipe folder can be overriden with the `-recipes` flag
   - The ingredient file can be overriden with the `-ingredients` flag
- [FEATURE] Add the recipe and ingredient parser. They validate the format of the field but also the data itself (ingredients referenced exist, etc...)
  Ingredient file format:
  ```yaml
  # file:ingredients.yaml
  # shouldn't contain any spaces only a-z and dashes
  tomatoes:
    title: Tomatoes
  spinach:
    title: Spinach
  ```
  Recipe file format:
  ```yaml
  # file:recipes/spinach-pasta.yaml
  # the name of the file is used to generate the URL
  title: Spinach Pasta
  cooking: 20m
  preparation: 10m
  difficulty: easy
  pricing: cheap
  guests: 1
  ingredients:
    pasta:
      quantity: 50g # can be one of `g`, `kg`, `ml`, `cl`, `l`, `cc`, `cs` or no suffix.
    parmesan:
      quantity: 10g
    spinach:
      quantity: 150g
  instructions:
    - Cook pasta in boiling water
    - A bit before the end of the pasta cooking add the spinach to the water
    - ....
  ```
