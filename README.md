# sqlchk

`sqlchk` is a Go static analyzer that validates raw SQL queries in your code by
running them on your local/test database. It basically does so by just 
preparing them.

> [!WARNING]
> It's a demonstration project and is not production ready AT ALL. If you like the idea and want to support the project, contributions are welcomed.

## Installation

```
go install github.com/thedenisnikulin/sqlchk/cmd/sqlchk@latest
```

## Example usage

Suppose we have a database table:
```
               table "public.products"
 Name  |          Type          | Nullable | Default
-------+------------------------+----------+---------
 id    | integer                | "YES"    |
 name  | character varying(100) | "YES"    |
 price | integer                | "YES"    |
```

And application code:
```go
// main.go
package main

func main() {
	query := `--sql
SELECT some_nonexistent_field FROM products`
	
    // run query here...
}
```

We can run `sqlchk` on it to check if all queries in the code are valid:
```bash
DATABASE_URL=postgres://user:pwd@host:port/db sqlchk main.go
```

And we get errors:
```
.../main.go:4:16: pq: column "some_nonexistent_field" does not exist
```

## Credits

This project is inspired by this Rust crate https://github.com/launchbadge/sqlx
