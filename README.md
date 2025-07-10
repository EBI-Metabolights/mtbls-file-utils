# MetaboLights File Utils


## Download

Please select appropriate executable (operating system and architecture) on 'build' folder.

- mtbls-compress
- mtbls-rename


## Compress utility usage

mtbls-compress compress all folders in a folder (not recursive) and moves original compressed folders into \<folder path\>_original folder.

```
Usage: mtbls-compress <folder_path> [--include=<folder pattern. e.g., *.d, *.raw, *] [--verbose]

mtbls-compress RAW_FILES --include=*.d

mtbls-compress RAW_FILES --include=*.raw

mtbls-compress RAW_FILES --include=*

mtbls-compress RAW_FILES --include=*.d --verbose

```


## Rename utility usage

Rename utility removed all unexpected characters in file names. Accepted characters 0-9 a-z A-Z -_ . characters. All invalid characters (except + ) will be replaces with \_ character. + character will be replaced with \_PLUS\_.

```

Usage:
  mtbls-rename [--dry=true] <folder_path>
  mtbls-rename [--dry=false] <folder_path>

Flags:
  -dry
        Perform a dry run without renaming (default true)
```

Example usages

```

🔧 Dry run mode: true
tests/invalid_files/z£$%sa+ ddd.&&s → tests/invalid_files/z______sa_PLUS___ddd.____s
tests/invalid_files/a+ ddd (copy)&&s → tests/invalid_files/a_PLUS___ddd____copy______s
🔎 Dry run complete. 2 item(s) would be renamed.
Now run again with --dry=false argument to rename all: mtbls-rename --dry=false <folder_path>

```