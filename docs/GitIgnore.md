# Comprehensive Documentation of .gitignore Rules

## Basic Patterns

1. **Specific Files**
    - Syntax: `filename.ext`
    - Example: `README.md`
    - Matches: Exact file name in any directory

2. **File Extensions**
    - Syntax: `*.extension`
    - Example: `*.log`
    - Matches: All files with the specified extension

3. **Directories**
    - Syntax: `directory_name/`
    - Example: `build/`
    - Matches: The entire directory and its contents

## Wildcards

4. **Single Asterisk (*)**
    - Syntax: `*`
    - Example: `*.txt`
    - Matches: Zero or more characters, except slash (/)
    - Note: Does not match across directory boundaries

5. **Double Asterisk (**)**
    - Syntax: `**/`
    - Example: `**/logs`
    - Matches: Zero or more directories
    - Note: Can be used to match across directory boundaries

6. **Question Mark (?)**
    - Syntax: `?`
    - Example: `file?.txt`
    - Matches: Exactly one character, except slash (/)

## Negation

7. **Negation (!)**
    - Syntax: `!pattern`
    - Example: `!important.txt`
    - Effect: Includes a file that would otherwise be ignored

## Character Classes

8. **Character Classes**
    - Syntax: `[characters]`
    - Example: `file[0-9].txt`
    - Matches: Any single character in the specified set

9. **Negated Character Classes**
    - Syntax: `[!characters]` or `[^characters]`
    - Example: `file[!a-z].txt`
    - Matches: Any single character not in the specified set

## Path Specifiers

10. **Leading Slash (/)**
    - Syntax: `/pattern`
    - Example: `/root.txt`
    - Matches: Patterns relative to the .gitignore file location

11. **Trailing Slash (/)**
    - Syntax: `pattern/`
    - Example: `logs/`
    - Matches: Only directories (not files) with that name

12. **Middle Slash (/)**
    - Syntax: `dir/file`
    - Example: `docs/build/`
    - Matches: Specifies a path relative to the .gitignore location

## Escaping Special Characters

13. **Backslash (\)**
    - Syntax: `\character`
    - Example: `\#important.txt`
    - Effect: Treats the following character as literal

## Comments

14. **Comments (#)**
    - Syntax: `# Comment text`
    - Example: `# This is a comment`
    - Effect: Gitignore ignores everything after #, unless escaped

## Advanced Patterns

15. **Combining Patterns**
    - Example: `**/logs/*.log`
    - Matches: All .log files in any logs directory

16. **Ignoring Files Only in Root Directory**
    - Syntax: `/filename.ext`
    - Example: `/README.md`
    - Matches: Only the file in the root directory, not in subdirectories

17. **Ignoring All Files in a Directory**
    - Syntax: `directory/*`
    - Example: `temp/*`
    - Matches: All files in the specified directory, but not the directory itself

18. **Ignoring All Files of a Type in a Specific Directory**
    - Syntax: `directory/*.ext`
    - Example: `logs/*.log`
    - Matches: All files with the specified extension in the given directory

## Special Cases

19. **Empty Lines**
    - Effect: Ignored by Git

20. **Trailing Spaces**
    - Effect: Ignored unless escaped with a backslash

21. **Duplicates**
    - Effect: Last duplicate pattern takes precedence

## Precedence Rules

- Patterns defined in a .gitignore file in a lower-level directory take precedence over higher-level ones.
- Within a single .gitignore file, later rules override earlier ones.
- Negation rules always take precedence over ignore rules.

## Best Practices

- Use specific patterns over broad ones to avoid unintended ignores.
- Place global ignore patterns in the root .gitignore file.
- Use local .gitignore files in subdirectories for project-specific ignores.
- Regularly review and update your .gitignore files as your project evolves.

This documentation covers the core functionality of .gitignore files, providing a comprehensive reference for creating
effective ignore patterns in Git repositories.