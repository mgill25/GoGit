# GoGit
An implementation of Git written in GoLang

## Motivation

To try and understand how Git works from the inside out by implementing some of its most commonly used features.

## ToDo

1. How does git work? What do I need to understand?
2. Implement all the foundational data structures and algorithms.
3. ggit init; ggit add; ggit commit -m "msg"


## Essential Links

1. [Git for computer scientists](http://eagain.net/articles/git-for-computer-scientists/)
2. [Git from the bottom up](ftp.newartisans.com/pub/git.from.bottom.up.pdf)
3. [Git Source Code](https://github.com/git/git)
4. [Dissecting Git's guts](https://www.youtube.com/watch?v=Y2Msq90ZknI)

## Broad concepts in Git

- Repository: Collection of Commits
- Commit: Snapshot of Working Tree at a given point in time
- Working Tree: Any directory on your filesystem with the Repository associated with it. Includes all sub-directories
- Index: Changes are first recorded from the Working Tree to the Index (Staging Area) before being committed.
- Branch: Just a name for a commit, also called a Reference
- HEAD: Used by the Repository to define what is currently checked out

## Workflow
- git init: Initializes the working tree
- git add: Adds all changes in the working tree to the index
- git commit: Changes are committed to the repository from the state of the index
- git checkout: Check out earlier stages of the working tree

## Guts of Git

*Blob* : SHA1 hash of a file being tracked by Git. Blob stores no metadata about its content, unlike a regular *nix file.
All the metadata info is kept in the tree that holds the blob. This helps to create multiple trees using same blob,
thus avoiding duplication of data.

- The .git makes a repo a git repo. If you ever want to back up/duplicate, just copy this. 
- In fact, git clone is basically just copying over .git
- "objects" directory - serves as git's "database". Can be described as a "content addressable file system". So that we can retrieve info
 based on the content.

### Plumbing commands

#### Basic object storage
Hash an object and write it out to the object db: `git hash-object filename`
```apacheconfig
> git hash-object -w hello_world.txt                                                                                          ◼ (git:master)
980a0d5f19a64b4b30a87d4206aade58726b60e3

# Hash is uniquely generated SHA1. Prepends some metadata (filetype and size) to the content. 
> printf "blob 13\000Hello World\041\n" | openssl sha1                                                                        ◼ (git:master)
980a0d5f19a64b4b30a87d4206aade58726b60e3

# First 2 characters are the sub-directory and the remaining 38 become the file name.
> find .git/objects -type f                                                                                                   ◼ (git:master)
.git/objects/98/0a0d5f19a64b4b30a87d4206aade58726b60e3 

# We can retreive the contents by using git cat-file
>  git cat-file -p 980a0d5f19a64b4b30a87d4206aade58726b60e3                                                                    ◼ (git:master)
Hello World!

# But the actual contents of the file are gibberish if we try to cat them! Since they are compressed using zlib! :)
> cat .git/objects/98/0a0d5f19a64b4b30a87d4206aade58726b60e3                                                                  ◼ (git:master)
xK��OR04f�H����/�IQ�Hak%

# Git uses the hash to detect when a file has changed and will be more selective about when to store new objects.
# If there's a file already existing and we copy it, it'll just spit back the same hash and won't duplicate the object.
```

### Versioning

Modify the contents of hello_world.txt and we see the following changes
```apacheconfig

# The hash has now changed because the contents have changed
> git hash-object -w hello_world.txt                                                                                          ◼ (git:master)
f907f535fae1314af17f142b0bf815cfaaf0050d

# Now we have another object in the object database
> find .git/objects -type f                                                                                                   ◼ (git:master)
.git/objects/98/0a0d5f19a64b4b30a87d4206aade58726b60e3
.git/objects/f9/07f535fae1314af17f142b0bf815cfaaf0050d

# NOTE: the new object has the content in its entirety, NOT the diff...initially.
> git cat-file -p f907f535fae1314af17f142b0bf815cfaaf0050d                                                                  ⏎ ◼ (git:master)
Hello World!
The world is beautiful

# What is this object called?
> git cat-file -t f907f535fae1314af17f142b0bf815cfaaf0050d                                                                  ⏎ ◼ (git:master)
blob
```

#### Tree objects
Blobs are only for single files. How do we group files? Using Tree objects!
A Tree object can be thought of as a complete snapshot of your project directory.
But first, we need an Index file! Lets stage some files aka put stuff into our index.

```apacheconfig
> git update-index --add hello_world.txt

# Create a new file and add it to index directly - hash + blob creation + happens under the hood.
> git update-index --add foobar.txt

# we now have an index file in .git. We can inspect the content with:
> git ls-files --stage
100644 f6a4b702215289bde4037bb76acf18d634628539 0	foobar.txt
100644 f907f535fae1314af17f142b0bf815cfaaf0050d 0	hello_world.txt

# Thus, the index --stage is just the running list of stuff in our staging area.
# We finally write that tree
> git write-tree
385c013d49b7be65e532fc9ccb718412310aebbf

# What is inside the tree object?
> git cat-file -p 385c013d49b7be65e532fc9ccb718412310aebbf                                                                    ✚ (git:master)
100644 blob f6a4b702215289bde4037bb76acf18d634628539	foobar.txt
100644 blob f907f535fae1314af17f142b0bf815cfaaf0050d	hello_world.txt

# Unlike the index file, tree objects are final. They can also refer other tree objects to support sub-directories
# Lets make sure that they are final. 
> find .git/objects -type f                                                                                                   ✚ (git:master)
.git/objects/38/5c013d49b7be65e532fc9ccb718412310aebbf
.git/objects/98/0a0d5f19a64b4b30a87d4206aade58726b60e3
.git/objects/f6/a4b702215289bde4037bb76acf18d634628539
.git/objects/f9/07f535fae1314af17f142b0bf815cfaaf0050d
```

#### Commit object

So we can save files and make snapshots, but we need metadata. Who saved these objects, when and why? Commit objects!
```apacheconfig
> echo 'first plumbing commit' | git commit-tree 385c013d49b7be65e532fc9ccb718412310aebbf                                     ✚ (git:master)
da9f4efa6d78917d721476c90386d89f23864338

# But what's inside a commit object?
> git cat-file -p da9f4efa6d78917d721476c90386d89f23864338                                                                    ✚ (git:master)
tree 385c013d49b7be65e532fc9ccb718412310aebbf
author Manish <manish.gill@wingify.com> 1489692140 +0530
committer Manish <manish.gill@wingify.com> 1489692140 +0530

first plumbing commit

# That looks pretty familiar! Note that it also has the hash of the tree it points to, as well as timestamps.
```

#### Relating commits to one another


## GoLang - What I need

1. Some basic knowledge about data structures builtin and https://github.com/emirpasic/gods
2. The ability to write command line applications - argument parsing and so on
3. How to interface with the Operating System - something like python's OS module
4. Using sha1 and zlib compression
