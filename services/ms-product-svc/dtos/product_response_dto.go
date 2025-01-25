package dtos

type ProductResponseDTO struct {
	ProductId        string `json:"productId"`        // ID produk (bebas digunakan)
	Name             string `json:"name"`             // Nama produk
	Category         string `json:"category"`         // Kategori produk
	Qty              int    `json:"qty"`              // Jumlah produk
	Price            int    `json:"price"`            // Harga produk
	Sku              string `json:"sku"`              // SKU produk
	FileId           string `json:"fileId"`           // ID file terkait
	FileUri          string `json:"fileUri"`          // URI file terkait
	FileThumbnailUri string `json:"fileThumbnailUri"` // URI thumbnail file terkait
	CreatedAt        string `json:"createdAt"`        // Waktu pembuatan
	UpdatedAt        string `json:"updatedAt"`        // Waktu pembaruan
}
