# Renamer

_Renaming media assets by hand so they can play well with Plex media server drove me crazy._

## TLDR

```
$ ls
Billions.1-1.mp4 Billions.1-3.mp4

$ cat source
Billions.1-1.mp4
Billions.1-3.mp4

$ cat target
Billions - S01E01.mp4
Billions - S01E03.mp4

$ go install
$ renamer source target

$ ls
Billions - S01E01.mp4 Billions - S01E03.mp4
```

## A bit more explanation

This little program takes two parameters: `source` and `target`. `source` is a plain text file containing all the **original** file names. `target` is the plain text file containing the file names you want to change them into.

Copy your media asset files and paste them into a text file, and you will have `source`. To get `target`, use whatever text editor (I use Sublime Text) to massage the content of source.

This program does some basic sanity check, it makes sure
- The file names in `source` does exist
- There are same number of entries in `source` and `target`