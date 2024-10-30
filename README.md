<h1 align="center">
  how many sxg sites?
</h1>

<p align="center">
  how many of the top million sites use cloudflare sxg
</p>

<div align="center">
  <a href="https://x.com/cryptograthor">
    <img src="https://img.shields.io/badge/made_by_cryptograthor-black?style=flat&logo=undertale&logoColor=hotpink" />
    <!-- ![](https://img.shields.io/badge/made_by_cryptograthor-black?style=flat&logo=undertale&logoColor=hotpink) -->
  </a>
  </div>

- Top million domain ranking obtained from [cloudflare](https://radar.cloudflare.com/domains).
- `dump-signedexchange` obtained from [webpackage](https://github.com/WICG/webpackage/blob/main/go/signedexchange/cmd/dump-signedexchange/main.go)

Usage:

```sh
# this will run for awhile, several minutes
go run main.go

# prints the number of sites using SXG plus 1
ls -l | wc -l 
```
