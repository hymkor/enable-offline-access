enable-offline-access
=====================

`enable-offline-access` is a command-line tool to preemptively trigger cloud storage downloads  
by opening and closing files under a given path, making them available offline before application access.

Purpose
-------

When using cloud storage services like Dropbox, OneDrive, or Google Drive, files may not be downloaded until the user or an application opens them.  
This can cause delays or errors in applications that try to access many files at once.

`enable-offline-access` walks through the specified paths, opens each file in read-only mode, and immediately closes it.  
This gentle access often triggers the cloud storage service to download the file in advance, ensuring that the file is available when actually needed.

Use Cases
---------

- Avoid timeouts in applications that batch-open many files.
- Warm up cloud-synced folders before launching a media manager, IDE, or build tool.
- Prepare offline availability on unstable or slow networks.

Usage
-----

```sh
enable-offline-access [OPTIONS] <path1> [<path2> ...]
```

### Example

```sh
enable-offline-access -c 8 ~/Dropbox/project/
```

This opens files under `~/Dropbox/project/` with a concurrency of 8, encouraging Dropbox to fetch them locally.

Options
-------

| Option   | Description                                              | Default |
| -------- | -------------------------------------------------------- | ------- |
| `-c <n>` | Number of files to open concurrently (recommended: 4–16) | 8       |

Notes
-----

* Directories are traversed recursively.
* Only regular files are touched — symbolic links, sockets, and directories are skipped.
* Files are **not modified**, only opened and closed.

Compatibility
-------------

`enable-offline-access` is designed to work with:

* **Dropbox**
* **OneDrive**
* **Google Drive**
* And other cloud storage systems that download files on access.

However, behavior may vary depending on the client software and platform.
No guarantees are made that a file will always be downloaded; this tool **suggests** access.

Installation
------------

If you have Go installed:

```sh
go install github.com/yourusername/enable-offline-access@latest
```

Benchmark (Reference)
---------------------

Here is a rough performance comparison based on a directory containing **735 files**.

| Condition                          | Option               | Elapsed Time |
| ---------------------------------- | -------------------- | ------------ |
| All files already online           | *(default = `-c 8`)* | \~91 ms      |
| All files offline, download needed | *(default = `-c 8`)* | \~1m 7.5s    |
| All files offline, download needed | `-c 16`              | \~58.4s      |

* The `-c` option controls the maximum number of files opened simultaneously.
  The default is `-c 8`.
* The difference between the first two cases reflects the state of cloud synchronization, not the concurrency setting.
* Increasing concurrency beyond 8 provides some improvement, but the effect is not linear.
* Performance may vary depending on your cloud storage backend, network speed, and local disk performance.

Author
------

[hymkor (HAYAMA Kaoru)](https://github.com/hymkor)

License
-------

MIT License
