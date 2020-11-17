
## Usage

> Note: Karai aims to always compile and run on **Linux** targetting the **AMD64** CPU architecture. Other operating systems and architectures may compile and run this software but should do so expecting some inconsistencies.

**Launch Karai Coordinator**

```
make postgres

make migrate

make karai
```

Deprecated> For optimal transaction speed as coordinator:

> Deprecated> When skipping the write process, you are taking some risk if Karai crashes before you write transactions to disk. You can write your transactions to disk with the `wt` command.

```
./go-karai -coordinator -write=false
```

To place graph objects in a different directory:

```
./go-karai -coordinator -graphDir="/literal/path/to/graph/"
```

**Launch Options**

```
 Usage of ./go-karai:
  -apiport int
        Port to run Karai Coordinator API on. (default 4200)
  -batchDir string
        Path where batched transactions should be saved (default "./graph/batch")
  -chunkSize int
        Number of transactions per batch on disk. (default 100)
  -clean
        Clear all peer certs and graph objects
  -consume
        Consume data from sources.
  -graphDir string
        Path where graph objects should be saved (default "./graph")
  -matrix
        Enable Matrix functions. Requires -matrixToken, -matrixURL, and -matrixRoomID
  -matrixRoomID string
        Room ID for matrix publishd events
  -matrixToken string
        Matrix homeserver token
  -matrixURL string
        Matrix homeserver URL
  -write
        Write each graph object to disk. (default true)
```

> Type `menu` to view a list of functions. Functions that are darkened are disabled.

## Dependencies

-   Golang 1.14+ https://golang.org
-   Docker https://www.docker.com/
-   Make http://www.gnu.org/software/make/
-   Migrate https://github.com/golang-migrate/migrate

## Operating System

Karai supports Linux on AMD64 architecture, but may compile in other settings. Differences between Linux and non-Linux installs should be expected.

**Optional:** Compile with all errors displayed, then run binary. Avoids "too many errors" from hiding error info.

`go build -gcflags="-e" && ./pythia`

## Contributing

This repo only receives stable version release updates, development happens in a private repo. Please make an issue before writing code for a PR.

-   MIT License
-   `gofmt`
-   go modules
-   stdlib > \*

## Thanks to:

[![turtlecoin](https://user-images.githubusercontent.com/34389545/80266529-fb0b6880-8661-11ea-9a75-4cb066834775.png)](https://turtlecoin.lol)
[![IPFS](https://user-images.githubusercontent.com/34389545/80266356-0c07aa00-8661-11ea-8308-84639318213a.png)](https://ipfs.io)
[![LibP2P](https://user-images.githubusercontent.com/34389545/80266502-e4651180-8661-11ea-8367-54bf59e26470.png)](https://libp2p.io)
[![GOLANG](https://user-images.githubusercontent.com/34389545/80266422-6b65ba00-8661-11ea-836a-d1904ec15b94.png)](https://golang.org)
