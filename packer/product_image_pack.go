package packer

import (
	"fmt"
	"time"

	"github.com/imega-teleport/db2file/imager"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/db2file/teleport"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

func (p *pkg) productImagePack(item storage.ProductImage) error {
	info, err := imager.GetImageInfo(fmt.Sprintf("%s/%s", p.Options.PathToImages, item.URL))
	if err != nil {
		return err
	}
	p.ThirdPack.AddItem(teleport.Post{
		ID:       teleport.UUID(fmt.Sprintf("%s_img_%s", item.ProductID, item.EntityID)),
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

	s, err := serializationAttachment(info)
	if err != nil {
		return err
	}
	p.ThirdPack.AddItem(teleport.PostMeta{
		PostID: teleport.UUID(fmt.Sprintf("%s_img", item.ProductID)),
		Key:    "_wp_attachment_metadata",
		Value:  s,
	})

	return nil
}

func serializationAttachment(i imager.ImageInfo) (string, error) {
	encoder := php_serialize.NewSerializer()
	source := map[php_serialize.PhpValue]php_serialize.PhpValue{}

	source = map[php_serialize.PhpValue]php_serialize.PhpValue{
		"width":  i.Width,
		"height": i.Height,
		"file":   i.Name,
		"image_meta": map[php_serialize.PhpValue]php_serialize.PhpValue{
			"aperture":          0,
			"credit":            "",
			"camera":            "",
			"caption":           "",
			"created_timestamp": 0,
			"copyright":         "",
			"focal_length":      0,
			"iso":               0,
			"shutter_speed":     0,
			"title":             "",
			"orientation":       0,
		},
	}
	return encoder.Encode(source)
}
