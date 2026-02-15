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

# Add/replace MPL 2.0 header in source files
mpl:
    #!/usr/bin/env bash
    set -euo pipefail

    # Process files with sed
    find . -type f \( \
        -name "*.go" -o \
        -name "*.rs" -o \
        -name "*.cs" \
    \) -exec sed -i '1i\
    // This Source Code Form is subject to the terms of the Mozilla Public\
    // License, v. 2.0. If a copy of the MPL was not distributed with this\
    // file, You can obtain one at https://mozilla.org/MPL/2.0/.
    ' {} \;

        find . -type f -name "*.py" -exec sed -i '1i\
    # This Source Code Form is subject to the terms of the Mozilla Public\
    # License, v. 2.0. If a copy of the MPL was not distributed with this\
    # file, You can obtain one at https://mozilla.org/MPL/2.0/.
    ' {} \;

        find . -type f -name "*.php" -exec sed -i '1i\
    /* This Source Code Form is subject to the terms of the Mozilla Public\
    * License, v. 2.0. If a copy of the MPL was not distributed with this\
    * file, You can obtain one at https://mozilla.org/MPL/2.0/. */
    ' {} \;
