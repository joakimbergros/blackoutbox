package documents

import "blackoutbox/internal/stores"

type DocumentHandler struct {
	Store stores.DocumentStoreInterface
}
