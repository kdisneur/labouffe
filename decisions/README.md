# Decisions

## Ingredient code plural / singular?

We've decided ingredients should be created using the singular form, always.
It means, even if we never consumes only one, we still use the singular form.

```diff
- olives
+ olive

- mushrooms
+ mushroom
```

## Ingredient specialization

When we need to create several form of the same ingredient (powder, seed, juice,...), we add the specialization as a prefix to keep the same ingredient's base together when alphabetically sorting the ingredient list.

```diff
- graine-cumin
+ cumin-graine

- poudre-cumin
+ cumin-poudre

- jus-citron
+ citron-jus
```
