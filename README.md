# Volcano Engine Certificate Center certificate management

Shared certificate upload code for the Certimate Volcano Engine deployment
plugins that store certificates in Certificate Center.

Consumers pin this repository as a Git submodule inside their Go module and
import it at that path. They build and test it with `go test ./...` from the
consumer checkout. The consumer's Go module supplies the Certimate protocol
and Volcano Engine SDK dependencies.
