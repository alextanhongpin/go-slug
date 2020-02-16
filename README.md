# go-slug
Slugify your text into readable string.

## Slugifying text

There are several concerns when dealing with slugification

- different languages, how to slugify chinese characters/accented characters?
- what is the best length for slug urls?
- how to deal with changing slug? (whenever we update the name of the resource)
- how to deal with unique slug (once used, cannot be reused)
- how to deal with profanity in urls?

## Dealing with broken slug

There are several reasons why a url broke - someone decided to rename the resource you had, and it generates a new slug, so when a user attempts to visit the older ones, they get a 404.

So how do you deal with changing urls?

- ensure that the slug remains the same whenever the original name changes. Note that in some cases, the requirement requires the slug to change, hence this is not always a viable option.
- append a unique id behind the slug, something like slug-<unique-id-with-fix-characters>, then when user pass in a slug, just extract the last n characters to query. 
  - pros: don't need additional storage
  - pros: SEO doesn't care about the additional id behind, as long as the readable parts still maintain the context
  - pros: Notion uses this technique. This is especially useful when the name and slug changes very often.
  - cons: requires handling to extract the id from the url
- create a history of all the slugs that are created on every update, delete etc.
  - similar like `friendly_id` gem for rails
  - cons: takes up a lot of storage
  - cons: require additional query, and also additional steps to insert this into the table
