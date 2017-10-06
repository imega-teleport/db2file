package packer

import (
	"fmt"
	"testing"
	"time"

	"github.com/imega-teleport/db2file/indexer"
	"github.com/imega-teleport/db2file/storage"
	"github.com/stretchr/testify/assert"
)

func offTest_PackProductImage(t *testing.T) {
	p := pkg{
		Indexer: indexer.NewIndexer(),
		Options: Options{
			PathToImages: "/go/src/github.com/imega-teleport/db2file",
		},
	}

	item := storage.ProductImage{
		ProductID: "2d551104-b0ef-11e3-91a0-94de8026f172",
		URL:       "imager/teleport.png",
	}

	err := p.productImagePack(item)
	assert.NoError(t, err)

	err = p.ThirdPackToContent(true)
	assert.NoError(t, err)

	expected := fmt.Sprintf("set @max_post_id=(select ifnull(max(id),0)from posts);"+
		"set @2d551104b0ef11e391a094de8026f172_img=1;"+
		"INSERT INTO posts (id,post_author,post_date,post_date_gmt,post_content,post_title,post_excerpt,post_status,post_name,post_modified,post_modified_gmt,post_parent,post_type,post_mime_type) "+
		"VALUES (@max_post_id+@2d551104b0ef11e391a094de8026f172_img,"+
		"1,"+
		"'%s',"+
		"'%s',"+
		"'',"+
		"'/go/src/github.com/imega-teleport/db2file/imager/teleport.png',"+
		"'',"+
		"'inherit',"+
		"'/go/src/github.com/imega-teleport/db2file/imager/teleport.png',"+
		"'%s',"+
		"'%s',"+
		"@2d551104b0ef11e391a094de8026f172,"+
		"'attachment',"+
		"'image/png');"+
		//teleport_item
		"set @2d551104b0ef11e391a094de8026f172_img=(select id from teleport_item where guid='2d551104b0ef11e391a094de8026f172_img');"+
		//PostMeta
		"INSERT INTO postmeta (post_id,meta_key,meta_value) "+
		"VALUES (@2d551104b0ef11e391a094de8026f172_img,"+
		"'_wp_attachment_metadata',"+
		"'a:4:{s:5:\"width\";i:200;s:6:\"height\";i:200;s:4:\"file\";s:61:\"/go/src/github.com/imega-teleport/db2file/imager/teleport.png\";s:10:\"image_meta\";a:11:{s:8:\"aperture\";i:0;s:6:\"credit\";s:0:\"\";s:6:\"camera\";s:0:\"\";s:7:\"caption\";s:0:\"\";s:17:\"created_timestamp\";i:0;s:9:\"copyright\";s:0:\"\";s:12:\"focal_length\";i:0;s:3:\"iso\";i:0;s:13:\"shutter_speed\";i:0;s:5:\"title\";s:0:\"\";s:11:\"orientation\";i:0;}}');",
		time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	assert.Equal(t, expected, p.Content)
}
