# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

migrate-up:
    migrate -database sqlite3://app.db -path migrations up

migrate-down:
    migrate -database sqlite3://app.db -path migrations down

test:
    go test -v ./...

test-coverage:
    go test -cover ./...

# Add/replace MPL 2.0 header in Go and SQL files
mpl:
    ./add_mpl_header.sh
