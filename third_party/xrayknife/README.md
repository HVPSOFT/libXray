This directory contains a minimal vendored subset of `xray-knife` used for share-link parsing.

Source:
- https://github.com/HVPSOFT/xray-knife

Notes:
- The copied code is based on the v9 codebase and keeps the original MIT license in [LICENSE](/Users/vlad/GitHub/libXray/third_party/xrayknife/LICENSE).
- libXray vendors only the packages needed to parse share links into Xray outbound configs.
- A small compatibility patch was applied so the code builds with the newer `github.com/xtls/xray-core` used by this repository.
