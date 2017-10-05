package packer

import (
	"fmt"
	"time"

	"github.com/imega-teleport/db2file/imager"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/db2file/teleport"
)

func (p *pkg) productImagePack(item storage.ProductImage) error {
	info, err := imager.GetImageInfo(fmt.Sprintf("%s/%s", p.Options.PathToImages, item.URL))
	if err != nil {
		return err
	}
	p.ThirdPack.AddItem(teleport.Post{
		ID:       teleport.UUID(fmt.Sprintf("%s_img", item.ProductID)),
		AuthorID: 1,
		Date:     time.Now(),
		Title:    info.Name,
		Excerpt:  "",
		Status:   "inherit",
		Name:     info.Name,
		Modified: time.Now(),
		ParentID: teleport.UUID(item.ProductID),
		Type:     "attachment",
		MimeType: info.Mime,
	})
	return nil
}
