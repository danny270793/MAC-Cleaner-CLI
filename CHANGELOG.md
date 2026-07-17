# Changelog

## [0.2.0](https://github.com/danny270793/MAC-Cleaner-CLI/compare/v0.1.0...v0.2.0) (2026-07-17)


### Features

* **install:** add install.sh to download and install the latest release ([416035f](https://github.com/danny270793/MAC-Cleaner-CLI/commit/416035fb170645ef65a39b5d0d8e81de6b120092))

## [0.1.0](https://github.com/danny270793/MAC-Cleaner-CLI/compare/v0.0.1...v0.1.0) (2026-07-16)


### Features

* **cargocache:** add Cargo registry cache cleaner ([dd4b7c8](https://github.com/danny270793/MAC-Cleaner-CLI/commit/dd4b7c8d0ba13b70acb3f5164815f8799b54a6a9))
* **cleaner:** add Cleaner interface ([ac8072a](https://github.com/danny270793/MAC-Cleaner-CLI/commit/ac8072a1e5be1a81fce43e6faa3d44f687689a8d))
* **cleaner:** add Name() to Cleaner interface ([4412865](https://github.com/danny270793/MAC-Cleaner-CLI/commit/4412865114255128de2bae55d6e03ef6638743ed))
* **cleaner:** add Size() to Cleaner interface ([dcf432a](https://github.com/danny270793/MAC-Cleaner-CLI/commit/dcf432aeade039eba7db79a6465f0191d78e5332))
* **coresimulatorcaches:** add CoreSimulator caches cleaner ([d5ca9c6](https://github.com/danny270793/MAC-Cleaner-CLI/commit/d5ca9c6e78e1bb2daadaefb8eb1dc9377636ff4e))
* **docker:** extract Docker cleaner into its own file ([9c456a2](https://github.com/danny270793/MAC-Cleaner-CLI/commit/9c456a25875f576fdf92b0c1338a14503902b116))
* **docker:** implement Name() ([d0dba00](https://github.com/danny270793/MAC-Cleaner-CLI/commit/d0dba00a4cef24340e194cbad6893258a27a5631))
* **docker:** implement Size() as unmeasurable ([4532bd6](https://github.com/danny270793/MAC-Cleaner-CLI/commit/4532bd6483db1bbe9757eb2c83432ad18cc15609))
* **gomodcache:** add Go module cache cleaner via go clean -modcache ([3c43f72](https://github.com/danny270793/MAC-Cleaner-CLI/commit/3c43f725b257491e50832ee200497a4092e0eb85))
* **gradle:** extract Gradle cleaner into its own file ([fbc01af](https://github.com/danny270793/MAC-Cleaner-CLI/commit/fbc01af4e89ea774aec6a787028d6e96d1852623))
* **gradle:** implement Name() ([722b2a1](https://github.com/danny270793/MAC-Cleaner-CLI/commit/722b2a13900d4fd123558eab5ee86134cc5ed84b))
* **gradle:** implement Size() and share paths with Clean() ([5ab6829](https://github.com/danny270793/MAC-Cleaner-CLI/commit/5ab682976b9f4ccc44165d69433f59728c63120c))
* **librarycaches:** extract Library Caches cleaner into its own file ([dad779d](https://github.com/danny270793/MAC-Cleaner-CLI/commit/dad779d91553cd01f736d998f0fc548e960ac984))
* **librarycaches:** implement Name() ([d23d14f](https://github.com/danny270793/MAC-Cleaner-CLI/commit/d23d14f223bc5288658a044736712372c0effa89))
* **librarycaches:** implement Size() and share path with Clean() ([01eceaf](https://github.com/danny270793/MAC-Cleaner-CLI/commit/01eceaf657aabcf5276b259b9c90d4ed3cf3cefa))
* **main:** add --docker --gradle --library-caches --pub-cache --all flags ([f0b8ab2](https://github.com/danny270793/MAC-Cleaner-CLI/commit/f0b8ab2c63283959f2d5c2dc2f7d6b1815cf4761))
* **main:** add --version flag ([393a3c4](https://github.com/danny270793/MAC-Cleaner-CLI/commit/393a3c4d42d57543d5634a1cace0dabd5d98e543))
* **main:** add friendly --help usage message ([df03d37](https://github.com/danny270793/MAC-Cleaner-CLI/commit/df03d371dcd0cbb82e0643b53084c76146393296))
* **main:** pass measured size to pending output ([38fd9c6](https://github.com/danny270793/MAC-Cleaner-CLI/commit/38fd9c6720ca69b9611daeb1506d1b8caf0fe08d))
* **main:** print checklist-style progress per cleaner ([d248eb4](https://github.com/danny270793/MAC-Cleaner-CLI/commit/d248eb46f116d2008ac87649746bd30b1f0f86c6))
* **main:** wire up --vscode-extensions flag ([dc51637](https://github.com/danny270793/MAC-Cleaner-CLI/commit/dc51637f196b211f7fba9febf42aeea81e0f14ea))
* **main:** wire up new cleaner flags ([3a99f42](https://github.com/danny270793/MAC-Cleaner-CLI/commit/3a99f42901e353032c26483bfd7b572a45bc573b))
* **npmcache:** add npm cache cleaner ([33bb186](https://github.com/danny270793/MAC-Cleaner-CLI/commit/33bb186f3b442933b370eb75d29d7617935f02b6))
* **pnpmstore:** add pnpm store cleaner via pnpm store prune ([bb9e37b](https://github.com/danny270793/MAC-Cleaner-CLI/commit/bb9e37bc0ff6c719f39e14185b3f9e61b4df008b))
* **pubcache:** extract Pub Cache cleaner into its own file ([96c83f7](https://github.com/danny270793/MAC-Cleaner-CLI/commit/96c83f7515a11be3fcea00ab86777206b01108b2))
* **pubcache:** implement Name() ([9256058](https://github.com/danny270793/MAC-Cleaner-CLI/commit/92560589c46e5f6c089b79da76e1ade3f270e3de))
* **pubcache:** implement Size() and share path with Clean() ([8254204](https://github.com/danny270793/MAC-Cleaner-CLI/commit/8254204badd2db0edc0b53affd4097d25eb6f53e))
* **ui:** add colored pending/done output helpers ([fcf4d7b](https://github.com/danny270793/MAC-Cleaner-CLI/commit/fcf4d7b821842c59a39f0e8420cd22ee77ff45fd))
* **ui:** show measured size in pending output ([a6d58ab](https://github.com/danny270793/MAC-Cleaner-CLI/commit/a6d58ab373e8fbc98da0d2794f3401697908d866))
* **util:** add dirSize and sizeOfPaths helpers ([91ae964](https://github.com/danny270793/MAC-Cleaner-CLI/commit/91ae96499a4de2e6768c91fc977a031b94ef60ba))
* **util:** add removeContents helper for clearing directory contents ([13500bf](https://github.com/danny270793/MAC-Cleaner-CLI/commit/13500bf2a75c27b0c2cb796f1bd22dd02cc25574))
* **version:** add version variable for ldflags injection ([913c61f](https://github.com/danny270793/MAC-Cleaner-CLI/commit/913c61f44ce1042e97c70002bc0d2ed712d90dab))
* **vscodeextensions:** remove outdated vscode extension versions ([c3fda5a](https://github.com/danny270793/MAC-Cleaner-CLI/commit/c3fda5a8659f06881d9533d4ed50c1f2b90402af))
* **xcodederiveddata:** add Xcode DerivedData cleaner ([1b216d6](https://github.com/danny270793/MAC-Cleaner-CLI/commit/1b216d6f93383a816c0e977052c913925dd1b359))
