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


## Scenarios 

- user visit old url - they should be redirected to the latest url
- item has been renamed from a to b to c, visiting a or b will redirect to c
- if calling the api from the frontend, the redirection will be handled on the client side, and a shallow routing (replacing url parts) or redirection can be performed
- item 1 has be renamed from a to b, and the back to a. Visiting a will go to a, visiting b will redirect the user back to a.
- item 1 has been renamed from a to b, and then item 2 has been renamed from c to a. Visiting a  will redirect to a, and visiting c will redirect to a, and visiting b will redirect to b. Basically, the named can be taken over by another item as long as it is no longer in used.


## Code

Pseudo code with golang to show how the api could work. 

TODO:
- create migration file
- create basic sluggable service
- user can just add the migration to their file, and use the service with their database provider

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, playground")
}

type Slug struct {
	ID            string
	Slug          string
	SluggableType string // UID
	SluggableID   string // UID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SlugService interface {
	Find(slug, sluggableType string) (*Slug, error)
	// This behaves like upsert, it will create if does not exists or update existing ones.
	Create(slug, sluggableID, sluggableType string) (*Slug, error)
	Delete(sluggableID, sluggableType string)
	WithTx(Tx) SlugService
}

type Comment struct {
	Title       string
	Description string
}

const CommentSlug = "comment"

// Find the comments that has been renamed...
func getComment(ctx context.Context, slug string) (*Comment, error) {
	// If the slug is not the latest, the comment table could return the latest slugs to reduce load.
	// Then, the frontend can check if the requested slug matches the returned slug - if not
	// a shallow routing can be performed on the front end to reduce load.
	slug, err := slugsvc.Find(slug, CommentSlug)
	if err != nil {
		return nil, err
	}
	comment, err := commentsvc.Find(slug.SluggableID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func createComment(ctx context.Context, title, description string) (*Comment, error) {
	var slug *Slug
	err := withTransaction(db, func(tx Tx) error {
		commentTx, err := commentsvc.WithTx(tx)
		comment, err := commentTx.Create(title, description)
		if err != nil {
			return err
		}

		slugTx, err := slugsvc.WithTX(tx)
		if err != nil {
			return err
		}

		// This will fail if the slug already exists. We can provide an implementation
		// that will append a unique id if the slug exists.
		var err error
		slug, err = slugTx.Create(slugify(title), comment.ID, CommentSlug)
		// TODO: Check err unique - retry by appending a unique id behind.
		return err
	})
	return slug, err
}

func updateComment(ctx context.Context, commentID, title, description string) (*Comment, error) {
	var slug *Slug
	err := withTransaction(db, func(tx Tx) error {
		commentTx, err := commentsvc.WithTx(tx)
		comment, err := commentTx.Update(commentID, title, description)
		if err != nil {
			return err
		}

		slugTx, err := slugsvc.WithTX(tx)
		if err != nil {
			return err
		}

		// This will create or update existing slug's timestamp.
		var err error
		slug, err = slugTx.Create(slugify(title), comment.ID, CommentSlug)
		return err
	})
	return slug, err
}

func deleteComment(ctx context.Context, commentID string) error {
	return withTransaction(db, func(tx Tx) error {
		commentTx, err := commentsvc.WithTx(tx)
		if err := commentTx.Delete(commentID); err != nil {
			return err
		}

		slugTx, err := slugsvc.WithTX(tx)
		if err != nil {
			return err
		}

		return slugTx.Delete(commentID, CommentSlug)
	})
}
```
