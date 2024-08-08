# KVLM (Key-Value List with Message) Format Documentation

## Introduction

KVLM is a simple yet powerful format used to store structured data, particularly in version control systems like Git. It's designed to handle key-value pairs along with a special multi-line message section.

## Basic Structure

A KVLM document consists of two main parts:
1. Key-Value Pairs
2. Message

### 1. Key-Value Pairs

- Each pair is on a separate line
- Format: `key value`
- The key and value are separated by a single space
- Keys cannot contain spaces
- Values can contain spaces and span multiple lines

### 2. Message

- Comes after all key-value pairs
- Separated from key-value pairs by a blank line
- Can span multiple lines
- Does not have a key

## Special Features

### Multi-line Values

- If a value spans multiple lines, each continuation line starts with a space
- This distinguishes between a new key-value pair and a continuation of the previous value

### Message Handling

- The message is treated as a special case
- It's typically associated with an empty string key (`""`) or `null` key in implementations

## Example

Here's an example of a KVLM document (similar to a Git commit object):

```
tree 29ff16c9c14e2652b22f8b78bb08a5a07930c147
parent 206941306e8a8af65b66eaaaea388a7ae24d49a0
author John Doe <john@example.com> 1527025023 +0200
committer Jane Doe <jane@example.com> 1527025044 +0200
gpgsig -----BEGIN PGP SIGNATURE-----
 
 iQIzBAABCAAdFiEExwXquOM8bWb4Q2zVGxM2FxoLkGQFAlsEjZQACgkQGxM2FxoL
 kGQdcBAAqPP+ln4nGDd2gETXjvOpOxLzIMEw4A9gU6CzWzm+oB8mEIKyaH0UFIPh
 -----END PGP SIGNATURE-----

This is the commit message.
It can span multiple lines.
The blank line above separates it from the key-value pairs.
```

## Parsing KVLM

When parsing KVLM:
1. Read line by line
2. For each line:
   - If it starts with a space, it's a continuation of the previous value
   - If it's a blank line, the rest is the message
   - Otherwise, it's a new key-value pair
3. For multi-line values, remove the leading space from continuation lines

## Serializing KVLM

When writing KVLM:
1. For each key-value pair:
   - Write the key, a space, then the value
   - If the value has multiple lines, add a space at the start of each line (except the first)
2. Add a blank line
3. Write the message (if any)

## Use Cases

KVLM is particularly useful for:
- Version control systems (like Git) to store commit information
- Configuration files that need a free-form message section
- Any scenario requiring structured data with a mix of single-line and multi-line values

## Benefits

- Simple to parse and generate
- Human-readable
- Flexible (handles both single-line and multi-line values)
- Separates structured data from free-form text

## Limitations

- Keys cannot contain spaces
- No nested structures (unlike JSON or YAML)
- No direct support for data types (everything is text)

## Conclusion

KVLM is a straightforward yet versatile format. Its simplicity makes it easy to implement and use, while its design allows it to handle complex data structures like Git commit objects efficiently.