# QuickBooks Online XSD schema

These files are Intuit's official XML schema for the QuickBooks Online v3 API. They are copied unmodified from
[intuit/QuickBooks-V3-Java-SDK](https://github.com/intuit/QuickBooks-V3-Java-SDK/tree/develop/ipp-v3-java-data/src/main/xsd)
(Apache-2.0, see `LICENSE` in this directory), commit `865ec27967b64db818a4b330f040b361aae81153`. `EntitlementsResponse.xsd` is left out: nothing
includes it, and its header marks it confidential.

Entity types (Account, Invoice, Line, ...) are defined in `Finance.xsd`. Shared base types are in `IntuitBaseTypes.xsd`.

To update, re-download the files from the `develop` branch and update the commit above.
