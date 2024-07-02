# ghashdeep

If you have a large library of files (music, images, videos, etc.), this utility will help you easily calculate or check previously calculated hash sums of the entire library recursively with a single command.

---

### How to calculate the hash:

If you want to calculate the hash of all files starting from the current directory and recursively for all subdirectories, run this command:
```bash
ghashdeep calculate
```
This command will create a file named "checksum.{md5}" in each file folder, the extension will depend on the checksum algorithm selected, the default is md5.

After running this command, you will see a warning message:
```
All existing checksums will be replaced with new checksums.
Are you sure? [YES/NO]:
```
Because if you already have files with calculated checksums, they will be overwritten. If you have never calculated checksums for these files, then type `YES` in the response and the checksums will be calculated.

For example, if you have a file hierarchy like this:
```bash
$ tree
.
├── a
│   └── 1.txt
├── b
│   └── b1
│       └── 2.txt
└── c
```

After executing the `calculate` command, you will see the following:
```bash
$ tree
.
├── a
│   ├── 1.txt
│   └── checksum.md5
├── b
│   └── b1
│       ├── 2.txt
│       └── checksum.md5
└── c
```

The resulting checksum.md5 files are fully compatible with the md5sum system utility, so you can manually verify checksums if you wish:
```bash
$ md5sum -c checksum.md5
1.txt: OK
```

---

### How to validate a hash:

If you have checksum files, you can just run this command and if everything is fine, you will see the following:
```bash
$ ghashdeep check
Jul  1 23:08:45.596 INF crawler/check.go:202 Success status=GOOD path=/tmp/mytmp/test folder=a duration=249.231µs
Jul  1 23:08:45.597 INF crawler/check.go:202 Success status=GOOD path=/tmp/mytmp/test/b folder=b1 duration=111.057µs
```

#### Invalid checksum:
If some files have changed their checksum, you will see the following message:
```bash
$ ghashdeep check
Jul  1 23:11:28.653 ERR crawler/check.go:169 folder have errors status=BAD path=/tmp/mytmp/test folder=a duration=149.804µs
Jul  1 23:11:28.653 ERR crawler/check.go:186 bad checksum status=BAD file=1.txt
```
This message means that the file named '1.txt' in folder 'a' has an invalid checksum.

#### File not found:
If a file has been deleted or a new file has been added but does not exist in checksum.md5, you will see these errors:
```bash
$ ghashdeep check
Jul  1 23:18:01.588 ERR crawler/check.go:169 folder have errors status=BAD path=/tmp/mytmp/test folder=a duration=104.575µs
Jul  1 23:18:01.588 ERR crawler/check.go:179 no checksum status=BAD file=2.txt
Jul  1 23:18:01.588 ERR crawler/check.go:193 not found status=BAD file=1.txt
```
